package simulation

import ("testing";"time")

func TestBehavioralState(t *testing.T){s:=NewStateService();if s.Get("npc1").Needs.Energy!=100{t.Fatal("initial needs incorrect")};s.Tick("npc1");m:=s.Remember("npc1","SOCIAL","npc2","met at harbor",5,2);if m.Event==""||len(s.Get("npc1").Memories)!=1{t.Fatal("memory missing")}}

func TestRestRecoversAndClamps(t *testing.T){
 s:=NewStateService()
 s.states["npc1"]=BehavioralState{Needs:Needs{Energy:20,Rest:10,Social:50,Satisfaction:50}}
 s.Apply("npc1",ActionRest,time.Unix(1,0))
 got:=s.Get("npc1")
 if got.Needs.Energy<=20||got.Needs.Rest<=10{t.Fatalf("rest did not recover needs: %+v",got.Needs)}
 s.states["npc1"]=BehavioralState{Needs:Needs{Energy:95,Rest:95,Social:50,Satisfaction:50}}
 s.Apply("npc1",ActionRest,time.Unix(2,0))
 got=s.Get("npc1")
 if got.Needs.Energy!=100||got.Needs.Rest!=100{t.Fatalf("rest exceeded max or failed clamp: %+v",got.Needs)}
}
