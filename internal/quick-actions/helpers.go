package quick_actions

import (
	"fmt"

	"github.com/google/go-github/v39/github"

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
