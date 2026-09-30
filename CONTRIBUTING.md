# Contributing to Worldbuster

## Development rules

- Keep authoritative game state on the server.
- Route currency and item changes through transactional domain services.
- Make randomness injectable and seedable.
- Add tests and documentation with every feature.
- Run `gofmt`, `go test -race ./...`, `go vet ./...`, and the frontend build before opening a change.
- Keep commits small and focused.
- Do not add secrets to the repository.

## Pull requests

Explain the user-visible change, persistence changes, API changes, tests, and operational impact. Include migration and rollback notes for database changes.
