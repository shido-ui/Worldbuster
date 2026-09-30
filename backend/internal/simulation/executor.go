package simulation

import "time"

// DomainExecutor is the single mutation boundary for simulated actions.
// Implementations must call the same authoritative services used by human players.
type DomainExecutor struct {
	Now func() time.Time
}

func (e DomainExecutor) Execute(action Action, character SimCharacter) ActionResult {
	if e.Now == nil { e.Now=time.Now }
	switch action.Type {
	case ActionWork, ActionTravel, ActionStudy, ActionSocialize, ActionRest:
		return ActionResult{Accepted:true,Action:action}
	default:
		return ActionResult{Accepted:false,Action:action,Reason:"unsupported action"}
	}
}
