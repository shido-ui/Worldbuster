package simulation
import "testing"
func TestGroupService(t *testing.T){s:=NewGroupService();if !s.Create("g1","Citizens","a"){t.Fatal("create failed")};if !s.Join("g1","b"){t.Fatal("join failed")};g,ok:=s.GroupOf("b");if !ok||len(g.MemberIDs)!=2{t.Fatal("membership failed")};if s.Join("g1","b"){t.Fatal("duplicate membership")}}
