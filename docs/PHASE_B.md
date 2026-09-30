# Phase B — Integration, Load Validation & Release

## Automated gates
GitHub Actions runs:
- Go unit tests.
- Go backend compilation.
- Android debug compilation.

A green CI run proves compilation/tests for those jobs only; it does not prove staging, database, physical-device or production deployment readiness.

## Staging gates
1. Start PostgreSQL and Redis from the infrastructure baseline.
2. Apply all migrations in order.
3. Start the API with staging secrets.
4. Exercise authentication/session lifecycle.
5. Exercise market inventory transfer, mission XP, combat/equipment and achievements.
6. Exercise world events/news and NPC simulation.
7. Run simulation load tests at 100, 1,000 and 10,000 characters.
8. Verify Android against staging API.
9. Build and install release APK/AAB.
10. Test on physical Android hardware.
11. Verify backup and restore.
12. Review logs and resource usage.

## Release gate
Do not mark Phase B as passed until the staging and release gates above have actually been executed and recorded.
