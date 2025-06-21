package main

import (
	"context"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/lnk.by/aws/common"
	"github.com/lnk.by/shared/service"
	"github.com/lnk.by/shared/service/organization"
)

func updateOrganization(ctx context.Context, organizationId string, requestedOrganization organization.Organization) (int, string) {
	return service.Update(ctx, organization.UpdateSQL, organizationId, &requestedOrganization)
}

func main() {
	lambda.Start(common.UpdateAdapter(updateOrganization, organization.IdParam, common.StringIDParser))
}
