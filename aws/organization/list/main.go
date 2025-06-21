package main

import (
	"context"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/lnk.by/aws/common"
	"github.com/lnk.by/shared/service/organization"
)

func listOrganizations(ctx context.Context, request events.APIGatewayProxyRequest) events.APIGatewayProxyResponse {
	return common.ListEntities(ctx, request, organization.ListSQL)
}

func main() {
	lambda.Start(listOrganizations)
}
