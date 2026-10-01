package simulation

import (
	"context"
	"time"
)

type ActionExecutor func(action Action, character SimCharacter) ActionResult
type EventSink func(eventType, actorID, targetID string, payload map[string]any)
type PersistenceErrorSink func(operation, characterID string, err error)
type Runner struct {
	Population       *Service
	State            *StateService
	BuildContext     ContextBuilder
	Execute          ActionExecutor
	Emit             EventSink
	Persistence      PersistenceSink
	PersistenceError PersistenceErrorSink
}

func (r *Runner) reportPersistenceError(operation, characterID string, err error) {
	if err != nil && r.PersistenceError != nil {
		r.PersistenceError(operation, characterID, err)
	}
}

func (r *Runner) Tick(now time.Time) int {
	if r.Population == nil {
		return 0
	}
	return r.TickCharacters(r.Population.List(1000), now)
}

func (r *Runner) TickCharacters(chars []SimCharacter, now time.Time) int {
	count := 0
	seen := map[string]struct{}{}
	for _, c := range chars {
		if !c.Active || c.ID == "" {
			continue
		}
		if _, duplicate := seen[c.ID]; duplicate {
			continue
		}
		seen[c.ID] = struct{}{}
		r.Population.Touch(c.ID, now)
		if r.Persistence != nil {
			if err := r.Persistence.EnsureCharacter(context.Background(), c); err != nil {
				r.reportPersistenceError("ensure_character", c.ID, err)
				continue
			}
		}
		ctx := r.BuildContext(c.ID)
		state := r.State.Get(c.ID)
		result := DecideAndValidateWithState(c, ctx, state)
		if !result.Accepted {
			if r.Persistence != nil {
				r.reportPersistenceError("record_action", c.ID, r.Persistence.RecordAction(context.Background(), c.ID, result.Action, false))
			}
			r.State.Tick(c.ID)
			continue
		}
		if !r.State.CanAffordAction(c.ID, result.Action.Type) {
			if r.Persistence != nil {
				r.reportPersistenceError("record_action", c.ID, r.Persistence.RecordAction(context.Background(), c.ID, result.Action, false))
			}
			r.State.Tick(c.ID)
			continue
		}
		executed := r.Execute(result.Action, c)
		if !executed.Accepted {
			if r.Persistence != nil {
				r.reportPersistenceError("record_action", c.ID, r.Persistence.RecordAction(context.Background(), c.ID, executed.Action, false))
			}
			r.State.Tick(c.ID)
			continue
		}
		_, economicOK := r.State.ApplyEconomy(c.ID, result.Action.Type)
		if !economicOK {
			if r.Persistence != nil {
				r.reportPersistenceError("record_action", c.ID, r.Persistence.RecordAction(context.Background(), c.ID, result.Action, false))
			}
			r.State.Tick(c.ID)
			continue
		}
		r.State.Apply(c.ID, result.Action.Type, now)
		if r.Persistence != nil {
			r.reportPersistenceError("record_action", c.ID, r.Persistence.RecordAction(context.Background(), c.ID, result.Action, true))
		}
		if result.Action.Type == ActionSocialize && result.Action.TargetID != "" {
			rel := state.Relationships[result.Action.TargetID]
			rel.FromID = c.ID
			rel.ToID = result.Action.TargetID
			rel.Kind = RelationshipAcquaintance
			rel.Familiarity = clamp(rel.Familiarity + 5)
			rel.Trust = clampSigned(rel.Trust + 2)
			rel.Affinity = clampSigned(rel.Affinity + 4)
			rel.Interactions++
			rel.LastInteraction = &now
			r.State.RecordRelationship(c.ID, rel, now)
			if rp, ok := r.Persistence.(RelationshipPersistence); ok {
				r.reportPersistenceError("record_relationship", c.ID, rp.RecordRelationship(context.Background(), rel))
			}
			if rep, ok := r.Persistence.(ReputationPersistence); ok {
				r.reportPersistenceError("propagate_social_reputation", c.ID, rep.PropagateSocialReputation(context.Background(), c.ID, 1, 2, "positive social interaction"))
			}
			if r.Emit != nil {
				r.Emit("SIMULATION_RELATIONSHIP_CHANGED", c.ID, rel.ToID, map[string]any{"affinity": rel.Affinity, "trust": rel.Trust, "familiarity": rel.Familiarity, "interactions": rel.Interactions})
			}
		}
		r.State.Tick(c.ID)
		if r.Persistence != nil {
			st := r.State.Get(c.ID)
			if len(st.Memories) > 0 {
				r.reportPersistenceError("record_memory", c.ID, r.Persistence.RecordMemory(context.Background(), c.ID, st.Memories[len(st.Memories)-1]))
			}
			if ep, ok := r.Persistence.(EconomyPersistence); ok {
				r.reportPersistenceError("record_economy", c.ID, ep.RecordEconomy(context.Background(), c.ID, st.Economy))
			}
		}
		if r.Emit != nil {
			st := r.State.Get(c.ID)
			r.Emit("SIMULATED_ACTION", c.ID, result.Action.TargetID, map[string]any{"action": string(result.Action.Type), "mood": st.Mood, "balance": st.Economy.Balance})
		}
		count++
	}
	return count
}
