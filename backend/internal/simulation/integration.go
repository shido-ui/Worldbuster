package simulation

import "time"

// WorldAdapter exposes authoritative domain state to the simulation.
// Implementations must delegate mutations to the corresponding domain service.
type WorldAdapter interface {
 Context(characterID string) WorldContext
 Execute(action Action, character SimCharacter) ActionResult
}

type IntegratedRunner struct {
 Population *Service
 State *StateService
 World WorldAdapter
 Emit EventSink
 Persistence PersistenceSink
}

func (r *IntegratedRunner) Tick(now time.Time) int {
 if r.Population == nil || r.State == nil || r.World == nil { return 0 }
 return r.TickCharacters(r.Population.List(10000), now)
}

func (r *IntegratedRunner) TickCharacters(chars []SimCharacter, now time.Time) int {
 runner := Runner{
  Population: r.Population,
  State: r.State,
  BuildContext: r.World.Context,
  Execute: r.World.Execute,
  Emit: r.Emit,
  Persistence: r.Persistence,
 }
 return runner.TickCharacters(chars, now)
}
