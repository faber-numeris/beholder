module github.com/faber-numeris/beholder/authn

go 1.26.1

require (
	github.com/agiledragon/gomonkey/v2 v2.14.1
	github.com/alexedwards/argon2id v1.0.0
	github.com/caarlos0/env/v11 v11.4.1
	github.com/faber-numeris/foundation/beholder v0.1.3
	github.com/faber-numeris/foundation/testutils v0.0.0-00010101000000-000000000000
	github.com/go-chi/chi/v5 v5.3.1
	github.com/go-chi/cors v1.2.2
	github.com/google/uuid v1.6.0
	github.com/jackc/pgerrcode v0.0.0-20250907135507-afb5586c32a6
	github.com/jackc/pgx/v5 v5.10.0
	github.com/jmoiron/sqlx v1.4.0
	github.com/lmittmann/tint v1.2.0
	github.com/oapi-codegen/runtime v1.6.0
	github.com/oaswrap/spec-ui v0.2.0
	github.com/stretchr/testify v1.11.1
	gopkg.in/yaml.v3 v3.0.1
)

require (
	github.com/apapsch/go-jsonmerge/v2 v2.0.0 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/dlclark/regexp2 v1.11.5 // indirect
	github.com/getkin/kin-openapi v0.140.0 // indirect
	github.com/go-openapi/jsonpointer v1.0.0 // indirect
	github.com/jackc/pgpassfile v1.0.0 // indirect
	github.com/jackc/pgservicefile v0.0.0-20240606120523-5a60cdf6a761 // indirect
	github.com/jackc/puddle/v2 v2.2.2 // indirect
	github.com/oasdiff/yaml v0.1.1 // indirect
	github.com/oasdiff/yaml3 v0.0.14 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/rogpeppe/go-internal v1.14.1 // indirect
	github.com/santhosh-tekuri/jsonschema/v6 v6.0.2 // indirect
	github.com/stretchr/objx v0.5.3 // indirect
	golang.org/x/crypto v0.54.0 // indirect
	golang.org/x/sync v0.22.0 // indirect
	golang.org/x/sys v0.47.0 // indirect
	golang.org/x/text v0.40.0 // indirect
)

replace github.com/faber-numeris/foundation/testutils => ../../foundation/testutils
