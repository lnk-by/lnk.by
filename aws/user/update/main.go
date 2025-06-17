package main

import (
	"context"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/lnk.by/aws/common"
	"github.com/lnk.by/shared/service"
	"github.com/lnk.by/shared/service/user"
)

func updateUser(ctx context.Context, userId string, requestedUser user.User) (int, string) {
	requestedUser.ID = userId
	return service.Update(ctx, user.UpdateSQL, &requestedUser)
}

func main() {
	lambda.Start(common.UpdateAdapter(updateUser, user.IdParam, common.StringIDParser))
}
