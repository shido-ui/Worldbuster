# Worldbuster Roadmap

## Phase 1 — Foundation Repair
1. Application bootstrap, routing, scheduler isolation, and graceful shutdown
2. PostgreSQL as production source of truth with explicit test fakes
3. Immutable/checksummed migration system and schema drift checks
4. Idempotent fixture-driven seeding
5. Redis integration for sessions, rate limits, presence, realtime fan-out, and short-lived caches
6. Structured logging, metrics, tracing, health/readiness, and operational dashboards
7. Security hardening and audit logging
8. Race, fuzz, integration, and economy-concurrency testing

## Phase 2 — Core Gameplay Loops
9. Resource bars and lazy regeneration
10. Character stats and training
11. Crime/progression loop
12. Combat and competitive systems
13. Jobs, education, and player companies
14. Transactional economy, market, trading, banking, property, and stock systems
15. Travel and world map
16. Factions and territory
17. Missions, achievements, rankings, and newspaper
18. Social graph, messaging, chat, and notifications

## Phase 3 — Differentiators
19. Living deterministic NPC world
20. Dynamic event director
21. Provider-agnostic AI gateway with deterministic state validation
22. Redis-backed realtime events and presence
23. First-class mobile client and push notifications
24. Public API, API keys, scopes, rate limits, and webhooks
25. Player-generated content and seasonal ladders

## Phase 4 — Frontend and UX
26. React/TypeScript application architecture and routing
27. Guided onboarding and progressive disclosure
28. Mobile-responsive design system, accessibility, localization, and live UI
29. Frontend unit/E2E tests and performance budget
30. In-game wiki and mentor tooling

## Phase 5 — Scale and Operations
31. Load testing and capacity model
32. Postgres query/index/partition tuning
33. Durable event/outbox architecture
34. Horizontal scaling and leader election
35. Kubernetes/Terraform staging and production infrastructure
36. CI/CD, feature flags, backups, restore testing, and incident runbooks

## Phase 6 — Fair Monetization and Anti-Cheat
37. Non-pay-to-win supporter/cosmetic/convenience model
38. Transparent premium-currency accounting
39. Anti-bot and multi-account abuse detection
40. Moderation, appeals, audit history, and compliance tooling

## Phase 7 — Community and Launch
41. Public changelog and roadmap
42. Closed beta telemetry and economy balancing
43. Retention/funnel/session analytics
44. Hotfixable configuration with audit history
45. Public launch, status, and support operations

The roadmap intentionally excludes restricted gambling mechanics and focuses monetization on non-gambling systems.
