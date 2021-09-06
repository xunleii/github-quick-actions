package quick_action_test

import (
	"context"
	"encoding/json"
	"fmt"
	"sort"
	"testing"

	"github.com/brianvoe/gofakeit"
	"github.com/google/go-github/v38/github"
	"github.com/hashicorp/go-multierror"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"

	"github.com/google/uuid"

	quick_action "xnku.be/github-quick-actions/pkg/gh-quick-action"
	"xnku.be/github-quick-actions/pkg/gh-quick-action/fixtures"
)

type GithubQuickActionsSuite struct {
	suite.Suite
	GithubQuickActions *quick_action.GithubQuickActions

	SimpleQA  *quick_action.MockGithubQuickActionHandler
	ComplexQA *quick_action.MockGithubQuickActionHandler
}

func (suite *GithubQuickActionsSuite) SetupTest() {
	suite.SimpleQA = &quick_action.MockGithubQuickActionHandler{}
	suite.ComplexQA = &quick_action.MockGithubQuickActionHandler{}

	suite.GithubQuickActions = quick_action.NewGithubQuickActions(&fixtures.MockClientCreator{})

	suite.GithubQuickActions.AddQuickAction(
		"simple",
		quick_action.GithubQuickAction{
			OnEvent: quick_action.GithubIssueCommentEvent,
			Handler: suite.SimpleQA.Fnc,
		},
		quick_action.GithubQuickAction{
			OnEvent: quick_action.GithubIssuesEvent,
			Handler: suite.SimpleQA.Fnc,
		},
	)

	suite.GithubQuickActions.AddQuickAction(
		"complex",
		quick_action.GithubQuickAction{
			OnEvent: quick_action.GithubIssuesEvent,
			Handler: suite.ComplexQA.Fnc,
		},
	)
}

func (suite *GithubQuickActionsSuite) TestHandles() {
	handles := suite.GithubQuickActions.Handles()
	suite.Assert().ElementsMatch(handles, []string{"issue_comment", "issues"})
}

// event failed on unknow event type
func (suite *GithubQuickActionsSuite) TestUnknownEvent() {
	err := suite.GithubQuickActions.Handle(
		context.Background(),
		"invalid_event",
		uuid.NewString(),
		nil,
	)
	suite.Assert().EqualError(err, "'invalid_event' event not handled... rejected")

	suite.SimpleQA.AssertNotCalled(suite.T(), "Fnc")
	suite.ComplexQA.AssertNotCalled(suite.T(), "Fnc")
}

// event ignored on event other than `created`
func (suite *GithubQuickActionsSuite) TestInvalidEvent() {
	event := &github.RepositoryEvent{Action: github.String(gofakeit.RandString([]string{"modified", "deleted"}))}
	payload, _ := json.Marshal(event)

	err := suite.GithubQuickActions.Handle(
		context.Background(),
		"issue_comment",
		uuid.NewString(),
		payload,
	)
	suite.Assert().NoError(err)

	suite.SimpleQA.AssertNotCalled(suite.T(), "Fnc")
}

// comment line are ignored if not starting with '/' or if an unknown action
func (suite *GithubQuickActionsSuite) TestInvalidLinesEvent() {
	event := &github.IssueCommentEvent{
		Action:  github.String("created"),
		Comment: &github.IssueComment{Body: github.String("not an action\n/unknown action\n/")},
	}
	payload, _ := json.Marshal(event)

	err := suite.GithubQuickActions.Handle(
		context.Background(),
		"issue_comment",
		uuid.NewString(),
		payload,
	)
	suite.Assert().NoError(err)

	suite.SimpleQA.AssertNotCalled(suite.T(), "Fnc")
	suite.ComplexQA.AssertNotCalled(suite.T(), "Fnc")
}

func (suite *GithubQuickActionsSuite) TestSimpleIssuesEvent() {
	event := &github.IssueCommentEvent{
		Action:  github.String("created"),
		Comment: &github.IssueComment{Body: github.String(`/simple`)},
	}
	payload, _ := json.Marshal(event)

	suite.SimpleQA.
		On("Fnc", mock.Anything, []string{}).
		Return(nil)

	err := suite.GithubQuickActions.Handle(
		context.Background(),
		"issue_comment",
		uuid.NewString(),
		payload,
	)
	suite.Assert().NoError(err)

	suite.SimpleQA.AssertExpectations(suite.T())
	suite.SimpleQA.AssertNumberOfCalls(suite.T(), "Fnc", 1)
	suite.ComplexQA.AssertNotCalled(suite.T(), "Fnc")
}

func (suite *GithubQuickActionsSuite) TestSimpleActionWithArgsEvent() {
	event := &github.IssueCommentEvent{
		Action: github.String("created"),
		Comment: &github.IssueComment{Body: github.String(`/simple arg1 		arg2   arg3`)},
	}
	payload, _ := json.Marshal(event)

	suite.SimpleQA.
		On("Fnc", mock.Anything, []string{"arg1", "arg2", "arg3"}).
		Return(nil)

	err := suite.GithubQuickActions.Handle(
		context.Background(),
		"issue_comment",
		uuid.NewString(),
		payload,
	)

	suite.Assert().NoError(err)

	suite.SimpleQA.AssertExpectations(suite.T())
	suite.SimpleQA.AssertNumberOfCalls(suite.T(), "Fnc", 1)
	suite.ComplexQA.AssertNotCalled(suite.T(), "Fnc")
}

func (suite *GithubQuickActionsSuite) TestSimpleActionError() {
	event := &github.IssueCommentEvent{
		Action:  github.String("created"),
		Comment: &github.IssueComment{Body: github.String(`/simple`)},
	}
	payload, _ := json.Marshal(event)

	actionErr := fmt.Errorf(gofakeit.HipsterSentence(10))
	suite.SimpleQA.
		On("Fnc", mock.Anything, []string{}).
		Return(actionErr)

	err := suite.GithubQuickActions.Handle(
		context.Background(),
		"issue_comment",
		uuid.NewString(),
		payload,
	)

	suite.SimpleQA.AssertExpectations(suite.T())
	suite.SimpleQA.AssertNumberOfCalls(suite.T(), "Fnc", 1)
	suite.ComplexQA.AssertNotCalled(suite.T(), "Fnc")

	suite.Assert().Error(err)
	suite.Assert().IsType(&multierror.Error{}, err)
	suite.Assert().NotPanics(func() {
		suite.Assert().Len(err.(*multierror.Error).WrappedErrors(), 1)
		suite.Assert().Equal(err.(*multierror.Error).WrappedErrors()[0], actionErr)
	})
}

func (suite *GithubQuickActionsSuite) TestMultiActionsEvent() {
	event := &github.IssuesEvent{
		Action: github.String("created"),
		Issue: &github.Issue{
			Body: github.String("/simple a b c\n/simple\n/complex\n/complex /simple\n\n/simple not complex"),
		},
	}
	payload, _ := json.Marshal(event)

	suite.SimpleQA.
		On("Fnc", mock.Anything, []string{"a", "b", "c"}).
		Return(nil)
	suite.SimpleQA.
		On("Fnc", mock.Anything, []string{}).
		Return(nil)
	suite.SimpleQA.
		On("Fnc", mock.Anything, []string{"not", "complex"}).
		Return(fmt.Errorf("yep, not complex at all"))

	suite.ComplexQA.
		On("Fnc", mock.Anything, []string{}).
		Return(nil)
	suite.ComplexQA.
		On("Fnc", mock.Anything, []string{"/simple"}).
		Return(fmt.Errorf("don't do that, please"))

	err := suite.GithubQuickActions.Handle(
		context.Background(),
		"issues",
		uuid.NewString(),
		payload,
	)

	suite.SimpleQA.AssertExpectations(suite.T())
	suite.SimpleQA.AssertNumberOfCalls(suite.T(), "Fnc", 3)
	suite.ComplexQA.AssertExpectations(suite.T())
	suite.ComplexQA.AssertNumberOfCalls(suite.T(), "Fnc", 2)

	suite.Assert().Error(err)
	suite.Assert().IsType(&multierror.Error{}, err)
	suite.Assert().NotPanics(func() {
		suite.Assert().Len(err.(*multierror.Error).WrappedErrors(), 2)
		sort.Sort(err.(*multierror.Error))

		suite.Assert().EqualError(err.(*multierror.Error).WrappedErrors()[0], "don't do that, please")
		suite.Assert().EqualError(err.(*multierror.Error).WrappedErrors()[1], "yep, not complex at all")
	})
}

func (suite *GithubQuickActionsSuite) TestMultiActionsOnSpecificEvent() {
	event := &github.IssueCommentEvent{
		Action:  github.String("created"),
		Comment: &github.IssueComment{Body: github.String("/simple a b c\n/simple\n/complex\n/complex /simple\n\n/simple not complex")},
	}
	payload, _ := json.Marshal(event)

	suite.SimpleQA.
		On("Fnc", mock.Anything, []string{"a", "b", "c"}).
		Return(nil)
	suite.SimpleQA.
		On("Fnc", mock.Anything, []string{}).
		Return(nil)
	suite.SimpleQA.
		On("Fnc", mock.Anything, []string{"not", "complex"}).
		Return(fmt.Errorf("yep, not complex at all"))

	err := suite.GithubQuickActions.Handle(
		context.Background(),
		"issue_comment",
		uuid.NewString(),
		payload,
	)

	suite.SimpleQA.AssertExpectations(suite.T())
	suite.SimpleQA.AssertNumberOfCalls(suite.T(), "Fnc", 3)
	suite.ComplexQA.AssertNotCalled(suite.T(), "Fnc")

	suite.Assert().Error(err)
	suite.Assert().IsType(&multierror.Error{}, err)
	suite.Assert().NotPanics(func() {
		suite.Assert().Len(err.(*multierror.Error).WrappedErrors(), 1)
		sort.Sort(err.(*multierror.Error))

		suite.Assert().EqualError(err.(*multierror.Error).WrappedErrors()[0], "yep, not complex at all")
	})
}

// TestIssueQuickActions starts the testing suite
func TestIssueQuickActions(t *testing.T) { suite.Run(t, new(GithubQuickActionsSuite)) }
