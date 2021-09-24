package quick_actions

import (
	"context"
	"fmt"
	"strings"

	"github.com/google/go-github/v39/github"
	"github.com/palantir/go-githubapp/githubapp"
	"github.com/rs/zerolog"
	"github.com/thoas/go-funk"

	"xnku.be/github-quick-actions/pkg/gh-quick-action/v1"
)

var UnassignIssueComment = v1.GithubQuickAction{
	OnEvent: v1.GithubIssueCommentEvent,
	Handler: Unassign,
}

func Unassign(ctx context.Context, cc githubapp.ClientCreator, event v1.GithubQuickActionEvent) error {
	logger := zerolog.Ctx(ctx).With().
		Str("quick_action", "unassign").
		Logger()

	issueEvent, valid := event.Payload.(*github.IssueCommentEvent)
	if !valid {
		return fmt.Errorf("invalid event type; only accept %T but get %T", issueEvent, event)
	}

	logger.Debug().Msgf("handle `/unassign` (args: %v)", event.Arguments)

	var assignees []string
	for _, assignee := range event.Arguments {
		switch {
		case assignee == "@" || assignee == "":
			// ignore empty assignees
		case assignee == "me":
			assignees = append(assignees, issueEvent.GetComment().GetUser().GetLogin())
		case strings.HasPrefix(assignee, "@"):
			assignees = append(assignees, assignee[1:])
		}
	}

	if len(assignees) == 0 {
		// do nothing if no assignees
		return nil
	}

	client, err := cc.NewInstallationClient(issueEvent.GetInstallation().GetID())
	if err != nil {
		return err
	}

	owner := issueEvent.GetRepo().GetOwner().GetLogin()
	repo := issueEvent.GetRepo().GetName()
	noIssue := issueEvent.GetIssue().GetNumber()

	_, _, err = client.Issues.RemoveAssignees(ctx, owner, repo, noIssue, funk.UniqString(assignees))

	return err
}
