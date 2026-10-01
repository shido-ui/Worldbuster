# Production Readiness

## Deployment checklist
- [ ] TLS termination configured
- [ ] Strong database credentials stored as secrets
- [ ] PostgreSQL backups configured and restore tested
- [ ] Shared-state/cache infrastructure selected when required
- [ ] API health checks wired to deployment platform
- [ ] CPU/memory/database metrics collected
- [ ] Error logs centralized
- [ ] Rate limits tuned per endpoint/account/IP
- [ ] Authentication and mutation endpoints tested under load
- [ ] Database migrations reviewed before release
- [ ] Android release signing configured
- [ ] Release build tested against staging API

No production deployment is implied by the presence of these files.
