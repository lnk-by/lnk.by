package main

import (
	"context"

	"github.com/aws/aws-lambda-go/events"
	"github.com/lnk.by/aws/adapter"
	"github.com/lnk.by/shared/service"
	"github.com/lnk.by/shared/service/short_url"
)

func updateShortURL(ctx context.Context, request events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {
	return adapter.Update(ctx, request, short_url.UpdateSQL, service.IdParam), nil
}

func main() {
	adapter.LambdaMain(updateShortURL)
}
