# Worldbuster Database

PostgreSQL is the authoritative persistent store when WORLDBUSTER_DATABASE_URL is configured.

| Migration | Purpose |
|---|---|
| 001_core.sql | accounts and player profiles |
| 002_auth.sql | persistent sessions |
| 003_integrity.sql | integrity constraints |
| 004_economy.sql | economy persistence |
| 005_organizations.sql | organizations |
| 006_jobs.sql | jobs |
| 007_progression.sql | progression |
| 008_social.sql | social data |
| 009_events.sql | world events |
| 010_simulation.sql | simulation state |
| 011_progression.sql | player progression stats |
| 012_skills.sql | player skills |
| 013_unlocks.sql | player unlocks |

At startup, ApplyMigrations creates schema_migrations, executes pending SQL migrations in lexical order, and records successful versions.

Identity flow: session -> account -> player profile. Authoritative operations resolve the player from the authenticated account rather than trusting a client-supplied player ID.

Without WORLDBUSTER_DATABASE_URL, local development uses in-memory authentication.


### Phase 79 migration
- `016_salary_settlement.sql` adds `player_employment.last_paid_at` for idempotent scheduled salary settlement.
- Salaries are settled transactionally with the economy ledger and player XP.
- Salary processing is capped per batch and uses row locks so concurrent workers do not double-pay the same employment.
