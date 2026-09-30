# Worldbuster Status

Current development phase: 97

NPCs now have the first end-to-end production pipeline: production recipes consume business inventory, charge deterministic labor costs, create outputs, record production runs, and update world supply signals. PostgreSQL-backed startup seeds a small set of NPC production businesses with safe generic inputs.

Phase 95 is focused on production persistence and supply propagation. Phase 97 adds authoritative territory consequences: organization membership is required to add influence, influence determines control, contested/stable control changes territory stability, and territory effects are persisted for downstream economy/world systems.
