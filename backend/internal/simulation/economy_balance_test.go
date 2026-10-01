package simulation

import "testing"

func TestCalculatePriceAdjustmentDemandPressure(t *testing.T) {
	a := CalculatePriceAdjustment(EconomySignal{Supply: 10, Demand: 90, BasePrice: 100, CurrentPrice: 100})
	if a.NewPrice <= 100 || a.Direction != 1 {
		t.Fatalf("expected upward pressure: %+v", a)
	}
}

func TestCalculatePriceAdjustmentSupplyPressure(t *testing.T) {
	a := CalculatePriceAdjustment(EconomySignal{Supply: 90, Demand: 10, BasePrice: 100, CurrentPrice: 100})
	if a.NewPrice >= 100 || a.Direction != -1 {
		t.Fatalf("expected downward pressure: %+v", a)
	}
}

func TestCalculatePriceAdjustmentCapsSingleStep(t *testing.T) {
	a := CalculatePriceAdjustment(EconomySignal{Supply: 0, Demand: 100000, BasePrice: 100, CurrentPrice: 100})
	if a.NewPrice > 110 {
		t.Fatalf("price moved too far in one cycle: %+v", a)
	}
}
