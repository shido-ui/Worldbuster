package simulation

import "testing"

func TestResolveTerritoryStableLead(t *testing.T) {
	o := ResolveTerritory(TerritoryPressure{TopInfluence: 500, SecondInfluence: 100, Stability: 50})
	if o.Stability <= 50 || o.CommerceModifierBPS <= 0 || o.SafetyModifier <= 0 {
		t.Fatalf("expected stable control effects: %+v", o)
	}
}

func TestResolveTerritoryContested(t *testing.T) {
	o := ResolveTerritory(TerritoryPressure{TopInfluence: 110, SecondInfluence: 100, Stability: 50})
	if o.Stability >= 50 || o.CommerceModifierBPS >= 0 {
		t.Fatalf("expected contested penalty: %+v", o)
	}
}
