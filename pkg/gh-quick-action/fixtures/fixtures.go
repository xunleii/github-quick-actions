package fixtures

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"strings"
	"testing"
	"text/template"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	quick_action "xnku.be/github-quick-actions/pkg/gh-quick-action"
)

type (
	// EventFixtures defines a common structure to tests quick actions for a
	// specific kind of events
	EventFixtures struct {
		QuickActionName string

		// EventType defines which type of event should be generated
		EventType

		// Quick action test cases
		Fixtures []EventFixture
	}

	// EventType generates the correct GithubQuickActionEvent using the given fixture
	EventType func(cc *MockClientCreator, fixture EventFixture) quick_action.GithubQuickActionEvent

	// EventFixture defines a single test case for a specific event
	EventFixture struct {
		Name      string
		Arguments []string

		APICalls      map[string]APICallMITM
		ExpectedError string
	}

	// APICallMITM defines how "Github" should respond to a specific call
	APICallMITM struct {
		Response *http.Response
		Error    error
		Validate func(t *testing.T, r *http.Request, payload []byte)
	}
)

// RunForQuickAction runs all defined fixtures for a specific quick action handler.
func (fixtures EventFixtures) RunForQuickAction(t *testing.T, handler quick_action.GithubQuickActionHandler) {
	for _, fixture := range fixtures.Fixtures {
		t.Run(fmt.Sprintf("%s@%s", fixtures.QuickActionName, fixture.Name), func(t *testing.T) {
			cc := &MockClientCreator{}
			cc.On("NewInstallationClient", mock.Anything).Return(MockGithubClient, nil)
			for api, mitm := range fixture.APICalls {
				cc.On("github.Request", api).Return(mitm.Response, mitm.Error)
				// NOTE: github.RequestValidation is used to inject RequestValidation for the given API
				cc.On("github.RequestValidation", api).Return(func(r *http.Request) {
					if mitm.Validate != nil {
						payload, err := io.ReadAll(r.Body)
						require.NoError(t, err)
						require.NoError(t, r.Body.Close())

						mitm.Validate(t, r, payload)
					}
				})
			}

			err := handler(
				context.TODO(),
				fixtures.EventType(cc, fixture),
			)

			if len(fixture.APICalls) > 0 {
				// if API call MITM are defined, we expect that they should be called one time each
				cc.AssertNumberOfCalls(t, "github.Request", len(fixture.APICalls))
			}

			if fixture.ExpectedError == "" {
				assert.NoError(t, err)
				return
			}

			assert.EqualError(t, err, fixture.ExpectedError)
		})
	}
}

// EventJsonFromTemplate creates a JSON event from the given template.
func EventJsonFromTemplate(tpl string) func(action, body string, _ map[string]interface{}) []byte {
	jsonTemplate := template.Must(template.New("").Parse(tpl))

	return func(action, body string, opts map[string]interface{}) []byte {
		if opts == nil {
			opts = map[string]interface{}{}
		}
		opts["action"] = action
		opts["body"] = body

		b := new(strings.Builder)
		if err := jsonTemplate.Execute(b, opts); err != nil {
			panic(fmt.Errorf("failed to generate JSON from template: %w", err))
		}

		return []byte(b.String())
	}
}
