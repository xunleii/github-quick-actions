package quick_actions

import (
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/rs/zerolog"
	"github.com/thoas/go-funk"

	. "xnku.be/github-quick-actions/pkg/gh_quick_action/v2"
)

type (
	labelsHelper struct{ githubEventHelper }

	// LabelQuickAction implements QuickAction interface for /label command.
	// This quick action adds one or several labels to an issue or a PR.
	LabelQuickAction struct{ labelsHelper }
	// UnlabelQuickAction implements QuickAction interface for /unlabel or /remove_label command.
	// This quick action removes one or several labels to an issue or a PR.
	UnlabelQuickAction struct{ labelsHelper }
)

func (qa LabelQuickAction) TriggerOnEvents() []EventType {
	// NOTE: adding label should be triggered on issues & pull requests description too
	return []EventType{EventTypeIssue, EventTypeIssueComment, EventTypePullRequest, EventTypePullRequestReviewComment}
}
func (qa LabelQuickAction) HandleCommand(ctx *EventContext, command *EventCommand) error {
	logger := zerolog.Ctx(ctx).With().
		Str("quick_action", "label").
		Logger()

	logger.Info().Msgf("handle `/label` (args: %v)", command.Arguments)

	labels := qa.getLabels(command)
	if len(labels) == 0 {
		logger.Debug().Msgf("no labels found; ignored")
		return nil
	}

	client, err := qa.newInstallationClient(ctx, command.Payload)
	if err != nil {
		return err
	}

	_, _, err = client.Issues.AddLabelsToIssue(
		ctx,
		command.Payload.RepositoryOwner(),
		command.Payload.RepositoryName(),
		command.Payload.IssueNumber(),
		labels,
	)

	return err
}

func (qa UnlabelQuickAction) HandleCommand(ctx *EventContext, command *EventCommand) error {
	logger := zerolog.Ctx(ctx).With().
		Str("quick_action", "unlabel").
		Logger()

	logger.Info().Msgf("handle `/unlabel` (args: %v)", command.Arguments)

	labels := qa.getLabels(command)

	client, err := qa.newInstallationClient(ctx, command.Payload)
	if err != nil {
		return err
	}

	if len(labels) > 0 {
		errs := multierror.Group{}

		for _, label := range labels {
			label := label
			errs.Go(func() error {
				_, err := client.Issues.RemoveLabelForIssue(
					ctx,
					command.Payload.RepositoryOwner(),
					command.Payload.RepositoryName(),
					command.Payload.IssueNumber(),
					label,
				)
				return err
			})
		}
		return errs.Wait().ErrorOrNil()
	}

	_, err = client.Issues.RemoveLabelsForIssue(
		ctx,
		command.Payload.RepositoryOwner(),
		command.Payload.RepositoryName(),
		command.Payload.IssueNumber(),
	)
	return err
}

func (labelsHelper) TriggerOnEvents() []EventType {
	// NOTE: all label changes should be triggered on comment
	return []EventType{EventTypeIssueComment, EventTypePullRequestReviewComment}
}
func (labelsHelper) getLabels(command *EventCommand) []string {
	var labels []string
	for _, label := range command.Arguments {
		switch {
		case label == "~" || label == "":
		// ignore empty labels
		case strings.HasPrefix(label, "~"):
			labels = append(labels, label[1:])
		default:
			labels = append(labels, label)
		}
	}

	return funk.UniqString(labels)
}

func init() {
	// NOTE: register quick actions
	registerQuickAction("label", &LabelQuickAction{})
	registerQuickAction("unlabel", &UnlabelQuickAction{})
	registerQuickAction("remove_label", &UnlabelQuickAction{})
}
