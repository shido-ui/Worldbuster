package simulation

type WorldContext struct {
 CharacterID string `json:"characterId"`
 HasJob bool `json:"hasJob"`
 CanStudy bool `json:"canStudy"`
 SocialOpportunity bool `json:"socialOpportunity"`
 RestNeeded bool `json:"restNeeded"`
 HasTravelDestination bool `json:"hasTravelDestination"`
}

type ActionResult struct {
 Accepted bool `json:"accepted"`
 Action Action `json:"action"`
 Reason string `json:"reason,omitempty"`
}

type ActionValidator interface {
 Validate(action Action, ctx WorldContext) bool
}

type ContextBuilder func(characterID string) WorldContext

func DecideAndValidate(c SimCharacter, ctx WorldContext) ActionResult { return DecideAndValidateWithState(c,ctx,BehavioralState{}) }

func DecideAndValidateWithState(c SimCharacter, ctx WorldContext, state BehavioralState) ActionResult {
 action, err := ChooseActionWithState(c, Context{HasJob:ctx.HasJob,CanStudy:ctx.CanStudy,SocialOpportunity:ctx.SocialOpportunity,RestNeeded:ctx.RestNeeded}, state)
 if err != nil { return ActionResult{Accepted:false,Reason:err.Error()} }
 if !validateAction(action,ctx) { return ActionResult{Accepted:false,Action:action,Reason:"action is not valid in current world context"} }
 return ActionResult{Accepted:true,Action:action}
}

func validateAction(a Action,ctx WorldContext) bool {
 switch a.Type {
 case ActionWork: return ctx.HasJob
 case ActionStudy: return ctx.CanStudy
 case ActionSocialize: return ctx.SocialOpportunity
 case ActionTravel: return ctx.HasTravelDestination
 case ActionRest: return ctx.RestNeeded
 default: return false
 }
}
