package simulation

import "testing"

func TestResolveCombatDeterministic(t *testing.T){
 a:=CombatInput{10,8,7,6,9,7,5,5,3,3}
 x:=ResolveCombat(a);y:=ResolveCombat(a)
 if x!=y{t.Fatalf("combat must be deterministic: %+v != %+v",x,y)}
 if x.EnergyCost<=0||x.XPReward<=0{t.Fatalf("expected bounded costs/reward: %+v",x)}
}
