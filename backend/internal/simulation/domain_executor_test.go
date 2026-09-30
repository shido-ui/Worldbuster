package simulation

import ("errors";"testing";"time")

func TestServiceExecutorUsesDomainServices(t *testing.T){
 called:=false
 e:=ServiceExecutor{Now:func()time.Time{return time.Unix(10,0)},Services:DomainServices{Work:func(id string)error{called=id=="npc-1";return nil}}}
 r:=e.Execute(Action{Type:ActionWork},SimCharacter{ID:"npc-1"})
 if !r.Accepted||!called{t.Fatal("work was not delegated")}
}
func TestServiceExecutorRejectsMissingService(t *testing.T){
 e:=ServiceExecutor{}
 r:=e.Execute(Action{Type:ActionWork},SimCharacter{ID:"npc-1"})
 if r.Accepted||r.Reason==""{t.Fatal("missing service should reject")}
}
func TestServiceExecutorPropagatesFailure(t *testing.T){
 e:=ServiceExecutor{Services:DomainServices{Study:func(string)error{return errors.New("not enrolled")}}}
 r:=e.Execute(Action{Type:ActionStudy},SimCharacter{ID:"npc-1"})
 if r.Accepted{t.Fatal("failed domain action should reject")}
}
