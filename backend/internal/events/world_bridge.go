package events

import "time"

type WorldEventSink struct {
	History   *History
	News      *NewsGenerator
	Published []NewsItem
}

func (s *WorldEventSink) Publish(event WorldEvent) (NewsItem, bool) {
	if s == nil || s.History == nil || event.ID == "" {
		return NewsItem{}, false
	}
	s.History.Record(event)
	if s.News == nil {
		return NewsItem{}, true
	}
	item := s.News.Generate(event)
	s.Published = append(s.Published, item)
	if len(s.Published) > 1000 {
		s.Published = s.Published[len(s.Published)-1000:]
	}
	return item, true
}

func (s *WorldEventSink) PublishSimulation(eventType, actorID, targetID, summary string, importance int, at time.Time) (NewsItem, bool) {
	return s.Publish(WorldEvent{
		ID:   eventType + "-" + actorID + "-" + at.UTC().Format("20060102150405.000000000"),
		Type: eventType, ActorID: actorID, TargetID: targetID,
		Summary: summary, Importance: importance, At: at,
	})
}
