package main

import (
	"net/http"
	"os"
	"time"

	"github.com/gregjones/httpcache"
	"github.com/palantir/go-githubapp/githubapp"
	"github.com/rs/zerolog"
	"github.com/scaleway/scaleway-functions-go/events"

	"xnku.be/github-quick-actions/pkg/cmd"
)

// ScalewayFunctionEntrypoint handles Github events on Scaleway serverless platform
func ScalewayFunctionEntrypoint(req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	// define logger before anything
	logger := zerolog.New(os.Stdout).With().Timestamp().Logger()
	config, err := cmd.FromEnvironment()
	if err != nil {
		logger.Error().Err(err).Send()
		return events.APIGatewayProxyResponse{StatusCode: http.StatusInternalServerError}, err
	}
	llvl, _ := zerolog.ParseLevel(config.LogLevel)
	logger = logger.Level(llvl)

	appConfig, err := config.GithubAppConfig()
	if err != nil {
		logger.Error().Err(err).Msgf("failed to configure Github Application")
		return events.APIGatewayProxyResponse{StatusCode: http.StatusInternalServerError}, err
	}

	_, err = githubapp.NewDefaultCachingClientCreator(
		appConfig,
		githubapp.WithClientUserAgent(config.Github.UserAgent),
		githubapp.WithClientTimeout(3*time.Second),
		githubapp.WithClientCaching(false, func() httpcache.Cache { return httpcache.NewMemoryCache() }),
	)
	if err != nil {
		logger.Error().Err(err).Send()
		return events.APIGatewayProxyResponse{StatusCode: http.StatusInternalServerError}, err
	}

	logger.Info().Interface("request", req).Send()
	return events.APIGatewayProxyResponse{StatusCode: http.StatusOK}, nil
}
