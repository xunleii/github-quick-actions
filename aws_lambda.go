//+build aws_lambda

package main

import (
	"github.com/aws/aws-lambda-go/lambda"

	"xnku.be/github-quick-actions/pkg/serverless"
)

func main() {
	lambda.Start(serverless.AWSLambdaHandler)
}
