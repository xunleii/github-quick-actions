package quick_actions

import (
	"context"
	"fmt"
	"strings"

	"github.com/rs/zerolog"

	. "xnku.be/github-quick-actions/pkg/gh-quick-action"
)

func Unassign(ctx context.Context, event GithubQuickActionEvent) error {
	logger := zerolog.Ctx(ctx).With().
		Str("quick_action", "unassign").
		Logger()

	issueEvent, valid := event.(*IssueCommentEvent)
	if !valid {
		return fmt.Errorf("invalid event type; only accept %T but get %T", issueEvent, event)
	}

	logger.Debug().Msgf("handle `/unassign` (args: %v)", issueEvent.Arguments())

	var assignees []string
	for _, assignee := range issueEvent.Arguments() {
		switch {
		case assignee == "@" || assignee == "":
			// ignore empty assignees
		case assignee == "me":
			assignees = append(assignees, issueEvent.GetComment().GetUser().GetLogin())
		case strings.HasPrefix(assignee, "@"):
			assignees = append(assignees, assignee[1:])
		}
	}

	client, err := issueEvent.NewInstallationClient(issueEvent.GetInstallation().GetID())
	if err != nil {
		zerolog.Ctx(ctx).Error().Err(err).Send()
		return err
	}

	owner := issueEvent.GetRepo().GetOwner().GetLogin()
	repo := issueEvent.GetRepo().GetName()
	noIssue := issueEvent.GetIssue().GetNumber()

	_, _, err = client.Issues.RemoveAssignees(ctx, owner, repo, noIssue, assignees)

	return err
}
