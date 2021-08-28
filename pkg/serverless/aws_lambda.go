package serverless

import (
	"context"
	"encoding/json"
	"net/http"
	"os"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/awslabs/aws-lambda-go-api-proxy/gorillamux"
	"github.com/gorilla/mux"
	"github.com/gregjones/httpcache"
	"github.com/hashicorp/go-multierror"
	"github.com/palantir/go-githubapp/githubapp"
	"github.com/rs/zerolog"

	quick_actions "xnku.be/github-quick-actions/internal/quick-actions"
	"xnku.be/github-quick-actions/pkg/cmd"
	quick_action "xnku.be/github-quick-actions/pkg/gh-quick-action"
)

var muxLambda *gorillamux.GorillaMuxAdapter

func init() {
	logger := zerolog.New(os.Stdout).With().Timestamp().Logger()
	config, err := cmd.FromEnvironment()
	if err != nil {
		logger.Fatal().Err(err).Msg("failed to configure AWS lambda from environment")
	}

	llvl, _ := zerolog.ParseLevel(config.LogLevel)
	logger = logger.Level(llvl)

	// warn if GQA_LISTEN_* vars are set up (useless)
	if config.ListenAddr != "" {
		logger.Warn().Msg("CGA_LISTEN_ADDR ignored in serverless mode")
	}
	if config.ListenPath != "" {
		logger.Warn().Msg("CGA_LISTEN_PATH ignored in serverless mode")
	}

	appConfig, err := config.GithubAppConfig()
	if err != nil {
		logger.Fatal().Err(err).Msg("failed to configure Github Application")
	}

	logger.Debug().
		Interface("config", config).
		Interface("appconfig", appConfig).
		Send()
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

	logger.Info().Msgf("prepare application event dispatcher")
	app := githubapp.NewEventDispatcher(
		[]githubapp.EventHandler{issueQuickActions},
		appConfig.App.WebhookSecret,
		githubapp.WithErrorCallback(githubappErrorCallback),
	)

	r := mux.NewRouter()
	r.Handle("/", app)
	r.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
			logctx := logger.WithContext(request.Context())
			next.ServeHTTP(writer, request.WithContext(logctx))
		})
	})

	muxLambda = gorillamux.New(r)
}

func githubappErrorCallback(w http.ResponseWriter, r *http.Request, err error) {
	logger := zerolog.Ctx(r.Context())

	errs, valid := err.(*multierror.Error)
	if !valid {
		// not handled errors
		return
	}

	var errors []string
	for _, err := range errs.WrappedErrors() {
		errors = append(errors, err.Error())
	}

	json, err := json.Marshal(map[string][]string{"errors": errors})
	logger.Debug().RawJSON("body", json).Err(err).Send()
	if err != nil {
		logger.Error().Err(err).Send()
	}

	w.Header().Add("Content-Type", "application/json")
	_, _ = w.Write(json)
	w.WriteHeader(http.StatusInternalServerError)
}

func AWSLambdaHandler(ctx context.Context, event events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	zerolog.Ctx(ctx).Trace().Interface("event", event).Send()
	return muxLambda.ProxyWithContext(ctx, event)
}