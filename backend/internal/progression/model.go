package progression

type Stats struct{Strength int `json:"strength"`;Defense int `json:"defense"`;Speed int `json:"speed"`;Intelligence int `json:"intelligence"`;Education int `json:"education"`}
type CharacterProgression struct{CharacterID string `json:"characterId"`;Level int `json:"level"`;XP int64 `json:"xp"`;Stats Stats `json:"stats"`}
type Course struct{ID string `json:"id"`;Name string `json:"name"`;DurationHours int `json:"durationHours"`;EducationGain int `json:"educationGain"`;RequiredLevel int `json:"requiredLevel"`}
type Enrollment struct{CharacterID string `json:"characterId"`;CourseID string `json:"courseId"`;StartedAt int64 `json:"startedAt"`}
