package main

import (
	"context"

	"github.com/aws/aws-lambda-go/events"
	"github.com/lnk.by/aws/adapter"
	"github.com/lnk.by/shared/service/organization"
)

func createOrganization(ctx context.Context, request events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {
	return adapter.Create(ctx, request, organization.CreateSQL), nil
}

func main() {
	adapter.LambdaMain(createOrganization)
}
