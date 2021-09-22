//go:build aws_lambda
// +build aws_lambda

package main

import (
	"github.com/aws/aws-lambda-go/lambda"

	"xnku.be/github-quick-actions/pkg/serverless"
)

var awsLambdaAdapter = serverless.NewAWSLambdaAdapter(
	serverless.LoggerFromEnvironment(),
	serverless.GithubApplicationFromEnvironment(),
)

func main() {
	lambda.Start(awsLambdaAdapter.ProxyWithContext)
}
