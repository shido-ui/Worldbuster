package world

import "time"

type ControllerType string

const (
	ControllerHuman     ControllerType = "HUMAN"
	ControllerSimulated ControllerType = "SIMULATED"
)

type Character struct {
	ID             string         `json:"id"`
	DisplayName    string         `json:"displayName"`
	ControllerType ControllerType `json:"controllerType"`
	Level          int            `json:"level"`
	XP             int64          `json:"xp"`
	Cash           int64          `json:"cash"`
	Energy         int            `json:"energy"`
	CreatedAt      time.Time      `json:"createdAt"`
}

type WorldState struct {
	Tick        int64  `json:"tick"`
	OnlineCount int    `json:"onlineCount"`
	Day         int    `json:"day"`
	Time        string `json:"time"`
	Status      string `json:"status"`
}
