package quick_actions

import (
	"context"
	"fmt"
	"strings"

	"github.com/google/go-github/v39/github"
	"github.com/palantir/go-githubapp/githubapp"
	"github.com/rs/zerolog"
	"github.com/thoas/go-funk"

	quick_action "xnku.be/github-quick-actions/pkg/gh-quick-action"
)

var AssignIssueComment = quick_action.GithubQuickAction{
	OnEvent: quick_action.GithubIssueCommentEvent,
	Handler: Assign,
}

func Assign(ctx context.Context, cc githubapp.ClientCreator, event quick_action.GithubQuickActionEvent) error {
	logger := zerolog.Ctx(ctx).With().
		Str("quick_action", "assign").
		Logger()

	issueEvent, valid := event.Payload.(*github.IssueCommentEvent)
	if !valid {
		return fmt.Errorf("invalid event type; only accept %T but get %T", issueEvent, event)
	}

	logger.Debug().Msgf("handle `/assign` (args: %v)", event.Arguments)

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

	_, _, err = client.Issues.AddAssignees(ctx, owner, repo, noIssue, funk.UniqString(assignees))

	return err
}
