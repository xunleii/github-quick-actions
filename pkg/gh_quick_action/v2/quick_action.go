package gh_quick_actions

import (
	"context"

	"github.com/palantir/go-githubapp/githubapp"
)

type (
	// QuickAction defines a Github quick action.
	QuickAction interface {
		TriggerOnEvents() []EventType
		HandleCommand(ctx *EventContext, command *EventCommand) error
	}

	// EventContext implement all tools required in order to handle a
	// Github event.
	EventContext struct {
		context.Context
		githubapp.ClientCreator
	}

	// EventCommand contains the command to be handled by the
	// quick action implementation.
	EventCommand struct {
		Command   string
		Arguments []string

		Payload EventPayload
	}

	// EventPayload wraps native Github events in order to use them more easily.
	EventPayload interface {
		Type() EventType
		// Action returns the action that was performed on the comment
		// (one of "created", "edited" or "deleted").
		Action() EventAction
		RepositoryName() string
		RepositoryOwner() string
		IssueNumber() int
		Body() string

		// Raw contains the raw events if it needed to be used by the
		// quick action implementation.
		Raw() interface{}
	}
)
