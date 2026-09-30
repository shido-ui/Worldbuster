package simulation

import ("testing";"time")

func TestRunnerExecutesOnlyValidatedActions(t *testing.T){
 p:=NewService();_ = p.Register(SimCharacter{ID:"npc1",Name:"Worker",Controller:ControllerSimulated,Goals:[]Goal{{Kind:"WORK",Priority:10}}})
 st:=NewStateService();ran:=0;events:=0
 r:=Runner{Population:p,State:st,BuildContext:func(string)WorldContext{return WorldContext{HasJob:true}},Execute:func(a Action,c SimCharacter)ActionResult{ran++;return ActionResult{Accepted:true,Action:a}},Emit:func(string,string,string,map[string]any){events++}}
 if n:=r.Tick(time.Unix(1,0));n!=1||ran!=1||events!=1{t.Fatalf("runner failed n=%d ran=%d events=%d",n,ran,events)}
}
