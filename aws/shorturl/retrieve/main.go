package main

import (
	"context"

	"github.com/aws/aws-lambda-go/events"
	"github.com/lnk.by/aws/adapter"
	"github.com/lnk.by/shared/service"
	"github.com/lnk.by/shared/service/shorturl"
)

func retrieveShortURL(ctx context.Context, request events.APIGatewayV2HTTPRequest) events.APIGatewayV2HTTPResponse {
	return adapter.Retrieve(ctx, request, shorturl.RetrieveSQL, service.IdParam)
}

func main() {
	adapter.LambdaMain(retrieveShortURL)
}
