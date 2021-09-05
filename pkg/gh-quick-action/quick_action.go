package quick_action

import (
	"context"

	"github.com/palantir/go-githubapp/githubapp"
)

// GithubQuickActions is the container that manage all GitHub quick actions and handle requests.
type GithubQuickActions interface {
	// EventHandler implements the EventHandler used by the githubapp package
	githubapp.EventHandler
	// AddQuickAction add a new Github quick action that should be handled
	AddQuickAction(command string, handler GithubQuickActionHandler)
}

// GithubQuickActionEvent contains all information and arguments from a quick action event.
type GithubQuickActionEvent interface {
	// ClientCreator inherits ClientCreator methods
	githubapp.ClientCreator

	// Arguments lists all action's arguments
	Arguments() []string
}

// GithubQuickActionHandler define how a specific command should be handled. This is the
// main logical part of this project.
type GithubQuickActionHandler func(ctx context.Context, payload GithubQuickActionEvent) error
