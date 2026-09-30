package progression

type SkillDefinition struct {
 ID string `json:"id"`
 Name string `json:"name"`
 Description string `json:"description"`
 Stat string `json:"stat"`
 MaxLevel int `json:"maxLevel"`
}

type SkillState struct {
 SkillID string `json:"skillId"`
 Level int `json:"level"`
 XP int64 `json:"xp"`
}

type SkillRegistry struct {
 definitions map[string]SkillDefinition
}

func NewSkillRegistry(defs []SkillDefinition)*SkillRegistry{
 r:=&SkillRegistry{definitions:map[string]SkillDefinition{}}
 for _,d:=range defs{if d.ID!=""&&d.MaxLevel>0{r.definitions[d.ID]=d}}
 return r
}

func(r *SkillRegistry) Get(id string)(SkillDefinition,bool){
 if r==nil{return SkillDefinition{},false};d,ok:=r.definitions[id];return d,ok
}

func(s *SkillState) AddXP(def SkillDefinition,amount int64)bool{
 if s==nil||amount<=0||s.SkillID!=def.ID||s.Level<1||s.Level>=def.MaxLevel{return false}
 s.XP+=amount
 for s.Level<def.MaxLevel {
  need:=int64(s.Level*s.Level*50)
  if s.XP<need{break}
  s.XP-=need;s.Level++
 }
 return true
}
