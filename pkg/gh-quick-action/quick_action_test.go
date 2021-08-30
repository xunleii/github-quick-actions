package quick_action

import (
	"context"

	"github.com/stretchr/testify/mock"
)

type MockGithubQuickActionHandler struct{ mock.Mock }

func (h *MockGithubQuickActionHandler) Fnc(ctx context.Context, payload GithubQuickActionEvent) error {
	return h.Called(ctx, payload.Arguments()).Error(0)
}
