package simulation

import "math"

type CombatInput struct {
 AttackerStrength,AttackerSpeed,AttackerDefense,AttackerEndurance int
 DefenderStrength,DefenderSpeed,DefenderDefense,DefenderEndurance int
 AttackerLevel,DefenderLevel int
}

type CombatResult struct {
 AttackerScore,DefenderScore int
 Winner int
 EnergyCost int
 XPReward int64
}

func ResolveCombat(in CombatInput) CombatResult {
 a:=in.AttackerStrength*3+in.AttackerSpeed*2+in.AttackerDefense+in.AttackerEndurance+in.AttackerLevel*5
 d:=in.DefenderStrength*3+in.DefenderSpeed*2+in.DefenderDefense+in.DefenderEndurance+in.DefenderLevel*5
 if a<1 {a=1};if d<1 {d=1}
 // Deterministic outcome: identical state always produces identical result.
 margin:=int(math.Abs(float64(a-d)))
 winner:=0;if a>d{winner=1}else if d>a{winner=2}
 xp:=int64(25+margin/10);if xp>250{xp=250}
 return CombatResult{AttackerScore:a,DefenderScore:d,Winner:winner,EnergyCost:10,XPReward:xp}
}
