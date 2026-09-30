package simulation

import ("errors";"sort")

var ErrNoAction=errors.New("no valid action")

type ActionType string
const(ActionWork ActionType="WORK";ActionTravel ActionType="TRAVEL";ActionStudy ActionType="STUDY";ActionSocialize ActionType="SOCIALIZE";ActionRest ActionType="REST")

type Action struct{Type ActionType;TargetID string;Priority int}
type Context struct{HasJob bool;CanStudy bool;SocialOpportunity bool;SocialTargetID string;RestNeeded bool;HasTravelDestination bool}

func ChooseAction(c SimCharacter,ctx Context)(Action,error){return ChooseActionWithState(c,ctx,defaultState())}

func ChooseActionWithState(c SimCharacter,ctx Context,state BehavioralState)(Action,error){
 state=normalizeDecisionState(state)
 criticalRest:=ctx.RestNeeded||state.Needs.Energy<25||state.Needs.Rest<20
 if criticalRest{return Action{Type:ActionRest,Priority:1000000},nil}
 actions:=[]Action{}
 for _,g:=range c.Goals{
  p:=g.Priority+personalityBias(c,g.Kind)
  switch g.Kind{
  case "WORK":
   if ctx.HasJob{actions=append(actions,Action{Type:ActionWork,TargetID:g.TargetID,Priority:p+stateBias(state,ActionWork)})}
  case "STUDY":
   if ctx.CanStudy{actions=append(actions,Action{Type:ActionStudy,TargetID:g.TargetID,Priority:p+stateBias(state,ActionStudy)})}
  case "SOCIAL":
   if ctx.SocialOpportunity&&ctx.SocialTargetID!=""{
    rp:=0
    if rel,ok:=state.Relationships[ctx.SocialTargetID];ok{rp=rel.SocialWeight()}
    actions=append(actions,Action{Type:ActionSocialize,TargetID:ctx.SocialTargetID,Priority:p+stateBias(state,ActionSocialize)+rp})
   }
  case "TRAVEL":
   if ctx.HasTravelDestination&&g.TargetID!=""{actions=append(actions,Action{Type:ActionTravel,TargetID:g.TargetID,Priority:p+stateBias(state,ActionTravel)})}
  }
 }
 if len(actions)==0{return Action{},ErrNoAction}
 sort.SliceStable(actions,func(i,j int)bool{return actions[i].Priority>actions[j].Priority})
 return actions[0],nil
}

func personalityBias(c SimCharacter,kind string)int{switch kind{case "SOCIAL":return c.Personality.Sociability/5;case "TRAVEL":return c.Personality.Curiosity/5;case "STUDY","WORK":return c.Personality.Discipline/5};return 0}
func stateBias(s BehavioralState,a ActionType)int{switch a{case ActionRest:return maxInt(0,100-s.Needs.Energy)/2+maxInt(0,100-s.Needs.Rest)/3;case ActionSocialize:return maxInt(0,50-s.Needs.Social)/2+s.Mood/10;case ActionWork:return s.Needs.Satisfaction/10-s.Stress/10;case ActionStudy:return s.Needs.Satisfaction/12-s.Stress/12;case ActionTravel:return s.Mood/10};return 0}
func maxInt(a,b int)int{if a>b{return a};return b}

func normalizeDecisionState(s BehavioralState) BehavioralState {
 if s.Mood==0&&s.Stress==0&&s.Needs.Energy==0&&s.Needs.Social==0&&s.Needs.Rest==0&&s.Needs.Satisfaction==0&&s.Economy.Balance==0&&len(s.Memories)==0&&len(s.Relationships)==0{return defaultState()}
 return s
}
