package quick_action

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/palantir/go-githubapp/githubapp"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"github.com/thoas/go-funk"
)

// GithubQuickActions manages all defined GitHub quick actions through
// a githubapp Handler.
type GithubQuickActions struct {
	// inherit ClientCreator methods
	githubapp.ClientCreator

	// actions contains all defined Github quick actions.
	actions map[string]map[string]GithubQuickActionHandler
	// eventWrappers contains all Github type definition handled by this instance.
	eventWrappers map[string]githubEventType
}

// NewGithubQuickActions creates a new instance of GithubQuickActions.
func NewGithubQuickActions(cc githubapp.ClientCreator) *GithubQuickActions {
	return &GithubQuickActions{
		ClientCreator: cc,
		actions:       map[string]map[string]GithubQuickActionHandler{},
		eventWrappers: map[string]githubEventType{},
	}
}

// AddQuickAction adds or replaces quick actions for the given command.
func (a GithubQuickActions) AddQuickAction(command string, handlers ...GithubQuickAction) {
	for _, handler := range handlers {
		a.eventWrappers[handler.OnEvent.Name()] = handler.OnEvent

		if a.actions[handler.OnEvent.Name()] == nil {
			a.actions[handler.OnEvent.Name()] = map[string]GithubQuickActionHandler{}
		}
		a.actions[handler.OnEvent.Name()][command] = handler.Handler
	}
}

// Handles implements githubapp.Handles
func (a GithubQuickActions) Handles() []string {
	var handles []string
	for k := range a.eventWrappers {
		handles = append(handles, k)
	}
	return funk.UniqString(handles)
}

// Handle implements githubapp.Handle
func (a GithubQuickActions) Handle(ctx context.Context, eventType, deliveryID string, payload []byte) error {
	eventWrapper, handled := a.eventWrappers[eventType]
	if !handled {
		return fmt.Errorf("'%s' event not handled... rejected", eventType)
	}

	event, err := eventWrapper.Wraps(payload)
	if err != nil {
		return errors.Wrapf(err, "failed to parse '%s' event payload", eventType)
	}

	// ignore if is not a new created comment
	if event.GetAction() != "created" {
		return nil
	}

	// update logger context with current event metadata
	logger := *zerolog.Ctx(ctx)
	logger.UpdateContext(func(ctx zerolog.Context) zerolog.Context {
		return ctx.
			Fields(map[string]interface{}{
				githubapp.LogKeyDeliveryID:      deliveryID,
				githubapp.LogKeyEventType:       eventType,
				githubapp.LogKeyRepositoryOwner: event.GetRepo().GetOwner().GetLogin(),
				githubapp.LogKeyRepositoryName:  event.GetRepo().GetName(),
			})
	})

	logger.Info().Msgf("new '%s' event handled", eventType)
	logger.Trace().
		RawJSON("payload", payload).
		Msgf("new '%s' event handled", eventType)

	handlers := a.actions[eventWrapper.Name()]

	var quickActions [][]string
	for n, line := range strings.Split(event.GetBody(), "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			// empty line, ignored
			continue
		}

		if line[0] != '/' {
			// not a command, ignored
			logger.Trace().Msgf("no command on line n°%d, ignored...", n)
			continue
		}

		commandLine := strings.Split(line, " ")
		if commandLine[0] == "/" {
			// empty command, ignored
			logger.Trace().Msgf("line n°%d has an empty command, ignored...", n)
			continue
		}

		action := commandLine[0][1:]
		if _, exists := handlers[action]; !exists {
			logger.Warn().Msgf("quick action '%s' doesn't exists, ignored", action)
			continue
		}

		var args []string
		for _, arg := range strings.Split(line, " ") {
			arg = strings.TrimSpace(arg)
			if arg != "" {
				args = append(args, arg)
			}
		}
		quickActions = append(quickActions, args)
	}
	if len(quickActions) == 0 {
		// no command line found, ignored
		return nil
	}

	var errsGroup = multierror.Group{}
	for _, quickAction := range quickActions {
		name := quickAction[0][1:]
		args := quickAction[1:]

		logger.Debug().Msgf("handle action '%s'", name)
		errsGroup.Go(func() error {
			err := handlers[name](ctx, a.ClientCreator, GithubQuickActionEvent{event.GetEventPayload(), args})
			if err != nil {
				logger.Err(err).Msgf("failed to run '%s': %s", name, err)
			}
			return err
		})
	}

	return errsGroup.Wait().ErrorOrNil()
}
