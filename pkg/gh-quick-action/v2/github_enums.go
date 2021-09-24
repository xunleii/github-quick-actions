package v1

// EventType enumerates all currently managed event types
type EventType string

const (
	EventTypeIssueComment EventType = "issue_comment"
)

// eventAction enumerates all possible action available on a event
type eventAction string

const (
	EventActionCreated eventAction = "created"
	EventActionEdited  eventAction = "edited"
	EventActionDeleted eventAction = "deleted"
)
