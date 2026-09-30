package simulation
import "testing"
func TestGroupRuntime(t *testing.T){r:=NewGroupRuntime();if !r.SetRole("g","a",GroupLeader){t.Fatal("role")};if v,_:=r.Role("g","a");v!=GroupLeader{t.Fatal("wrong role")};r.AddResource("g",100,5);if v,_:=r.Resource("g");v.Currency!=100||v.Influence!=5{t.Fatal("resource")};r.SetGoal("g","grow");if v,_:=r.Goal("g");v!="grow"{t.Fatal("goal")}}
