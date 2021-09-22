package serverless

import (
	"context"

	"github.com/aws/aws-lambda-go/events"
	"github.com/awslabs/aws-lambda-go-api-proxy/gorillamux"
	"github.com/gorilla/mux"
)

type (
	AWSLambdaAdapter struct {
		*adapter
		mux *gorillamux.GorillaMuxAdapter
	}
)

func NewAWSLambdaAdapter(opts ...AdapterOption) Adapter {
	awsLambda := AWSLambdaAdapter{&adapter{}, nil}

	for _, opt := range opts {
		opt(awsLambda)
	}

	if awsLambda.adapter.app == nil {
		awsLambda.logger.Fatal().
			Msgf("NewAWSLambdaAdapter should be called with at least `serverless.GithubApplicationFromEnvironment` as parameter")
	}

	r := mux.NewRouter()
	r.Handle("/", awsLambda.app)
	awsLambda.mux = gorillamux.New(r)

	return awsLambda
}

func (a AWSLambdaAdapter) ProxyWithContext(ctx context.Context, event events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	a.logger.Trace().Interface("event", event).Send()
	return a.mux.ProxyWithContext(a.logger.WithContext(ctx), event)
}
