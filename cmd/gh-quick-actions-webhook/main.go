package main

import (
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
	appv2 "xnku.be/github-quick-actions/pkg/gh_quick_action/v2"
)

func main() {
	// define logger before anything
	logger := zerolog.New(os.Stdout).With().Timestamp().Logger()

	config := cmd.CLIConfig{}
	kong.Parse(&config, kong.Vars{"version": version.Version})

	llvl, _ := zerolog.ParseLevel(config.LogLevel)
	logger = logger.Level(llvl)

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
	githubQuickActions := appv2.NewGithubQuickActions(cc)
	quick_actions.InjectAll(githubQuickActions)

	app := githubapp.NewEventDispatcher(
		[]githubapp.EventHandler{githubQuickActions},
		appConfig.App.WebhookSecret,
		githubapp.WithErrorCallback(appv2.HttpErrorCallback),
	)

	r := mux.NewRouter()
	r.Handle(config.ListenPath, app)
	r.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
			logctx := logger.WithContext(request.Context())
			next.ServeHTTP(writer, request.WithContext(logctx))
		})
	})

	srv := &http.Server{
		Handler: r,
		Addr:    config.ListenAddr,

		// Good practice: enforce timeouts for servers you create!
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	logger.Info().Msgf("start listening on %s%s", config.ListenAddr, config.ListenPath)
	logger.Fatal().Err(srv.ListenAndServe())
}
