package main

import (
	"context"

	"github.com/aws/aws-lambda-go/events"
	"github.com/lnk.by/aws/adapter"
	"github.com/lnk.by/shared/service/shorturl"
)

func listShortURLs(ctx context.Context, request events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {
	return adapter.List(ctx, request, shorturl.ListSQL), nil
}

func main() {
	adapter.LambdaMain(listShortURLs)
}
