package simulation

import "time"

type SocialWorldEvent struct {
	ActorID string
	GroupID string
	Action  GroupActionType
	At      time.Time
}

type SocialWorldRuntime struct {
	Relationships *RelationshipService
	Groups        *GroupService
	GroupRuntime  *GroupRuntime
	Events        []SocialWorldEvent
}

func (r *SocialWorldRuntime) ApplyGroupEvent(event SocialWorldEvent) {
	if r == nil || r.GroupRuntime == nil {
		return
	}
	r.GroupRuntime.ApplyDecision(event.GroupID, GroupAction{Type: event.Action, TargetID: event.ActorID})
	r.Events = append(r.Events, event)
}

func (r *SocialWorldRuntime) RecentEvents(limit int) []SocialWorldEvent {
	if r == nil || limit <= 0 {
		return nil
	}
	start := len(r.Events) - limit
	if start < 0 {
		start = 0
	}
	return append([]SocialWorldEvent(nil), r.Events[start:]...)
}
