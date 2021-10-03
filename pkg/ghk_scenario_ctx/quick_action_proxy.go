package gqa_scenario_context

import (
	"encoding/json"
	"net/http"

	"github.com/google/go-github/v39/github"
	"github.com/improbable-eng/go-httpwares"
	"github.com/shurcooL/githubv4"
	"golang.org/x/oauth2"

	gh_quick_actions "xnku.be/github-quick-actions/pkg/gh_quick_action/v2"
)

type (
	// ProxyQuickAction injects some information on context in order to get
	// information on server side.
	ProxyQuickAction struct {
		gh_quick_actions.QuickAction
		client *http.Client
	}

	ProxyQuickActionErr struct {
		error
		ctx *gh_quick_actions.EventCommand
	}
)

const (
	EventHeader   string = "QuickAction-EventType"
	CommandHeader string = "QuickAction-Command"
	ArgsHeader    string = "QuickAction-Args"
)

func (action *ProxyQuickAction) HandleCommand(ctx *gh_quick_actions.EventContext, command *gh_quick_actions.EventCommand) error {
	// NOTE: arguments are converted into JSON in order to easily access from Gherkin rules
	jsonArgs, _ := json.Marshal(command.Arguments)

	// Use custom tripperware for the current client
	client := httpwares.WrapClient(action.client, func(next http.RoundTripper) http.RoundTripper {
		return httpwares.RoundTripperFunc(func(req *http.Request) (*http.Response, error) {
			req.Header.Add(EventHeader, string(command.Payload.Type()))
			req.Header.Add(CommandHeader, command.Command)
			req.Header.Add(ArgsHeader, string(jsonArgs))

			return next.RoundTrip(req)
		})
	})

	// Prepare the EventContext
	ctx.ClientCreator = &ClientCreator{client}
	_, _ = client.Get("quick-action://localhost/triggered")

	err := action.QuickAction.HandleCommand(ctx, command)
	if err != nil {
		return &ProxyQuickActionErr{
			error: err,
			ctx:   command,
		}
	}
	return nil
}


type ClientCreator struct{ *http.Client }

func (cc *ClientCreator) NewAppClient() (*github.Client, error) 							 { return github.NewClient(cc.Client), nil }
func (cc *ClientCreator) NewInstallationClient(installationID int64) (*github.Client, error) { return cc.NewAppClient() }
func (cc *ClientCreator) NewTokenSourceClient(ts oauth2.TokenSource) (*github.Client, error) { return cc.NewAppClient() }
func (cc *ClientCreator) NewTokenClient(token string) (*github.Client, error) 				 { return cc.NewAppClient() }

func (cc *ClientCreator) NewAppV4Client() (*githubv4.Client, error) 							 { panic("not implemented") }
func (cc *ClientCreator) NewInstallationV4Client(installationID int64) (*githubv4.Client, error) { panic("not implemented") }
func (cc *ClientCreator) NewTokenSourceV4Client(ts oauth2.TokenSource) (*githubv4.Client, error) { panic("not implemented") }
func (cc *ClientCreator) NewTokenV4Client(token string) (*githubv4.Client, error) 				 { panic("not implemented") }
