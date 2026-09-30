# Jobs / Organizations Foundation

Worldbuster now has a shared organization domain for companies, factions and guild-like groups.

Implemented foundation:
- organization identity and type
- owner identity
- member capacity
- membership roles
- join/leave lifecycle
- server-side validation
- PostgreSQL schema
- unit tests

The domain intentionally separates membership from future payroll, production, faction reputation, territory, permissions and progression systems. Those systems will consume organization membership rather than duplicating it.
