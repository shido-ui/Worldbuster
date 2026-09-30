package simulation
import "testing"
func TestClassifyPopulation(t *testing.T){p:=LifecyclePolicy{ActivePercent:10,RecentPercent:20,BackgroundPercent:30};if ClassifyPopulation(5,100,p)!=TierActive{t.Fatal("active")};if ClassifyPopulation(20,100,p)!=TierRecent{t.Fatal("recent")};if ClassifyPopulation(40,100,p)!=TierBackground{t.Fatal("background")};if ClassifyPopulation(90,100,p)!=TierDormant{t.Fatal("dormant")}}
func TestSelectRegion(t *testing.T){r:=[]PopulationRegion{{Name:"A",Weight:1},{Name:"B",Weight:3}};if SelectRegion(0,r)!="A"{t.Fatal("region")};if SelectRegion(3,r)!="B"{t.Fatal("region")}}
