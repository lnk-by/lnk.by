package main

import (
	"context"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/lnk.by/shared/service"
	"github.com/lnk.by/shared/service/short_url"
)

func handler(ctx context.Context, req events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {
	status, err_str, url := service.RetrieveValueAndMarshalError(ctx, short_url.RetrieveSQL, short_url.IdParam)

	switch {
	case status == http.StatusOK:
		return events.APIGatewayV2HTTPResponse{
			StatusCode: http.StatusMovedPermanently, // TODO: in future we can return 302 if the URL TTL is short or 307 or  308 if we will support methods other then GET
			Headers: map[string]string{
				"Location": url.Target,
			},
		}, nil
	default:
		return events.APIGatewayV2HTTPResponse{StatusCode: status, Body: err_str}, nil
	}
}

func main() {
	lambda.Start(handler)
}
