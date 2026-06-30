# Sync the Bruno collection with the OpenAPI spec file using interactive diffs
sync-api:
    npx @sayedameer/bruno-openapi-sync -s authn/internal/adapters/inbound/httpapi/openapi.yaml -o "authn/tests/"

mod authn "authn/justfile"
