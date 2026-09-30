package simulation

import ("testing";"time")

func TestTierScheduler(t *testing.T){
 s:=TierScheduler{Policy:LifecyclePolicy{ActivePercent:10,RecentPercent:20,BackgroundPercent:30}}
 now:=time.Unix(100,0)
 if !s.ShouldTick(1,100,now){t.Fatal("active should tick")}
 if s.ShouldTick(90,100,now){t.Fatal("dormant should not tick")}
}
