package main

import (
	"context"
	"log/slog"

	"github.com/lnk.by/aws/adapter"
	"github.com/lnk.by/shared/service/stats"
	"github.com/lnk.by/shared/service/stats/maxmind"
)

func acceptStatistics(ctx context.Context, event stats.Event) error {
	return stats.Process(ctx, event)
}

func main() {
	if err := maxmind.Init(); err != nil {
		slog.Error("Failed to intialize mixmind", "error", err)
	}
	adapter.LambdaMain(acceptStatistics)
}
