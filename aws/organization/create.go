package main

import (
	"context"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/lnk.by/aws/common"
	"github.com/lnk.by/shared/service"
	"github.com/lnk.by/shared/service/organization"
)

func createOrganization(ctx context.Context, requestedOrganization organization.Organization) (int, string) {
	return service.Create(ctx, organization.CreateSQL, &requestedOrganization)
}

func main() {
	lambda.Start(common.CreateAdapter(createOrganization))
}
