package simulation

import ("errors";"sort")

var ErrNoAction=errors.New("no valid action")

type ActionType string
const(ActionWork ActionType="WORK";ActionTravel ActionType="TRAVEL";ActionStudy ActionType="STUDY";ActionSocialize ActionType="SOCIALIZE";ActionRest ActionType="REST")

type Action struct{Type ActionType `json:"type"`;TargetID string `json:"targetId,omitempty"`;Priority int `json:"priority"`}

type Context struct{HasJob bool `json:"hasJob"`;CanStudy bool `json:"canStudy"`;SocialOpportunity bool `json:"socialOpportunity"`;RestNeeded bool `json:"restNeeded"`}

func ChooseAction(c SimCharacter,ctx Context)(Action,error){
 actions:=make([]Action,0,6)
 for _,g:=range c.Goals{
  p:=g.Priority+personalityBias(c,g.Kind)
  switch g.Kind{
  case "WORK":if ctx.HasJob{actions=append(actions,Action{Type:ActionWork,TargetID:g.TargetID,Priority:p+stateBias(c,ActionWork)})}
  case "STUDY":if ctx.CanStudy{actions=append(actions,Action{Type:ActionStudy,TargetID:g.TargetID,Priority:p+stateBias(c,ActionStudy)})}
  case "SOCIAL":if ctx.SocialOpportunity{actions=append(actions,Action{Type:ActionSocialize,TargetID:g.TargetID,Priority:p+stateBias(c,ActionSocialize)})}
  case "TRAVEL":if g.TargetID!=""{actions=append(actions,Action{Type:ActionTravel,TargetID:g.TargetID,Priority:p+stateBias(c,ActionTravel)})}
  }
 }
 if ctx.RestNeeded||stateNeedsRest(c){actions=append(actions,Action{Type:ActionRest,Priority:110})}
 if len(actions)==0{return Action{},ErrNoAction}
 sort.SliceStable(actions,func(i,j int)bool{return actions[i].Priority>actions[j].Priority})
 return actions[0],nil
}

func personalityBias(c SimCharacter,kind string)int{
 switch kind{case "SOCIAL":return c.Personality.Sociability/5;case "TRAVEL":return c.Personality.Curiosity/5;case "STUDY","WORK":return c.Personality.Discipline/5}
 return 0
}

// SimCharacter carries a compact behavioral snapshot so decision making remains deterministic.
func stateBias(c SimCharacter,a ActionType)int{
 switch a{
 case ActionRest:return 20
 case ActionSocialize:return c.Personality.Sociability/5
 case ActionTravel:return c.Personality.Curiosity/5
 case ActionStudy,ActionWork:return c.Personality.Discipline/5
 }
 return 0
}

// These fields are supplied by the runner when a behavioral state exists.
// The fallback keeps legacy callers deterministic.
func stateNeedsRest(c SimCharacter)bool{return false}
