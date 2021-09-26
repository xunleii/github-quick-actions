package serverless

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/gregjones/httpcache"
	"github.com/palantir/go-githubapp/githubapp"
	"github.com/prometheus/common/version"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/pkgerrors"

	quick_actions "xnku.be/github-quick-actions/internal/quick-actions"
	"xnku.be/github-quick-actions/pkg/cmd"
	"xnku.be/github-quick-actions/pkg/gh_quick_action/v1"
)

func init() {
	logger := zerolog.New(os.Stderr).With().
		Timestamp().
		Str("version", version.Version).
		Logger()
	zerolog.DefaultContextLogger = &logger
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack

	// warn if GQA_LISTEN_* vars are set up (useless)
	if _, exists := os.LookupEnv(cmd.EnvVarListenAddr); exists {
		logger.Warn().Msgf("'%s' ignored in serverless mode", cmd.EnvVarListenAddr)
	}
	if _, exists := os.LookupEnv(cmd.EnvVarListenPath); exists {
		logger.Warn().Msgf("'%s' ignored in serverless mode", cmd.EnvVarListenPath)
	}
}

type (
	// Adapter interfaces all serverless provider with the same proxy definition
	Adapter interface {
		ProxyWithContext(ctx context.Context, event events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error)

		injectLogger(logger zerolog.Logger)
		injectGithubApp(app http.Handler)
	}
	// AdapterOption extends the configuration of an Adapter
	AdapterOption func(adapter Adapter)

	// adapter implements shared properties between Adapters
	adapter struct {
		logger *zerolog.Logger
		app    http.Handler
	}
)

func (a *adapter) injectLogger(logger zerolog.Logger) { a.logger = &logger }
func (a *adapter) injectGithubApp(app http.Handler)   { a.app = app }

func LoggerFromEnvironment() AdapterOption {
	return func(adapter Adapter) {
		level, err := zerolog.ParseLevel(os.Getenv(cmd.EnvVarLogLevel))
		if err != nil {
			level = zerolog.InfoLevel
		}

		adapter.injectLogger(
			zerolog.New(os.Stdout).With().
				Timestamp().
				Str("version", version.Version).
				Logger().Level(level),
		)
	}
}

func GithubApplicationFromEnvironment() AdapterOption {
	return func(adapter Adapter) {
		var exists bool
		var err error
		cliConfig := cmd.CLIConfig{}

		// NOTE: We don't use Kong to build the cli configuration in order
		//		 to reduce time consumption on serverless platform.
		{
			cliConfig.Github.APIVersion, exists = os.LookupEnv(cmd.EnvVarAPIVersion)
			if !exists {
				cliConfig.Github.APIVersion = "v3"
			}

			ghUrl, exists := os.LookupEnv(cmd.EnvVarAPIUrl)
			if !exists {
				ghUrl = "https://api.github.com/"
			}
			if cliConfig.Github.APIUrl, err = url.Parse(ghUrl); err != nil {
				zerolog.DefaultContextLogger.
					Fatal().Err(err).
					Msgf("invalid '%s' environment variable: %s", cmd.EnvVarAPIUrl, err)
			}

			cliConfig.Github.UserAgent, exists = os.LookupEnv(cmd.EnvVarUserAgent)
			if !exists {
				cliConfig.Github.UserAgent = fmt.Sprintf("github-quick-action/%s", version.Version)
			}

			integrationID, exists := os.LookupEnv(cmd.EnvVarIntegrationID)
			if !exists {
				zerolog.DefaultContextLogger.
					Fatal().
					Msgf("environment variable '%s' is required", cmd.EnvVarIntegrationID)
			}
			if cliConfig.Github.Application.IntegrationID, err = strconv.Atoi(integrationID); err != nil {
				zerolog.DefaultContextLogger.
					Fatal().Err(err).
					Msgf("invalid '%s' environment variable: %s", cmd.EnvVarIntegrationID, err)
			}

			cliConfig.Github.Application.Pkey, exists = os.LookupEnv(cmd.EnvVarPkey)
			if !exists {
				zerolog.DefaultContextLogger.
					Fatal().
					Msgf("environment variable '%s' is required", cmd.EnvVarPkey)
			}

			cliConfig.Github.Application.WebhookSecret, exists = os.LookupEnv(cmd.EnvVarWebhookSecret)
			if !exists {
				zerolog.DefaultContextLogger.
					Fatal().
					Msgf("environment variable '%s' is required", cmd.EnvVarWebhookSecret)
			}

		}

		appConfig, err := cliConfig.GithubAppConfig()
		if err != nil {
			zerolog.DefaultContextLogger.
				Fatal().Err(err).
				Msgf("failed to build Github Application configuration: %s", err)
		}

		cc, err := githubapp.NewDefaultCachingClientCreator(
			appConfig,
			githubapp.WithClientUserAgent(cliConfig.Github.UserAgent),
			githubapp.WithClientTimeout(3*time.Second),
			githubapp.WithClientCaching(false, func() httpcache.Cache { return httpcache.NewMemoryCache() }),
		)
		if err != nil {
			zerolog.DefaultContextLogger.
				Fatal().Err(err).
				Msgf("failed to create Github client cache: %s", err)
		}

		zerolog.DefaultContextLogger.WithLevel(zerolog.InfoLevel).
			Msgf("prepare issues/pull_requests quick actions handlers")
		// NOTE: need something to automatically injects all actions
		issueQuickActions := v1.NewGithubQuickActions(cc)
		issueQuickActions.AddQuickAction("assign", quick_actions.AssignIssueComment)
		issueQuickActions.AddQuickAction("unassign", quick_actions.UnassignIssueComment)
		issueQuickActions.AddQuickAction("label", quick_actions.LabelIssueComment)

		zerolog.DefaultContextLogger.WithLevel(zerolog.InfoLevel).
			Msgf("prepare application event dispatcher")

		adapter.injectGithubApp(githubapp.NewEventDispatcher(
			[]githubapp.EventHandler{issueQuickActions},
			appConfig.App.WebhookSecret,
			githubapp.WithErrorCallback(v1.HttpErrorCallback),
		))
	}
}
