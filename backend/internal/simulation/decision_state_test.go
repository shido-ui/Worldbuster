package simulation
import "testing"
func TestChooseActionWithStateRestPressure(t *testing.T){
 c:=SimCharacter{Goals:[]Goal{{Kind:"WORK",Priority:100}},Personality:Personality{Discipline:100}}
 a,err:=ChooseActionWithState(c,Context{HasJob:true},BehavioralState{Needs:Needs{Energy:10,Rest:10}})
 if err!=nil||a.Type!=ActionRest{t.Fatal("low needs should trigger rest")}
}
func TestChooseActionWithStateSocialPressure(t *testing.T){
 c:=SimCharacter{Goals:[]Goal{{Kind:"SOCIAL",Priority:50}},Personality:Personality{Sociability:100}}
 a,err:=ChooseActionWithState(c,Context{SocialOpportunity:true,SocialTargetID:"npc2"},BehavioralState{Needs:Needs{Energy:100,Rest:100,Social:0,Satisfaction:50}})
 if err!=nil||a.Type!=ActionSocialize{t.Fatal("social pressure should influence decision")}
}
