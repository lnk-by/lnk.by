package main

import (
	"context"

	"github.com/aws/aws-lambda-go/events"
	"github.com/gofrs/uuid"
	"github.com/lnk.by/aws/adapter"
	"github.com/lnk.by/aws/s3client"
	"github.com/lnk.by/shared/service"
	"github.com/lnk.by/shared/service/landingpage"
)

func deleteLandingPage(ctx context.Context, request events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {
	return adapter.DeleteAndFinalize(ctx, request, landingpage.DeleteSQL, service.IdParam, func(id uuid.UUID) error {
		return landingpage.DeleteConfiguration(ctx, id)
	}), nil
}

func main() {
	s3client.InitializeFromEnvironment(context.Background())
	adapter.LambdaMain(deleteLandingPage)
}
