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
  switch g.Kind {
  case "WORK":
   if ctx.HasJob { actions=append(actions,Action{Type:ActionWork,TargetID:g.TargetID,Priority:g.Priority}) }
  case "STUDY":
   if ctx.CanStudy { actions=append(actions,Action{Type:ActionStudy,TargetID:g.TargetID,Priority:g.Priority}) }
  case "SOCIAL":
   if ctx.SocialOpportunity { actions=append(actions,Action{Type:ActionSocialize,TargetID:g.TargetID,Priority:g.Priority}) }
  }
 }
 if ctx.RestNeeded { actions=append(actions,Action{Type:ActionRest,Priority:100}) }
 if len(actions)==0 { return Action{},ErrNoAction }
 sort.SliceStable(actions,func(i,j int)bool{return actions[i].Priority>actions[j].Priority})
 return actions[0],nil
}
