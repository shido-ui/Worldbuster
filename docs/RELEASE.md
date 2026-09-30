# Worldbuster Release Checklist

## Backend
- [ ] Run all database migrations against a staging PostgreSQL instance.
- [ ] Verify authentication, sessions, permissions and rate limits.
- [ ] Exercise economy, organizations, jobs, progression, missions, combat, equipment, events, news, achievements and simulation.
- [ ] Run load tests at 100 / 1,000 / 10,000 simulated characters.
- [ ] Verify backups and restore procedure.
- [ ] Configure TLS, secrets and production observability.

## Android
- [ ] Point the app at the staging API.
- [ ] Verify login/session lifecycle.
- [ ] Verify world, market, missions and profile data.
- [ ] Test API failures and offline/retry behavior.
- [ ] Configure release signing.
- [ ] Build and install the signed release APK/AAB.
- [ ] Test on physical Android hardware.

## Release gate
A release is ready only after the unchecked items above have been executed and verified. Repository configuration alone is not proof that those checks passed.
