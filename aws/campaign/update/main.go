package main

import (
	"context"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/lnk.by/aws/common"
	"github.com/lnk.by/shared/service"
	"github.com/lnk.by/shared/service/customer"
)

func updateCustomer(ctx context.Context, customerId string, requestedCustomer customer.Customer) (int, string) {
	return service.Update(ctx, customer.UpdateSQL, customerId, &requestedCustomer)
}

func main() {
	lambda.Start(common.UpdateAdapter(updateCustomer, customer.IdParam, common.StringIDParser))
}
