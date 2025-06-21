package main

import (
	"context"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/lnk.by/shared/service"
	"github.com/lnk.by/shared/service/organization"
)

func deleteOrganization(ctx context.Context, request events.APIGatewayProxyRequest) events.APIGatewayProxyResponse {
	status, body := service.Delete(ctx, organization.DeleteSQL, request.PathParameters[organization.IdParam])
	return events.APIGatewayProxyResponse{StatusCode: status, Body: body}
}

func main() {
	lambda.Start(deleteOrganization)
}
