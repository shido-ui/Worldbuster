package simulation

import "time"

type ActionExecutor func(action Action, character SimCharacter) ActionResult
type EventSink func(eventType, actorID, targetID string, payload map[string]any)

type Runner struct {
 Population *Service
 State *StateService
 BuildContext ContextBuilder
 Execute ActionExecutor
 Emit EventSink
}

func (r *Runner) Tick(now time.Time) int {
 if r.Population==nil||r.State==nil||r.BuildContext==nil||r.Execute==nil{return 0}
 count:=0
 for _,c:=range r.Population.List(1000){
  if !c.Active{continue}
  ctx:=r.BuildContext(c.ID)
  state:=r.State.Get(c.ID)
  result:=DecideAndValidateWithState(c,ctx,state)
  if !result.Accepted{r.State.Tick(c.ID);continue}
  result=r.Execute(result.Action,c)
  if result.Accepted{
   r.State.Apply(c.ID,result.Action.Type,now)
   r.State.Tick(c.ID)
   if r.Emit!=nil{r.Emit("SIMULATED_ACTION",c.ID,result.Action.TargetID,map[string]any{"action":string(result.Action.Type),"mood":r.State.Get(c.ID).Mood})}
   count++
  }
 }
 return count
}
