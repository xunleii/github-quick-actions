package quick_actions

import (
	"fmt"
	"testing"

	"github.com/cucumber/godog"
	"github.com/stretchr/testify/assert"

	. "xnku.be/github-quick-actions/pkg/gh_quick_action/v2"
	gqa_scenario_context "xnku.be/github-quick-actions/pkg/ghk_scenario_ctx"
)

func TestAssigneesHelper_getAssignees(t *testing.T) {
	noErr := func(e EventPayload, _ error) EventPayload { return e }

	payloads := []EventPayload{
		noErr(PayloadFactory(EventTypeIssue, []byte(`{"issue":{"user":{"login":"xunleii"}}}`))),
		noErr(PayloadFactory(EventTypeIssueComment, []byte(`{"comment":{"user":{"login":"xunleii"}}}`))),
		noErr(PayloadFactory(EventTypePullRequest, []byte(`{"pull_request":{"user":{"login":"xunleii"}}}`))),
		noErr(PayloadFactory(EventTypePullRequestReviewComment, []byte(`{"comment":{"user":{"login":"xunleii"}}}`))),
	}

	ts := map[string]struct {
		command   EventCommand
		assignees []string
	}{
		"single assignee": {
			command:   EventCommand{Arguments: []string{"@mojombo"}},
			assignees: []string{"mojombo"},
		},
		"single invalid assignee": {
			command:   EventCommand{Arguments: []string{"mojombo"}},
			assignees: []string{},
		},
		"multiple assignee": {
			command:   EventCommand{Arguments: []string{"@mojombo", "@defunkt"}},
			assignees: []string{"mojombo", "defunkt"},
		},
		"multiple assignee with invalid": {
			command:   EventCommand{Arguments: []string{"@mojombo", "defunkt"}},
			assignees: []string{"mojombo"},
		},
		"self assignee": {
			command:   EventCommand{Arguments: []string{"me"}},
			assignees: []string{"xunleii"},
		},
		"multiple assignee with self": {
			command:   EventCommand{Arguments: []string{"@mojombo", "me"}},
			assignees: []string{"mojombo", "xunleii"},
		},
		"multiple assignee duplication": {
			command:   EventCommand{Arguments: []string{"@mojombo", "@mojombo"}},
			assignees: []string{"mojombo"},
		},
		"invalid assignee": {
			command:   EventCommand{Arguments: []string{"@"}},
			assignees: []string{},
		},
		"empty assignee": {
			command:   EventCommand{Arguments: []string{""}},
			assignees: []string{},
		},
	}

	for _, payload := range payloads {
		t.Run(string(payload.Type()), func(t *testing.T) {
			for name, tc := range ts {
				t.Run(name, func(t *testing.T) {
					command := tc.command
					command.Payload = payload

					assignees := assigneesHelper{}.getAssignees(&command)
					assert.ElementsMatch(t, tc.assignees, assignees)
				})
			}
		})
	}
}

func TestAssign_TriggerOnEvents(t *testing.T) {
	assert.ElementsMatch(t,
		[]EventType{EventTypeIssue, EventTypeIssueComment, EventTypePullRequest, EventTypePullRequestReviewComment},
		AssignQuickAction{}.TriggerOnEvents(),
	)
}

func TestAssignFeature(t *testing.T) {
	events := AssignQuickAction{}.TriggerOnEvents()

	for _, event := range events {
		t.Run(string(event), func(t *testing.T) {
			suite := godog.TestSuite{
				ScenarioInitializer: gqa_scenario_context.ScenarioInitializer(map[string]QuickAction{"assign": &AssignQuickAction{}}),
				Options: &godog.Options{
					Format:   "pretty",
					Paths:    []string{"ghk::features"},
					Tags:     fmt.Sprintf("assign && %s", event),
					TestingT: t,
				},
			}

			if suite.Run() != 0 {
				t.Fatal("non-zero status returned, failed to run feature tests")
			}
		})
	}
}

func TestUnassign_TriggerOnEvents(t *testing.T) {
	assert.ElementsMatch(t,
		[]EventType{EventTypeIssueComment, EventTypePullRequestReviewComment},
		UnassignQuickAction{}.TriggerOnEvents(),
	)
}

func TestUnassignFeature(t *testing.T) {
	suite := godog.TestSuite{
		ScenarioInitializer: gqa_scenario_context.ScenarioInitializer(map[string]QuickAction{"unassign": &UnassignQuickAction{}}),
		Options: &godog.Options{
			Format:   "pretty",
			Paths:    []string{"ghk::features"},
			Tags:     "unassign",
			TestingT: t,
		},
	}

	if suite.Run() != 0 {
		t.Fatal("non-zero status returned, failed to run feature tests")
	}
}
