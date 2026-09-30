# Security Policy

## Reporting a vulnerability

Please do not disclose security vulnerabilities in public issues. Report them privately to the repository maintainers through the project's GitHub security reporting channel.

Include the affected component, reproduction steps, impact, and any relevant logs without including secrets or personal data.

## Security principles

Worldbuster treats the server as authoritative. Sensitive actions should be validated server-side, persisted transactionally, rate-limited, and recorded in audit logs.
