package simulation
import "testing"
func TestBehavioralState(t *testing.T){s:=NewStateService();if s.Get("npc1").Needs.Energy!=100{t.Fatal("initial needs incorrect")};s.Tick("npc1");m:=s.Remember("npc1","SOCIAL","npc2","met at harbor",5,2);if m.Event==""||len(s.Get("npc1").Memories)!=1{t.Fatal("memory missing")}}
