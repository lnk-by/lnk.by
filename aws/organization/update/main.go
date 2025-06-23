package main

import (
	"context"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/lnk.by/aws/adapter"
	"github.com/lnk.by/shared/service/organization"
)

func updateOrganization(ctx context.Context, request events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {
	return adapter.Update(ctx, request, organization.UpdateSQL, organization.IdParam), nil
}

func main() {
	lambda.Start(updateOrganization)
}
