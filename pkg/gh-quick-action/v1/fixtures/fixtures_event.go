package fixtures

import (
	"bytes"
	"encoding/json"
	"text/template"

	"github.com/Masterminds/sprig"
	"github.com/google/go-github/v39/github"

	quick_action "xnku.be/github-quick-actions/pkg/gh-quick-action/v1"
)

const (
	IssueCommentEventJSON = `{"action":"created","issue":{"number":1},"comment":{"user":{"login":"xunleii"},"body":{{ .body | quote }}},"repository":{"name":"github-quick-actions","full_name":"xunleii/github-quick-actions","owner":{"login":"xunleii"}},"installation":{"id":1234567890}}`
)

var (
	// IssueCommentEventType manages handler using IssueComment events
	IssueCommentEventType EventGenerator = func(cc *MockClientCreator, fixture EventFixture) quick_action.GithubQuickActionEvent {
		var issueCommentEvent github.IssueCommentEvent

		buffer := &bytes.Buffer{}
		payloadTpl, _ := template.New("").Funcs(sprig.TxtFuncMap()).Parse(IssueCommentEventJSON)
		_ = payloadTpl.Execute(buffer, map[string]interface{}{"body": ""})
		_ = json.Unmarshal(buffer.Bytes(), &issueCommentEvent)

		return quick_action.GithubQuickActionEvent{
			Payload:   &issueCommentEvent,
			Arguments: fixture.Arguments,
		}
	}
)
