package main

import (
	"context"

	"github.com/aws/aws-lambda-go/events"
	"github.com/lnk.by/aws/adapter"
	"github.com/lnk.by/shared/service/short_url"
)

func createShortURL(ctx context.Context, request events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {
	return adapter.Create(ctx, request, short_url.CreateSQL), nil
}

func main() {
	adapter.LambdaMain(createShortURL)
}
