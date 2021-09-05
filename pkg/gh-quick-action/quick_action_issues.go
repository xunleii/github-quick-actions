package quick_action

import (
	"context"
	"encoding/json"
	"strings"

	"github.com/google/go-github/v38/github"
	"github.com/hashicorp/go-multierror"
	"github.com/palantir/go-githubapp/githubapp"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
)

// IssueCommentQuickActions implements quick actions handling issues and pull requests comments.
type IssueCommentQuickActions struct {
	githubapp.ClientCreator

	handlers map[string]GithubQuickActionHandler
}

// NewIssueCommentQuickActions creates a new quick action manager handling issues and pull requests comments.
func NewIssueCommentQuickActions(cc githubapp.ClientCreator) *IssueCommentQuickActions {
	return &IssueCommentQuickActions{
		ClientCreator: cc,
		handlers:      map[string]GithubQuickActionHandler{},
	}
}
func (issue IssueCommentQuickActions) Handles() []string { return []string{"issue_comment"} }
func (issue IssueCommentQuickActions) AddQuickAction(command string, handler GithubQuickActionHandler) {
	issue.handlers[command] = handler
}

func (issue IssueCommentQuickActions) Handle(ctx context.Context, eventType, deliveryID string, payload []byte) error {
	var event github.IssueCommentEvent
	if err := json.Unmarshal(payload, &event); err != nil {
		return errors.Wrap(err, "failed to parse issue comment event payload")
	}

	// ignore if is not a new created comment
	if event.GetAction() != "created" {
		return nil
	}

	// update logger context with current event metadata
	logger := *zerolog.Ctx(ctx)
	logger.UpdateContext(func(ctx zerolog.Context) zerolog.Context {
		return ctx.
			Fields(map[string]interface{}{
				githubapp.LogKeyDeliveryID:      deliveryID,
				githubapp.LogKeyEventType:       eventType,
				githubapp.LogKeyRepositoryOwner: event.GetRepo().GetOwner().GetLogin(),
				githubapp.LogKeyRepositoryName:  event.GetRepo().GetName(),
			})
	})

	logger.Info().Msgf("new event %s handled", eventType)
	logger.Trace().
		RawJSON("payload", payload).
		Msgf("new event %s handled", eventType)

	var quickActions [][]string
	for n, line := range strings.Split(*event.Comment.Body, "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			// empty line, ignored
			continue
		}

		if line[0] != '/' {
			// not a command, ignored
			logger.Trace().Msgf("no command on line n°%d, ignored...", n)
			continue
		}

		commandLine := strings.Split(line, " ")
		if commandLine[0] == "/" {
			// empty command, ignored
			logger.Trace().Msgf("line n°%d has an empty command, ignored...", n)
			continue
		}

		action := commandLine[0][1:]
		if _, exists := issue.handlers[action]; !exists {
			logger.Warn().Msgf("quick action '%s' doesn't exists, ignored", action)
			continue
		}

		var args []string
		for _, arg := range strings.Split(line, " ") {
			arg = strings.TrimSpace(arg)
			if arg != "" {
				args = append(args, arg)
			}
		}
		quickActions = append(quickActions, args)
	}
	if len(quickActions) == 0 {
		// no command line found, ignored
		return nil
	}

	var errsGroup = multierror.Group{}
	for _, quickAction := range quickActions {
		name := quickAction[0][1:]
		args := quickAction[1:]

		logger.Debug().Msgf("handle action '%s'", name)
		errsGroup.Go(func() error {
			err := issue.handlers[name](ctx, &IssueCommentEvent{issue.ClientCreator, &event, args})
			if err != nil {
				logger.Err(err).Msgf("failed to run '%s': %s", name, err)
			}
			return err
		})
	}

	return errsGroup.Wait().ErrorOrNil()
}

// IssueCommentEvent implements GithubQuickActionEvent interface for
// issue/pull_request events.
type IssueCommentEvent struct {
	githubapp.ClientCreator
	*github.IssueCommentEvent

	Args []string
}

func (event IssueCommentEvent) Arguments() []string { return event.Args }
