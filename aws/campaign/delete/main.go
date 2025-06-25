package main

import (
	"context"

	"github.com/aws/aws-lambda-go/events"
	"github.com/lnk.by/aws/adapter"
	"github.com/lnk.by/shared/service"
	"github.com/lnk.by/shared/service/campaign"
)

func deleteCampaign(ctx context.Context, request events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {
	return adapter.Delete(ctx, request, campaign.DeleteSQL, service.IdParam), nil
}

func main() {
	adapter.LambdaMain(deleteCampaign)
}
