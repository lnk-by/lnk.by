package main

import (
	"context"

	"github.com/aws/aws-lambda-go/events"
	"github.com/lnk.by/aws/adapter"
	"github.com/lnk.by/aws/s3client"
	"github.com/lnk.by/shared/service/landingpage"
)

func listLandingPages(ctx context.Context, request events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {
	return adapter.ListAndTransform(ctx, request, landingpage.ListSQL, func(p *landingpage.LandingPage) (*landingpage.LandingPage, error) {
		return landingpage.SetConfiguration(ctx, p)
	}), nil
}

func main() {
	s3client.InitializeFromEnvironment(context.Background())
	adapter.LambdaMain(listLandingPages)
}
