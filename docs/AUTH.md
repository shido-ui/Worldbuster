# Identity & Authentication

Worldbuster now has an authentication boundary built around bcrypt password hashing, normalized unique usernames, random 256-bit session identifiers, HttpOnly SameSite=Lax cookies, server-side session resolution/revocation, and generic login failure responses.

API:
- POST /api/v1/auth/register
- POST /api/v1/auth/login
- GET /api/v1/auth/me
- POST /api/v1/auth/logout

The current auth service is deliberately in-memory. PostgreSQL persistence is the next infrastructure step; the SQL session schema is already prepared.

The client never becomes authoritative for account identity, session validity, progression, cash, inventory, or world state.
