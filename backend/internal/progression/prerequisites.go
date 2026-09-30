package progression

type SkillPrerequisite struct {
 SkillID string `json:"skillId"`
 RequiredLevel int `json:"requiredLevel"`
}

type SkillPath struct {
 SkillID string `json:"skillId"`
 Prerequisites []SkillPrerequisite `json:"prerequisites"`
}

func CanUnlock(path SkillPath, owned map[string]SkillState) bool {
 for _,req:=range path.Prerequisites {
  state,ok:=owned[req.SkillID]
  if !ok || state.Level<req.RequiredLevel { return false }
 }
 return true
}

func ValidatePaths(paths []SkillPath) bool {
 seen:=map[string]bool{}
 for _,p:=range paths {
  if p.SkillID=="" || seen[p.SkillID] { return false }
  seen[p.SkillID]=true
  for _,req:=range p.Prerequisites {
   if req.SkillID=="" || req.RequiredLevel<1 { return false }
  }
 }
 return true
}
