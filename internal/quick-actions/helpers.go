package quick_actions

import (
	"fmt"

	"github.com/google/go-github/v39/github"
	"github.com/rs/zerolog"

	. "xnku.be/github-quick-actions/pkg/gh_quick_action/v2"
)

type (
	githubEventHelper struct{}

	githubInstallationInterface interface{ GetInstallation() *github.Installation }
)

func (githubEventHelper) newInstallationClient(ctx *EventContext, payload EventPayload) (*github.Client, error) {
	switch event := payload.Raw().(type) {
	case githubInstallationInterface:
		return ctx.NewInstallationClient(event.GetInstallation().GetID())
	default:
		return nil, fmt.Errorf("invalid event type %T", event)
	}
}

type (
	autoLogMiddleware struct{QuickAction}
)

func (mw autoLogMiddleware) HandleCommand(ctx *EventContext, command *EventCommand) error {
	logger := zerolog.Ctx(ctx).With().
		Str("quick_action", command.Command).
		Logger()
	logger.Info().Msgf("handling /%s (args: %v)", command.Command, command.Arguments)

	ctx.Context = logger.WithContext(ctx)
	err := mw.QuickAction.HandleCommand(ctx, command)

	return err
}