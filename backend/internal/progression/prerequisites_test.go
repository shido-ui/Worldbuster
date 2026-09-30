package progression

import "testing"

func TestSkillPrerequisites(t *testing.T){
 path:=SkillPath{SkillID:"advanced",Prerequisites:[]SkillPrerequisite{{SkillID:"focus",RequiredLevel:2}}}
 owned:=map[string]SkillState{"focus":{SkillID:"focus",Level:2}}
 if !CanUnlock(path,owned){t.Fatal("valid path rejected")}
 owned["focus"]=SkillState{SkillID:"focus",Level:1}
 if CanUnlock(path,owned){t.Fatal("insufficient prerequisite accepted")}
}
