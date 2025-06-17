package main

import (
	"context"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/lnk.by/aws/common"
	"github.com/lnk.by/shared/service"
	"github.com/lnk.by/shared/service/user"
)

func createUser(ctx context.Context, requestedUser user.User) (int, string) {
	return service.Create(ctx, user.CreateSQL, &requestedUser)
}

func main() {
	lambda.Start(common.CreateAdapter(createUser))
}
