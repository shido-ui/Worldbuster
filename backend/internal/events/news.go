package events

import "strings"

type NewsItem struct {
	ID         string `json:"id"`
	EventID    string `json:"event_id"`
	Headline   string `json:"headline"`
	Body       string `json:"body"`
	Importance int    `json:"importance"`
}

type NewsGenerator struct{}

func (NewsGenerator) Generate(event WorldEvent) NewsItem {
	headline := strings.TrimSpace(event.Summary)
	if headline == "" {
		headline = "World event reported"
	}
	return NewsItem{
		ID:         "news-" + event.ID,
		EventID:    event.ID,
		Headline:   headline,
		Body:       event.Summary,
		Importance: event.Importance,
	}
}
