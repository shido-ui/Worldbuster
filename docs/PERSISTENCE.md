# Persistent State & Transaction Layer

Phase 4 establishes the durable-state boundary.

- PostgreSQL is authoritative for durable account/player state.
- Repository code owns SQL access.
- Multi-record authoritative mutations use one transaction.
- Transactions commit only after the complete operation succeeds; otherwise they roll back.
- Cash, XP and energy are server-owned values.
- Expired sessions are rejected at the database boundary.
- Repositories return domain models rather than leaking SQL rows.

Repositories now exist for accounts, sessions and player profiles. The existing in-memory services remain the development runtime until the application wiring switches to PostgreSQL.
