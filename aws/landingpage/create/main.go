package main

import (
	"context"

	"github.com/aws/aws-lambda-go/events"
	"github.com/lnk.by/aws/adapter"
	"github.com/lnk.by/aws/s3client"
	"github.com/lnk.by/shared/service"
	"github.com/lnk.by/shared/service/landingpage"
)

func createLandingPage(ctx context.Context, request events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {
	userID := service.GetUUIDFromAuthorization(request.Headers["authorization"])
	status, body := landingpage.CreateLandingPage(ctx, request.Body, userID)
	return events.APIGatewayV2HTTPResponse{StatusCode: status, Body: body, Headers: adapter.StandardHeaders}, nil
}

func main() {
	s3client.InitializeFromEnvironment(context.Background())
	adapter.LambdaMain(createLandingPage)
}
