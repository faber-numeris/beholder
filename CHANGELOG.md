## v0.1.0 (2026-06-12)

### Feat

- implement hashing

### Fix

- migration script for local docker compose setup (#61)
- credentials verification (#59)
- existing user message (#58)
- email message (#48)

### Refactor

- implement hexagonal architecture for postgres persistence using sqlx
- dependency injection
- restructure to Hexagonal Architecture (Ports and Adapters) (#49)

## v0.0.3 (2026-03-07)

### Fix

- email templates

## v0.0.2 (2026-03-07)

### Fix

- email confirmation (#46)
- email confirmation (#45)
- extract password hashing to dedicated service (#38)
- wrong statement on sqlc query

### Refactor

- point out some potential refactoring (#8)
