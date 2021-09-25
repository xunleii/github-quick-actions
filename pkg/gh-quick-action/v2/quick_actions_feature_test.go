package gh_quick_actions_test

import (
	"fmt"
	"strings"
	"testing"
	"unicode"

	"github.com/cucumber/godog"
	"github.com/google/go-github/v39/github"
	"github.com/thoas/go-funk"

	gh_quick_actions "xnku.be/github-quick-actions/pkg/gh_quick_action/v2"
	gqa_scenario_context "xnku.be/github-quick-actions/pkg/ghk_scenario_ctx"
)

func TestGithubQuickActionsFeatures(t *testing.T) {
	suite := godog.TestSuite{
		ScenarioInitializer: gqa_scenario_context.ScenarioInitializer(map[string]gh_quick_actions.QuickAction{"hello_world": &HelloWorldQuickAction{}}),
		Options: &godog.Options{
			Format:   "pretty",
			Paths:    []string{"ghk::features"},
			TestingT: t,
		},
	}

	if suite.Run() != 0 {
		t.Fatal("non-zero status returned, failed to run feature tests")
	}
}

// HelloWorldQuickAction implements /hello_world command that creates labels
// using the given argument and prefixing them with `hello_world@`. If no
// argument are given, it will only create the label `hello_world`
// NOTE: this quick action is only here for example purpose; do not use it
type HelloWorldQuickAction struct{}

func (h HelloWorldQuickAction) TriggerOnEvents() []gh_quick_actions.EventType {
	return []gh_quick_actions.EventType{gh_quick_actions.EventTypeIssueComment}
}

func (h HelloWorldQuickAction) HandleCommand(ctx *gh_quick_actions.EventContext, command *gh_quick_actions.EventCommand) error {
	payload, valid := command.Payload.Raw().(*github.IssueCommentEvent)
	if !valid {
		return fmt.Errorf("invalid event type")
	}

	var labels []string
	for _, argument := range command.Arguments {
		argument = strings.Map(func(r rune) rune {
			switch {
			case unicode.IsSpace(r):
				return '-'
			case unicode.IsDigit(r):
			case unicode.IsLetter(r):
			default:
				return -1
			}
			return r
		}, argument)
		labels = append(labels, "hello_world@"+strings.ToLower(argument))
	}
	if len(labels) == 0 {
		labels = append(labels, "hello_world")
	}

	client, err := ctx.NewInstallationClient(payload.GetInstallation().GetID())
	if err != nil {
		return err
	}

	owner := payload.GetRepo().GetOwner().GetLogin()
	repo := payload.GetRepo().GetName()
	noIssue := payload.GetIssue().GetNumber()

	_, _, err = client.Issues.AddLabelsToIssue(ctx, owner, repo, noIssue, funk.UniqString(labels))

	return err
}
