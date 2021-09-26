package quick_actions_test

import (
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	quick_actions "xnku.be/github-quick-actions/internal/quick-actions"
	fixtures2 "xnku.be/github-quick-actions/pkg/gh_quick_action/v1/fixtures"
)

var UnassignIssueCommentFixtures = fixtures2.EventFixtures{
	QuickActionName: "unassign",
	EventGenerator:  fixtures2.IssueCommentEventType,

	Fixtures: []fixtures2.EventFixture{
		{
			Name:      "simple assignment",
			Arguments: []string{"@mojombo"},
			APICalls: map[string]fixtures2.APICallMITM{
				"DELETE https://api.github.com/repos/xunleii/github-quick-actions/issues/1/assignees": {
					Response: &http.Response{StatusCode: http.StatusOK},
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
			APICalls: map[string]fixtures2.APICallMITM{
				"DELETE https://api.github.com/repos/xunleii/github-quick-actions/issues/1/assignees": {
					Response: &http.Response{StatusCode: http.StatusOK},
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
			APICalls: map[string]fixtures2.APICallMITM{
				"DELETE https://api.github.com/repos/xunleii/github-quick-actions/issues/1/assignees": {
					Response: &http.Response{StatusCode: http.StatusOK},
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
			APICalls: map[string]fixtures2.APICallMITM{
				"DELETE https://api.github.com/repos/xunleii/github-quick-actions/issues/1/assignees": {
					Response: &http.Response{StatusCode: http.StatusOK},
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
			APICalls: map[string]fixtures2.APICallMITM{
				"DELETE https://api.github.com/repos/xunleii/github-quick-actions/issues/1/assignees": {
					Response: &http.Response{StatusCode: http.StatusOK},
					Validate: func(t *testing.T, r *http.Request, payload []byte) {
						assignees, err := assigneesFromPayload(payload)
						require.NoError(t, err)
						assert.ElementsMatch(t, assignees, []string{"mojombo"})
					},
				},
			},
		},

		{
			Name:      "invalid assignee",
			Arguments: []string{"@"},
			APICalls: map[string]fixtures2.APICallMITM{},
		},
		{
			Name:      "empty assignee",
			Arguments: []string{""},
			APICalls: map[string]fixtures2.APICallMITM{},
		},
		{
			Name:      "repository/owner not found",
			Arguments: []string{"@mojombo"},
			APICalls: map[string]fixtures2.APICallMITM{
				"DELETE https://api.github.com/repos/xunleii/github-quick-actions/issues/1/assignees": {
					Response: func() *http.Response {
						req, _ := http.NewRequest(http.MethodDelete, "https://api.github.com/repos/xunleii/github-quick-actions/issues/1/assignees", nil)
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
			ExpectedError: "DELETE https://api.github.com/repos/xunleii/github-quick-actions/issues/1/assignees: 404 Not Found []",
		},
		{
			Name:      "assignee doesn't have permission to this repo",
			Arguments: []string{"@mojombo"},
			APICalls: map[string]fixtures2.APICallMITM{
				"DELETE https://api.github.com/repos/xunleii/github-quick-actions/issues/1/assignees": {
					Response: &http.Response{StatusCode: http.StatusOK}, // NOTE: don't know why ...
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

func TestUnassign(t *testing.T) {
	UnassignIssueCommentFixtures.RunForQuickAction(t, quick_actions.Unassign)
}
