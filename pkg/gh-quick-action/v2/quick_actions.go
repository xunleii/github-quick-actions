package gh_quick_actions

import (
	"context"
	"encoding/csv"
	"fmt"
	"strings"
	"unicode"

	"github.com/hashicorp/go-multierror"
	"github.com/palantir/go-githubapp/githubapp"
	"github.com/rs/zerolog"
	"github.com/thoas/go-funk"
)

type (
	// quickActionRegistry represents the registry that contains the
	// implemented quick actions.
	// The first key is the even type (issue, comment, ...) and the second
	// one is the command name. This allows to easily access to the correct
	// quick action giving the context (event type + command name).
	quickActionRegistry map[EventType]map[string]QuickAction

	// GithubQuickActions manages all defined GitHub quick actions through
	// a githubapp Handler.
	GithubQuickActions struct {
		// ClientCreator used by quick actions to create Github API clients.
		cc githubapp.ClientCreator

		// registry contains all Github quick actions implementations
		// that will be handled.
		registry quickActionRegistry
	}
)

// NewGithubQuickActions creates a new instance of GithubQuickActions.
func NewGithubQuickActions(cc githubapp.ClientCreator) *GithubQuickActions {
	return &GithubQuickActions{cc: cc, registry: quickActionRegistry{}}
}

// AddQuickAction add quick action for the given command.
func (a GithubQuickActions) AddQuickAction(command string, action QuickAction) {
	if action == nil {
		// NOTE: panic is used to avoid unknown issue like unexpected nil
		//		 action (so without panic, the action will be just ignored,
		//		 increasing the complexity to debug)
		panic(fmt.Errorf("quick action cannot be nil for command '/%s'", command))
	}

	eventTypes := action.TriggerOnEvents()

	for _, eventType := range eventTypes {
		if a.registry[eventType] == nil {
			a.registry[eventType] = map[string]QuickAction{}
		}

		if a.registry[eventType][command] != nil {
			// NOTE: panic to avoid unexpected overwrite of an existing action
			panic(fmt.Errorf("quick action already defined for command '/%s'", command))
		}
		a.registry[eventType][command] = action
	}
}

// Handles implements githubapp.Handles
func (a GithubQuickActions) Handles() []string {
	var handles []string
	for eventType := range a.registry {
		handles = append(handles, string(eventType))
	}
	return funk.UniqString(handles)
}

// Handle implements githubapp.Handle
func (a GithubQuickActions) Handle(ctx context.Context, evntType, deliveryID string, json []byte) error {
	logger := zerolog.Ctx(ctx)
	logger.UpdateContext(func(ctx zerolog.Context) zerolog.Context {
		return ctx.
			Str(githubapp.LogKeyEventType, evntType).
			Str(githubapp.LogKeyDeliveryID, deliveryID)
	})

	payload, err := PayloadFactory(EventType(evntType), json)
	if err != nil {
		logger.Error().Err(err).Msgf("failed to generate payload from JSON")
		return err
	}

	if payload.Action() != EventActionCreated {
		// NOTE: ignore all event if not "created"
		return nil
	}

	logger.UpdateContext(func(ctx zerolog.Context) zerolog.Context {
		return ctx.
			Str(githubapp.LogKeyRepositoryOwner, payload.RepositoryOwner()).
			Str(githubapp.LogKeyRepositoryName, payload.RepositoryName())
	})
	logger.Info().Send()
	logger.Trace().RawJSON("payload", json).Send()

	commands := a.payloadToCommands(ctx, payload)
	if len(commands) == 0 {
		logger.Info().Msgf("no command found, aborted")
		return nil
	}
	eventCtx := &EventContext{Context: ctx, ClientCreator: a.cc}

	errors := &multierror.Error{}
	for _, command := range commands {
		action := a.registry[command.Payload.Type()][command.Command]
		// TODO: in order to preserve user command order, all calls are sequential,
		// 		 increasing the execution time. Find a way to detect conflicts between
		//		 commands (like /label & /remove_label) and group them.
		err := action.HandleCommand(eventCtx, command)
		if err != nil {
			logger.Error().Err(err).Msgf("failed to run quick action: %s", err)
			errors = multierror.Append(errors, err)
		}
	}

	return errors.ErrorOrNil()
}

// payloadToCommands extracts all command defined by the user in the event body.
func (a GithubQuickActions) payloadToCommands(ctx context.Context, event EventPayload) []*EventCommand {
	logger := zerolog.Ctx(ctx)
	actions := a.registry[event.Type()]

	var commands []*EventCommand
	for n, line := range strings.Split(event.Body(), "\n") {
		if line == "" || line[0] != '/' {
			// not a command, ignored
			logger.Trace().Msgf("no command on line n°%d, ignored...", n)
			continue
		}

		// NOTE: in order to keep to CPU time, we avoid creating the CSV and
		// 		 parse the line if the action doesn't exist.
		idx := strings.IndexFunc(line, unicode.IsSpace)
		command := line[1:]
		switch idx {
		case 1: // NOTE: if idx == 1 means that le first "word" is only `/` and should be ignored
			logger.Trace().Msgf("no command on line n°%d, ignored...", n)
			continue
		case -1:
			// ignore because no space found means that the full line is the command
		default:
			command = line[1:idx]
		}

		if _, exists := actions[command]; !exists {
			logger.Warn().Msgf("quick action '/%s' doesn't exists, ignored", command)
			continue
		}

		reader := csv.NewReader(strings.NewReader(line))
		reader.Comma = ' '

		record, err := reader.Read()
		if err != nil {
			logger.Error().
				Err(err).
				Str("quick_action", command).
				Str("line", line).
				Msgf("failed to parse command line '%s', ignored", line)
			continue
		}

		var args []string
		for _, item := range record {
			item := strings.TrimSpace(item)
			if len(item) > 0 {
				args = append(args, item)
			}
		}

		commands = append(commands, &EventCommand{
			Command:   command,
			Arguments: args[1:],
			Payload:   event,
		})
	}
	return commands
}
