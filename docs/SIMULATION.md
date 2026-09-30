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

The simulation layer does not bypass inventory, economy, jobs, travel, organizations or other authoritative systems.

This architecture lets simulated characters participate in the same world rather than living in a separate fake database.
