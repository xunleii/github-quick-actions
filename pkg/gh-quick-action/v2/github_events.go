package v1

import (
	"encoding/json"
	"fmt"

	"github.com/google/go-github/v39/github"
)

// IssueCommentEvent wraps *github.com/google/go-github/v39/github.IssueCommentEvent
// by implement EventPayload interface.
type IssueCommentEvent struct{ *github.IssueCommentEvent }

func newIssueCommentEvent(payload []byte) (EventPayload, error) {
	var event github.IssueCommentEvent
	err := json.Unmarshal(payload, &event)
	return &IssueCommentEvent{&event}, err
}

func (i *IssueCommentEvent) Type() EventType         { return EventTypeIssueComment }
func (i *IssueCommentEvent) Action() eventAction     { return eventAction(i.GetAction()) }
func (i *IssueCommentEvent) RepositoryName() string  { return i.GetRepo().GetName() }
func (i *IssueCommentEvent) RepositoryOwner() string { return i.GetRepo().GetOwner().GetLogin() }
func (i *IssueCommentEvent) Body() string            { return i.GetComment().GetBody() }
func (i *IssueCommentEvent) Raw() interface{}        { return i.IssueCommentEvent }

// payloadFactory is only used internally to generate EventPayload from raw JSON.
var payloadFactory = map[EventType]func([]byte) (EventPayload, error){
	EventTypeIssueComment: newIssueCommentEvent,
}

// PayloadFactory generates an EventPayload using the given JSON.
func PayloadFactory(eventType EventType, payload []byte) (EventPayload, error) {
	builder, exists := payloadFactory[eventType]
	if !exists {
		return nil, fmt.Errorf("event type '%s' not currently managed", eventType)
	}

	eventPayload, err := builder(payload)
	if err != nil {
		return nil, fmt.Errorf("failed to extract data from JSON for event '%s': %w", eventType, err)
	}
	return eventPayload, nil
}

//@formatter:off
