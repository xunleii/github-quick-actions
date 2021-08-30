package mock

import (
	"github.com/google/go-github/v38/github"
	"github.com/shurcooL/githubv4"
	"github.com/stretchr/testify/mock"
	"golang.org/x/oauth2"
)

type ClientCreator struct{ mock.Mock }

func (cc *ClientCreator) NewAppClient() (*github.Client, error) {
	return cc.mock(cc.Called())
}
func (cc *ClientCreator) NewAppV4Client() (*githubv4.Client, error) {
	return cc.mockv4(cc.Called())
}
func (cc *ClientCreator) NewInstallationClient(installationID int64) (*github.Client, error) {
	return cc.mock(cc.Called(installationID))
}
func (cc *ClientCreator) NewInstallationV4Client(installationID int64) (*githubv4.Client, error) {
	return cc.mockv4(cc.Called(installationID))
}
func (cc *ClientCreator) NewTokenSourceClient(ts oauth2.TokenSource) (*github.Client, error) {
	return cc.mock(cc.Called(ts))
}
func (cc *ClientCreator) NewTokenSourceV4Client(ts oauth2.TokenSource) (*githubv4.Client, error) {
	return cc.mockv4(cc.Called(ts))
}
func (cc *ClientCreator) NewTokenClient(token string) (*github.Client, error) {
	return cc.mock(cc.Called(token))
}
func (cc *ClientCreator) NewTokenV4Client(token string) (*githubv4.Client, error) {
	return cc.mockv4(cc.Called(token))
}

func (cc *ClientCreator) mock(args mock.Arguments) (*github.Client, error) {
	return args.Get(0).(*github.Client), args.Error(1)
}
func (cc *ClientCreator) mockv4(args mock.Arguments) (*githubv4.Client, error) {
	return args.Get(0).(*githubv4.Client), args.Error(1)
}
