# Repository Guidelines

## Project Structure & Module Organization
- Entry point `main.go` wires `internal/cmd`; domain code in `internal/controller`, `internal/logic`, `internal/service`, `internal/dao`, `internal/model`, `internal/consts`, and generated `internal/packed`.
- HTTP/RPC API contracts live under `api/...`; endpoint notes in `docs/API.md`.
- Configuration defaults in `config/config.yaml`; override via environment/GoFrame config sources.
- Deployment assets in `manifest/{deploy,docker,config,i18n,protobuf}`; database schema in `hack/schema.sql`; static files in `resource/{public,template}`.
- Makefile helpers in `hack/*.mk` drive generators; outputs land back in `internal/*`.

## Build, Test, and Development Commands
- `make up`: update GoFrame CLI before other tasks.
- `make build`: run `gf build -ew` using `hack/config.yaml`.
- `make ctrl|dao|service|enums|pb|pbentity`: regenerate controllers, DAO, services, enums, or protobuf bindings from `api`/`manifest` definitions.
- `make image TAG=<tag> [PUSH=-p]`: build (and optionally push) Docker image; `make deploy TAG=<tag> _ENV=<overlay>` applies kustomize manifests via kubectl.
- Quick checks: `go test ./...` for unit tests; add `-race` for concurrency issues.

## Coding Style & Naming Conventions
- Go 1.23 code formatted with `gofmt`; keep imports gofmt-ordered.
- Package names are lowercase without underscores; files in `api`/`internal` use descriptive nouns/verbs (e.g., `file_upload.go`).
- Request/response structs mirror API definitions; keep validation near controllers and business rules in logic/service layers.
- Use contextual errors (`fmt.Errorf("...: %w", err)`) and avoid global state outside config wiring.

## Testing Guidelines
- Place `_test.go` files next to implementation packages; name tests `TestXxx` and subtests with `t.Run`.
- Prefer table-driven tests and helper builders; mock external IO where possible.
- Run `go test ./...` before PRs; include coverage notes when touching core logic.

## Commit & Pull Request Guidelines
- Repository history is not available; default to concise, imperative messages (Conventional Commits style recommended, e.g., `feat: add telegram webhook handler`).
- PRs should describe intent, link issues/tasks, list commands run, and note config/db/migration changes; include screenshots or sample payloads for API-visible changes.
