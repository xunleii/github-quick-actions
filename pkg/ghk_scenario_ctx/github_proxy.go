package gqa_scenario_context

import (
	"net/http"
	"net/url"
	"sync"

	"github.com/gorilla/mux"
)

type GithubAPIProxy struct {
	*mux.Router

	requests APIRequests
	mx       sync.Mutex
}

func NewGithubAPIProxy() *GithubAPIProxy {
	apiProxy := &GithubAPIProxy{
		Router: mux.NewRouter(),
	}

	apiProxy.Use(apiProxy.proxyMiddleware)
	apiProxy.NotFoundHandler = apiProxy.proxyMiddleware(http.HandlerFunc(func(wr http.ResponseWriter, _ *http.Request) { wr.WriteHeader(http.StatusOK) }))

	return apiProxy
}

func (g *GithubAPIProxy) proxyMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(wr http.ResponseWriter, req *http.Request) {
		apiRequest := APIRequest{
			Request:   req,
			Method:    req.Method,
			URL:       req.URL,
			EventType: req.Header.Get(EventHeader),
			Command:   req.Header.Get(CommandHeader),
			Arguments: req.Header.Get(ArgsHeader),
		}

		g.mx.Lock()
		g.requests = append(g.requests, apiRequest)
		g.mx.Unlock()

		next.ServeHTTP(wr, req)
	})
}

func (g *GithubAPIProxy) HandledRequests() APIRequests { return g.requests }

type (
	APIRequests []APIRequest
	APIRequest  struct {
		*http.Request

		Method    string
		URL       *url.URL
		EventType string
		Command   string
		Arguments string
	}
)

func (r APIRequest) IsMetadataRequest() bool  { return r.URL.Scheme == "quick-action" }
func (r APIRequest) IsGithubAPIRequest() bool { return !r.IsMetadataRequest() }

func (r APIRequests) WithEventType(eventType string) (requests APIRequests) {
	return r.With(func(r APIRequest) bool { return r.EventType == eventType })
}
func (r APIRequests) WithCommand(command string) (requests APIRequests) {
	return r.With(func(r APIRequest) bool { return r.Command == command })
}
func (r APIRequests) WithArguments(args string) (requests APIRequests) {
	return r.With(func(r APIRequest) bool { return r.Arguments == args })
}

func (r APIRequests) With(predicate func(r APIRequest) bool) (requests APIRequests) {
	for _, request := range r {
		if predicate(request) {
			requests = append(requests, request)
		}
	}
	return requests
}
