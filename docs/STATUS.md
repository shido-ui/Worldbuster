# Worldbuster Status

Current development phase: 93

The simulation now covers autonomous NPC decisions, behavioral needs, persistent memories, relationships, reputation, and a deterministic economic wallet.

Phase 93 adds autonomous economic participation: NPC work earns simulated currency; selected activities consume it; insufficient funds block the action; balances and lifetime flows persist to PostgreSQL.

Known production hardening:
- Rehydrate simulation state from PostgreSQL after restart.
- Replace per-tick persistence with batching/queues at scale.
- Integrate NPC wallets with the authoritative player/economy ledger where appropriate.
- Market inventory escrow/ownership integration remains required.
- Territory influence authorization and aggregate-control updates remain required.
- Add broad load/integration tests.
