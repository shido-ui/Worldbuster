# Worldbuster Architecture

## Core rule
**Code owns truth. Simulation creates consequences. AI adds contextual intelligence.**

## Stack
- Go: authoritative backend and deterministic simulation
- TypeScript/React: web client
- Python: isolated AI/analytics services
- PostgreSQL: authoritative persistence
- Redis: cache/session/rate limiting/hot state

## Runtime
Client -> API -> auth -> validation -> rules -> transaction -> domain events -> secondary systems -> response.

Human and simulated characters share the same authoritative state model; only their controller differs.
