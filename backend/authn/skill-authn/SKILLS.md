---
name: skill-authn
description: What it does + when to trigger it
---

# Beholder — Architecture & Conventions

## Hexagonal Architecture (Ports & Adapters)

```
internal/
├── core/domain/       # Enterprise business rules (entities, value objects)
├── core/services/     # Application business rules (use cases)
├── ports/inbound/     # Inbound port interfaces (driving side)
├── ports/outbound/    # Outbound port interfaces (driven side)
├── adapters/inbound/  # Adapters implementing inbound ports (HTTP API, etc.)
│   └── httpapi/
├── adapters/outbound/ # Adapters implementing outbound ports (Postgres, Mail, etc.)
│   └── postgres/
│       ├── gen/       # Generated sqlc code
│       └── migrations/
├── app/bootstrap/     # Composition root (DI wiring)
├── infrastructure/    # Cross-cutting concerns (config, DB pool, etc.)
└── platform/          # Shared utilities, mappers
    └── mapper/
        └── generated/ # Generated goverter mappers
```

**Rules:**
- `core/` must never import from `adapters/` or `infrastructure/` — it depends only on `ports/`
- `core/` must never have external library dependencies, except for go standard library components. 
- `adapters/` depends on `ports/` and `core/domain/`
- `ports/` depends on `core/domain/`
- `app/bootstrap/` wires everything together

## Generated Files — Never Edit Manually

The following files are auto-generated. Always use the corresponding `just` recipe to regenerate them after changing their source.

| Directory | Generator | Source | Recipe |
|---|---|---|---|
| `adapters/inbound/httpapi/gen/` | ogen | `openapi.yaml` | `just authn generate-oas` |
| `adapters/outbound/postgres/gen/` | sqlc | `.sql` queries + migrations | `just authn generate-sqlc` |
| `platform/mapper/generated/` | goverter | `goverter.go` config | `just authn generate-mappers` |
| `internal/mocks/` | mockery | `.mockery.yaml` config | _Run `go run github.com/vektra/mockery/v2@latest --config .mockery.yaml` from `authn/`_ |

**Never** edit files in `gen/` or `generated/` directories or `mock_*.go` files by hand. Change the source and regenerate.

## Build & Test Discipline

1. **Always** run `go build ./...` after any code change.
2. **Always** run `go test ./... -count=1` before committing.
3. If tests fail after a change, assume the new code is breaking them and **ask for confirmation** before altering the tests or rolling back.

## Code quality

### REST API

This project is API First, so any changes on the http inbound should be added first on the Openapi.yaml file, have the 
components generated via just authn generate-oas.

If the REST API is changed, Bruno test files should be synchronized via just authn sync-api.


## Documentation

### README.md

Always keep the documentation up to date if :

* The directory structure has changed
* We are changed any automation tool
* We changed anything on the technical stack
    