package quick_action

import (
	"context"

	"github.com/google/go-github/v38/github"
	"github.com/palantir/go-githubapp/githubapp"
)

type (
	// githubEventPayload defines required methods from all Github event payload.
	githubEventPayload interface{}

	// githubWrappedEvent wraps a githubEventPayload in order to simplify it usage.
	githubWrappedEvent interface {
		GetAction() string
		GetInstallation() *github.Installation
		GetRepo() *github.Repository

		GetEventPayload() githubEventPayload
		GetBody() string
	}

	// githubEventType defines how to handle a specific Github event.
	// NOTE: this is not exported to force usage of definition implemented
	// 		 in this package.
	githubEventType interface {
		Name() string
		Wraps(payload []byte) (githubWrappedEvent, error)
	}
)

type (
	// GithubQuickActionHandler define how a specific command should be handled.
	// This is the main logical part of this project.
	GithubQuickActionHandler func(ctx context.Context, cc githubapp.ClientCreator, event GithubQuickActionEvent) error

	// GithubQuickActionEvent contains all information and arguments from a
	// quick action event.
	GithubQuickActionEvent struct {
		// Payload represents the Github event payload
		Payload githubEventPayload

		// Arguments lists all parsed action's arguments
		Arguments []string
	}

	// GithubQuickAction is a single definition of a Github quick action linked to
	// a specific handle.
	GithubQuickAction struct {
		OnEvent githubEventType
		Handler GithubQuickActionHandler
	}
)
