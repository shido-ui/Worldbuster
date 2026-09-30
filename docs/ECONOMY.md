# Economy Foundation

Worldbuster now has a server-authoritative currency ledger.

- Money changes only through server operations.
- Every change creates a ledger entry.
- Balances are authoritative derived state and can be reconciled against history.
- Debits reject insufficient funds.
- Entries contain reason and optional reference ID.
- Currency is explicit (`WBX`) for future expansion.

Production will add PostgreSQL repositories, double-entry transfers, marketplace settlement, fees, business accounts and idempotency controls.
