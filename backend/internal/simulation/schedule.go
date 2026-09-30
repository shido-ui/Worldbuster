package simulation

type ScheduleBlock struct {
 StartHour int `json:"startHour"`
 EndHour int `json:"endHour"`
 Action ActionType `json:"action"`
}

type Schedule struct { Blocks []ScheduleBlock `json:"blocks"` }

func (s Schedule) ActionAt(hour int) (ActionType,bool) {
 for _,b:=range s.Blocks {
  if hour>=b.StartHour && hour<b.EndHour { return b.Action,true }
 }
 return "",false
}
