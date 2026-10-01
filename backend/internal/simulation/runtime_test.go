package simulation

import (
	"testing"
	"time"
)

func TestRuntimeTickWithoutRunner(t *testing.T) {
	r := NewRuntime(nil, TierScheduler{})
	if got := r.Tick(time.Unix(1, 0)); got != 0 {
		t.Fatal("expected zero")
	}
}

func BenchmarkRuntimeTick1000Characters(b *testing.B) {
	population := NewService()
	generator := NewPopulationGenerator(42, PopulationProfile{Names: []string{"A", "B", "C", "D"}})
	for _, character := range generator.Generate(1000) {
		if err := population.Register(character); err != nil {
			b.Fatal(err)
		}
	}
	runner := &IntegratedRunner{Population: population, State: NewStateService(), World: &LiveAdapter{Locations: map[string]string{}}}
	runtime := NewRuntime(runner, TierScheduler{Policy: LifecyclePolicy{ActivePercent: 20, RecentPercent: 30, BackgroundPercent: 30}})
	now := time.Unix(1000, 0)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		runtime.Tick(now.Add(time.Duration(i) * time.Second))
	}
}
