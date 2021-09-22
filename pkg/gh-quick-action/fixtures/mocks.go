package fixtures

import (
	"fmt"
	"net/http"

	"github.com/google/go-github/v39/github"
	"github.com/shurcooL/githubv4"
	"github.com/stretchr/testify/mock"
	"golang.org/x/oauth2"
)

type MockClientCreator struct{ mock.Mock }

const MockGithubClient string = "fixtures:github.Client"

func (cc *MockClientCreator) NewAppClient() (*github.Client, error) {
	return cc.mock(cc.Called())
}
func (cc *MockClientCreator) NewAppV4Client() (*githubv4.Client, error) {
	return cc.mockv4(cc.Called())
}
func (cc *MockClientCreator) NewInstallationClient(installationID int64) (*github.Client, error) {
	return cc.mock(cc.Called(installationID))
}
func (cc *MockClientCreator) NewInstallationV4Client(installationID int64) (*githubv4.Client, error) {
	return cc.mockv4(cc.Called(installationID))
}
func (cc *MockClientCreator) NewTokenSourceClient(ts oauth2.TokenSource) (*github.Client, error) {
	return cc.mock(cc.Called(ts))
}
func (cc *MockClientCreator) NewTokenSourceV4Client(ts oauth2.TokenSource) (*githubv4.Client, error) {
	return cc.mockv4(cc.Called(ts))
}
func (cc *MockClientCreator) NewTokenClient(token string) (*github.Client, error) {
	return cc.mock(cc.Called(token))
}
func (cc *MockClientCreator) NewTokenV4Client(token string) (*githubv4.Client, error) {
	return cc.mockv4(cc.Called(token))
}

func (cc *MockClientCreator) mock(args mock.Arguments) (*github.Client, error) {
	switch v := args.Get(0).(type) {
	case string:
		if v == MockGithubClient {
			return github.NewClient(&http.Client{Transport: &MockRoundTripper{cc}}), args.Error(1)
		}
	case *github.Client:
		return v, args.Error(1)
	}
	return nil, args.Error(1)
}
func (cc *MockClientCreator) mockv4(args mock.Arguments) (*githubv4.Client, error) {
	switch v := args.Get(0).(type) {
	case string:
		if v == MockGithubClient {
			return githubv4.NewClient(&http.Client{Transport: &MockRoundTripper{cc}}), args.Error(1)
		}
	case *githubv4.Client:
		return v, args.Error(1)
	}
	return nil, args.Error(1)
}

type MockRoundTripper struct{ *MockClientCreator }

func (r *MockRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	request := fmt.Sprintf("%s %s", req.Method, req.URL)

	args := r.MethodCalled("github.RequestValidation", request)
	if args != nil {
		args.Get(0).(func(r *http.Request))(req)
	}

	args = r.MethodCalled("github.Request", request)
	return args.Get(0).(*http.Response), args.Error(1)
}
