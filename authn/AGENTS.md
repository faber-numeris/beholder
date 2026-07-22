# AGENTS.md — authn service

## Database access

All SQL access in this service goes through **sqlc**-generated code in
`internal/adapters/outbound/postgres/gen`. Never write hardcoded SQL strings
in repository/adapter files (e.g. `internal/adapters/outbound/postgres/*.go`).

Workflow for any new or changed query:

1. Add or edit the query in `internal/adapters/outbound/postgres/queries/*.sql`
   (sqlc query syntax, using `sqlc.arg`/`sqlc.narg` for parameters).
2. Regenerate with `just generate-sqlc` (runs `sqlc generate`). This is the
   only way files under `internal/adapters/outbound/postgres/gen` should
   ever change — never hand-edit that directory.
3. Call the generated method from the repository (e.g.
   `internal/adapters/outbound/postgres/user_repository.go`), depending on
   the `gen.Querier` interface rather than a concrete `*gen.Queries` or a
   `database/sql`/`sqlx` handle, so tests can inject a fake implementing
   `gen.Querier`.

If a nullable `timestamptz` column needs to round-trip through Go, prefer the
`*time.Time` override in `sqlc.yaml` (see the `overrides` section) over
scanning into a bare `time.Time` — pgx cannot scan SQL `NULL` into a
non-pointer `time.Time` and will error at runtime.
