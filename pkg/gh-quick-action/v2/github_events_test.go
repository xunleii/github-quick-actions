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
	issueCommentEventFixture = &github.IssueCommentEvent{
		Action:  github.String("created"),
		Repo:    &github.Repository{Name: github.String("github-quick-actions"), Owner: &github.User{Login: github.String("xunleii")}},
		Comment: &github.IssueComment{Body: github.String("...")},
	}
	issueCommentEventFixtureJSON, _ = json.Marshal(issueCommentEventFixture)
)

func TestPayloadFactory(t *testing.T) {
	tts := []struct {
		name      string
		eventType EventType
		eventJSON []byte
		expectedPayload EventPayload
		err             error
	}{
		{
			name:      "unknown event type",
			eventType: EventType("unknown"),
			err:       fmt.Errorf("event type 'unknown' not managed"),
		},

		{
			name:            "github.IssueCommentEvent",
			eventType:       EventTypeIssueComment,
			eventJSON:       issueCommentEventFixtureJSON,
			expectedPayload: mockEventPayload{EventTypeIssueComment, EventActionCreated, "github-quick-actions", "xunleii", "..."},
		},
		{
			name:      "github.IssueCommentEvent@invalid",
			eventType: EventTypeIssueComment,
			eventJSON: []byte("...invalid..."),
			err:       fmt.Errorf("failed to extract data from JSON for event 'issue_comment': invalid character '.' looking for beginning of value"),
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
	eventType EventType
	action    eventAction
	repoName  string
	repoOwner string
	body      string
}

func (m mockEventPayload) Type() EventType     { return m.eventType }
func (m mockEventPayload) Action() eventAction { return m.action }
func (m mockEventPayload) RepositoryName() string  { return m.repoName }
func (m mockEventPayload) RepositoryOwner() string { return m.repoOwner }
func (m mockEventPayload) Body() string            { return m.body }
func (m mockEventPayload) Raw() interface{}        { panic("not implemented") }

//@format=off
