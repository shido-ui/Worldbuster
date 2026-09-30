package simulation

type TerritoryPressure struct {
 TotalInfluence int
 TopInfluence int
 SecondInfluence int
 Stability int
}

type TerritoryOutcome struct {
 ControllerChanged bool
 Stability int
 CommerceModifierBPS int
 SafetyModifier int
}

func ResolveTerritory(p TerritoryPressure) TerritoryOutcome {
 stability:=p.Stability
 if stability<0 {stability=0}; if stability>100 {stability=100}
 gap:=p.TopInfluence-p.SecondInfluence
 if p.TopInfluence<=0 { return TerritoryOutcome{Stability:stability,CommerceModifierBPS:-500,SafetyModifier:-10} }
 if gap<25 {stability-=2} else if gap>150 {stability+=2}
 if stability<0 {stability=0}; if stability>100 {stability=100}
 commerce:=0
 safety:=0
 if gap>=150 {commerce=500;safety=10}
 if gap>=300 {commerce=1000;safety=20}
 if gap<25 {commerce=-500;safety=-10}
 return TerritoryOutcome{Stability:stability,CommerceModifierBPS:commerce,SafetyModifier:safety}
}
