package leaderboard_worker

import (
	"context"
	"time"

	core_logger "github.com/kupr666/Orange_Team/internal/core/logger"
)

type SnapshotService interface {
	FinalizeClosedPeriods(ctx context.Context, now time.Time) error
}

type SnapshotWorker struct {
	service  SnapshotService
	log      *core_logger.Logger
	interval time.Duration
	now      func() time.Time
}

func NewSnapshotWorker(
	service SnapshotService,
	log *core_logger.Logger,
	interval time.Duration,
) *SnapshotWorker {
	return &SnapshotWorker{
		service:  service,
		log:      log,
		interval: interval,
		now:      time.Now,
	}
}

func (w *SnapshotWorker) Run(ctx context.Context) {
	w.finalize(ctx)

	ticker := time.NewTicker(w.interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			w.finalize(ctx)
		}
	}
}

func (w *SnapshotWorker) finalize(ctx context.Context) {
	if err := w.service.FinalizeClosedPeriods(ctx, w.now()); err != nil {
		w.log.Error("finalize leaderboard snapshots", "error", err)
	}
}
