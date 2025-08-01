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

func updateOrganization(ctx context.Context, request events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {
	return adapter.UpdateAndFinalize(ctx, request, landingpage.UpdateSQL, service.IdParam, func(id uuid.UUID, p *landingpage.LandingPage) error {
		_, err := landingpage.SetConfiguration(ctx, p)
		return err
	}), nil
}

func main() {
	s3client.InitializeFromEnvironment(context.Background())
	adapter.LambdaMain(updateOrganization)
}
