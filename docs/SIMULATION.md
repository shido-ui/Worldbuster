# Living Population Simulation

The simulated population is a persistent domain, not a collection of disposable request-time bots.

Each simulated character has identity, controller type, personality, goals, last simulation tick and active/background state.

The authoritative game services remain responsible for inventory, money, jobs, organizations and other state. The simulation layer chooses when a simulated character should act and calls validated domain operations.

Simulation frequency is tiered: active, recently active, background and dormant. Future phases add schedules, needs, jobs, relationships, economy behavior, travel, missions, memory and event-driven reactions.


## Decision architecture

A simulated character reads validated world context, evaluates goals and schedules, selects an allowed action, then submits it to the normal authoritative game service. The controller cannot directly edit money, inventory, jobs, relationships or other authoritative state.

Foundation actions include work, travel, study, socialization and rest. Simulation frequency remains tiered: active, recently active, background and dormant.
