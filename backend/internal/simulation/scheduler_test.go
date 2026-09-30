package simulation

import("testing";"time")

func TestTierScheduler(t *testing.T){
 s:=TierScheduler{Policy:LifecyclePolicy{ActivePercent:10,RecentPercent:20,BackgroundPercent:30}}
 now:=time.Unix(100,0)
 if !s.ShouldTick(1,100,now){t.Fatal("active should tick")}
 if s.ShouldTick(90,100,now){t.Fatal("dormant should not tick")}
}

func TestTierSchedulerCadence(t *testing.T){
 now:=time.Unix(1000,0)
 chars:=make([]SimCharacter,10)
 for i:=range chars{chars[i]=SimCharacter{ID:formatID(i+1),LastTick:now.Unix()}}
 s:=TierScheduler{Policy:LifecyclePolicy{ActivePercent:20,RecentPercent:30,BackgroundPercent:30}}
 if got:=len(s.Due(chars,now.Add(time.Second)));got!=2{t.Fatalf("1s due=%d want 2",got)}
 if got:=len(s.Due(chars,now.Add(5*time.Second)));got!=5{t.Fatalf("5s due=%d want 5",got)}
 if got:=len(s.Due(chars,now.Add(30*time.Second)));got!=8{t.Fatalf("30s due=%d want 8",got)}
 if got:=len(s.Due(chars,now.Add(5*time.Minute)));got!=10{t.Fatalf("5m due=%d want 10",got)}
}
