package liveness

import (
	"context"
	"log/slog"
	"time"

	"github.com/Southclaws/fault/ftag"
	"github.com/rs/xid"
	"go.uber.org/fx"

	jobs "github.com/Southclaws/storyden/app/resources/background_job"
	"github.com/Southclaws/storyden/app/resources/link/link_check"
	"github.com/Southclaws/storyden/app/resources/link/link_url"
	"github.com/Southclaws/storyden/app/services/background_job"
)

const (
	sweepInterval  = 4 * time.Hour
	maxDailyChecks = 600
	maxBatchSize   = 100
)

type Payload struct {
	Limit int `json:"limit"`
}

type Checker struct {
	links  *link_check.Repository
	jobs   *jobs.Repository
	prober Prober
	logger *slog.Logger
}

func Build() fx.Option {
	return fx.Options(
		fx.Provide(NewProber, New),
		fx.Invoke(func(*Checker) {}),
	)
}

func New(ctx context.Context, lc fx.Lifecycle, links *link_check.Repository, jobs *jobs.Repository, prober Prober, runner *background_job.Runner, logger *slog.Logger) (*Checker, error) {
	c := &Checker{links: links, jobs: jobs, prober: prober, logger: logger}
	if err := runner.Register(link_check.Kind, 1, background_job.Handler[Payload]{Run: c.run}); err != nil {
		return nil, err
	}

	run, cancel := context.WithCancel(ctx)
	done := make(chan struct{})
	lc.Append(fx.Hook{
		OnStart: func(context.Context) error {
			go func() {
				defer close(done)
				c.schedule(run)
			}()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			cancel()
			select {
			case <-done:
				return nil
			case <-ctx.Done():
				return ctx.Err()
			}
		},
	})

	return c, nil
}

func (c *Checker) enqueue(ctx context.Context, now time.Time) error {
	total, err := c.links.Count(ctx)
	if err != nil {
		return err
	}

	limit := batchBudget(total, now)
	if limit == 0 {
		return nil
	}

	candidates, err := c.links.Due(ctx, now, 1)
	if err != nil {
		return err
	}

	if len(candidates) == 0 {
		return nil
	}

	key := now.UTC().Truncate(sweepInterval).Format(time.RFC3339)
	return c.jobs.Ensure(ctx, xid.New(), link_check.Kind, key, Payload{Limit: limit})
}

func dailyBudget(total int) int {
	days := int(link_check.CheckInterval / (24 * time.Hour))
	checks := (total + days - 1) / days
	return min(maxDailyChecks, checks)
}

func batchBudget(total int, now time.Time) int {
	daily := dailyBudget(total)
	windows := int(24 * time.Hour / sweepInterval)
	window := now.UTC().Hour() / int(sweepInterval/time.Hour)
	allocatedBefore := daily * window / windows
	allocatedThrough := daily * (window + 1) / windows
	return allocatedThrough - allocatedBefore
}

func (c *Checker) schedule(ctx context.Context) {
	// TODO: Replace this ticker with recurring background-job schedules.
	ticker := time.NewTicker(sweepInterval)
	defer ticker.Stop()

	for {
		if err := c.enqueue(ctx, time.Now()); err != nil && ctx.Err() == nil {
			c.logger.WarnContext(ctx, "failed to enqueue link liveness check", "error", err)
		}

		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
		}
	}
}

func (c *Checker) run(ctx context.Context, job *jobs.Job[Payload]) (background_job.Result, error) {
	limit := min(maxBatchSize, job.Payload.Limit)
	candidates, err := c.links.Due(ctx, time.Now(), limit)
	if err != nil {
		return background_job.Result{}, err
	}

	for _, candidate := range candidates {
		if err := c.check(ctx, candidate); err != nil {
			return background_job.Result{}, err
		}
	}

	return background_job.Result{}, nil
}

func (c *Checker) check(ctx context.Context, candidate link_check.Candidate) error {
	if err := ctx.Err(); err != nil {
		return err
	}

	u, err := link_url.Parse(candidate.URL)
	if err != nil {
		return c.links.Remove(ctx, candidate.ID)
	}

	status, err := c.prober.Check(ctx, u)
	if ctx.Err() != nil {
		return ctx.Err()
	}

	if ftag.Get(err) == ftag.InvalidArgument {
		return c.links.Remove(ctx, candidate.ID)
	}

	if err != nil {
		status = link_check.Offline
	}

	return c.links.Record(ctx, candidate, status, time.Now())
}
