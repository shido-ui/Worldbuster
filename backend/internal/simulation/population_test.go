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
