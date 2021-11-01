package quick_actions

import (
	"fmt"
	"testing"

	"github.com/cucumber/godog"
	"github.com/stretchr/testify/assert"
	. "xnku.be/github-quick-actions/pkg/gh_quick_action/v2"
	gqa_scenario_context "xnku.be/github-quick-actions/pkg/ghk_scenario_ctx"
)

func TestDuplicateQuickAction_TriggerOnEvents(t *testing.T) {
	assert.ElementsMatch(t,
		[]EventType{EventTypeIssue, EventTypeIssueComment, EventTypePullRequest},
		DuplicateQuickAction{}.TriggerOnEvents(),
	)
}

func TestDuplicateFeature(t *testing.T) {
	events := DuplicateQuickAction{}.TriggerOnEvents()

	for _, event := range events {
		t.Run(string(event), func(t *testing.T) {
			suite := godog.TestSuite{
				ScenarioInitializer: gqa_scenario_context.ScenarioInitializer(map[string]QuickAction{"duplicate": &DuplicateQuickAction{}}),
				Options: &godog.Options{
					Format:   "pretty",
					Paths:    []string{"ghk::features"},
					Tags:     fmt.Sprintf("duplicate && %s", event),
					TestingT: t,
				},
			}

			if suite.Run() != 0 {
				t.Fatal("non-zero status returned, failed to run feature tests")
			}
		})
	}
}