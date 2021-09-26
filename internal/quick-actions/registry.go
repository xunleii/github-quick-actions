package quick_actions

import (
	v2 "xnku.be/github-quick-actions/pkg/gh_quick_action/v2"
)

// registry is a shared registry containing all default Github quick actions
var registry = map[string]v2.QuickAction{}

// registerQuickAction add quick action to the internal registry.
// NOTE: this is for internal use only
func registerQuickAction(command string, quickAction v2.QuickAction) {
	registry[command] = quickAction
}

// InjectAll adds all defined Github quick actions to
// the given GithubQuickActions instance.
func InjectAll(gh *v2.GithubQuickActions) {
	for command, action := range registry {
		gh.AddQuickAction(command, action)
	}
}
