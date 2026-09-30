package progression

import "testing"

func TestSkillFramework(t *testing.T){
 r:=NewSkillRegistry([]SkillDefinition{{ID:"focus",Name:"Focus",MaxLevel:10}})
 d,ok:=r.Get("focus")
 if !ok{t.Fatal("skill missing")}
 s:=SkillState{SkillID:"focus",Level:1}
 if !s.AddXP(d,50)||s.Level!=2{t.Fatal("skill progression failed")}
}
