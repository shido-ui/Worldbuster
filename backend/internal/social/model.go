package social

type Relationship struct {
	CharacterID string `json:"characterId"`
	TargetID    string `json:"targetId"`
	Kind        string `json:"kind"`
	Score       int    `json:"score"`
}
type Message struct {
	ID        string `json:"id"`
	FromID    string `json:"fromId"`
	ToID      string `json:"toId"`
	Body      string `json:"body"`
	CreatedAt int64  `json:"createdAt"`
	Read      bool   `json:"read"`
}
