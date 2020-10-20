package scheduler

import (
	"context"
	"time"

	"github.com/mlhamel/rougecombien/pkg/scraper"

	"github.com/mlhamel/rougecombien/pkg/config"
	"github.com/mlhamel/rougecombien/pkg/running"
	"github.com/pior/runnable"
)

const DURATION = 15 * time.Second

type Scheduler struct {
	cfg     *config.Config
	timeout time.Duration
}

func NewScheduler(cfg *config.Config) Scheduler {
	return Scheduler{cfg, DURATION}
}

func (s *Scheduler) Run(ctx context.Context) error {
	periodic := running.Periodic(s.cfg, s.timeout, scraper.NewRiviereRouge(s.cfg))

	return runnable.
		Signal(periodic).
		Run(ctx)
}
