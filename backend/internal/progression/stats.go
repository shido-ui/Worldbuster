package progression

type StatBlock struct {
	Strength     int `json:"strength"`
	Defense      int `json:"defense"`
	Speed        int `json:"speed"`
	Intelligence int `json:"intelligence"`
	Endurance    int `json:"endurance"`
}

func NewStatBlock() StatBlock {
	return StatBlock{Strength: 1, Defense: 1, Speed: 1, Intelligence: 1, Endurance: 1}
}

func (s StatBlock) Total() int {
	return s.Strength + s.Defense + s.Speed + s.Intelligence + s.Endurance
}

func (s StatBlock) Valid() bool {
	return s.Strength >= 0 && s.Defense >= 0 && s.Speed >= 0 && s.Intelligence >= 0 && s.Endurance >= 0
}

func (s *StatBlock) Add(stat string, amount int) bool {
	if s == nil || amount == 0 {
		return false
	}
	switch stat {
	case "strength":
		s.Strength += amount
	case "defense":
		s.Defense += amount
	case "speed":
		s.Speed += amount
	case "intelligence":
		s.Intelligence += amount
	case "endurance":
		s.Endurance += amount
	default:
		return false
	}
	return s.Valid()
}
