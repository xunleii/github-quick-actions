package main

import (
	"context"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
	"github.com/gregjones/httpcache"
	"github.com/palantir/go-githubapp/githubapp"
	"github.com/prometheus/common/version"
	"github.com/rs/zerolog"
	"gopkg.in/alecthomas/kingpin.v2"

	quick_actions "xnku.be/github-quick-actions/internal/quick-actions"
	quick_action "xnku.be/github-quick-actions/pkg/gh-quick-action"
)

func main() {
	var (
		githubAPIVersion = kingpin.Flag("github.api_version", "Github API version").
					Default("v3").
					Enum("v3", "v4")
		githubAPIUrl = kingpin.Flag("github.addr", "Github API url").
				Default("https://api.github.com/").
				URL()
		githubIntegrationID = kingpin.Flag("github.app_id", "Github Application ID").
					Required().
					Int()
		githubWebhookSecret = kingpin.Flag("github.webhook_secret", "Github Webhook secret").
					Required().
					String()
		githubPkey = kingpin.Flag("github.pkey", "Github Application private key path").
				Required().
				OpenFile(os.O_RDONLY, 0644)

		listenAddr = kingpin.Flag("listen.addr", "Webhook listening address").
				Default("localhost:3000").
				TCP()
		listenPath = kingpin.Flag("listen.path", "Webhook listening path").
				Default("/api/v1/webhook").
				String()
		logLevel = kingpin.Flag("log.level", "Log level verbosity").
				Default("info").
				Enum("trace", "debug", "info", "warn", "error", "fatal", "panic")
		userAgent = kingpin.Flag("user_agent", "User agent used on Github requests").
				Default(fmt.Sprintf("github-quick-action/%s", version.Version)).
				String()
	)
	kingpin.Version(version.Info())
	kingpin.CommandLine.UsageWriter(os.Stderr)
	kingpin.HelpFlag.Short('h')

	kingpin.Parse()

	// NOTE: define logger before logging possible errors
	llvl, _ := zerolog.ParseLevel(*logLevel)
	logger := zerolog.New(os.Stdout).
		With().Timestamp().Logger().
		Level(llvl)

	config := githubapp.Config{}
	switch *githubAPIVersion {
	case "v3":
		config.V3APIURL = (*githubAPIUrl).String()
	case "v4":
		config.V4APIURL = (*githubAPIUrl).String()
	}

	config.App.IntegrationID = int64(*githubIntegrationID)
	config.App.WebhookSecret = *githubWebhookSecret

	pkey, err := io.ReadAll(*githubPkey)
	if err != nil {
		logger.Fatal().Err(err).Msgf("failed to read the given Github private key '%s'", (*githubPkey).Name())
	}
	config.App.PrivateKey = string(pkey)

	cc, err := githubapp.NewDefaultCachingClientCreator(
		config,
		githubapp.WithClientUserAgent(*userAgent),
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

	webhookHandler := githubapp.NewDefaultEventDispatcher(config, issueQuickActions)

	r := mux.NewRouter()
	r.Handle(*listenPath, webhookHandler)

	ctx := logger.WithContext(context.Background())
	srv := &http.Server{
		BaseContext: func(_ net.Listener) context.Context { return ctx },
		Handler:     http.HandlerFunc(func(wr http.ResponseWriter, rq *http.Request) { r.ServeHTTP(wr, rq.WithContext(ctx)) }),
		Addr:        (*listenAddr).String(),

		// Good practice: enforce timeouts for servers you create!
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	logger.Info().Msgf("start listening on %s%s", (*listenAddr).String(), *listenPath)
	logger.Fatal().Err(srv.ListenAndServe())
}
