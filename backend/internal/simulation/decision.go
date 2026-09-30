package simulation

import("context";"encoding/json";"errors";"sort";"time")

type Decision struct{ActionType string `json:"actionType"`;TargetID string `json:"targetId,omitempty"`;Score int `json:"score"`}
type DecisionContext struct{Now time.Time;Energy int;Cash int64;Level int;Nearby []string}
var ErrNoDecision=errors.New("no valid simulation decision")
func ChooseDecision(c SimCharacter,ctx DecisionContext)(Decision,error){
 if !c.Active{return Decision{},ErrNoDecision}
 var choices []Decision
 for _,g:=range c.Goals{score:=g.Priority;switch g.Kind{case "WORK":score+=c.Personality.Discipline;case "SOCIAL":score+=c.Personality.Sociability;case "EXPLORE":score+=c.Personality.Curiosity;case "RISK":score+=c.Personality.RiskTolerance;case "PROGRESS":score+=c.Personality.Discipline+c.Personality.Curiosity/2};if ctx.Energy<=0&&g.Kind!="SOCIAL"{score-=100};if ctx.Cash<=0&&g.Kind=="RISK"{score-=25};choices=append(choices,Decision{ActionType:g.Kind,TargetID:g.TargetID,Score:score})}
 if len(choices)==0{return Decision{},ErrNoDecision};sort.SliceStable(choices,func(i,j int)bool{return choices[i].Score>choices[j].Score});return choices[0],nil
}
func Normalize(c SimCharacter)SimCharacter{if c.Personality.Curiosity<0{c.Personality.Curiosity=0};if c.Personality.Curiosity>100{c.Personality.Curiosity=100};if c.Personality.Sociability<0{c.Personality.Sociability=0};if c.Personality.Sociability>100{c.Personality.Sociability=100};if c.Personality.RiskTolerance<0{c.Personality.RiskTolerance=0};if c.Personality.RiskTolerance>100{c.Personality.RiskTolerance=100};if c.Personality.Discipline<0{c.Personality.Discipline=0};if c.Personality.Discipline>100{c.Personality.Discipline=100};return c}
func Encode(c SimCharacter)([]byte,error){return json.Marshal(c)}
func TickCharacter(c SimCharacter,ctx DecisionContext)(SimCharacter,Decision,error){c=Normalize(c);d,err:=ChooseDecision(c,ctx);if err!=nil{return c,Decision{},err};c.LastTick=ctx.Now.Unix();return c,d,nil}
func LoadGoals(raw []byte)([]Goal,error){var g []Goal;if len(raw)==0{return g,nil};return g,json.Unmarshal(raw,&g)}
var _=context.Background
