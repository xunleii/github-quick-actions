package gqa_scenario_context

import (
	"context"
	"fmt"
	"io"
	"net/http/httptest"
	"strings"

	"github.com/cucumber/godog"
	"github.com/google/uuid"
	"github.com/hashicorp/go-multierror"
	"github.com/nsf/jsondiff"
	"github.com/thoas/go-funk"

	gh_quick_actions "xnku.be/github-quick-actions/pkg/gh_quick_action/v2"
)

type (
	// QuickActionScenarioContext implements a feature context for Godog allowing us
	// to test quick actions directly through Gherkin scenario.
	QuickActionScenarioContext struct {
		ghQuickActions *gh_quick_actions.GithubQuickActions

		sharedProxy *ProxyRoundTripper
		registry    map[CommandEventTypeKey]*ProxyQuickAction
		errs        *multierror.Error
	}

	CommandEventTypeKey struct {
		Command   string
		EventType string
	}
)

func ScenarioInitializer(quickActions map[string]gh_quick_actions.QuickAction) func(ctx *godog.ScenarioContext) {
	return func(ctx *godog.ScenarioContext) {
		scenario := &QuickActionScenarioContext{
			ghQuickActions: gh_quick_actions.NewGithubQuickActions(nil),
			sharedProxy:    NewProxyRoundTripper(),
			registry:       map[CommandEventTypeKey]*ProxyQuickAction{},
		}

		// NOTE: generates some steps dynamically using registered quick actions
		for command, quickAction := range quickActions {
			for _, eventType := range quickAction.TriggerOnEvents() {
				ctx.Step(fmt.Sprintf("^quick action \"/%s\" is registered for \"%s\" events$", command, eventType), func() {
					proxy := &ProxyQuickAction{
						QuickAction: quickAction,
						proxies:     map[string]*ProxiesCommand{},
						scenario:    scenario,
					}

					scenario.registry[CommandEventTypeKey{command, string(eventType)}] = proxy
					scenario.ghQuickActions.AddQuickAction(command, proxy)
				})
			}
		}

		// WHEN steps
		ctx.Step(`^Github sends an event "([^"]*)" with$`, scenario.simulateGithubEvent)
		ctx.Step(`^Github replies to '([A-Z]+) ([^']+)' with '(\d{3}) (.+)'$`, scenario.simulateGithubAPIReply)

		// THEN steps
		ctx.Step(`^Github Quick Actions shouldn't do anything$`, scenario.assertNoQuickActionsCalled)
		ctx.Step(`^Github Quick Actions should return these errors$`, scenario.assertErrorsHasBeenReturned)
		ctx.Step(`^Github Quick Actions should handle command "/([^"]+)" for "([^"]+)" event with arguments (\[.+\]) by sending these following requests$`, scenario.assertCommandTriggeredSuccessfully)
		ctx.Step(`^Github Quick Actions should handle command "/([^"]+)" for "([^"]+)" event with no argument by sending these following requests$`, scenario.assertNoArgCommandTriggeredSuccessfully)
		ctx.Step(`^Github Quick Actions should handle command "/([^"]+)" for "([^"]+)" event with arguments (\[.+\]) without sending anything$`, scenario.assertCommandTriggeredSuccessfullyWithoutRequest)
		ctx.Step(`^Github Quick Actions should handle command "/([^"]+)" for "([^"]+)" event with no argument without sending anything$`, scenario.assertNoArgCommandTriggeredSuccessfullyWithoutRequest)
		ctx.Step(`^Github Quick Actions should handle command "/([^"]+)" for "([^"]+)" event with arguments (\[.+\]) but returns this error: '(.+)'$`, scenario.assertCommandTriggeredWithError)
		ctx.Step(`^Github Quick Actions should handle command "/([^"]+)" for "([^"]+)" event with no argument but returns this error: '(.+)'$`, scenario.assertNoArgCommandTriggeredWithError)
	}
}

// simulateGithubEvent simulates an event sent by Github using the Gherkin rule's
// arguments (event type and JSON payload).
func (ctx *QuickActionScenarioContext) simulateGithubEvent(eventType string, json *godog.DocString) {
	err := ctx.ghQuickActions.Handle(context.Background(), eventType, uuid.New().String(), []byte(json.Content))
	ctx.errs = multierror.Append(ctx.errs, err)
}

// simulateGithubAPIReply simulates a Github reply for the given request.
func (ctx *QuickActionScenarioContext) simulateGithubAPIReply(method, url string, code int, response string) {
	wr := httptest.NewRecorder()
	_, _ = wr.WriteString(response)
	wr.Code = code
	ctx.sharedProxy.InjectResponse(method, url, wr)
}

// assertNoQuickActionsCalled asserts that Github Quick Actions didn't use any
// Quick Actions during the current scenario.
func (ctx *QuickActionScenarioContext) assertNoQuickActionsCalled() error {
	for key, quickAction := range ctx.registry {
		if quickAction.called > 0 {
			return fmt.Errorf(`Command "/%s" has been trigger on event "%s"`, key.Command, key.EventType)
		}
	}
	return nil
}

// assertErrorsHasBeenReturned asserts that the given errors hase been return
// from GithubQuickActions during the scenario.
// WARN: only internal error should be checked here; if you want to check
//		 specific QuickAction error, you must use the other rules.
func (ctx *QuickActionScenarioContext) assertErrorsHasBeenReturned(errors *godog.DocString) error {
	expectedErrors := strings.Split(errors.Content, "\n")
	var actualErrors []string
	for _, err := range ctx.errs.WrappedErrors() {
		actualErrors = append(actualErrors, err.Error())
	}

	for _, expectedError := range expectedErrors {
		if !funk.ContainsString(actualErrors, expectedError) {
			return fmt.Errorf(`Github Quick Actions didn't returned the required error "%s"`, expectedError)
		}
	}

	return nil
}

// assertCommandTriggeredSuccessfully asserts that the specified command
// should be triggered, has sent the given requests and didn't have returned anything
func (ctx *QuickActionScenarioContext) assertCommandTriggeredSuccessfully(command, eventType, argumentsJSON string, requests *godog.Table) error {
	proxy := ctx.registry[CommandEventTypeKey{Command: command, EventType: eventType}]
	if proxy == nil {
		return fmt.Errorf(`Command "/%s" for "%s" events is not registered`, command, eventType)
	}

	errs := proxy.InterceptedErrors(argumentsJSON)
	if len(errs) != 0 {
		return fmt.Errorf(`Command "/%s" on "%s" event has returned the following error(s): %v`, command, eventType, errs)
	}

	if proxy.called == 0 {
		return fmt.Errorf(`Command "/%s" on "%s" event hasn't been called'`, command, eventType)
	}

	// check row validity
	if len(requests.Rows) < 2 {
		return fmt.Errorf("At least 1 request should be defined")
	}

	switch {
	case len(requests.Rows[0].Cells) != 3:
		fallthrough
	case requests.Rows[0].Cells[0].Value != "API request method":
		fallthrough
	case requests.Rows[0].Cells[1].Value != "API request URL":
		fallthrough
	case requests.Rows[0].Cells[2].Value != "API request payload":
		return fmt.Errorf(`Invalid table definition; it must contain these 3 columns: "API request method", "API request URL" and "API request payload"`)
	}

tableIterator:
	for _, row := range requests.Rows[1:] {
		method := row.Cells[0].Value
		url := row.Cells[1].Value
		expectedPayload := row.Cells[2].Value

		requests := proxy.InterceptedRequests(argumentsJSON, method, url)
		if requests == nil {
			errs = append(errs, fmt.Sprintf(`Request %s on "%s" not found for command "/%s" (with %s) on "%s" event`, method, url, command, argumentsJSON, eventType))
			continue
		}

		for _, request := range requests {
			actual, _ := io.ReadAll(request.Body)
			_ = request.Body.Close()

			diff, _ := jsondiff.Compare(actual, []byte(expectedPayload), &jsondiff.Options{})
			if diff == jsondiff.FullMatch || diff == jsondiff.SupersetMatch {
				continue tableIterator
			}
		}

		errs = append(errs, fmt.Sprintf(`No valid payload found for request %s on "%s" for command "/%s" (with %s) on "%s" event`, method, url, command, argumentsJSON, eventType))
	}

	switch len(errs) {
	case 0:
		return nil
	case 1:
		return fmt.Errorf(errs[0])
	default:
		return fmt.Errorf("Several requests are not validated: [%s]", strings.Join(errs, ", "))
	}
}

// assertNoArgCommandTriggeredSuccessfully asserts the same things that
// assertCommandTriggeredSuccessfully but for command without arguments.
func (ctx *QuickActionScenarioContext) assertNoArgCommandTriggeredSuccessfully(command, eventType string, requests *godog.Table) error {
	return ctx.assertCommandTriggeredSuccessfully(command, eventType, "[]", requests)
}

func (ctx *QuickActionScenarioContext) assertCommandTriggeredSuccessfullyWithoutRequest(command, eventType, argumentsJSON string) error {
	proxy := ctx.registry[CommandEventTypeKey{Command: command, EventType: eventType}]
	if proxy == nil {
		return fmt.Errorf(`Command "/%s" for "%s" events is not registered`, command, eventType)
	}

	errs := proxy.InterceptedErrors(argumentsJSON)
	if len(errs) != 0 {
		return fmt.Errorf(`Command "/%s" on "%s" event has returned the following error(s): %v`, command, eventType, errs)
	}

	if proxy.called == 0 {
		return fmt.Errorf(`Command "/%s" on "%s" event hasn't been called`, command, eventType)
	}

	if _, exists := proxy.proxies[argumentsJSON]; !exists {
		return fmt.Errorf(`Command "/%s" (with %s) on "%s" event hasn't been called`, command, argumentsJSON, eventType)
	}

	var requests []string
	for _, request := range proxy.proxies[argumentsJSON].Requests {
		for key := range request.interceptedRequests {
			requests = append(requests, fmt.Sprintf("%s %s", key.Method, key.URL))
		}
	}

	if len(requests) > 0 {
		return fmt.Errorf(`Command "/%s" on "%s" has sent some requests: [%s]`, command, eventType, strings.Join(requests, ", "))
	}

	return nil
}

func (ctx *QuickActionScenarioContext) assertNoArgCommandTriggeredSuccessfullyWithoutRequest(command, eventType string) error {
	return ctx.assertCommandTriggeredSuccessfullyWithoutRequest(command, eventType, "[]")
}

// assertCommandTriggeredWithError asserts that the specified command
// should be triggered but returned the given error
func (ctx *QuickActionScenarioContext) assertCommandTriggeredWithError(command, eventType, argumentsJSON, error string) error {
	proxy := ctx.registry[CommandEventTypeKey{Command: command, EventType: eventType}]
	if proxy == nil {
		return fmt.Errorf(`Command "/%s" for "%s" events is not registered`, command, eventType)
	}

	errs := proxy.InterceptedErrors(argumentsJSON)
	if len(errs) == 0 {
		return fmt.Errorf(`Command "/%s" on "%s" event didn't have returned any error`, command, eventType)
	}

	for _, err := range errs {
		if err == error {
			return nil
		}
	}
	return fmt.Errorf(`Command "/%s" on "%s" event didn't have returned the specified error: %s`, command, eventType, errs)
}

// assertNoArgCommandTriggeredWithError asserts the same things that
// assertCommandTriggeredWithError but for command without arguments.
func (ctx *QuickActionScenarioContext) assertNoArgCommandTriggeredWithError(command, eventType, error string) error {
	return ctx.assertCommandTriggeredWithError(command, eventType, "[]", error)
}
