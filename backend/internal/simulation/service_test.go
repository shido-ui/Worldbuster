package simulation
import("testing";"time")
func TestSimulationPopulation(t *testing.T){s:=NewService();if err:=s.Register(SimCharacter{ID:"npc1",Name:"Aster Citizen",Controller:ControllerSimulated,Goals:[]Goal{{ID:"g1",Kind:"WORK",Priority:1}}});err!=nil{t.Fatal(err)};if n:=s.Tick(time.Unix(123,0));n!=1{t.Fatalf("tick count=%d",n)};c,ok:=s.Get("npc1");if !ok||c.LastTick!=123{t.Fatalf("character was not ticked: %+v",c)}}
