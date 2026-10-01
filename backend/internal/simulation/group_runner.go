package simulation

import "time"

type GroupContextBuilder func(group Group) GroupContext
type GroupEventSink func(groupID string, action GroupAction, at time.Time)

type GroupRunner struct {
	Groups       *GroupService
	Runtime      *GroupRuntime
	BuildContext GroupContextBuilder
	Emit         GroupEventSink
}

func (r *GroupRunner) Tick(now time.Time) int {
	if r == nil || r.Groups == nil || r.Runtime == nil || r.BuildContext == nil {
		return 0
	}
	count := 0
	for _, g := range r.Groups.List() {
		action := ChooseGroupAction(r.BuildContext(g))
		if action.Type == "" || !r.Runtime.ApplyDecision(g.ID, action) {
			continue
		}
		if r.Emit != nil {
			r.Emit(g.ID, action, now)
		}
		count++
	}
	return count
}
