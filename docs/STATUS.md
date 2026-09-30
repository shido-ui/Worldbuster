# Worldbuster Status

Current development phase: 92

The simulation now has persistent NPC identities, goals, deterministic decisions, behavioral needs, action history, memory, relationships, and simulated reputation. Human and simulated characters continue to use the same core action model while simulated controllers provide autonomous decisions.

Phase 92 added persistent relationship state with familiarity, trust, affinity, interaction history, social-memory creation, and deterministic simulated reputation propagation.

Known production hardening still required:
- Relationship state rehydration from PostgreSQL on process restart.
- Avoid per-tick character upserts for large populations.
- Batch/queue simulation persistence for scale.
- Market inventory escrow/ownership integration.
- Territory influence authorization and aggregate-control updates.
- Broader test coverage and load testing.
