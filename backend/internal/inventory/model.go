package inventory

type Item struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	Category  string `json:"category"`
	Stackable bool   `json:"stackable"`
	MaxStack  int    `json:"maxStack"`
}

type Stack struct {
	ItemID   string `json:"itemId"`
	Quantity int    `json:"quantity"`
}

type Inventory struct {
	CharacterID string  `json:"characterId"`
	Capacity    int     `json:"capacity"`
	Stacks      []Stack `json:"stacks"`
}
