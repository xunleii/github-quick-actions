package quick_actions_test

import (
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	quick_actions "xnku.be/github-quick-actions/internal/quick-actions"
	fixtures2 "xnku.be/github-quick-actions/pkg/gh-quick-action/v1/fixtures"
)

var LabelIssueCommentFixtures = fixtures2.EventFixtures{
	QuickActionName: "label",
	EventGenerator:  fixtures2.IssueCommentEventType,

	Fixtures: []fixtures2.EventFixture{
		{
			Name:      "simple label",
			Arguments: []string{"~feature"},
			APICalls: map[string]fixtures2.APICallMITM{
				"POST https://api.github.com/repos/xunleii/github-quick-actions/issues/1/labels": {
					Response: &http.Response{StatusCode: http.StatusCreated},
					Validate: func(t *testing.T, r *http.Request, payload []byte) {
						labels, err := labelsFromPayload(payload)
						require.NoError(t, err)
						assert.ElementsMatch(t, labels, []string{"feature"})
					},
				},
			},
		},
		{
			Name:      "multi label",
			Arguments: []string{"feature", "~doc"},
			APICalls: map[string]fixtures2.APICallMITM{
				"POST https://api.github.com/repos/xunleii/github-quick-actions/issues/1/labels": {
					Response: &http.Response{StatusCode: http.StatusCreated},
					Validate: func(t *testing.T, r *http.Request, payload []byte) {
						labels, err := labelsFromPayload(payload)
						require.NoError(t, err)
						assert.ElementsMatch(t, labels, []string{"feature", "doc"})
					},
				},
			},
		},
		{
			Name:      "multi label duplication",
			Arguments: []string{"feature", "feature"},
			APICalls: map[string]fixtures2.APICallMITM{
				"POST https://api.github.com/repos/xunleii/github-quick-actions/issues/1/labels": {
					Response: &http.Response{StatusCode: http.StatusCreated},
					Validate: func(t *testing.T, r *http.Request, payload []byte) {
						labels, err := labelsFromPayload(payload)
						require.NoError(t, err)
						assert.ElementsMatch(t, labels, []string{"feature"})
					},
				},
			},
		},

		{
			Name:      "empty label",
			Arguments: []string{""},
			APICalls:  map[string]fixtures2.APICallMITM{},
		},
		{
			Name:      "repository/owner not found",
			Arguments: []string{"feature"},
			APICalls: map[string]fixtures2.APICallMITM{
				"POST https://api.github.com/repos/xunleii/github-quick-actions/issues/1/labels": {
					Response: func() *http.Response {
						req, _ := http.NewRequest(http.MethodPost, "https://api.github.com/repos/xunleii/github-quick-actions/issues/1/labels", nil)
						body := strings.NewReader(`{"message": "Not Found", "documentation_url": "https://docs.github.com/en/rest/reference/issues#add-labels-to-an-issue"}`)
						resp := &http.Response{Request: req, StatusCode: http.StatusNotFound, Body: io.NopCloser(body)}
						return resp
					}(),
					Validate: func(t *testing.T, r *http.Request, payload []byte) {
						labels, err := labelsFromPayload(payload)
						require.NoError(t, err)
						assert.ElementsMatch(t, labels, []string{"feature"})
					},
				},
			},
			ExpectedError: "POST https://api.github.com/repos/xunleii/github-quick-actions/issues/1/labels: 404 Not Found []",
		},
	},
}

func TestLabel(t *testing.T) { LabelIssueCommentFixtures.RunForQuickAction(t, quick_actions.Label) }

// labelsFromPayload extracts label field from JSON payload
func labelsFromPayload(payload []byte) ([]string, error) {
	var labels []string

	err := json.Unmarshal(payload, &labels)
	return labels, err
}
