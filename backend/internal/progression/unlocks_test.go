package progression

import "testing"

func TestUnlockStore(t *testing.T){
 s:=NewUnlockStore()
 r:=ProgressionReward{ID:"r1",SkillID:"focus",RequiredLevel:5,Type:RewardAbility,Value:"advanced_focus"}
 if !s.Grant("p1",r){t.Fatal("grant failed")}
 if s.Grant("p1",r){t.Fatal("duplicate grant accepted")}
 if !s.Has("p1","r1"){t.Fatal("unlock missing")}
 if len(s.List("p1"))!=1{t.Fatal("unlock list incorrect")}
}
