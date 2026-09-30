package simulation

import (
	"fmt"
	"time"
)

type DomainServices struct {
	Work func(characterID string) error
	Travel func(characterID, destination string, now time.Time) error
	Study func(characterID string) error
	Socialize func(characterID, targetID string) error
	Rest func(characterID string) error
}

type ServiceExecutor struct {
	Services DomainServices
	Now func() time.Time
}

func (e ServiceExecutor) Execute(action Action, character SimCharacter) ActionResult {
	now:=time.Now
	if e.Now!=nil { now=e.Now }
	switch action.Type {
	case ActionWork:
		if e.Services.Work==nil{return reject(action,"work service unavailable")}
		if err:=e.Services.Work(character.ID);err!=nil{return reject(action,err.Error())}
	case ActionTravel:
		if e.Services.Travel==nil{return reject(action,"travel service unavailable")}
		if err:=e.Services.Travel(character.ID,action.TargetID,now());err!=nil{return reject(action,err.Error())}
	case ActionStudy:
		if e.Services.Study==nil{return reject(action,"study service unavailable")}
		if err:=e.Services.Study(character.ID);err!=nil{return reject(action,err.Error())}
	case ActionSocialize:
		if e.Services.Socialize==nil{return reject(action,"social service unavailable")}
		if err:=e.Services.Socialize(character.ID,action.TargetID);err!=nil{return reject(action,err.Error())}
	case ActionRest:
		if e.Services.Rest==nil{return reject(action,"rest service unavailable")}
		if err:=e.Services.Rest(character.ID);err!=nil{return reject(action,err.Error())}
	default:
		return reject(action,"unsupported action")
	}
	return ActionResult{Accepted:true,Action:action}
}

func reject(action Action, reason string) ActionResult {
	return ActionResult{Accepted:false,Action:action,Reason:fmt.Sprintf("domain execution rejected: %s",reason)}
}
