package quick_action

import (
	"context"
	"encoding/json"
	"fmt"
	"testing"

	"github.com/brianvoe/gofakeit"
	"github.com/google/go-github/v38/github"
	"github.com/hashicorp/go-multierror"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"

	"github.com/google/uuid"
)

type IssueQuickActionsSuite struct {
	suite.Suite
	*IssueQuickActions

	simpleHdl  *MockGithubQuickActionHandler
	complexHdl *MockGithubQuickActionHandler
}

func (suite *IssueQuickActionsSuite) SetupTest() {
	suite.simpleHdl = &MockGithubQuickActionHandler{}
	suite.complexHdl = &MockGithubQuickActionHandler{}

	suite.IssueQuickActions = NewIssueQuickActions(nil)
	suite.IssueQuickActions.AddQuickAction("simple", suite.simpleHdl.Fnc)
	suite.IssueQuickActions.AddQuickAction("complex", suite.complexHdl.Fnc)
}

// event ignored on event other than `created`
func (suite *IssueQuickActionsSuite) TestInvalidEvent() {
	event := &github.RepositoryEvent{Action: github.String(gofakeit.RandString([]string{"modified", "deleted"}))}
	payload, _ := json.Marshal(event)

	err := suite.IssueQuickActions.Handle(
		context.Background(),
		"commented",
		uuid.NewString(),
		payload,
	)

	suite.simpleHdl.AssertNotCalled(suite.T(), "Fnc")

	suite.Assert().NoError(err)
}

// comment line are ignored if not starting with '/' or if an unknown action
func (suite *IssueQuickActionsSuite) TestInvalidLinesEvent() {
	event := &github.IssueCommentEvent{
		Action:  github.String(gofakeit.RandString([]string{"created", "modified", "deleted"})),
		Comment: &github.IssueComment{Body: github.String("not an action\n/unknown action\n/")},
	}
	payload, _ := json.Marshal(event)

	err := suite.IssueQuickActions.Handle(
		context.Background(),
		"commented",
		uuid.NewString(),
		payload,
	)

	suite.simpleHdl.AssertNotCalled(suite.T(), "Fnc")

	suite.Assert().NoError(err)
}

func (suite *IssueQuickActionsSuite) TestSimpleActionEvent() {
	event := &github.IssueCommentEvent{
		Action:  github.String("created"),
		Comment: &github.IssueComment{Body: github.String(`/simple`)},
	}
	payload, _ := json.Marshal(event)

	suite.simpleHdl.
		On("Fnc", mock.Anything, []string{}).
		Return(nil)

	err := suite.IssueQuickActions.Handle(
		context.Background(),
		"commented",
		uuid.NewString(),
		payload,
	)

	suite.simpleHdl.AssertExpectations(suite.T())
	suite.simpleHdl.AssertNumberOfCalls(suite.T(), "Fnc", 1)

	suite.Assert().NoError(err)
}

func (suite *IssueQuickActionsSuite) TestSimpleActionWithArgsEvent() {
	event := &github.IssueCommentEvent{
		Action: github.String("created"),
		Comment: &github.IssueComment{Body: github.String(`/simple arg1 		arg2   arg3`)},
	}
	payload, _ := json.Marshal(event)

	suite.simpleHdl.
		On("Fnc", mock.Anything, []string{"arg1", "arg2", "arg3"}).
		Return(nil)

	err := suite.IssueQuickActions.Handle(
		context.Background(),
		"commented",
		uuid.NewString(),
		payload,
	)

	suite.simpleHdl.AssertExpectations(suite.T())
	suite.simpleHdl.AssertNumberOfCalls(suite.T(), "Fnc", 1)

	suite.Assert().NoError(err)
}

func (suite *IssueQuickActionsSuite) TestSimpleActionError() {
	event := &github.IssueCommentEvent{
		Action:  github.String("created"),
		Comment: &github.IssueComment{Body: github.String(`/simple`)},
	}
	payload, _ := json.Marshal(event)

	actionErr := fmt.Errorf(gofakeit.HipsterSentence(10))
	suite.simpleHdl.
		On("Fnc", mock.Anything, []string{}).
		Return(actionErr)

	err := suite.IssueQuickActions.Handle(
		context.Background(),
		"commented",
		uuid.NewString(),
		payload,
	)

	suite.simpleHdl.AssertExpectations(suite.T())
	suite.simpleHdl.AssertNumberOfCalls(suite.T(), "Fnc", 1)

	suite.Assert().Error(err)
	suite.Assert().IsType(&multierror.Error{}, err)
	suite.Assert().NotPanics(func() {
		suite.Assert().Len(err.(*multierror.Error).WrappedErrors(), 1)
		suite.Assert().Equal(err.(*multierror.Error).WrappedErrors()[0], actionErr)
	})
}

func (suite *IssueQuickActionsSuite) TestMultiActionsEvent() {
	event := &github.IssueCommentEvent{
		Action:  github.String("created"),
		Comment: &github.IssueComment{Body: github.String("/simple a b c\n/simple\n/complex\n/complex /simple\n\n/simple not complex")},
	}
	payload, _ := json.Marshal(event)

	suite.simpleHdl.
		On("Fnc", mock.Anything, []string{"a", "b", "c"}).
		Return(nil)
	suite.simpleHdl.
		On("Fnc", mock.Anything, []string{}).
		Return(nil)
	suite.simpleHdl.
		On("Fnc", mock.Anything, []string{"not", "complex"}).
		Return(fmt.Errorf("yep, not complex at all"))

	suite.complexHdl.
		On("Fnc", mock.Anything, []string{}).
		Return(nil)
	suite.complexHdl.
		On("Fnc", mock.Anything, []string{"/simple"}).
		Return(fmt.Errorf("don't do that, please"))

	err := suite.IssueQuickActions.Handle(
		context.Background(),
		"commented",
		uuid.NewString(),
		payload,
	)

	suite.simpleHdl.AssertExpectations(suite.T())
	suite.simpleHdl.AssertNumberOfCalls(suite.T(), "Fnc", 3)
	suite.complexHdl.AssertExpectations(suite.T())
	suite.complexHdl.AssertNumberOfCalls(suite.T(), "Fnc", 2)

	suite.Assert().Error(err)
	suite.Assert().IsType(&multierror.Error{}, err)
	suite.Assert().NotPanics(func() {
		suite.Assert().Len(err.(*multierror.Error).WrappedErrors(), 2)
		suite.Assert().EqualError(err.(*multierror.Error).WrappedErrors()[0], "yep, not complex at all")
		suite.Assert().EqualError(err.(*multierror.Error).WrappedErrors()[1], "don't do that, please")
	})
}

// TestIssueQuickActions starts the testing suite
func TestIssueQuickActions(t *testing.T) { suite.Run(t, new(IssueQuickActionsSuite)) }
