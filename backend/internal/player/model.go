package player

import "time"

type Profile struct {
	ID          string    `json:"id"`
	AccountID   string    `json:"accountId"`
	DisplayName string    `json:"displayName"`
	Level       int       `json:"level"`
	XP          int64     `json:"xp"`
	Cash        int64     `json:"cash"`
	Energy      int       `json:"energy"`
	UpdatedAt   time.Time `json:"updatedAt"`
}
