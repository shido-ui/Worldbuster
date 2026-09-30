# Living Population Simulation

The simulated population is a persistent domain, not a collection of disposable request-time bots.

Each simulated character has identity, controller type, personality, goals, last simulation tick and active/background state.

## Decision architecture

A simulated character reads validated world context, evaluates goals and schedules, selects an allowed action, then submits that action to the normal authoritative game service. The controller cannot directly edit money, inventory, jobs, relationships or other authoritative state.

The decision layer now validates its selected action against a world context before it can proceed. This is the bridge between simulation and authoritative gameplay systems.

Foundation actions include work, travel, study, socialization and rest. Simulation frequency remains tiered: active, recently active, background and dormant.

## Behavioral state

Simulated characters have energy, social need, rest need, satisfaction, mood, stress and structured memories. Memory is bounded and structured rather than an unbounded chat transcript.

Future connectors will build world context from actual job, progression, travel, inventory, economy, social and organization services, then route accepted actions into those authoritative services.
