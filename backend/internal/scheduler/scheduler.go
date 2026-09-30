package scheduler

import (
	"context"
	"log/slog"
	"sync"
	"time"
)

type Job struct {
	Name     string
	Interval time.Duration
	Run      func(context.Context, time.Time) error
}

type Runner struct {
	logger *slog.Logger
	wg     sync.WaitGroup
}

func New(logger *slog.Logger) *Runner {
	if logger == nil {
		logger = slog.Default()
	}
	return &Runner{logger: logger}
}

func (r *Runner) Start(ctx context.Context, jobs ...Job) {
	for _, job := range jobs {
		if job.Name == "" || job.Interval <= 0 || job.Run == nil {
			r.logger.Warn("skipping invalid scheduler job", "name", job.Name)
			continue
		}
		j := job
		r.wg.Add(1)
		go r.run(ctx, j)
	}
}

func (r *Runner) Wait() {
	r.wg.Wait()
}

func (r *Runner) run(ctx context.Context, job Job) {
	defer r.wg.Done()

	ticker := time.NewTicker(job.Interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case now := <-ticker.C:
			if err := runSafely(ctx, job, now); err != nil {
				r.logger.Error("scheduler job failed", "job", job.Name, "error", err)
			}
		}
	}
}

func runSafely(ctx context.Context, job Job, now time.Time) (err error) {
	defer func() {
		if recovered := recover(); recovered != nil {
			err = panicError{value: recovered}
		}
	}()
	return job.Run(ctx, now)
}

type panicError struct {
	value any
}

func (e panicError) Error() string {
	return "scheduler job panicked"
}
