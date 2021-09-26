package gh_quick_actions

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/uuid"
	"github.com/hashicorp/go-multierror"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHttpErrorCallbackSimple(t *testing.T) {
	strErr := uuid.New().String()
	wr := httptest.NewRecorder()
	HttpErrorCallback(wr, nil, fmt.Errorf(strErr))

	assert.Equal(t, http.StatusInternalServerError, wr.Code)
	assert.Equal(t, strErr+"\n", wr.Body.String())
}

func TestHttpErrorCallbackMultiError(t *testing.T) {
	strErrs := []string{uuid.New().String(), uuid.New().String(), uuid.New().String()}

	errs := &multierror.Error{}
	for _, err := range strErrs {
		errs = multierror.Append(errs, fmt.Errorf(err))
	}

	wr := httptest.NewRecorder()
	HttpErrorCallback(wr, nil, errs)

	var jsonResponse struct{ Errors []string }
	err := json.Unmarshal(wr.Body.Bytes(), &jsonResponse)

	require.NoError(t, err)
	assert.Equal(t, http.StatusInternalServerError, wr.Code)
	assert.Equal(t, strErrs, jsonResponse.Errors)
}
