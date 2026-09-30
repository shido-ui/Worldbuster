# Simulated Population

The simulation layer models persistent non-player characters using the same character-domain concepts as human players.

Design:
- deterministic rules remain authoritative
- each simulated character has goals, personality and a schedule
- controllers choose from validated actions
- active population receives higher simulation frequency
- background population is processed in batches
- simulation emits domain events instead of directly mutating unrelated systems
- simulated accounts remain auditable as SIMULATED internally

This phase adds the population/controller foundation; detailed professions, economies, relationships and behaviors will plug into the existing domain services.
