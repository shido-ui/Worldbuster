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
 Strength    int       `json:"strength"`
 Defense     int       `json:"defense"`
 Speed       int       `json:"speed"`
 Intelligence int      `json:"intelligence"`
 Endurance   int       `json:"endurance"`
 Education   int       `json:"education"`
 LocationID  string    `json:"locationId,omitempty"`
	UpdatedAt   time.Time `json:"updatedAt"`
}
