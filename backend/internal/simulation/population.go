package simulation

import (
	"crypto/sha256"
	"fmt"
	mrand "math/rand"
	"sync"
)

type PopulationProfile struct {
	Names           []string
	JobWeights      map[string]int
	PersonalityBias Personality
}

type PopulationGenerator struct {
	mu      sync.Mutex
	next    int
	seed    int64
	rng     *mrand.Rand
	profile PopulationProfile
}

func NewPopulationGenerator(seed int64, profile PopulationProfile) *PopulationGenerator {
	return &PopulationGenerator{
		seed:    seed,
		rng:     mrand.New(mrand.NewSource(seed)),
		profile: profile,
	}
}

func (g *PopulationGenerator) Generate(count int) []SimCharacter {
	g.mu.Lock()
	defer g.mu.Unlock()

	out := make([]SimCharacter, 0, count)
	for i := 0; i < count; i++ {
		g.next++
		name := "Citizen"
		if len(g.profile.Names) > 0 {
			name = g.profile.Names[g.rng.Intn(len(g.profile.Names))]
		}
		out = append(out, SimCharacter{
			ID:           deterministicCharacterID(g.seed, g.next),
			Name:         name + " " + formatID(g.next),
			Controller:   ControllerSimulated,
			Personality:  g.profile.PersonalityBias,
			Goals:        []Goal{{ID: "work", Kind: "WORK", Priority: 50}},
			Active:       false,
		})
	}
	return out
}

func deterministicCharacterID(seed int64, sequence int) string {
	input := fmt.Sprintf("worldbuster:npc:%d:%d", seed, sequence)
	sum := sha256.Sum256([]byte(input))
	b := sum[:16]
	b[6] = (b[6] & 0x0f) | 0x40
	b[8] = (b[8] & 0x3f) | 0x80
	return fmt.Sprintf("%08x-%04x-%04x-%04x-%012x", b[0:4], b[4:6], b[6:8], b[8:10], b[10:16])
}

func formatID(n int) string {
	digits := "0123456789"
	s := ""
	if n == 0 {
		return "0"
	}
	for n > 0 {
		s = string(digits[n%10]) + s
		n /= 10
	}
	return s
}
