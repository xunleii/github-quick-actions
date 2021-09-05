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
	"xnku.be/github-quick-actions/pkg/gh-quick-action/fixtures"
)

var AssignIssueCommentFixtures = fixtures.EventFixtures{
	QuickActionName: "assign",
	EventType:       fixtures.IssueCommentEventType,

	Fixtures: []fixtures.EventFixture{
		{
			Name:      "simple assignment",
			Arguments: []string{"@mojombo"},
			APICalls: map[string]fixtures.APICallMITM{
				"POST https://api.github.com/repos/xunleii/github-quick-actions/issues/1/assignees": {
					Response: &http.Response{StatusCode: http.StatusCreated},
					Validate: func(t *testing.T, r *http.Request, payload []byte) {
						assignees, err := assigneesFromPayload(payload)
						require.NoError(t, err)
						assert.ElementsMatch(t, assignees, []string{"mojombo"})
					},
				},
			},
		},
		{
			Name:      "multi assignment",
			Arguments: []string{"@mojombo", "@defunkt"},
			APICalls: map[string]fixtures.APICallMITM{
				"POST https://api.github.com/repos/xunleii/github-quick-actions/issues/1/assignees": {
					Response: &http.Response{StatusCode: http.StatusCreated},
					Validate: func(t *testing.T, r *http.Request, payload []byte) {
						assignees, err := assigneesFromPayload(payload)
						require.NoError(t, err)
						assert.ElementsMatch(t, assignees, []string{"mojombo", "defunkt"})
					},
				},
			},
		},
		{
			Name:      "self assignment",
			Arguments: []string{"me"},
			APICalls: map[string]fixtures.APICallMITM{
				"POST https://api.github.com/repos/xunleii/github-quick-actions/issues/1/assignees": {
					Response: &http.Response{StatusCode: http.StatusCreated},
					Validate: func(t *testing.T, r *http.Request, payload []byte) {
						assignees, err := assigneesFromPayload(payload)
						require.NoError(t, err)
						assert.ElementsMatch(t, assignees, []string{"xunleii"})
					},
				},
			},
		},
		{
			Name:      "multi assignment with self",
			Arguments: []string{"@mojombo", "@defunkt", "me"},
			APICalls: map[string]fixtures.APICallMITM{
				"POST https://api.github.com/repos/xunleii/github-quick-actions/issues/1/assignees": {
					Response: &http.Response{StatusCode: http.StatusNoContent},
					Validate: func(t *testing.T, r *http.Request, payload []byte) {
						assignees, err := assigneesFromPayload(payload)
						require.NoError(t, err)
						assert.ElementsMatch(t, assignees, []string{"mojombo", "defunkt", "xunleii"})
					},
				},
			},
		},
		{
			Name:      "multi assignment duplication",
			Arguments: []string{"@mojombo", "@mojombo"},
			APICalls: map[string]fixtures.APICallMITM{
				"POST https://api.github.com/repos/xunleii/github-quick-actions/issues/1/assignees": {
					Response: &http.Response{StatusCode: http.StatusCreated},
					Validate: func(t *testing.T, r *http.Request, payload []byte) {
						assignees, err := assigneesFromPayload(payload)
						require.NoError(t, err)
						assert.ElementsMatch(t, assignees, []string{"mojombo"})
					},
				},
			},
		},

		{
			Name:      "repository/owner not found",
			Arguments: []string{"@mojombo"},
			APICalls: map[string]fixtures.APICallMITM{
				"POST https://api.github.com/repos/xunleii/github-quick-actions/issues/1/assignees": {
					Response: func() *http.Response {
						req, _ := http.NewRequest(http.MethodPost, "https://api.github.com/repos/xunleii/github-quick-actions/issues/1/assignees", nil)
						body := strings.NewReader(`{"message": "Not Found", "documentation_url": "https://docs.github.com/rest/reference/issues#add-assignees-to-an-issue"}`)
						resp := &http.Response{Request: req, StatusCode: http.StatusNotFound, Body: io.NopCloser(body)}
						return resp
					}(),
					Validate: func(t *testing.T, r *http.Request, payload []byte) {
						assignees, err := assigneesFromPayload(payload)
						require.NoError(t, err)
						assert.ElementsMatch(t, assignees, []string{"mojombo"})
					},
				},
			},
			ExpectedError: "POST https://api.github.com/repos/xunleii/github-quick-actions/issues/1/assignees: 404 Not Found []",
		},
		{
			Name:      "assignee doesn't have permission to this repo",
			Arguments: []string{"@mojombo"},
			APICalls: map[string]fixtures.APICallMITM{
				"POST https://api.github.com/repos/xunleii/github-quick-actions/issues/1/assignees": {
					Response: &http.Response{StatusCode: http.StatusCreated}, // NOTE: don't know why ...
					Validate: func(t *testing.T, r *http.Request, payload []byte) {
						assignees, err := assigneesFromPayload(payload)
						require.NoError(t, err)
						assert.ElementsMatch(t, assignees, []string{"mojombo"})
					},
				},
			},
		},
	},
}

func TestAssign(t *testing.T) { AssignIssueCommentFixtures.RunForQuickAction(t, quick_actions.Assign) }

// assigneesFromPayload extracts assignees field from JSON payload
func assigneesFromPayload(payload []byte) ([]string, error) {
	assignees := struct {
		Assignees []string `json:"assignees,omitempty"`
	}{}

	err := json.Unmarshal(payload, &assignees)
	return assignees.Assignees, err
}
