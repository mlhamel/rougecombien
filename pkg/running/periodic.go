package running

import (
	"context"
	"time"

	"github.com/mlhamel/rougecombien/pkg/config"
	"github.com/pior/runnable"
)

// Periodic is a constructor for creating a periodic runnable
func Periodic(cfg *config.Config, duration time.Duration, runner runnable.Runnable) runnable.Runnable {
	return &periodicRunnable{cfg, duration, runner}
}

type periodicRunnable struct {
	cfg      *config.Config
	duration time.Duration
	runner   runnable.Runnable
}

func (p *periodicRunnable) Run(ctx context.Context) error {
	p.cfg.Logger().Info().Msg("Initializing periodic runner")

	errs := make(chan error, 1)
	done := make(chan bool)

	defer close(errs)
	defer close(done)

	go p.runForever(ctx, errs, done)

	select {
	case <-ctx.Done():
		done <- true
		return nil
	case err := <-errs:
		return err
	}
}

func (p *periodicRunnable) runForever(ctx context.Context, errs chan<- error, done chan bool) {
	ticker := time.NewTicker(p.duration)
	defer ticker.Stop()

	for {
		p.cfg.Logger().Info().Str("duration", fmtDuration(p.duration)).Msg("Sleeping")
		select {
		case <-done:
			return
		case <-ticker.C:
			p.cfg.Logger().Info().Msg("Running")
			if err := p.runner.Run(ctx); err != nil {
				errs <- err
				return
			}
		}
	}
}
