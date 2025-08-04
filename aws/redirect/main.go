package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"strconv"
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

	now := time.Now()
	day := fmt.Sprintf("day%03d", now.YearDay())
	hour := fmt.Sprintf("hour%02d", now.Hour())

	status, url, errStr := service.RetrieveValueAndMarshalError(ctx, shorturl.RetrieveValidSQL, key, day, hour)
	if errStr != "" {
		return events.APIGatewayV2HTTPResponse{StatusCode: status, Body: errStr}, nil
	}
	if limitExceeded, retryAfter := shorturl.GetLimitExceededMessage(url); limitExceeded != "" {
		headers := map[string]string{"Content-Type": "application/json"}
		if retryAfter > 0 {
			headers["Retry-After"] = strconv.Itoa(retryAfter)
		}

		return events.APIGatewayV2HTTPResponse{
			StatusCode: http.StatusTooManyRequests,
			Headers:    headers,
			Body:       string(mustJSON(map[string]string{"error": limitExceeded})),
		}, nil
	}

	if err := sendStatistics(ctx, key, req); err != nil {
		slog.Warn("Failed to send stats", "error", err)
	}

	return events.APIGatewayV2HTTPResponse{
		StatusCode: http.StatusMovedPermanently, // TODO: in future we can return 302 if the URL TTL is short or 307 or  308 if we will support methods other then GET
		Headers: map[string]string{
			"Location":      url.Target,
			"Cache-Control": "no-store, no-cache, must-revalidate, max-age=0",
		},
	}, nil
}

func mustJSON(v any) []byte {
	b, err := json.Marshal(v)
	if err != nil {
		panic("failed to marshal JSON: " + err.Error())
	}
	return b
}

func sendStatistics(ctx context.Context, key string, req events.APIGatewayV2HTTPRequest) error {
	event := stats.Event{
		Key:       key,
		IP:        req.RequestContext.HTTP.SourceIP,
		UserAgent: req.Headers["user-agent"],
		Referer:   req.Headers["referer"],
		Timestamp: time.Now().UTC(),
		Language:  req.Headers["accept-language"],
	}

	payload, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("failed to marshal stats event %v: %w", event, err)
	}

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	_, err = lambdaClient.Invoke(ctx, &lambdasdk.InvokeInput{
		FunctionName:   aws.String("aws_stats_record"),
		InvocationType: types.InvocationTypeEvent,
		Payload:        payload,
	})

	if err != nil {
		return fmt.Errorf("failed to invoke stats lambda: %w", err)
	}

	return nil
}

func main() {
	adapter.LambdaMain(redirect)
}
