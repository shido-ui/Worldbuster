package events
import "testing"
func TestEventsAndNotifications(t *testing.T){s:=NewService();e:=s.Publish("TRAVEL_COMPLETED","p1","","map");if e.ID==""{t.Fatal("event id missing")};n:=s.Notify("p1","world","Arrival","You arrived.",e.ID);if n.EventID!=e.ID{t.Fatal("event link missing")};if len(s.Recent(10))!=1||len(s.Notifications("p1"))!=1{t.Fatal("event data missing")}}
