package main

import (
	"context"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/lnk.by/aws/common"
	"github.com/lnk.by/shared/service"
	"github.com/lnk.by/shared/service/short_url"
)

func createShortURL(ctx context.Context, requestedShortURL short_url.ShortURL) (int, string) {
	attempts := 10
	if requestedShortURL.Key != "" {
		attempts = 1
		key := requestedShortURL.Key
		short_url.Custom.Put(key, true)
		defer short_url.Custom.Delete(key)
	}
	return service.CreateWithRetries(ctx, short_url.CreateSQL, &requestedShortURL, attempts)
}

func main() {
	lambda.Start(common.CreateAdapter(createShortURL))
}
