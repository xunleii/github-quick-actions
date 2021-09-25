package gqa_scenario_context

import (
	"net/http"
	"net/http/httptest"

	"github.com/google/go-github/v39/github"
	"github.com/shurcooL/githubv4"
	"golang.org/x/oauth2"
)

type ClientCreator struct{ *ProxyRoundTripper }

func (cc *ClientCreator) NewAppClient() (*github.Client, error) 							 { return github.NewClient(&http.Client{Transport: cc.ProxyRoundTripper}), nil }
func (cc *ClientCreator) NewInstallationClient(installationID int64) (*github.Client, error) { return cc.NewAppClient() }
func (cc *ClientCreator) NewTokenSourceClient(ts oauth2.TokenSource) (*github.Client, error) { return cc.NewAppClient() }
func (cc *ClientCreator) NewTokenClient(token string) (*github.Client, error) 				 { return cc.NewAppClient()}

func (cc *ClientCreator) NewAppV4Client() (*githubv4.Client, error) 							 { panic("not implemented") }
func (cc *ClientCreator) NewInstallationV4Client(installationID int64) (*githubv4.Client, error) { panic("not implemented")}
func (cc *ClientCreator) NewTokenSourceV4Client(ts oauth2.TokenSource) (*githubv4.Client, error) { panic("not implemented")}
func (cc *ClientCreator) NewTokenV4Client(token string) (*githubv4.Client, error) 				 { panic("not implemented")}

type (
	// ProxyRoundTripper implements a round tripper used to intercept
	// all requests and inject responses for specific ones.
	// By default, it returns a 200.
	ProxyRoundTripper struct {
		injectedResponses   map[RoundTripperRequestKey]*httptest.ResponseRecorder
		interceptedRequests map[RoundTripperRequestKey][]*http.Request
	}
	// RoundTripperRequestKey represent a specific request
	RoundTripperRequestKey struct {
		Method string
		URL    string
	}
)

func NewProxyRoundTripper() *ProxyRoundTripper {
	return &ProxyRoundTripper{
		injectedResponses:   map[RoundTripperRequestKey]*httptest.ResponseRecorder{},
		interceptedRequests: map[RoundTripperRequestKey][]*http.Request{},
	}
}

// RoundTrip implements RoundTrip interface
func (r ProxyRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	requestKey := RoundTripperRequestKey{req.Method, req.URL.String()}

	r.interceptedRequests[requestKey] = append(r.interceptedRequests[requestKey], req)

	response, exists := r.injectedResponses[requestKey]
	if exists {
		response := response.Result()
		response.Request = req
		return response, nil
	}

	wr := httptest.NewRecorder()
	wr.Code = http.StatusOK
	return wr.Result(), nil
}

// Copy returns a copy a the current ProxyRoundTripper instance, without
// the intercepted requests.
func (r ProxyRoundTripper) Copy() *ProxyRoundTripper {
	return &ProxyRoundTripper{
		injectedResponses:   r.injectedResponses,
		interceptedRequests: map[RoundTripperRequestKey][]*http.Request{},
	}
}

// InjectResponse configures the roundtriper to use a specific response for the given request.
func (r *ProxyRoundTripper) InjectResponse(method, url string, response *httptest.ResponseRecorder) {
	requestKey := RoundTripperRequestKey{method, url}
	r.injectedResponses[requestKey] = response
}

// InterceptedRequests returns all intercepted requests for a specific method and URL.
func (r *ProxyRoundTripper) InterceptedRequests(method, url string) []*http.Request {
	requestKey := RoundTripperRequestKey{method, url}
	return r.interceptedRequests[requestKey]
}

//@format=off
