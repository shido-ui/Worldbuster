# Living Population Simulation

The simulated population is a persistent domain, not a collection of disposable request-time bots.

Each simulated character has identity, controller type, personality, goals, last simulation tick and active/background state.

The authoritative game services remain responsible for inventory, money, jobs, organizations and other state. The simulation layer chooses when a simulated character should act and calls validated domain operations.

Simulation frequency is tiered: active, recently active, background and dormant. Future phases add schedules, needs, jobs, relationships, economy behavior, travel, missions, memory and event-driven reactions.
