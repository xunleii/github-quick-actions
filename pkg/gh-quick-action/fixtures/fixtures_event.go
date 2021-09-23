package fixtures

import (
	"encoding/json"

	"github.com/google/go-github/v39/github"

	quick_action "xnku.be/github-quick-actions/pkg/gh-quick-action"
)

const (
	IssueCommentEventJSON = `{"action":"created","issue":{"number":1},"comment":{"user":{"login":"xunleii"},"body":"{{ .body }}"},"repository":{"name":"github-quick-actions","full_name":"xunleii/github-quick-actions","owner":{"login":"xunleii"}},"installation":{"id":1234567890}}`
)

var (
	// IssueCommentEventType manages handler using IssueComment events
	IssueCommentEventType EventGenerator = func(cc *MockClientCreator, fixture EventFixture) quick_action.GithubQuickActionEvent {
		var issueCommentEvent github.IssueCommentEvent
		_ = json.Unmarshal([]byte(IssueCommentEventJSON), &issueCommentEvent)

		return quick_action.GithubQuickActionEvent{
			Payload:   &issueCommentEvent,
			Arguments: fixture.Arguments,
		}
	}
)
