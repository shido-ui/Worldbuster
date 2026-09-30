package simulation

import "sort"

type GroupActionType string

const (
 GroupRecruit GroupActionType = "RECRUIT"
 GroupCooperate GroupActionType = "COOPERATE"
 GroupCompete GroupActionType = "COMPETE"
 GroupInvest GroupActionType = "INVEST"
)

type GroupContext struct {
 MemberCount int
 Influence int
 Currency int64
 Goal string
 NearbyCandidates int
 RivalInfluence int
}

type GroupAction struct {
 Type GroupActionType
 TargetID string
 Priority int
}

func ChooseGroupAction(ctx GroupContext) GroupAction {
 actions:=make([]GroupAction,0,4)
 if ctx.NearbyCandidates>0 && ctx.MemberCount<20 {
  actions=append(actions,GroupAction{Type:GroupRecruit,Priority:50+ctx.NearbyCandidates})
 }
 if ctx.Goal=="grow" && ctx.Currency>=50 {
  actions=append(actions,GroupAction{Type:GroupInvest,Priority:60+ctx.MemberCount})
 }
 if ctx.RivalInfluence>ctx.Influence && ctx.MemberCount>2 {
  actions=append(actions,GroupAction{Type:GroupCompete,Priority:55+ctx.RivalInfluence-ctx.Influence})
 }
 if ctx.MemberCount>1 {
  actions=append(actions,GroupAction{Type:GroupCooperate,Priority:40+ctx.MemberCount})
 }
 if len(actions)==0{return GroupAction{}}
 sort.SliceStable(actions,func(i,j int)bool{return actions[i].Priority>actions[j].Priority})
 return actions[0]
}

func ApplyGroupAction(runtime *GroupRuntime,groupID string,action any){
 if runtime==nil||groupID==""{return}
 var typ GroupActionType
 switch a:=action.(type){case GroupAction: typ=a.Type; case GroupActionType: typ=a; default: return}
 switch typ{
 case GroupInvest:
  runtime.AddResource(groupID,-50,10)
 case GroupRecruit:
  runtime.AddResource(groupID,0,2)
 case GroupCooperate:
  runtime.AddResource(groupID,0,1)
 case GroupCompete:
  runtime.AddResource(groupID,0,3)
 }
}
