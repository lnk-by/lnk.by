package main

import (
	"context"

	"github.com/lnk.by/aws/adapter"
	"github.com/lnk.by/shared/service/stats"
)

func acceptStatistics(ctx context.Context, event stats.StatsEvent) error {
	return stats.ProcessStatistics(ctx, event)
}

func main() {
	adapter.LambdaMain(acceptStatistics)
}
