package quick_actions

import (
	"context"
	"strings"

	"github.com/rs/zerolog"

	. "xnku.be/github-quick-actions/pkg/gh-quick-action"
)

func Assign(ctx context.Context, event GithubQuickActionEvent) error {
	logger := zerolog.Ctx(ctx).With().
		Str("quick_action", "assign").
		Logger()

	logger.Debug().Msgf("handle `/assign` (args: %v)", event.Args)

	var assignees []string
	for _, assignee := range event.Args {
		switch {
		case assignee == "@" || assignee == "":
			// ignore empty assignees
		case assignee == "me":
			assignees = append(assignees, event.GetComment().GetUser().GetLogin())
		case strings.HasPrefix(assignee, "@"):
			assignees = append(assignees, assignee[1:])
		}
	}

	client, err := event.NewInstallationClient(event.GetInstallation().GetID())
	if err != nil {
		zerolog.Ctx(ctx).Error().Err(err).Send()
		return err
	}

	owner := event.GetRepo().GetOwner().GetLogin()
	repo := event.GetRepo().GetName()
	noIssue := event.GetIssue().GetNumber()

	_, _, err = client.Issues.AddAssignees(ctx, owner, repo, noIssue, assignees)

	return err
}
