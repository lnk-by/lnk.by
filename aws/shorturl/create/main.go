package main

import (
	"context"

	"github.com/aws/aws-lambda-go/events"
	"github.com/lnk.by/aws/adapter"
)

func createShortURL(ctx context.Context, request events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {
	return adapter.CreateShortURL(ctx, request), nil
}

func main() {
	adapter.LambdaMain(createShortURL)
}
