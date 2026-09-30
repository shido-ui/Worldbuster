package simulation

import (
 "errors"
 "sort"
)

var ErrNoAction = errors.New("no valid action")

type ActionType string
const (
 ActionWork ActionType = "WORK"
 ActionTravel ActionType = "TRAVEL"
 ActionStudy ActionType = "STUDY"
 ActionSocialize ActionType = "SOCIALIZE"
 ActionRest ActionType = "REST"
)

type Action struct {
 Type ActionType `json:"type"`
 TargetID string `json:"targetId,omitempty"`
 Priority int `json:"priority"`
}

type Context struct {
 HasJob bool `json:"hasJob"`
 CanStudy bool `json:"canStudy"`
 SocialOpportunity bool `json:"socialOpportunity"`
 RestNeeded bool `json:"restNeeded"`
}

func ChooseAction(c SimCharacter, ctx Context) (Action, error) {
 actions := make([]Action,0,5)
 for _, g := range c.Goals {
  priority := g.Priority
  switch g.Kind {
  case "WORK":
   if ctx.HasJob { actions=append(actions,Action{Type:ActionWork,TargetID:g.TargetID,Priority:priority + needBias(c,ActionWork)}) }
  case "STUDY":
   if ctx.CanStudy { actions=append(actions,Action{Type:ActionStudy,TargetID:g.TargetID,Priority:priority + needBias(c,ActionStudy)}) }
  case "SOCIAL":
   if ctx.SocialOpportunity { actions=append(actions,Action{Type:ActionSocialize,TargetID:g.TargetID,Priority:priority + needBias(c,ActionSocialize)}) }
  case "TRAVEL":
   if g.TargetID!="" { actions=append(actions,Action{Type:ActionTravel,TargetID:g.TargetID,Priority:priority + needBias(c,ActionTravel)}) }
  }
 }
 if ctx.RestNeeded || needsRest(c) {
  actions=append(actions,Action{Type:ActionRest,Priority:100 + needBias(c,ActionRest)})
 }
 if len(actions)==0 { return Action{},ErrNoAction }
 sort.SliceStable(actions,func(i,j int)bool{return actions[i].Priority>actions[j].Priority})
 return actions[0],nil
}

func needBias(c SimCharacter, action ActionType) int {
 // Behavioral state is intentionally bounded; extreme needs override routine goals.
 switch action {
 case ActionRest:
  return 30
 case ActionSocialize:
  return c.Personality.Sociability/10
 case ActionStudy:
  return c.Personality.Discipline/10
 case ActionWork:
  return c.Personality.Discipline/10
 case ActionTravel:
  return c.Personality.Curiosity/10
 default:
  return 0
 }
}

func needsRest(c SimCharacter) bool {
 return false
}
