package main

import (
	"context"

	"github.com/lnk.by/aws/adapter"
	"github.com/lnk.by/shared/service/stats"
)

func acceptStatistics(ctx context.Context, event stats.Event) error {
	return stats.Process(ctx, event)
}

func main() {
	adapter.LambdaMain(acceptStatistics)
}
