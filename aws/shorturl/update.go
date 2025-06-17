package main

import (
	"context"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/lnk.by/aws/common"
	"github.com/lnk.by/shared/service"
	"github.com/lnk.by/shared/service/shortURL"
)

func updateShortURL(ctx context.Context, key string, requestedShortURL shortURL.ShortURL) (int, string) {
	return service.Update(ctx, shortURL.UpdateSQL, key, &requestedShortURL)
}

func main() {
	lambda.Start(common.UpdateAdapter(updateShortURL, shortURL.IdParam, common.StringIDParser))
}
