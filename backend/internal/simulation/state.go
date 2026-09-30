package simulation

type Needs struct{Energy int `json:"energy"`;Social int `json:"social"`;Rest int `json:"rest"`;Satisfaction int `json:"satisfaction"`}
type Memory struct{ID string `json:"id"`;Type string `json:"type"`;TargetID string `json:"targetId,omitempty"`;Event string `json:"event"`;Importance int `json:"importance"`;Sentiment int `json:"sentiment"`;CreatedAt int64 `json:"createdAt"`}
type BehavioralState struct{CurrentAction ActionType `json:"currentAction"`;Stress int `json:"stress"`;Mood int `json:"mood"`;Needs Needs `json:"needs"`;Memories []Memory `json:"memories"`}
