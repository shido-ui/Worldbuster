# Worldbuster Status

Current development phase: 96

NPCs now have the first end-to-end production pipeline: production recipes consume business inventory, charge deterministic labor costs, create outputs, record production runs, and update world supply signals. PostgreSQL-backed startup seeds a small set of NPC production businesses with safe generic inputs.

Phase 95 is focused on production persistence and supply propagation. Phase 96 adds deterministic supply/demand price pressure, bounded price movement, price history, demand/supply decay, and a scheduled rebalance cycle. Production outputs can now feed matching market goods.
