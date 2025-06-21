package main

import (
	"context"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/lnk.by/shared/service"
	"github.com/lnk.by/shared/service/campaign"
)

func retrieveCampaign(ctx context.Context, request events.APIGatewayProxyRequest) events.APIGatewayProxyResponse {
	status, body := service.Retrieve(ctx, campaign.RetrieveSQL, request.PathParameters[campaign.IdParam])
	return events.APIGatewayProxyResponse{StatusCode: status, Body: body}
}

func main() {
	lambda.Start(retrieveCampaign)
}
