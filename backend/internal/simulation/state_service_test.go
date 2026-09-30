package simulation

import("testing";"time")

func TestStateServicePersistsAndEvolves(t *testing.T){
 s:=NewStateService()
 before:=s.Get("npc")
 s.Apply("npc",ActionWork,time.Unix(1,0))
 after:=s.Get("npc")
 if after.Needs.Energy>=before.Needs.Energy||len(after.Memories)!=1{t.Fatal("state was not persisted")}
 s.Tick("npc")
 if s.Get("npc").Needs.Energy>=after.Needs.Energy{t.Fatal("needs did not drift")}
}
