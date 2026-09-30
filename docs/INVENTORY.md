# Inventory Foundation

Inventory is server-authoritative and item definitions are data-driven.

The first foundation supports:
- item definitions
- stackable/non-stackable items
- maximum stack sizes
- per-character slot capacity
- atomic add/remove operations under the service lock
- API retrieval and item addition
- unit tests

The inventory layer deliberately does not decide economic value, combat effects, crafting effects, or marketplace behavior. Those systems will consume validated inventory operations later.

Production persistence will move the same domain rules behind PostgreSQL repositories; the client must never become the authority for quantities.
