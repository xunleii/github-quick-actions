package gh_quick_actions

import (
	"encoding/json"
	"fmt"

	"github.com/google/go-github/v39/github"
)

// IssueEvent wraps *github.com/google/go-github/v39/github.IssuesEvent
// by implement EventPayload interface.
type IssueEvent struct{ *github.IssuesEvent }

func newIssueEvent(payload []byte) (EventPayload, error) {
	var event github.IssuesEvent
	err := json.Unmarshal(payload, &event)
	return &IssueEvent{&event}, err
}

func (i *IssueEvent) Type() EventType         { return EventTypeIssue }
func (i *IssueEvent) Action() EventAction     { return EventAction(i.GetAction()) }
func (i *IssueEvent) RepositoryName() string  { return i.GetRepo().GetName() }
func (i *IssueEvent) RepositoryOwner() string { return i.GetRepo().GetOwner().GetLogin() }
func (i *IssueEvent) IssueNumber() int        { return i.GetIssue().GetNumber() }
func (i *IssueEvent) Body() string            { return i.GetIssue().GetBody() }
func (i *IssueEvent) Raw() interface{}        { return i.IssuesEvent }

// IssueCommentEvent wraps *github.com/google/go-github/v39/github.IssueCommentEvent
// by implement EventPayload interface.
type IssueCommentEvent struct{ *github.IssueCommentEvent }

func newIssueCommentEvent(payload []byte) (EventPayload, error) {
	var event github.IssueCommentEvent
	err := json.Unmarshal(payload, &event)
	return &IssueCommentEvent{&event}, err
}

func (i *IssueCommentEvent) Type() EventType         { return EventTypeIssueComment }
func (i *IssueCommentEvent) Action() EventAction     { return EventAction(i.GetAction()) }
func (i *IssueCommentEvent) RepositoryName() string  { return i.GetRepo().GetName() }
func (i *IssueCommentEvent) RepositoryOwner() string { return i.GetRepo().GetOwner().GetLogin() }
func (i *IssueCommentEvent) IssueNumber() int        { return i.GetIssue().GetNumber() }
func (i *IssueCommentEvent) Body() string            { return i.GetComment().GetBody() }
func (i *IssueCommentEvent) Raw() interface{}        { return i.IssueCommentEvent }

// PullRequestEvent wraps *github.com/google/go-github/v39/github.PullRequestEvent
// by implement EventPayload interface.
type PullRequestEvent struct { *github.PullRequestEvent }

func newPullRequestEvent(payload []byte) (EventPayload, error) {
	var event github.PullRequestEvent
	err := json.Unmarshal(payload, &event)
	return &PullRequestEvent{&event}, err
}

func (i *PullRequestEvent) Type() EventType         { return EventTypePullRequest }
func (i *PullRequestEvent) Action() EventAction     { return EventAction(i.GetAction()) }
func (i *PullRequestEvent) RepositoryName() string  { return i.GetRepo().GetName() }
func (i *PullRequestEvent) RepositoryOwner() string { return i.GetRepo().GetOwner().GetLogin() }
func (i *PullRequestEvent) IssueNumber() int        { return i.GetPullRequest().GetNumber() }
func (i *PullRequestEvent) Body() string            { return i.GetPullRequest().GetBody() }
func (i *PullRequestEvent) Raw() interface{}        { return i.PullRequestEvent }

// PullRequestReviewCommentEvent wraps *github.com/google/go-github/v39/github.PullRequestReviewCommentEvent
// by implement EventPayload interface.
type PullRequestReviewCommentEvent struct { *github.PullRequestReviewCommentEvent }

func newPullRequestReviewCommentEvent(payload []byte) (EventPayload, error) {
	var event github.PullRequestReviewCommentEvent
	err := json.Unmarshal(payload, &event)
	return &PullRequestReviewCommentEvent{&event}, err
}

func (i *PullRequestReviewCommentEvent) Type() EventType         { return EventTypePullRequestReviewComment }
func (i *PullRequestReviewCommentEvent) Action() EventAction     { return EventAction(i.GetAction()) }
func (i *PullRequestReviewCommentEvent) RepositoryName() string  { return i.GetRepo().GetName() }
func (i *PullRequestReviewCommentEvent) RepositoryOwner() string { return i.GetRepo().GetOwner().GetLogin() }
func (i *PullRequestReviewCommentEvent) IssueNumber() int        { return i.GetPullRequest().GetNumber() }
func (i *PullRequestReviewCommentEvent) Body() string            { return i.GetComment().GetBody() }
func (i *PullRequestReviewCommentEvent) Raw() interface{}        { return i.PullRequestReviewCommentEvent }

// payloadFactory is only used internally to generate EventPayload from raw JSON.
var payloadFactory = map[EventType]func([]byte) (EventPayload, error){
	EventTypeIssue: newIssueEvent,
	EventTypeIssueComment: newIssueCommentEvent,
	EventTypePullRequest: newPullRequestEvent,
	EventTypePullRequestReviewComment: newPullRequestReviewCommentEvent,
}

// PayloadFactory generates an EventPayload using the given JSON.
func PayloadFactory(eventType EventType, json []byte) (EventPayload, error) {
	builder, exists := payloadFactory[eventType]
	if !exists {
		return nil, fmt.Errorf("event type '%s' not managed", eventType)
	}

	payload, err := builder(json)
	if err != nil {
		return nil, fmt.Errorf("failed to extract data from JSON for event '%s': %w", eventType, err)
	}
	return payload, nil
}
