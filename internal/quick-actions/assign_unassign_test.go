package quick_actions

import (
	"testing"

	"github.com/cucumber/godog"
	"github.com/stretchr/testify/assert"

	. "xnku.be/github-quick-actions/pkg/gh_quick_action/v2"
	gqa_scenario_context "xnku.be/github-quick-actions/pkg/ghk_scenario_ctx"
)

func TestAssigneesHelper_TriggerOnEvents(t *testing.T) {
	assert.ElementsMatch(t, []EventType{EventTypeIssueComment}, assigneesHelper{}.TriggerOnEvents())
}

func TestAssigneesHelper_getAssignees(t *testing.T) {
	noErr := func(e EventPayload, _ error) EventPayload { return e }

	payloads := []EventPayload{
		noErr(PayloadFactory(EventTypeIssueComment, []byte(`{"comment":{"user":{"login":"xunleii"}}}`))),
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

func TestAssignFeature(t *testing.T) {
	suite := godog.TestSuite{
		ScenarioInitializer: gqa_scenario_context.ScenarioInitializer(map[string]QuickAction{"assign": &AssignQuickAction{}}),
		Options: &godog.Options{
			Format:   "pretty",
			Paths:    []string{"ghk::features"},
			Tags:     "assign",
			TestingT: t,
		},
	}

	if suite.Run() != 0 {
		t.Fatal("non-zero status returned, failed to run feature tests")
	}
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
