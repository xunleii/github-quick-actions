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

// IssueQuickActions implements quick actions handling issues and pull requests comments.
type IssueQuickActions struct {
	githubapp.ClientCreator

	handlers map[string]GithubQuickActionHandler
}

// NewIssueQuickActions creates a new quick action manager handling issues and pull requests comments.
func NewIssueQuickActions(cc githubapp.ClientCreator) *IssueQuickActions {
	return &IssueQuickActions{
		ClientCreator: cc,
		handlers:      map[string]GithubQuickActionHandler{},
	}
}
func (issue IssueQuickActions) Handles() []string { return []string{"issue_comment"} }
func (issue IssueQuickActions) AddQuickAction(command string, handler GithubQuickActionHandler) {
	issue.handlers[command] = handler
}

func (issue IssueQuickActions) Handle(ctx context.Context, eventType, deliveryID string, payload []byte) error {
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
	for _, line := range strings.Split(*event.Comment.Body, "\n") {
		line = strings.TrimSpace(line)
		if line[0] != '/' {
			// not a command, ignored
			logger.Trace().Msgf("line %d not a command, ignored...", line)
			continue
		}

		commandLine := strings.Split(line, " ")
		if commandLine[0] == "/" {
			// empty command, ignored
			logger.Trace().Msgf("line %d is an empty command, ignored...", line)
			continue
		}

		action := commandLine[0][1:]
		if _, exists := issue.handlers[action]; !exists {
			logger.Warn().Msgf("quick action '%s' doesn't exists, ignored", action)
			continue
		}

		quickActions = append(quickActions, strings.Split(line, " "))
	}
	if len(quickActions) == 0 {
		// no command line found, ignored
		return nil
	}

	var errs = &multierror.Group{}
	for _, quickAction := range quickActions {
		name := quickAction[0][1:]
		args := quickAction[1:]

		logger.Debug().Msgf("handle action '%s'", name)
		errs.Go(func() error {
			err := issue.handlers[name](ctx, GithubQuickActionEvent{issue.ClientCreator, event, args})
			if err != nil {
				logger.Err(err).Msgf("failed to run '%s': %s", name, err)
			}
			return err
		})
	}

	return errs.Wait().ErrorOrNil()
}
