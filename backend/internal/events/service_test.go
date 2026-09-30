package events

import ("testing";"time")

func TestEventsAndNotifications(t *testing.T){s:=NewService();e,err:=s.Publish("TRAVEL_COMPLETED","p1","",map[string]any{"source":"test"});if err!=nil{t.Fatal(err)};if e.ID==""{t.Fatal("event id missing")};n,nerr:=s.Notify("p1","world","Arrival","You arrived.",e.ID);if nerr!=nil{t.Fatal(nerr)};if n.EventID!=e.ID{t.Fatal("event link missing")};if len(s.Recent(10))!=1||len(s.Notifications("p1"))!=1{t.Fatal("event data missing")}}


func TestEventSubscribers(t *testing.T){
 s:=NewService()
 ch:=s.Subscribe()
 e,err:=s.Publish("LIVE","actor","target",nil);if err!=nil{t.Fatal(err)}
 select{case got:=<-ch: if got.ID!=e.ID{t.Fatal("subscriber received wrong event")}
 case <-time.After(time.Second): t.Fatal("subscriber did not receive event")}
 s.Unsubscribe(ch)
 select{case _,ok:=<-ch: if ok{t.Fatal("subscriber channel should be closed")}
 case <-time.After(time.Second): t.Fatal("unsubscribe did not close channel")}
}

func TestSlowSubscriberDoesNotBlockPublish(t *testing.T){
 s:=NewService()
 ch:=s.Subscribe()
 defer s.Unsubscribe(ch)
 for i:=0;i<64;i++{
  done:=make(chan struct{})
  go func(){s.Publish("LIVE","","",nil);close(done)}()
  select{case <-done:case <-time.After(time.Second):t.Fatal("publish blocked on slow subscriber")}
 }
}
