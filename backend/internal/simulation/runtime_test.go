package simulation

import (
 "testing"
 "time"
)

func TestRuntimeTickWithoutRunner(t *testing.T){
 r:=NewRuntime(nil,TierScheduler{})
 if got:=r.Tick(time.Unix(1,0));got!=0{t.Fatal("expected zero")}
}
