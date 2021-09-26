package gqa_scenario_context

import (
	"encoding/json"
	"net/http"
	"strings"
	"sync"
	"unicode"

	gh_quick_actions "xnku.be/github-quick-actions/pkg/gh_quick_action/v2"
)

type (
	// ProxyQuickAction is a sort of middleware for quick actions. It intercepts
	// all commands to a specific quick action and saves errors and responses
	// in order to use them lately.
	ProxyQuickAction struct {
		gh_quick_actions.QuickAction

		called   int
		proxies  map[string]*ProxiesCommand
		scenario *QuickActionScenarioContext
		mx       sync.Mutex
	}

	ProxiesCommand struct {
		Requests []*ProxyRoundTripper
		Errors   []string
	}
)

func (action *ProxyQuickAction) HandleCommand(ctx *gh_quick_actions.EventContext, command *gh_quick_actions.EventCommand) error {
	action.mx.Lock()
	defer action.mx.Unlock()
	if action.proxies == nil {
		action.proxies = map[string]*ProxiesCommand{}
	}

	// NOTE: arguments are converted into JSON in order to easily access from Gherkin rules
	jsonArgs, _ := json.Marshal(command.Arguments)
	if action.proxies[string(jsonArgs)] == nil {
		action.proxies[string(jsonArgs)] = &ProxiesCommand{}
	}
	proxy := action.proxies[string(jsonArgs)]

	roundTripper := action.scenario.sharedProxy.Copy()
	proxy.Requests = append(proxy.Requests, roundTripper)
	ctx.ClientCreator = &ClientCreator{roundTripper}

	action.called++
	err := action.QuickAction.HandleCommand(ctx, command)
	if err != nil {
		strErr := strings.Map(func(r rune) rune {
			switch {
			case unicode.IsSpace(r): // NOTE: replace all \t, \n ... with simple space
				return ' '
			default:
				return r
			}
		}, err.Error())

		proxy.Errors = append(proxy.Errors, strings.TrimSpace(strErr))
	}

	return err
}

// InterceptedRequests returns all intercepted http.Request for the given
// command and method/url combo.
func (action *ProxyQuickAction) InterceptedRequests(arguments, method, url string) []*http.Request {
	interceptedCmd := action.interceptedCommand(arguments)
	if interceptedCmd == nil {
		return nil
	}

	var requests []*http.Request
	for _, roundTripper := range interceptedCmd.Requests {
		requests = append(requests, roundTripper.InterceptedRequests(method, url)...)
	}

	return requests
}

// InterceptedErrors returns all intercepted errors for the given command.
func (action *ProxyQuickAction) InterceptedErrors(arguments string) []string {
	interceptedCmd := action.interceptedCommand(arguments)
	if interceptedCmd == nil {
		return nil
	}

	return interceptedCmd.Errors
}

func (action *ProxyQuickAction) interceptedCommand(arguments string) *ProxiesCommand {
	if action.proxies == nil {
		return nil
	}

	if action.proxies[arguments] == nil {
		return nil
	}
	return action.proxies[arguments]
}
