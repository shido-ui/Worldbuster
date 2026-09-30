package simulation
import("time";"github.com/shido-ui/Worldbuster/backend/internal/job";"github.com/shido-ui/Worldbuster/backend/internal/progression";"github.com/shido-ui/Worldbuster/backend/internal/social";"github.com/shido-ui/Worldbuster/backend/internal/world")
func reject(action Action, reason string) ActionResult{return ActionResult{Accepted:false,Action:action,Reason:reason}}

type LiveAdapter struct{Jobs *job.Service;Progression *progression.Service;Social *social.Service;Travel *world.TravelService;Locations map[string]string}
func(a *LiveAdapter)Context(characterID string)WorldContext{
 _,employed:=a.Jobs.GetEmployment(characterID)
 canStudy:=len(a.Progression.ListCourses())>0
 socialTarget:=""
 for id:=range a.Locations{if id!=characterID{socialTarget=id;break}}
 now:=time.Now().UTC()
 if a.Travel!=nil{
  if t,ok:=a.Travel.Get(characterID,now);ok&&!now.Before(t.ArrivesAt){
   if completed,done:=a.Travel.Complete(characterID,now);done{a.Locations[characterID]=completed.To}
  }
 }
 travelDestination:=""
 current:=a.Locations[characterID]
 traveling:=false
 if a.Travel!=nil{_,traveling=a.Travel.Get(characterID,now)}
 if !traveling&&current!=""{
  for _,route:=range world.Routes(){if route.From==current&&route.To!=""{travelDestination=route.To;break}}
 }
 return WorldContext{CharacterID:characterID,HasJob:employed,CanStudy:canStudy,SocialOpportunity:socialTarget!="",SocialTargetID:socialTarget,RestNeeded:false,HasTravelDestination:travelDestination!=""}
}
func(a *LiveAdapter)Execute(action Action,character SimCharacter)ActionResult{switch action.Type{case ActionWork:if _,ok:=a.Jobs.GetEmployment(character.ID);!ok{return reject(action,"character is not employed")};a.Progression.AddXP(character.ID,25);case ActionStudy:if len(a.Progression.ListCourses())==0{return reject(action,"no courses available")};a.Progression.AddXP(character.ID,15);case ActionSocialize:if action.TargetID==""{return reject(action,"social target required")};if err:=a.Social.SetRelationship(character.ID,action.TargetID,"acquaintance",1);err!=nil{return reject(action,err.Error())};case ActionTravel:from:=a.Locations[character.ID];if from==""||action.TargetID==""{return reject(action,"travel origin or destination missing")};if _,err:=a.Travel.Start(character.ID,from,action.TargetID,time.Now());err!=nil{return reject(action,err.Error())};case ActionRest:default:return reject(action,"unsupported action")};return ActionResult{Accepted:true,Action:action}}
