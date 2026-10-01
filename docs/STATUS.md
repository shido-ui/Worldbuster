# Worldbuster Status

## Foundation repair in progress

The repository has entered the foundation-repair track. The application entrypoint is now isolated from dependency wiring, background jobs are independently scheduled with panic/error isolation, graceful shutdown is enabled, and the HTTP address is configurable through `WORLDBUSTER_ADDR`.

CI now verifies Go formatting and module cleanliness, and project contribution/security/license policies are checked into the repository.

### Completed in this track

- Application bootstrap moved into `backend/internal/app/bootstrap.go`.
- Route dependency wiring moved into `backend/internal/app/routes.go`.
- Background simulation, production, territory, economy, salary, and education work are independent scheduler jobs.
- Scheduler jobs stop through the application context and isolate failures from one another.
- Graceful HTTP shutdown is wired to SIGINT/SIGTERM.
- PostgreSQL is required for production startup; there is no production in-memory downgrade.
- `WORLDBUSTER_ADDR` controls the HTTP listen address.
- `go.sum`, LICENSE, CONTRIBUTING.md, SECURITY.md, CODE_OF_CONDUCT.md, and golangci-lint configuration are present.
- CI runs golangci-lint and verifies module tidiness without a standalone database-job tidy step.
- Production bootstrap now requires PostgreSQL; persistent routes no longer silently downgrade to in-memory behavior.
- Job persistence is injected through a `store.JobStore` interface, with the PostgreSQL repository as the production implementation.
- World content and NPC population configuration are loaded from `database/seed/content.json` through a validated seed loader.
- CI checks formatting and fails when `go mod tidy` would modify tracked module files.
- Redis has a production cache boundary at `backend/internal/cache`; configured deployments validate connectivity during startup.
- `/health`, `/ready`, and `/metrics` provide liveness, dependency readiness, and Prometheus-compatible request telemetry.
- Request IDs are propagated through `X-Request-ID`, and metrics middleware preserves `http.Flusher` for SSE.

### Next foundation work

1. Validate the refactor with the full Go race/vet/build/lint matrix.
2. Add fresh-DB migration/schema verification and explicit legacy migration compatibility.
3. Extract interfaces for player/economy/inventory persistence and complete dependency-container wiring.
5. Continue persistence/domain interface extraction and observability.
6. Then build the gameplay loops end-to-end.

## Current phase

Phase 1 — Foundation Repair
