package simulation

import "testing"

func TestChooseActionUsesValidatedContext(t *testing.T) {
 c:=SimCharacter{Goals:[]Goal{{Kind:"WORK",Priority:5},{Kind:"SOCIAL",Priority:3}}}
 a,err:=ChooseAction(c,Context{HasJob:true,SocialOpportunity:true})
 if err!=nil||a.Type!=ActionWork {t.Fatalf("unexpected action: %+v err=%v",a,err)}
 a,err=ChooseAction(c,Context{HasJob:false,SocialOpportunity:true})
 if err!=nil||a.Type!=ActionSocialize {t.Fatalf("invalid work action selected: %+v",a)}
}
