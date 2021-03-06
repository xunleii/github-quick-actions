package gh_quick_actions

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/google/go-github/v39/github"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	issueEventFixture = &github.IssuesEvent{
		Action: github.String("created"),
		Repo:   &github.Repository{Name: github.String("github-quick-actions"), Owner: &github.User{Login: github.String("xunleii")}},
		Issue:  &github.Issue{Body: github.String("..."), Number: github.Int(0)},
	}
	issueEventFixtureJSON, _ = json.Marshal(issueEventFixture)

	issueCommentEventFixture = &github.IssueCommentEvent{
		Action:  github.String("created"),
		Repo:    &github.Repository{Name: github.String("github-quick-actions"), Owner: &github.User{Login: github.String("xunleii")}},
		Comment: &github.IssueComment{Body: github.String("...")},
		Issue:   &github.Issue{Number: github.Int(0)},
	}
	issueCommentEventFixtureJSON, _ = json.Marshal(issueCommentEventFixture)

	pullRequestEventFixture = &github.PullRequestEvent{
		Action:      github.String("created"),
		Repo:        &github.Repository{Name: github.String("github-quick-actions"), Owner: &github.User{Login: github.String("xunleii")}},
		PullRequest: &github.PullRequest{Body: github.String("..."), Number: github.Int(0)},
	}
	pullRequestEventFixtureJSON, _ = json.Marshal(pullRequestEventFixture)

	pullRequestReviewCommentEventFixture = &github.PullRequestReviewCommentEvent{
		Action:      github.String("created"),
		Repo:        &github.Repository{Name: github.String("github-quick-actions"), Owner: &github.User{Login: github.String("xunleii")}},
		Comment:     &github.PullRequestComment{Body: github.String("...")},
		PullRequest: &github.PullRequest{Number: github.Int(0)},
	}
	pullRequestReviewCommentEventFixtureJSON, _ = json.Marshal(pullRequestReviewCommentEventFixture)
)

func TestPayloadFactory(t *testing.T) {
	tts := []struct {
		name            string
		eventType       EventType
		eventJSON       []byte
		expectedPayload EventPayload
		err             error
	}{
		{
			name:      "unknown event type",
			eventType: EventType("unknown"),
			err:       fmt.Errorf("event type 'unknown' not managed"),
		},

		{
			name:            "github.IssueEvent",
			eventType:       EventTypeIssue,
			eventJSON:       issueEventFixtureJSON,
			expectedPayload: mockEventPayload{EventTypeIssue, EventActionCreated, "github-quick-actions", "xunleii", 0, "..."},
		},
		{
			name:      "github.IssueEvent@invalid",
			eventType: EventTypeIssue,
			eventJSON: []byte("...invalid..."),
			err:       fmt.Errorf("failed to extract data from JSON for event 'issue': invalid character '.' looking for beginning of value"),
		},

		{
			name:            "github.IssueCommentEvent",
			eventType:       EventTypeIssueComment,
			eventJSON:       issueCommentEventFixtureJSON,
			expectedPayload: mockEventPayload{EventTypeIssueComment, EventActionCreated, "github-quick-actions", "xunleii", 0, "..."},
		},
		{
			name:      "github.IssueCommentEvent@invalid",
			eventType: EventTypeIssueComment,
			eventJSON: []byte("...invalid..."),
			err:       fmt.Errorf("failed to extract data from JSON for event 'issue_comment': invalid character '.' looking for beginning of value"),
		},

		{
			name:            "github.PullRequestEvent",
			eventType:       EventTypePullRequest,
			eventJSON:       pullRequestEventFixtureJSON,
			expectedPayload: mockEventPayload{EventTypePullRequest, EventActionCreated, "github-quick-actions", "xunleii", 0, "..."},
		},
		{
			name:      "github.PullRequestEvent@invalid",
			eventType: EventTypePullRequest,
			eventJSON: []byte("...invalid..."),
			err:       fmt.Errorf("failed to extract data from JSON for event 'pull_request': invalid character '.' looking for beginning of value"),
		},

		{
			name:            "github.PullRequestReviewCommentEvent",
			eventType:       EventTypePullRequestReviewComment,
			eventJSON:       pullRequestReviewCommentEventFixtureJSON,
			expectedPayload: mockEventPayload{EventTypePullRequestReviewComment, EventActionCreated, "github-quick-actions", "xunleii", 0, "..."},
		},
		{
			name:      "github.PullRequestReviewCommentEvent@invalid",
			eventType: EventTypePullRequestReviewComment,
			eventJSON: []byte("...invalid..."),
			err:       fmt.Errorf("failed to extract data from JSON for event 'pull_request_review_comment': invalid character '.' looking for beginning of value"),
		},
	}

	for _, tt := range tts {
		t.Run(tt.name, func(t *testing.T) {
			payload, err := PayloadFactory(tt.eventType, tt.eventJSON)

			switch {
			case tt.err != nil && tt.expectedPayload != nil:
				t.Error("expectedPayload and err cannot be defined together")
			case tt.expectedPayload != nil:
				require.NoError(t, err)
				require.NotNil(t, payload)

				assert.Equal(t, tt.expectedPayload.Type(), payload.Type())
				assert.Equal(t, tt.expectedPayload.Action(), payload.Action())
				assert.Equal(t, tt.expectedPayload.RepositoryName(), payload.RepositoryName())
				assert.Equal(t, tt.expectedPayload.RepositoryOwner(), payload.RepositoryOwner())
				assert.Equal(t, tt.expectedPayload.IssueNumber(), payload.IssueNumber())
				assert.Equal(t, tt.expectedPayload.Body(), payload.Body())
				assert.NotNil(t, payload.Raw())
			case tt.err != nil:
				assert.EqualError(t, err, tt.err.Error())
				assert.Nil(t, payload)
			default:
				t.Error("expectedPayload or err should be specified")
			}
		})
	}
}

type mockEventPayload struct {
	eventType   EventType
	action      EventAction
	repoName    string
	repoOwner   string
	issueNumber int
	body        string
}

func (m mockEventPayload) Type() EventType         { return m.eventType }
func (m mockEventPayload) Action() EventAction     { return m.action }
func (m mockEventPayload) RepositoryName() string  { return m.repoName }
func (m mockEventPayload) RepositoryOwner() string { return m.repoOwner }
func (m mockEventPayload) IssueNumber() int        { return m.issueNumber }
func (m mockEventPayload) Body() string            { return m.body }
func (m mockEventPayload) Raw() interface{}        { panic("not implemented") }

//@format=off
