# Living Population Simulation

The simulated population is a persistent domain, not a collection of disposable request-time bots.

Each simulated character has identity, controller type, personality, goals, last simulation tick and active/background state.

## Decision architecture

A simulated character reads validated world context, evaluates goals and schedules, selects an allowed action, then submits it to the normal authoritative game service. The controller cannot directly edit money, inventory, jobs, relationships or other authoritative state.

Foundation actions include work, travel, study, socialization and rest. Simulation frequency remains tiered: active, recently active, background and dormant.

## Behavioral state

Simulated characters have energy, social need, rest need, satisfaction, mood, stress and structured memories. Memory is bounded and structured rather than an unbounded chat transcript. Later systems can connect these memories to authoritative events and relationships.
