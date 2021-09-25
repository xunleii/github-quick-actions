package gh_quick_actions

import (
	"context"
	"testing"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/suite"
)

type quickActionsTestSuite struct {
	suite.Suite
	*GithubQuickActions
}

func (ts *quickActionsTestSuite) SetupTest() { ts.GithubQuickActions = NewGithubQuickActions(nil) }

// GithubQuickActions.AddQuickAction
func (ts *quickActionsTestSuite) TestAddQuickAction_valid() {
	ts.GithubQuickActions.AddQuickAction("cmd#1", &mockQuickAction{onEvents: []EventType{"aaa", "bbb"}})
	ts.GithubQuickActions.AddQuickAction("cmd#2", &mockQuickAction{onEvents: []EventType{"aaa", "bbb", "ccc"}})

	// current registry should be
	// 	- aaa
	//		- cmd#1
	//			- *mockQuickAction#1
	//		- cmd#2
	//			- *mockQuickAction#3
	// 	- bbb
	//		- cmd#1
	//			- *mockQuickAction#1
	//		- cmd#2
	//			- *mockQuickAction#2
	// 	- ccc
	//		- cmd#2
	//			- *mockQuickAction#2
	//			- *mockQuickAction#3

	ts.Assert().NotNil(ts.GithubQuickActions.registry["aaa"]["cmd#1"])
	ts.Assert().NotNil(ts.GithubQuickActions.registry["aaa"]["cmd#2"])
	ts.Assert().NotNil(ts.GithubQuickActions.registry["bbb"]["cmd#1"])
	ts.Assert().NotNil(ts.GithubQuickActions.registry["bbb"]["cmd#2"])
	ts.Assert().NotNil(ts.GithubQuickActions.registry["ccc"]["cmd#2"])
}

func (ts *quickActionsTestSuite) TestAddQuickAction_nil() {
	ts.Assert().PanicsWithError("quick action cannot be nil for command '/cmd#1'", func() {
		ts.GithubQuickActions.AddQuickAction("cmd#1", nil)
	})
}

func (ts *quickActionsTestSuite) TestAddQuickAction_alreadyExists() {
	ts.GithubQuickActions.AddQuickAction("cmd#1", &mockQuickAction{onEvents: []EventType{"aaa", "bbb"}})
	ts.Assert().PanicsWithError("quick action already defined for command '/cmd#1'", func() {
		ts.GithubQuickActions.AddQuickAction("cmd#1", &mockQuickAction{onEvents: []EventType{"aaa", "bbb"}})
	})
}

// GithubQuickActions.Handles
func (ts *quickActionsTestSuite) TestHandles() {
	ts.GithubQuickActions.AddQuickAction("cmd#1", &mockQuickAction{onEvents: []EventType{"aaa", "bbb"}})
	ts.GithubQuickActions.AddQuickAction("cmd#2", &mockQuickAction{onEvents: []EventType{"aaa", "bbb", "ccc"}})

	ts.Assert().ElementsMatch([]string{"aaa", "bbb", "ccc"}, ts.GithubQuickActions.Handles())
}

// GithubQuickActions.payloadToCommands
func (ts *quickActionsTestSuite) TestPayloadToCommands() {
	ts.GithubQuickActions.AddQuickAction("cmd#1", &mockQuickAction{onEvents: []EventType{"aaa", "bbb"}})
	ts.GithubQuickActions.AddQuickAction("cmd#2", &mockQuickAction{onEvents: []EventType{"aaa", "bbb", "ccc"}})

	payload := mockEventPayload{body: `
  
/
// invalid command
/unknown "command"
/cmd#1
/cmd#1 simple
/cmd#2 "quoted arguments"
/cmd#2 mixed "arguments" with simple and "quoted arguments"
  /cmd#1 even with spaces ?
`}

	noLog := zerolog.Nop()
	ctx := noLog.WithContext(context.Background())

	payload.eventType = "aaa"
	commands := ts.GithubQuickActions.payloadToCommands(ctx, payload)

	ts.Require().Len(commands, 4)
	ts.Assert().Equal(EventCommand{Command: "cmd#1", Arguments: []string{}, Payload: payload}, *commands[0])
	ts.Assert().Equal(EventCommand{Command: "cmd#1", Arguments: []string{"simple"}, Payload: payload}, *commands[1])
	ts.Assert().Equal(EventCommand{Command: "cmd#2", Arguments: []string{"quoted arguments"}, Payload: payload}, *commands[2])
	ts.Assert().Equal(EventCommand{Command: "cmd#2", Arguments: []string{"mixed", "arguments", "with", "simple", "and", "quoted arguments"}, Payload: payload}, *commands[3])

	payload.eventType = "ccc"
	commands = ts.GithubQuickActions.payloadToCommands(ctx, payload)

	ts.Require().Len(commands, 2) // NOTE: `cmd#1` is only available for events `aaa` and `bbb`
	ts.Assert().Equal(EventCommand{Command: "cmd#2", Arguments: []string{"quoted arguments"}, Payload: payload}, *commands[0])
	ts.Assert().Equal(EventCommand{Command: "cmd#2", Arguments: []string{"mixed", "arguments", "with", "simple", "and", "quoted arguments"}, Payload: payload}, *commands[1])
}

func TestGithubQuickActionsSuite(t *testing.T) { suite.Run(t, new(quickActionsTestSuite)) }

// mockQuickAction implements a simple QuickAction
type mockQuickAction struct {
	onEvents []EventType
	retErr   error
}

func (m mockQuickAction) TriggerOnEvents() []EventType { return m.onEvents }
func (m mockQuickAction) HandleCommand(ctx *EventContext, command *EventCommand) error {
	return m.retErr
}
