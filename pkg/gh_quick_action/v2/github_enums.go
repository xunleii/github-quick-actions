package gh_quick_actions

// EventType enumerates all currently managed event types
type EventType string

const (
	EventTypeIssueComment EventType = "issue_comment"
)

// EventAction enumerates all possible action available on a event
type EventAction string

const (
	EventActionCreated EventAction = "created"
	EventActionEdited  EventAction = "edited"
	EventActionDeleted EventAction = "deleted"
)
