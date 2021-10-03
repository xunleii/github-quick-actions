package gqa_scenario_context

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/cucumber/godog"
	"github.com/google/uuid"
	"github.com/hashicorp/go-multierror"
	"github.com/nsf/jsondiff"
	"github.com/thoas/go-funk"

	gh_quick_actions "xnku.be/github-quick-actions/pkg/gh_quick_action/v2"
	"xnku.be/github-quick-actions/pkg/ghk_scenario_ctx/httptest"
)

type (
	// QuickActionScenarioContext implements a feature context for Godog allowing us
	// to test quick actions directly through Gherkin scenario.
	QuickActionScenarioContext struct {
		ghQuickActions *gh_quick_actions.GithubQuickActions
		ghAPIProxy     *GithubAPIProxy

		errs []error
	}
)

func ScenarioInitializer(quickActions map[string]gh_quick_actions.QuickAction) func(ctx *godog.ScenarioContext) {
	return func(ctx *godog.ScenarioContext) {
		scenario := &QuickActionScenarioContext{
			ghQuickActions: gh_quick_actions.NewGithubQuickActions(nil),
			ghAPIProxy:     NewGithubAPIProxy(),
		}

		srv := httptest.NewServer(scenario.ghAPIProxy)

		// NOTE: generates some steps dynamically using registered quick actions
		client := srv.Client()
		for command, quickAction := range quickActions {
			for _, eventType := range quickAction.TriggerOnEvents() {
				ctx.Step(fmt.Sprintf("^quick action \"/%s\" is registered for \"%s\" events$", command, eventType), func() {
					proxy := &ProxyQuickAction{
						QuickAction: quickAction,
						client:      client,
					}

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
		ctx.Step(`^Github Quick Actions should handle command "/([^"]+)" for "([^"]+)" event without argument by sending these following requests$`, scenario.assertNoArgCommandTriggeredSuccessfully)
		ctx.Step(`^Github Quick Actions should handle command "/([^"]+)" for "([^"]+)" event with arguments (\[.+\]) without sending anything$`, scenario.assertCommandTriggeredSuccessfullyWithoutRequest)
		ctx.Step(`^Github Quick Actions should handle command "/([^"]+)" for "([^"]+)" event without argument without sending anything$`, scenario.assertNoArgCommandTriggeredSuccessfullyWithoutRequest)
		ctx.Step(`^Github Quick Actions should handle command "/([^"]+)" for "([^"]+)" event with arguments (\[.+\]) but returns this error: '(.+)'$`, scenario.assertCommandTriggeredWithError)
		ctx.Step(`^Github Quick Actions should handle command "/([^"]+)" for "([^"]+)" event without argument but returns this error: '(.+)'$`, scenario.assertNoArgCommandTriggeredWithError)
	}
}

// simulateGithubEvent simulates an event sent by Github using the Gherkin rule's
// arguments (event type and JSON payload).
func (ctx *QuickActionScenarioContext) simulateGithubEvent(eventType string, json *godog.DocString) {
	err := ctx.ghQuickActions.Handle(context.Background(), eventType, uuid.New().String(), []byte(json.Content))
	ctx.errs = append(ctx.errs, err)
}

// simulateGithubAPIReply simulates a Github reply for the given request.
func (ctx *QuickActionScenarioContext) simulateGithubAPIReply(method, rawURL string, code int, response string) error {
	url, err := url.Parse(rawURL)
	if err != nil {
		return fmt.Errorf("failed to add response for %s %s: %w", method, url, err)
	}

	ctx.ghAPIProxy.NewRoute().
		Methods(method).Host(url.Host).Path(url.Path).
		HandlerFunc(func(writer http.ResponseWriter, req *http.Request) {
			writer.WriteHeader(code)
			_, _ = writer.Write([]byte(response))
		})
	return nil
}

// assertNoQuickActionsCalled asserts that Github Quick Actions didn't use any
// Quick Actions during the current scenario.
func (ctx *QuickActionScenarioContext) assertNoQuickActionsCalled() error {
	var commands []string
	for _, apiRequest := range ctx.ghAPIProxy.HandledRequests() {
		if apiRequest.IsMetadataRequest() {
			commands = append(commands, fmt.Sprintf("%s/%s", apiRequest.EventType, apiRequest.Command))
		}
	}

	if len(commands) > 0 {
		return fmt.Errorf(`One or several commands has been triggered: %s`, strings.Join(funk.UniqString(commands), ", "))
	}
	return nil
}

// assertErrorsHasBeenReturned asserts that the given errors hase been return
// from GithubQuickActions during the scenario.
// WARN: only internal error should be checked here; if you want to check
//		 specific QuickAction error, you must use the other rules.
func (ctx *QuickActionScenarioContext) assertErrorsHasBeenReturned(errorsDoc *godog.DocString) error {
	expectedErrors := strings.Split(errorsDoc.Content, "\n")
	var actualErrors []string

	for _, err := range ctx.errs {
		switch err.(type) {
		case *ProxyQuickActionErr:
			// NOTE: only get non ProxyQuickActionErr
		default:
			actualErrors = append(actualErrors, err.Error())
		}
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
func (ctx *QuickActionScenarioContext) assertCommandTriggeredSuccessfully(command, eventType, argumentsJSON string, requestsTable *godog.Table) error {
	// check if no error was returned
	errs := ctx.errorsForCommand(eventType, command, argumentsJSON)
	if len(errs) != 0 {
		return fmt.Errorf(`Command "/%s" (with %s) on "%s" event has returned the following error(s): %v`, command, argumentsJSON, eventType, errs)
	}

	// check if the current command has been called at least one time
	requests := ctx.ghAPIProxy.HandledRequests().
		WithEventType(eventType).
		WithCommand(command).
		WithArguments(argumentsJSON)
	if len(requests.With(func(r APIRequest) bool { return r.IsMetadataRequest() })) == 0 {
		return fmt.Errorf(`Command "/%s" (with %s) on "%s" event hasn't been called'`, command, argumentsJSON, eventType)
	}

	// check row validity
	if len(requestsTable.Rows) < 2 {
		return fmt.Errorf("At least 1 request should be defined")
	}

	if len(requestsTable.Rows[0].Cells) != 3 ||
		requestsTable.Rows[0].Cells[0].Value != "API request method" ||
		requestsTable.Rows[0].Cells[1].Value != "API request URL" ||
		requestsTable.Rows[0].Cells[2].Value != "API request payload" {
		return fmt.Errorf(`Invalid table definition; it must contain these 3 columns: "API request method", "API request URL" and "API request payload"`)
	}

	// check each row
tableIterator:
	for _, row := range requestsTable.Rows[1:] {
		method := row.Cells[0].Value
		url := row.Cells[1].Value
		expectedPayload := row.Cells[2].Value

		for _, request := range requests.With(func(r APIRequest) bool { return r.IsGithubAPIRequest() }) {
			if request.Body == nil {
				if expectedPayload == "" {
					continue tableIterator
				}
			} else {
				actual, _ := io.ReadAll(request.Body)
				_ = request.Body.Close()

				diff, _ := jsondiff.Compare(actual, []byte(expectedPayload), &jsondiff.Options{})
				if diff == jsondiff.FullMatch || diff == jsondiff.SupersetMatch {
					continue tableIterator
				}

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
	// check if no error was returned
	errs := ctx.errorsForCommand(eventType, command, argumentsJSON)
	if len(errs) != 0 {
		return fmt.Errorf(`Command "/%s" (with %s) on "%s" event has returned the following error(s): %v`, command, argumentsJSON, eventType, errs)
	}

	// check if the current command has been called at least one time
	requests := ctx.ghAPIProxy.HandledRequests().
		WithEventType(eventType).
		WithCommand(command).
		WithArguments(argumentsJSON)
	if len(requests.With(func(r APIRequest) bool { return r.IsMetadataRequest() })) == 0 {
		return fmt.Errorf(`Command "/%s" (with %s) on "%s" event hasn't been called'`, command, argumentsJSON, eventType)
	}

	var apiCalls []string
	for _, request := range requests.With(func(r APIRequest) bool { return r.IsGithubAPIRequest() }) {
		apiCalls = append(apiCalls, fmt.Sprintf("%s %s", request.Method, request.URL))
	}

	if len(apiCalls) > 0 {
		return fmt.Errorf(`Command "/%s" (with %s) on "%s" has sent some requests: [%s]`, command, argumentsJSON, eventType, strings.Join(apiCalls, ", "))
	}

	return nil
}

func (ctx *QuickActionScenarioContext) assertNoArgCommandTriggeredSuccessfullyWithoutRequest(command, eventType string) error {
	return ctx.assertCommandTriggeredSuccessfullyWithoutRequest(command, eventType, "[]")
}

// assertCommandTriggeredWithError asserts that the specified command
// should be triggered but returned the given error
func (ctx *QuickActionScenarioContext) assertCommandTriggeredWithError(command, eventType, argumentsJSON, error string) error {
	errs := ctx.errorsForCommand(eventType, command, argumentsJSON)
	if len(errs) == 0 {
		return fmt.Errorf(`Command "/%s" (with %s) on "%s" event didn't have returned any error`, command, argumentsJSON, eventType)
	}

	for _, err := range errs {
		if err == error {
			return nil
		}
	}
	return fmt.Errorf(`Command "/%s" (with %s) on "%s" event didn't have returned the specified error: %s`, command, argumentsJSON, eventType, error)
}

// assertNoArgCommandTriggeredWithError asserts the same things that
// assertCommandTriggeredWithError but for command without arguments.
func (ctx *QuickActionScenarioContext) assertNoArgCommandTriggeredWithError(command, eventType, error string) error {
	return ctx.assertCommandTriggeredWithError(command, eventType, "[]", error)
}

// errorsForCommand returns all formatted errors for a specific command.
func (ctx *QuickActionScenarioContext) errorsForCommand(eventType, command, argumentsJSON string) []string {
	var errs []string
	for _, err := range flattenErrors(ctx.errs) {
		if _, validErr := err.(*ProxyQuickActionErr); !validErr {
			continue
		}

		err := err.(*ProxyQuickActionErr)
		if string(err.ctx.Payload.Type()) == eventType && err.ctx.Command == command {
			currentArgumentsJSON, _ := json.Marshal(err.ctx.Arguments)
			if string(currentArgumentsJSON) == argumentsJSON {
				errs = append(errs, err.Error())
			}
		}
	}
	return errs
}

// flattenErrors extract all errors from multiple.Error recursively.
func flattenErrors(errs []error) []error {
	var flatErrs []error
	for _, err := range errs {
		switch err := err.(type) {
		case *multierror.Error:
			flatErrs = append(flatErrs, flattenErrors(err.WrappedErrors())...)
		case *ProxyQuickActionErr:
			ctx := err.ctx
			for _, err := range flattenErrors([]error{err.error}) {
				flatErrs = append(flatErrs, &ProxyQuickActionErr{error: err, ctx: ctx})
			}
		default:
			flatErrs = append(flatErrs, err)
		}
	}
	return flatErrs
}
