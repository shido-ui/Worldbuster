package social
import "testing"
func TestSocial(t *testing.T){s:=NewService();if err:=s.SetRelationship("a","b","friend",10);err!=nil{t.Fatal(err)};if _,ok:=s.Relationship("a","b");!ok{t.Fatal("relationship missing")};if _,err:=s.Send("a","b","hello");err!=nil{t.Fatal(err)};if len(s.Inbox("b"))!=1{t.Fatal("message missing")}}
