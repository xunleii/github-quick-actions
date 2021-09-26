package quick_actions

import (
	"strings"

	"github.com/google/go-github/v39/github"
	"github.com/rs/zerolog"
	"github.com/thoas/go-funk"

	. "xnku.be/github-quick-actions/pkg/gh_quick_action/v2"
)

type (
	assigneesHelper struct{ githubEventHelper }

	// AssignQuickAction implements QuickAction interface for /assign command.
	// This quick action adds one or several assignees to an issue or a PR.
	AssignQuickAction struct{ assigneesHelper }
	// UnassignQuickAction implements QuickAction interface for /unassign command.
	// This quick action removes one or several assignees to an issue or a PR.
	UnassignQuickAction struct{ assigneesHelper }
)

func (qa AssignQuickAction) HandleCommand(ctx *EventContext, command *EventCommand) error {
	logger := zerolog.Ctx(ctx).With().
		Str("quick_action", "assign").
		Logger()

	logger.Info().Msgf("handle `/assign` (args: %v)", command.Arguments)

	assignees := qa.getAssignees(command)
	if len(assignees) == 0 {
		logger.Debug().Msgf("no assignees found; ignored")
		return nil
	}

	client, err := qa.newInstallationClient(ctx, command.Payload)
	if err != nil {
		return err
	}

	_, _, err = client.Issues.AddAssignees(
		ctx,
		command.Payload.RepositoryOwner(),
		command.Payload.RepositoryName(),
		command.Payload.IssueNumber(),
		assignees,
	)

	return err
}

func (qa UnassignQuickAction) HandleCommand(ctx *EventContext, command *EventCommand) error {
	logger := zerolog.Ctx(ctx).With().
		Str("quick_action", "unassign").
		Logger()

	logger.Info().Msgf("handle `/unassign` (args: %v)", command.Arguments)

	assignees := qa.getAssignees(command)
	if len(assignees) == 0 {
		logger.Debug().Msgf("no assignees found; ignored")
		return nil
	}

	client, err := qa.newInstallationClient(ctx, command.Payload)
	if err != nil {
		return err
	}

	_, _, err = client.Issues.RemoveAssignees(
		ctx,
		command.Payload.RepositoryOwner(),
		command.Payload.RepositoryName(),
		command.Payload.IssueNumber(),
		assignees,
	)

	return err
}

func (assigneesHelper) TriggerOnEvents() []EventType { return []EventType{EventTypeIssueComment} }
func (assigneesHelper) getAssignees(command *EventCommand) []string {
	var author string
	switch event := command.Payload.Raw().(type) {
	case *github.IssueCommentEvent:
		author = event.GetComment().GetUser().GetLogin()
	}

	var assignees []string
	for _, assignee := range command.Arguments {
		switch {
		case assignee == "@" || assignee == "":
			// ignore empty assignees
		case assignee == "me" && author != "":
			assignees = append(assignees, author)
		case strings.HasPrefix(assignee, "@"):
			assignees = append(assignees, assignee[1:])
		}
	}

	return funk.UniqString(assignees)
}

func init() {
	// NOTE: register quick actions
	registerQuickAction("assign", &AssignQuickAction{})
	registerQuickAction("unassign", &UnassignQuickAction{})
}
