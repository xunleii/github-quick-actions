package v1

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/hashicorp/go-multierror"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHttpErrorCallback(t *testing.T) {
	errs := &multierror.Error{}
	errs = multierror.Append(errs, fmt.Errorf("error 01"), fmt.Errorf("error 02"))

	req := httptest.NewRequest("POST", "/", nil)
	wr := httptest.NewRecorder()
	HttpErrorCallback(wr, req, errs)

	res := wr.Result()
	body, err := io.ReadAll(res.Body)
	require.NoError(t, err)

	assert.JSONEq(t, `{"errors":  ["error 01", "error 02"]}`, string(body))
	assert.Equal(t, http.StatusInternalServerError, res.StatusCode)
	assert.Equal(t, res.Header.Get("Content-Type"), "application/json")
}

func TestHttpErrorCallback_invalid(t *testing.T) {
	req := httptest.NewRequest("POST", "/", nil)
	wr := httptest.NewRecorder()
	HttpErrorCallback(wr, req, fmt.Errorf("unhandled error type"))

	res := wr.Result()
	body, err := io.ReadAll(res.Body)
	require.NoError(t, err)

	assert.Equal(t, "unhandled error type\n", string(body))
	assert.Equal(t, http.StatusInternalServerError, res.StatusCode)
	assert.Equal(t, res.Header.Get("Content-Type"), "text/plain; charset=utf-8")
}
