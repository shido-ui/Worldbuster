package world

import("testing";"time")

func TestTravelLifecycle(t *testing.T){
 s:=NewTravelService();now:=time.Unix(1000,0)
 tr,err:=s.Start("c1","central","harbor",now);if err!=nil{t.Fatal(err)}
 if tr.ArrivesAt.Sub(now)!=45*time.Second{t.Fatalf("unexpected duration: %v",tr.ArrivesAt.Sub(now))}
 if _,ok:=s.Get("c1",now.Add(44*time.Second));!ok{t.Fatal("travel should remain active")}
 if _,ok:=s.Get("c1",now.Add(45*time.Second));!ok{t.Fatal("completed travel should resolve once")}
 if _,ok:=s.Get("c1",now.Add(46*time.Second));ok{t.Fatal("completed travel should be cleared")}
}
