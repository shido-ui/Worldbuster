package simulation

import "time"

func UpdateBehavior(state *BehavioralState, action ActionType, now time.Time) {
	if state == nil { return }
	state.CurrentAction = action
	switch action {
	case ActionWork:
		state.Needs.Energy -= 8
		state.Stress += 3
		state.Needs.Satisfaction += 2
	case ActionStudy:
		state.Needs.Energy -= 5
		state.Stress += 2
		state.Needs.Satisfaction += 1
	case ActionSocialize:
		state.Needs.Energy -= 2
		state.Needs.Social += 8
		state.Stress -= 4
		state.Needs.Satisfaction += 3
	case ActionRest:
		state.Needs.Energy += 12
		state.Stress -= 6
		state.Needs.Rest += 10
	case ActionTravel:
		state.Needs.Energy -= 3
		state.Stress += 1
	}
	clampNeeds(&state.Needs)
	if state.Stress < 0 { state.Stress = 0 }
	if state.Stress > 100 { state.Stress = 100 }
	state.Mood = state.Needs.Satisfaction - state.Stress/2
	if state.Mood < -100 { state.Mood = -100 }
	if state.Mood > 100 { state.Mood = 100 }
	state.Memories = append(state.Memories, Memory{
		ID: now.UTC().Format("20060102150405.000000000"),
		Type: "ACTION",
		Event: string(action),
		Importance: 1,
		CreatedAt: now,
	})
	if len(state.Memories) > 50 {
		state.Memories = state.Memories[len(state.Memories)-50:]
	}
}

func clampNeeds(n *Needs) {
	if n.Energy < 0 { n.Energy = 0 }; if n.Energy > 100 { n.Energy = 100 }
	if n.Social < 0 { n.Social = 0 }; if n.Social > 100 { n.Social = 100 }
	if n.Rest < 0 { n.Rest = 0 }; if n.Rest > 100 { n.Rest = 100 }
	if n.Satisfaction < 0 { n.Satisfaction = 0 }; if n.Satisfaction > 100 { n.Satisfaction = 100 }
}
