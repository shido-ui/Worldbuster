# Worldbuster Production Infrastructure

This directory defines a local production-like baseline.

## Services
- PostgreSQL: authoritative persistent state.
- Worldbuster API: Go server.

## Production rules
1. Replace all example credentials through secrets/environment management.
2. Use TLS at the edge and keep the database private.
3. Do not expose PostgreSQL publicly.
4. Back up PostgreSQL and periodically verify restoration.
5. Set resource limits and monitor CPU, memory, database latency and API errors.
6. Scale API instances only after shared state is moved to production database infrastructure.
7. Never treat client state as authoritative.

## Required secrets
- WORLDBUSTER_DATABASE_URL
- session/signing secrets when introduced
- external service credentials when introduced

The compose file is a development/staging baseline, not a claim of production readiness.
