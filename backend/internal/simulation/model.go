package simulation

type ControllerType string

const (
	ControllerHuman     ControllerType = "HUMAN"
	ControllerSimulated ControllerType = "SIMULATED"
)

type Goal struct {
	ID       string `json:"id"`
	Kind     string `json:"kind"`
	Priority int    `json:"priority"`
	TargetID string `json:"targetId,omitempty"`
}
type Personality struct {
	Curiosity     int `json:"curiosity"`
	Sociability   int `json:"sociability"`
	RiskTolerance int `json:"riskTolerance"`
	Discipline    int `json:"discipline"`
}
type SimCharacter struct {
	ID          string         `json:"id"`
	Name        string         `json:"name"`
	Controller  ControllerType `json:"controllerType"`
	Personality Personality    `json:"personality"`
	Goals       []Goal         `json:"goals"`
	LastTick    int64          `json:"lastTick"`
	Active      bool           `json:"active"`
}
