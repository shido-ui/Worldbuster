package social
import ("testing";"errors")
func TestSocial(t *testing.T){s:=NewService();if err:=s.SetRelationship("a","b","friend",10);err!=nil{t.Fatal(err)};if _,ok:=s.Relationship("a","b");!ok{t.Fatal("relationship missing")};if _,err:=s.Send("a","b","hello");err!=nil{t.Fatal(err)};if len(s.Inbox("b"))!=1{t.Fatal("message missing")}}


func TestMessageIDEntropyFailureIsReturned(t *testing.T){
 old:=randomRead
 defer func(){randomRead=old}()
 randomRead=func([]byte)(int,error){return 0,errors.New("entropy unavailable")}
 s:=NewService()
 if _,err:=s.Send("a","b","hello");err!=ErrIDEntropyFailure{t.Fatalf("expected entropy failure, got %v",err)}
 if len(s.Inbox("b"))!=0{t.Fatal("message stored after ID failure")}
}
