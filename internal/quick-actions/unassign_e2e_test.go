package quick_actions_test

import (
	"bytes"
	"context"
	"net/http"
	"testing"
	"text/template"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	quick_actions "xnku.be/github-quick-actions/internal/quick-actions"
	quick_action "xnku.be/github-quick-actions/pkg/gh-quick-action"
	"xnku.be/github-quick-actions/pkg/gh-quick-action/fixtures"
)

func TestUnassignE2E(t *testing.T) {
	comments := []string{
		"/unassign me",
		"/unassign @mojombo",
		"/unassign me @mojombo",
		`/unassign me\n/assign @mojombo`,
	}

	cc := &fixtures.MockClientCreator{}
	cc.On("NewInstallationClient", mock.Anything).Return(fixtures.MockGithubClient, nil)
	cc.On("github.Request", mock.Anything).Return(&http.Response{StatusCode: http.StatusCreated}, nil)
	cc.On("github.RequestValidation", mock.Anything).Return(func(*http.Request) {})

	ghApp := quick_action.NewGithubQuickActions(cc)
	ghApp.AddQuickAction("unassign", quick_actions.UnassignIssueComment)

	payloadTpl, _ := template.New("").Parse(fixtures.IssueCommentEventJSON)
	for _, comment := range comments {
		t.Run(comment, func(t *testing.T) {
			buffer := &bytes.Buffer{}
			err := payloadTpl.Execute(buffer, map[string]interface{}{"body": comment})
			require.NoError(t, err)

			err = ghApp.Handle(context.TODO(), quick_actions.UnassignIssueComment.OnEvent.Name(), "", buffer.Bytes())
			assert.NoError(t, err)
		})
	}
}
