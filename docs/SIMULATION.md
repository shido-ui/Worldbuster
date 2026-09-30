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
Population generation is deterministic from a supplied seed and data profile. Generated characters receive stable simulated identity, personality defaults and initial goals.

## Lifecycle and regional distribution
Population uses active, recent, background and dormant tiers. Regional weights control deterministic population density.

## Behavioral simulation
NPC needs, stress, mood, satisfaction and bounded memories evolve after actions and over time. Decision priorities consume current behavioral state.

## Social simulation
Relationships persist with familiarity and trust. Social actions can deepen acquaintances into friendships or preserve rivalries.

## Group formation
A lightweight persistent group layer now supports NPC community formation: a leader can create a group, characters can join one group, and membership is queryable. This is the foundation for later organizations, factions, companies, gangs, clubs and other larger social structures.
