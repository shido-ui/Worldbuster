package simulation
type WorldContext struct{CharacterID string;HasJob bool;CanStudy bool;SocialOpportunity bool;SocialTargetID string;RestNeeded bool;HasTravelDestination bool}
type ActionResult struct{Accepted bool;Action Action;Reason string}
type ActionValidator func(action Action,ctx WorldContext)bool
type ContextBuilder func(characterID string)WorldContext
func DecideAndValidate(c SimCharacter,ctx WorldContext)ActionResult{return DecideAndValidateWithState(c,ctx,BehavioralState{})}

func DecideAndValidateWithState(c SimCharacter,ctx WorldContext,state BehavioralState)ActionResult{action,err:=ChooseActionWithState(c,Context{HasJob:ctx.HasJob,CanStudy:ctx.CanStudy,SocialOpportunity:ctx.SocialOpportunity,SocialTargetID:ctx.SocialTargetID,RestNeeded:ctx.RestNeeded},state);if err!=nil{return ActionResult{Accepted:false,Reason:err.Error()}};if action.Type==ActionWork&&!ctx.HasJob{return ActionResult{Action:action,Reason:"work unavailable"}};if action.Type==ActionStudy&&!ctx.CanStudy{return ActionResult{Action:action,Reason:"study unavailable"}};if action.Type==ActionSocialize&&(!ctx.SocialOpportunity||action.TargetID==""){return ActionResult{Action:action,Reason:"social unavailable"}};if action.Type==ActionTravel&&!ctx.HasTravelDestination{return ActionResult{Action:action,Reason:"travel unavailable"}};return ActionResult{Accepted:true,Action:action}}
