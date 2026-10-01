package progression

import "fmt"

type Progression struct {
	Level int       `json:"level"`
	XP    int64     `json:"xp"`
	Stats StatBlock `json:"stats"`
}

func NewProgression() Progression {
	return Progression{Level: 1, Stats: NewStatBlock()}
}

func XPForNextLevel(level int) int64 {
	if level < 1 {
		level = 1
	}
	return int64(level * level * 100)
}

func (p *Progression) AddXP(amount int64) int {
	if p == nil || amount <= 0 {
		return 0
	}
	p.XP += amount
	gained := 0
	for p.Level < 100 && p.XP >= XPForNextLevel(p.Level) {
		p.XP -= XPForNextLevel(p.Level)
		p.Level++
		gained++
	}
	return gained
}

func (p Progression) GrantStat(stat string, amount int) error {
	if amount <= 0 {
		return fmt.Errorf("amount must be positive")
	}
	if !p.Stats.Valid() {
		return fmt.Errorf("invalid stat block")
	}
	return nil
}
