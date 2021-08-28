package main

import (
	"context"
	"net"
	"net/http"
	"os"
	"time"

	"github.com/alecthomas/kong"
	"github.com/gorilla/mux"
	"github.com/gregjones/httpcache"
	"github.com/palantir/go-githubapp/githubapp"
	"github.com/prometheus/common/version"
	"github.com/rs/zerolog"

	quick_actions "xnku.be/github-quick-actions/internal/quick-actions"
	"xnku.be/github-quick-actions/pkg/cmd"
	quick_action "xnku.be/github-quick-actions/pkg/gh-quick-action"
)

func main() {
	config := cmd.CLIConfig{}
	kong.Parse(&config, kong.Vars{"version": version.Info()})

	// NOTE: define logger before logging possible errors
	llvl, _ := zerolog.ParseLevel(config.LogLevel)
	logger := zerolog.New(os.Stdout).
		With().Timestamp().Logger().
		Level(llvl)

	appConfig, err := config.GithubAppConfig()
	if err != nil {
		logger.Fatal().Err(err).Msgf("failed to configure Github Application")
	}

	cc, err := githubapp.NewDefaultCachingClientCreator(
		appConfig,
		githubapp.WithClientUserAgent(config.Github.UserAgent),
		githubapp.WithClientTimeout(3*time.Second),
		githubapp.WithClientCaching(false, func() httpcache.Cache { return httpcache.NewMemoryCache() }),
	)
	if err != nil {
		logger.Fatal().Err(err).Send()
	}

	logger.Info().Msgf("prepare issues/pull_requests quick actions handlers")
	issueQuickActions := quick_action.NewIssueQuickActions(cc)
	issueQuickActions.AddQuickAction("assign", quick_actions.Assign)
	issueQuickActions.AddQuickAction("unassign", quick_actions.Unassign)

	webhookHandler := githubapp.NewDefaultEventDispatcher(appConfig, issueQuickActions)

	r := mux.NewRouter()
	r.Handle(config.ListenPath, webhookHandler)

	ctx := logger.WithContext(context.Background())
	srv := &http.Server{
		BaseContext: func(_ net.Listener) context.Context { return ctx },
		Handler:     http.HandlerFunc(func(wr http.ResponseWriter, rq *http.Request) { r.ServeHTTP(wr, rq.WithContext(ctx)) }),
		Addr:        config.ListenAddr,

		// Good practice: enforce timeouts for servers you create!
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	logger.Info().Msgf("start listening on %s%s", config.ListenAddr, config.ListenPath)
	logger.Fatal().Err(srv.ListenAndServe())
}
