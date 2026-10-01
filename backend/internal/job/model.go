package job

type Position struct {
	ID                string `json:"id"`
	JobID             string `json:"jobId"`
	Name              string `json:"name"`
	Level             int    `json:"level"`
	BaseSalary        int64  `json:"baseSalary"`
	RequiredEducation int    `json:"requiredEducation"`
	RequiredStat      int    `json:"requiredStat"`
}
type Job struct {
	ID          string     `json:"id"`
	Name        string     `json:"name"`
	Department  string     `json:"department"`
	Description string     `json:"description"`
	Positions   []Position `json:"positions"`
}
type Employment struct {
	CharacterID string `json:"characterId"`
	JobID       string `json:"jobId"`
	PositionID  string `json:"positionId"`
	Experience  int    `json:"experience"`
}
