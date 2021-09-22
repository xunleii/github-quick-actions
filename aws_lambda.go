//go:build aws_lambda
// +build aws_lambda

package main

import (
	"os"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/prometheus/common/version"
	"github.com/rs/zerolog"

	"xnku.be/github-quick-actions/pkg/serverless"
)

func init() {
	logger := zerolog.New(os.Stdout)
	logger.Info().
		Timestamp().
		Dict("version", zerolog.Dict().
			Str("info", version.Info()).
			Str("build_context", version.BuildContext()),
		).
		Msg("Github quick actions starting on AWS Lambda")
}

var awsLambdaAdapter = serverless.NewAWSLambdaAdapter(
	serverless.LoggerFromEnvironment(),
	serverless.GithubApplicationFromEnvironment(),
)

func main() {
	lambda.Start(awsLambdaAdapter.ProxyWithContext)
}
