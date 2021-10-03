package httptest

import (
	"fmt"
	"net/http"
	"net/http/httptest"
)

// NoNetServer is an HTTP server working like other http.Server
// implementations without any interaction with any interfaces,
// for use in unit tests.
// In fact, we only want to detect want the http.Client send and not
// how it sends. This also allows us to easily intercepts http.Requests
// from a client inside an external package that we can't change the
// target URL.
type NoNetServer struct {
	// Handler to invoke on client requests
	Handler http.Handler
}

func NewServer(handler http.Handler) *NoNetServer {
	return &NoNetServer{Handler: handler}
}

// Close do nothing, just keep compatibility with other http.Server implementations.
func (s NoNetServer) Close() error { return nil }

// Client returns a http.Client linked to the NoNetServer; it will not send
// any TCP packet but will use the defined handler directly.
func (s *NoNetServer) Client() *http.Client {
	return &http.Client{Transport: &NoNetRoundTripper{s}}
}

type (
	// NoNetRoundTripper is an HTTP round tripper working likes other http.RoundTripper
	// implementations without any interactions with any interface.
	// Instead, it will directly communicate with the linked NoNetServer to send
	// the requests.
	NoNetRoundTripper struct {
		backend *NoNetServer
	}
)

// RoundTrip implements RoundTripper interface by sending the request directly to
// the server without TCP requests.
func (n NoNetRoundTripper) RoundTrip(request *http.Request) (*http.Response, error) {
	if n.backend.Handler == nil {
		return nil, fmt.Errorf("no handler defined on server side")
	}

	wr := httptest.NewRecorder()
	n.backend.Handler.ServeHTTP(wr, request)

	resp := wr.Result()
	resp.Request = request
	return resp, nil
}
