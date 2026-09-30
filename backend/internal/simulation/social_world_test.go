package simulation

import (
 "testing"
 "time"
)

func TestSocialWorldRuntime(t *testing.T){
 r:=&SocialWorldRuntime{GroupRuntime:NewGroupRuntime()}
 r.GroupRuntime.AddResource("g",100,0)
 r.ApplyGroupEvent(SocialWorldEvent{ActorID:"npc",GroupID:"g",Action:GroupInvest,At:time.Now()})
 v,_:=r.GroupRuntime.Resource("g")
 if v.Currency!=50 || v.Influence!=10 { t.Fatal("group event was not applied") }
 if len(r.RecentEvents(10))!=1 { t.Fatal("event history missing") }
}
