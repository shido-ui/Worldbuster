package world

import (
	"sync"
	"time"
)

type Service struct {
	mu          sync.RWMutex
	tick        int64
	onlineCount int
	startedAt   time.Time
}

func NewService() *Service {
	return &Service{startedAt: time.Now()}
}

func (s *Service) Tick() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.tick++
}

func (s *Service) Snapshot() WorldState {
	s.mu.RLock()
	defer s.mu.RUnlock()
	elapsed := time.Since(s.startedAt)
	hours := int(elapsed.Hours()) % 24
	minutes := int(elapsed.Minutes()) % 60
	return WorldState{
		Tick: s.tick, OnlineCount: s.onlineCount,
		Day: int(elapsed.Hours()/24) + 1,
		Time: time.Date(2000,1,1,hours,minutes,0,0,time.UTC).Format("15:04"),
		Status: "ONLINE",
	}
}
