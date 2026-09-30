# Worldbuster Phases

## Completed
- 87 — Advanced economy and market foundation
- 88 — NPC/bot deterministic decision engine
- 89 — NPC population lifecycle scheduler
- 90 — NPC action execution and world reactions
- 91 — Persistent NPC identity, memory, and action history
- 92 — NPC relationships, social memory, and simulated reputation propagation
- 93 — NPC autonomous economy foundation
- 94 — NPC businesses and employment foundation

## Phase 94 delivered
- Persistent NPC-owned business records.
- Business types: retail, service, production.
- Deterministic startup costs and baseline revenue/wage formulas.
- Business cash/revenue/expense/reputation/employee state.
- Business ledger and employment tables.
- Repository operations for creating businesses and recording expenses.

## Phase 95 delivered
- Deterministic production recipes with safe generic commodities.
- Persistent business inventories for raw inputs and produced outputs.
- Production-run history for auditability.
- Labor-cost consumption tied to business cash.
- World supply signals tracking generated supply.
- Periodic production-cycle execution from the server runtime.
- Seeded NPC production businesses with starter inputs so the pipeline is observable when PostgreSQL is enabled.
- Production remains server/DB authoritative; no AI-generated state changes.

## Phase 96 delivered
- Supply/demand-driven deterministic price adjustment.
- Per-cycle price movement capped to prevent runaway market shocks.
- Price history and economy rebalance audit records.
- Supply/demand pressure decays after each rebalance cycle.
- Production outputs now feed market asset supply when matching goods exist.
- Seeded safe commodity market assets for the production chain.
- Server runtime executes periodic economy rebalancing.
- Balance calculations remain deterministic and server authoritative.

## Next
Phase 97 — Territory consequences.
