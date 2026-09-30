package scheduler

import (
	"context"
	"log/slog"
	"sync/atomic"
	"testing"
	"time"
)

func TestRunnerIsolatesPanicsAndErrors(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	runner := New(slog.Default())
	var healthyCalls atomic.Int32

	runner.Start(ctx,
		Job{
			Name:     "panic-job",
			Interval: time.Millisecond,
			Run: func(context.Context, time.Time) error {
				panic("boom")
			},
		},
		Job{
			Name:     "healthy-job",
			Interval: time.Millisecond,
			Run: func(context.Context, time.Time) error {
				healthyCalls.Add(1)
				return nil
			},
		},
	)

	time.Sleep(10 * time.Millisecond)
	cancel()
	runner.Wait()

	if healthyCalls.Load() == 0 {
		t.Fatal("healthy scheduler job did not continue after another job panicked")
	}
}
