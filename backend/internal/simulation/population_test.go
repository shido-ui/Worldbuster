package simulation

import "testing"

func TestPopulationGenerator(t *testing.T) {
	g := NewPopulationGenerator(1, PopulationProfile{Names: []string{"Aster"}})
	p := g.Generate(3)
	if len(p) != 3 {
		t.Fatal("wrong population size")
	}
	if p[0].Controller != ControllerSimulated {
		t.Fatal("wrong controller")
	}
	if p[0].ID == p[1].ID {
		t.Fatal("duplicate ids")
	}
}

func TestPopulationGeneratorIsDeterministic(t *testing.T) {
	profile := PopulationProfile{
		Names:           []string{"Aster", "Mira", "Nox"},
		PersonalityBias: Personality{Curiosity: 7, RiskTolerance: 3},
	}
	first := NewPopulationGenerator(42, profile).Generate(20)
	second := NewPopulationGenerator(42, profile).Generate(20)

	for i := range first {
		if first[i].ID != second[i].ID ||
			first[i].Name != second[i].Name ||
			first[i].Personality != second[i].Personality {
			t.Fatalf("population differs at index %d", i)
		}
	}
}
