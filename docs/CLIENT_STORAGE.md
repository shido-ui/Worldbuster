# Client Storage Architecture

The future Android client may reserve up to several GB for local game data when the player permits it.

## Authority
PostgreSQL/server state is authoritative. Local data is a cache/snapshot, never the source of truth for protected gameplay state.

## Local layers
- Room/SQLite: indexed structured cache
- App-private files: large static datasets, downloaded content and snapshots
- Cache: disposable assets
- Sync journal: queued client actions with idempotency keys

## Sync model
1. Server publishes a versioned world/content revision.
2. Client downloads only missing chunks/deltas.
3. Local database applies the snapshot transactionally.
4. Client records the server revision.
5. Gameplay mutations are submitted to the server.
6. Server validates and returns authoritative state/deltas.

The client must tolerate interrupted downloads, schema upgrades, corrupted cache entries and server revisions being ahead of local state.

## Storage budget
The client should expose a configurable storage budget. Large content packs are optional and can be evicted/re-downloaded without losing authoritative progress.
