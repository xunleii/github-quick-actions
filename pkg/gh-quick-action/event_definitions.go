package quick_action

import (
	"encoding/json"

	"github.com/google/go-github/v38/github"
)

// GithubIssueCommentEvent defines Github issue comment event.
var GithubIssueCommentEvent = genericEventDefinition{
	name: "issue_comment",
	wraps: func(payload []byte) (githubWrappedEvent, error) {
		var event github.IssueCommentEvent
		err := json.Unmarshal(payload, &event)
		return &GithubWrappedIssueCommentEvent{event}, err
	},
}

// GithubWrappedIssueCommentEvent wraps github.IssueCommentEvent structure
type GithubWrappedIssueCommentEvent struct {
	github.IssueCommentEvent
}

func (g GithubWrappedIssueCommentEvent) GetEventPayload() githubEventPayload {
	return g.IssueCommentEvent
}
func (g GithubWrappedIssueCommentEvent) GetBody() string {
	return g.IssueCommentEvent.GetComment().GetBody()
}

// GithubIssuesEvent defines Github issue comment event.
var GithubIssuesEvent = genericEventDefinition{
	name: "issues",
	wraps: func(payload []byte) (githubWrappedEvent, error) {
		var event github.IssuesEvent
		err := json.Unmarshal(payload, &event)
		return &GithubWrappedIssueEvent{event}, err
	},
}

// GithubWrappedIssueEvent wraps github.IssuesEvent structure
type GithubWrappedIssueEvent struct {
	github.IssuesEvent
}

func (g GithubWrappedIssueEvent) GetEventPayload() githubEventPayload {
	return g.IssuesEvent
}
func (g GithubWrappedIssueEvent) GetBody() string {
	return g.IssuesEvent.GetIssue().GetBody()
}

// genericEventDefinition implements githubEventType to easily implement them as
// simple variable.
type genericEventDefinition struct {
	name  string
	wraps func(payload []byte) (githubWrappedEvent, error)
}

func (g genericEventDefinition) Name() string { return g.name }
func (g genericEventDefinition) Wraps(payload []byte) (githubWrappedEvent, error) {
	return g.wraps(payload)
}
