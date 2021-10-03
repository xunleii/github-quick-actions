package httptest_test

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"xnku.be/github-quick-actions/pkg/ghk_scenario_ctx/httptest"
)

func TestNewServer(t *testing.T) {
	srv := httptest.NewServer(nil)
	require.NotNil(t, srv)

	client := srv.Client()
	require.NotNil(t, client)

	_, err := client.Get("http://0.0.0.0")
	assert.EqualError(t, err, `Get "http://0.0.0.0": no handler defined on server side`)

	err = srv.Close()
	assert.Nil(t, err)
}

func TestNoNetServer_Client(t *testing.T) {
	handler := http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		assert.Equal(t, http.MethodGet, request.Method)
		assert.Equal(t, "http://0.0.0.0", request.URL.String())

		writer.WriteHeader(999)
	})

	srv := httptest.NewServer(handler)
	require.NotNil(t, srv)

	client := srv.Client()
	require.NotNil(t, client)

	resp, err := client.Get("http://0.0.0.0")
	require.NoError(t, err)

	assert.Equal(t, 999, resp.StatusCode)
}
