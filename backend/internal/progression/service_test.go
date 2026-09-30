package progression
import "testing"
func TestProgression(t *testing.T){s:=NewService();c:=s.AddXP("p1",1000);if c.Level!=2||c.XP!=0{t.Fatalf("unexpected level: %+v",c)};if err:=s.RegisterCourse(Course{ID:"basic",Name:"Foundation Studies",DurationHours:10,EducationGain:5,RequiredLevel:1});err!=nil{t.Fatal(err)};if err:=s.Enroll("p1","basic");err!=nil{t.Fatal(err)}}
