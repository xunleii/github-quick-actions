package quick_actions

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/google/go-github/v39/github"
	"github.com/hashicorp/go-multierror"
	"github.com/rs/zerolog"
	. "xnku.be/github-quick-actions/pkg/gh_quick_action/v2"
)

type (
	// DuplicateQuickAction implements QuickAction interface for /duplicate command.
	// This quick action closes an issue or a PR and marks as duplicate of
	// another issue.
	DuplicateQuickAction struct{ githubEventHelper }
)

func (qa DuplicateQuickAction) TriggerOnEvents() []EventType {
	// NOTE: assign should be triggered on issues & pull requests description
	return []EventType{EventTypeIssueComment, EventTypePullRequestReviewComment}
}

func (qa DuplicateQuickAction) HandleCommand(ctx *EventContext, command *EventCommand) error {
	logger := zerolog.Ctx(ctx).With().
		Str("quick_action", "duplicate").
		Logger()

	logger.Info().Msgf("handle `/duplicate` (args: %v)", command.Arguments)

	var issues []int
	for _, issue := range command.Arguments {
		if !strings.HasPrefix(issue, "#") {
			logger.Debug().Msgf("invalid issue '%s' provided; ignored", issue)
			continue
		}

		n, err := strconv.Atoi(strings.TrimPrefix(issue, "#"))
		if err != nil {
			logger.Error().Err(err).Msgf("invalid issue '%s' provided; ignored", issue)
			continue
		}
		issues = append(issues, n)
	}

	if len(issues) == 0 {
		logger.Error().Msgf("no valid issue provided; ignored")
		return nil
	}

	client, err := qa.newInstallationClient(ctx, command.Payload)
	if err != nil {
		return err
	}

	var errs *multierror.Error
	for _, issue := range issues {
		if issue == command.Payload.IssueNumber() {
			logger.Debug().Msgf("cannot mark current issue as duplicate of itself; ignored")
			continue
		}

		_, _, err := client.Issues.Get(
			ctx,
			command.Payload.RepositoryOwner(),
			command.Payload.RepositoryName(),
			issue,
		)
		if err != nil {
			// NOTE: invalid issue are ignored
			continue
		}

		_, _, err = client.Issues.CreateComment(
			ctx,
			command.Payload.RepositoryOwner(),
			command.Payload.RepositoryName(),
			command.Payload.IssueNumber(),
			&github.IssueComment{Body: github.String(fmt.Sprintf("Duplicate of #%d", issue))},
		)
		errs = multierror.Append(errs, err)
	}

	return errs.ErrorOrNil()
}

func init() {
	// NOTE: register quick actions
	registerQuickAction("duplicate", &DuplicateQuickAction{})
}
