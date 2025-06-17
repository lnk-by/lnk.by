package main

import (
	"context"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/lnk.by/aws/common"
	"github.com/lnk.by/shared/service"
	"github.com/lnk.by/shared/service/customer"
)

func createCustomer(ctx context.Context, requestedCustomer customer.Customer) (int, string) {
	return service.Create(ctx, customer.CreateSQL, &requestedCustomer)
}

func main() {
	lambda.Start(common.CreateAdapter(createCustomer))
}
