package main

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/lnk.by/aws/adapter"
	"github.com/lnk.by/shared/service"
	"github.com/lnk.by/shared/service/short_url"
)

func handler(ctx context.Context, req events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {
	slog.Debug("Handling redirect", "RawPath", req.RawPath, "param[key]", req.PathParameters[service.IdParam])
	status, url, errStr := service.RetrieveValueAndMarshalError(ctx, short_url.RetrieveSQL, req.PathParameters[service.IdParam])
	if errStr != "" {
		return events.APIGatewayV2HTTPResponse{StatusCode: status, Body: errStr}, nil
	}

	return events.APIGatewayV2HTTPResponse{
		StatusCode: http.StatusMovedPermanently, // TODO: in future we can return 302 if the URL TTL is short or 307 or  308 if we will support methods other then GET
		Headers: map[string]string{
			"Location": url.Target,
		},
	}, nil
}

func main() {
	adapter.LambdaMain(handler)
}
