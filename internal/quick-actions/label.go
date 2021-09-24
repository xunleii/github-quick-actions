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

var LabelIssueComment = v1.GithubQuickAction{
	OnEvent: v1.GithubIssueCommentEvent,
	Handler: Label,
}

func Label(ctx context.Context, cc githubapp.ClientCreator, event v1.GithubQuickActionEvent) error {
	logger := zerolog.Ctx(ctx).With().
		Str("quick_action", "label").
		Logger()

	issueEvent, valid := event.Payload.(*github.IssueCommentEvent)
	if !valid {
		return fmt.Errorf("invalid event type; only accept %T but get %T", issueEvent, event)
	}

	logger.Debug().Msgf("handle `/label` (args: %v)", event.Arguments)

	var labels []string
	for _, label := range event.Arguments {
		switch {
		case label == "~" || label == "":
		// ignore empty labels
		case strings.HasPrefix(label, "~"):
			labels = append(labels, label[1:])
		default:
			labels = append(labels, label)
		}
	}

	if len(labels) == 0 {
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

	_, _, err = client.Issues.AddLabelsToIssue(ctx, owner, repo, noIssue, funk.UniqString(labels))

	return err
}
