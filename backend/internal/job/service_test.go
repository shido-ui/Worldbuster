package job
import "testing"
func TestJobPositions(t *testing.T){s:=NewService();if err:=s.Register(Job{ID:"medical",Name:"Medical",Department:"Healthcare",Positions:[]Position{{ID:"intern",JobID:"medical",Name:"Intern",Level:1,BaseSalary:100,RequiredEducation:1,RequiredStat:5},{ID:"physician",JobID:"medical",Name:"Physician",Level:4,BaseSalary:500,RequiredEducation:4,RequiredStat:30}}});err!=nil{t.Fatal(err)};if err:=s.Employ("p1","medical","physician",4,30);err!=nil{t.Fatal(err)};if e,ok:=s.GetEmployment("p1");!ok||e.PositionID!="physician"{t.Fatal("employment missing")}}
