# Worldbuster Status

Current development phase: 106

NPCs now have the first end-to-end production pipeline: production recipes consume business inventory, charge deterministic labor costs, create outputs, record production runs, and update world supply signals. PostgreSQL-backed startup seeds a small set of NPC production businesses with safe generic inputs.

Phase 95 is focused on production persistence and supply propagation. Phase 106 adds database-backed world operator roles, audited world-control flags, and authenticated control APIs.
