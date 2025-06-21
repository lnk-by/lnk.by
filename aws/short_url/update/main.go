package main

import (
	"context"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/lnk.by/aws/common"
	"github.com/lnk.by/shared/service"
	"github.com/lnk.by/shared/service/short_url"
)

func updateShortURL(ctx context.Context, key string, requestedShortURL short_url.ShortURL) (int, string) {
	return service.Update(ctx, short_url.UpdateSQL, key, &requestedShortURL)
}

func main() {
	lambda.Start(common.UpdateAdapter(updateShortURL, short_url.IdParam, common.StringIDParser))
}
