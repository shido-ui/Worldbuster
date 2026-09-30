package events

import "time"

type Event struct{ID string `json:"id"`;Type string `json:"type"`;ActorID string `json:"actorId,omitempty"`;TargetID string `json:"targetId,omitempty"`;Payload map[string]any `json:"payload,omitempty"`;CreatedAt time.Time `json:"createdAt"`}
type Notification struct{ID string `json:"id"`;CharacterID string `json:"characterId"`;Type string `json:"type"`;Title string `json:"title"`;Body string `json:"body"`;EventID string `json:"eventId,omitempty"`;CreatedAt time.Time `json:"createdAt"`;Read bool `json:"read"`}
