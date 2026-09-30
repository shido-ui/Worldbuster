# Living Population Simulation

The simulated population is persistent and uses the same authoritative game services as players.

## Runtime flow
1. Population scheduler selects active simulated characters.
2. World context is built from authoritative services.
3. Goals/schedules produce a candidate action.
4. The action is validated against current world context.
5. The normal domain executor performs the action.
6. Behavioral state updates.
7. A structured event is emitted.

## Population generation
Population generation is deterministic from a supplied seed and data profile. Generated characters receive stable simulated identity, personality defaults and initial goals. Generation is a provisioning operation, not an uncontrolled per-request bot spawn.

Future balancing will use configurable population distributions, regional density, profession weights, activity tiers and lifecycle rules.


## Lifecycle and regional distribution

Population is divided into active, recent, background and dormant tiers so simulation compute can be budgeted instead of ticking every simulated character equally. Regional weights provide deterministic population density across locations. These controls are configuration-driven and remain separate from authoritative action validation.


## Domain execution boundary

Simulated actions now have an explicit domain-executor boundary. The executor is intentionally the single mutation point for simulated behavior and is designed to delegate each action to the same authoritative services used by human players. Unsupported actions are rejected rather than silently mutating state.
