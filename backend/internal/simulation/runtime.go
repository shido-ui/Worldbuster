package simulation

import (
	"sync"
	"time"
)

type Runtime struct {
	mu sync.Mutex
	Runner *IntegratedRunner
	Scheduler TierScheduler
	LastTick time.Time
}

func NewRuntime(runner *IntegratedRunner, scheduler TierScheduler) *Runtime {
	return &Runtime{Runner: runner, Scheduler: scheduler}
}

func (r *Runtime) Tick(now time.Time) int {
	r.mu.Lock()
	defer r.mu.Unlock()
	if r.Runner == nil || r.Runner.Population == nil {
		return 0
	}
	r.LastTick = now
	return r.Runner.Tick(now)
}
