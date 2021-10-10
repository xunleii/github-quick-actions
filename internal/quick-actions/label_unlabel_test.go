package quick_actions

import (
	"fmt"
	"testing"

	"github.com/cucumber/godog"
	"github.com/stretchr/testify/assert"

	. "xnku.be/github-quick-actions/pkg/gh_quick_action/v2"
	gqa_scenario_context "xnku.be/github-quick-actions/pkg/ghk_scenario_ctx"
)

func TestLabelsHelper_getLabels(t *testing.T) {
	ts := map[string]struct {
		command EventCommand
		labels  []string
	}{
		"simple label": {
			command: EventCommand{Arguments: []string{"label"}},
			labels:  []string{"label"},
		},
		"tilded label": {
			command: EventCommand{Arguments: []string{"~label"}},
			labels:  []string{"label"},
		},
		"multi labels": {
			command: EventCommand{Arguments: []string{"label1", "label2"}},
			labels:  []string{"label1", "label2"},
		},
		"mixed labels": {
			command: EventCommand{Arguments: []string{"label1", "~label2"}},
			labels:  []string{"label1", "label2"},
		},
		"duplicated labels": {
			command: EventCommand{Arguments: []string{"label1", "label1"}},
			labels:  []string{"label1"},
		},
		"duplicated mixed labels": {
			command: EventCommand{Arguments: []string{"~label1", "label1"}},
			labels:  []string{"label1"},
		},
		"invalid label": {
			command: EventCommand{Arguments: []string{"~"}},
			labels:  []string{},
		},
		"empty label": {
			command: EventCommand{Arguments: []string{""}},
			labels:  []string{},
		},
	}

	for name, tc := range ts {
		t.Run(name, func(t *testing.T) {
			command := tc.command

			labels := labelsHelper{}.getLabels(&command)
			assert.ElementsMatch(t, tc.labels, labels)

		})
	}
}

func TestLabel_TriggerOnEvents(t *testing.T) {
	assert.ElementsMatch(t,
		[]EventType{EventTypeIssue, EventTypeIssueComment, EventTypePullRequest, EventTypePullRequestReviewComment},
		LabelQuickAction{}.TriggerOnEvents(),
	)
}

func TestLabelFeature(t *testing.T) {
	events := LabelQuickAction{}.TriggerOnEvents()

	for _, event := range events {
		t.Run(string(event), func(t *testing.T) {
			suite := godog.TestSuite{
				ScenarioInitializer: gqa_scenario_context.ScenarioInitializer(map[string]QuickAction{"label": &LabelQuickAction{}}),
				Options: &godog.Options{
					Format:   "pretty",
					Paths:    []string{"ghk::features"},
					Tags:     fmt.Sprintf("label && %s", event),
					TestingT: t,
				},
			}

			if suite.Run() != 0 {
				t.Fatal("non-zero status returned, failed to run feature tests")
			}
		})
	}
}

func TestUnlabel_TriggerOnEvents(t *testing.T) {
	assert.ElementsMatch(t,
		[]EventType{EventTypeIssueComment, EventTypePullRequestReviewComment},
		UnlabelQuickAction{}.TriggerOnEvents(),
	)
}

func TestUnlabelFeature(t *testing.T) {
	suite := godog.TestSuite{
		ScenarioInitializer: gqa_scenario_context.ScenarioInitializer(map[string]QuickAction{"unlabel": &UnlabelQuickAction{}}),
		Options: &godog.Options{
			Format:   "pretty",
			Paths:    []string{"ghk::features"},
			Tags:     "unlabel",
			TestingT: t,
		},
	}

	if suite.Run() != 0 {
		t.Fatal("non-zero status returned, failed to run feature tests")
	}
}
