package events

import (
	"sort"
	"sync"
	"time"
)

type WorldEvent struct {
	ID         string    `json:"id"`
	Type       string    `json:"type"`
	ActorID    string    `json:"actor_id,omitempty"`
	TargetID   string    `json:"target_id,omitempty"`
	Summary    string    `json:"summary"`
	At         time.Time `json:"at"`
	Importance int       `json:"importance"`
}

type History struct {
	mu     sync.RWMutex
	events []WorldEvent
}

func NewHistory() *History { return &History{events: make([]WorldEvent, 0, 256)} }

func (h *History) Record(event WorldEvent) {
	if h == nil || event.ID == "" {
		return
	}
	h.mu.Lock()
	defer h.mu.Unlock()
	h.events = append(h.events, event)
	if len(h.events) > 5000 {
		h.events = h.events[len(h.events)-5000:]
	}
}

func (h *History) Recent(limit int) []WorldEvent {
	if h == nil || limit <= 0 {
		return nil
	}
	h.mu.RLock()
	defer h.mu.RUnlock()
	n := limit
	if n > len(h.events) {
		n = len(h.events)
	}
	out := append([]WorldEvent(nil), h.events[len(h.events)-n:]...)
	sort.SliceStable(out, func(i, j int) bool { return out[i].At.Before(out[j].At) })
	return out
}

func (h *History) Important(limit int) []WorldEvent {
	if h == nil || limit <= 0 {
		return nil
	}
	h.mu.RLock()
	defer h.mu.RUnlock()
	out := make([]WorldEvent, 0, limit)
	for i := len(h.events) - 1; i >= 0 && len(out) < limit; i-- {
		if h.events[i].Importance >= 70 {
			out = append(out, h.events[i])
		}
	}
	sort.SliceStable(out, func(i, j int) bool { return out[i].At.Before(out[j].At) })
	return out
}
