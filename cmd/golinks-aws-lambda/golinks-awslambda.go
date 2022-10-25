package main

import (
	"github.com/aws/aws-lambda-go/lambda"

	"github.com/belljustin/golinks/internal/golinks"
	_ "github.com/belljustin/golinks/internal/storage/dynamodb"
)

func main() {
	handler := golinks.NewLambdaHandler()
	lambda.Start(handler.Handle)
}
