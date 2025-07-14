package main

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/lnk.by/aws/adapter"
	"github.com/lnk.by/shared/service"
	"github.com/lnk.by/shared/service/shorturl"
	"github.com/lnk.by/shared/service/stats"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	lambdasdk "github.com/aws/aws-sdk-go-v2/service/lambda"
	"github.com/aws/aws-sdk-go-v2/service/lambda/types"
)

var lambdaClient *lambdasdk.Client

func init() {
	cfg, err := config.LoadDefaultConfig(context.Background())
	if err != nil {
		panic("failed to load AWS config: " + err.Error())
	}
	lambdaClient = lambdasdk.NewFromConfig(cfg)
}

func redirect(ctx context.Context, req events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {
	key := req.PathParameters[service.IdParam]
	slog.Info("Handling redirect", "RawPath", req.RawPath, "param[key]", key)
	status, url, errStr := service.RetrieveValueAndMarshalError(ctx, shorturl.RetrieveSQL, key)
	if errStr != "" {
		return events.APIGatewayV2HTTPResponse{StatusCode: status, Body: errStr}, nil
	}

	sendStatistics(ctx, key, req)

	return events.APIGatewayV2HTTPResponse{
		StatusCode: http.StatusMovedPermanently, // TODO: in future we can return 302 if the URL TTL is short or 307 or  308 if we will support methods other then GET
		Headers: map[string]string{
			"Location":      url.Target,
			"Cache-Control": "no-store, no-cache, must-revalidate, max-age=0",
		},
	}, nil
}

func sendStatistics(ctx context.Context, key string, req events.APIGatewayV2HTTPRequest) {
	event := stats.StatsEvent{
		Key:       key,
		IP:        req.RequestContext.HTTP.SourceIP,
		UserAgent: req.Headers["user-agent"],
		Referer:   req.Headers["referer"],
		Timestamp: time.Now().UTC(),
		Language:  req.Headers["accept-language"],
	}

	payload, err := json.Marshal(event)
	if err != nil {
		slog.Warn("Failed to marshal stats event", "error", err)
		return
	}

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	_, err = lambdaClient.Invoke(ctx, &lambdasdk.InvokeInput{
		FunctionName:   aws.String("aws_stats_record"),
		InvocationType: types.InvocationTypeEvent,
		Payload:        payload,
	})
	if err != nil {
		slog.Warn("Failed to invoke stats lambda", "error", err)
	}
}

func main() {
	adapter.LambdaMain(redirect)
}
