package main

import (
	"context"

	"github.com/aws/aws-lambda-go/events"
	"github.com/lnk.by/aws/adapter"
	"github.com/lnk.by/shared/service/short_url"
)

func listShortURLs(ctx context.Context, request events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {
	return adapter.List(ctx, request, short_url.ListSQL), nil
}

func main() {
	adapter.LambdaMain(listShortURLs)
}
