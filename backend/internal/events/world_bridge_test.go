package events

import (
 "testing"
 "time"
)

func TestWorldEventSink(t *testing.T){
 n:=NewsGenerator{}
 s:=&WorldEventSink{History:NewHistory(),News:&n}
 item,ok:=s.PublishSimulation("SIMULATED_ACTION","npc1","g1","NPC group action occurred",75,time.Now())
 if !ok || item.EventID=="" { t.Fatal("event was not published") }
 if len(s.History.Recent(10))!=1 || len(s.Published)!=1 { t.Fatal("event pipeline incomplete") }
}
