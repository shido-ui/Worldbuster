package simulation

import "testing"

func TestDecisionValidatedAgainstWorldContext(t *testing.T){
 c:=SimCharacter{ID:"npc1",Goals:[]Goal{{Kind:"WORK",Priority:5}}}
 r:=DecideAndValidate(c,WorldContext{CharacterID:"npc1",HasJob:false})
 if r.Accepted{t.Fatal("invalid work action was accepted")}
 r=DecideAndValidate(c,WorldContext{CharacterID:"npc1",HasJob:true})
 if !r.Accepted||r.Action.Type!=ActionWork{t.Fatalf("valid action rejected: %+v",r)}
}
