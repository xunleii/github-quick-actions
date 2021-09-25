package quick_actions_test

import (
	"bytes"
	"context"
	"net/http"
	"testing"
	"text/template"

	"github.com/Masterminds/sprig"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	quick_actions "xnku.be/github-quick-actions/internal/quick-actions"
	quick_action "xnku.be/github-quick-actions/pkg/gh-quick-action/v1"
	fixtures2 "xnku.be/github-quick-actions/pkg/gh-quick-action/v1/fixtures"
)

func TestLabelE2E(t *testing.T) {
	comments := []string{
		"/label feature",
		"/label bug feature",
		`/label feature\n/label "bug 'n feature"`,
	}

	cc := &fixtures2.MockClientCreator{}
	cc.On("NewInstallationClient", mock.Anything).Return(fixtures2.MockGithubClient, nil)
	cc.On("github.Request", mock.Anything).Return(&http.Response{StatusCode: http.StatusCreated}, nil)
	cc.On("github.RequestValidation", mock.Anything).Return(func(*http.Request) {})

	ghApp := quick_action.NewGithubQuickActions(cc)
	ghApp.AddQuickAction("label", quick_actions.LabelIssueComment)

	payloadTpl, _ := template.New("").Funcs(sprig.TxtFuncMap()).Parse(fixtures2.IssueCommentEventJSON)
	for _, comment := range comments {
		t.Run(comment, func(t *testing.T) {
			buffer := &bytes.Buffer{}
			err := payloadTpl.Execute(buffer, map[string]interface{}{"body": comment})
			require.NoError(t, err)

			err = ghApp.Handle(context.TODO(), quick_actions.LabelIssueComment.OnEvent.Name(), "", buffer.Bytes())
			assert.NoError(t, err)
		})
	}
}