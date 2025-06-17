package main

import (
	"context"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/lnk.by/aws/common"
	"github.com/lnk.by/shared/service"
	"github.com/lnk.by/shared/service/shortURL"
)

func createShortURL(ctx context.Context, requestedShortURL shortURL.ShortURL) (int, string) {
	attempts := 10
	if requestedShortURL.Key != "" {
		attempts = 1
		key := requestedShortURL.Key
		shortURL.CustomShortURLs.Put(key, true)
		defer shortURL.CustomShortURLs.Delete(key)
	}
	return service.CreateWithRetries(ctx, shortURL.CreateSQL, &requestedShortURL, attempts)
}

func main() {
	lambda.Start(common.CreateAdapter(createShortURL))
}
