<!--
Guidance for AI coding agents working on the `aoc2024` repository.
Keep this short, specific, and focused on discoverable, actionable patterns.
-->

# Quick orientation

- This repo implements an Advent-of-Code 2024 solver service written in Go. There are two runnable programs in `cmd/`:
  - `cmd/cli` — small CLI for running a solver against a local input file (see `cmd/cli/main.go`).
  - `cmd/web` — HTTP server exposing the solvers via REST and a tiny web UI (see `cmd/web/main.go`).
- Solver implementations live in `pkg/dN` (e.g. `pkg/d1`, `pkg/d2`, ...). They register with the runtime via `pkg/solver`.

## Big picture (files to read first)
- `cmd/web/main.go` — how the server is wired (muxes, middleware, TLS toggle, swagger).
- `cmd/web/middleware/middleware.go` — logging, rate-limiting, request context keys and helpers.
- `cmd/web/api/api.go` and `cmd/web/web/web.go` — API handlers and web templates; shows request handling conventions.
- `pkg/solver/solver.go` — solver registry contract (Register, New, ListRegister) and expected interface.
- `pkg/d*/` — example solver implementations. Follow these for new days.
- `Makefile` and `environments/local/*` — CI/CD, local run, Docker and Terraform/Ansible helper targets.

## Key conventions and patterns (explicit, code-backed)
- Solver contract: implement PuzzleSolver with methods Init(io.Reader) error and Solve(part int) (string, error). See `pkg/solver/solver.go` and an example in `pkg/d1/d1.go`.
- Registration: each day package should call `solver.Register("dX", func() PuzzleSolver { return &PuzzleStruct{} })` in its init(). The programs load solvers via blank imports, e.g. `_ "advent2024/pkg/d1"` in `cmd/cli/main.go` and `cmd/web/main.go`.
- Routing and path params: handlers are registered with patterns like `webMux.HandleFunc("GET /solve/{day}/{part}", ...)` and handlers access params via `r.PathValue("day")`.
  As of the repository toolchain (see `go.work` — `go 1.25.0`), `PathValue` is available on `*http.Request` from the standard library's `net/http` when using the newer method+path `HandleFunc` patterns. Check `cmd/web/main.go` for examples of the patterns used (e.g. `"GET /solve/{day}/{part}"`) and `go.work` for the Go version.
- Context logger: middleware injects a logger into the request context. Retrieve it with `middleware.GetLogger(r.Context())` (see `cmd/web/middleware/middleware.go`). Use it for request-scoped logs.

## Common developer workflows (how to build, run, test)
- Use the repository Makefile for high-level workflows:
  - `make localci` — build local Docker images (`environments/local/build_docker_images.sh`).
  - `make localrun` — build and run the web image locally (calls docker run).
  - `make localcd` — spin up a local deployment (calls `environments/local/*` scripts).
  - `make bootstrap`, `make init`, `make apply` — terraform/bootstrap targets for cloud workflows.
- Unit tests: the Makefile has a `test` target which runs `go work sync` then `go test` for each `pkg/*` module. Preferred commands:

```bash
go work sync        # ensure go.work modules are consistent
make test           # runs tests module-by-module as repo expects
```

- Running server locally (fast):

```bash
go run ./cmd/web     # runs the web server (reads flags and env vars)
```

- CLI usage example:

```bash
go run ./cmd/cli --filename=inputs/d1_input.txt --day=d1 --part=1
```

## Important config and environment variables
- `cmd/web/config/config.go` reads flags and environment variables. Relevant keys:
  - `PORT` or `--port` (default 8080)
  - `ENABLE_HTTPS` / `--https` (set to "true" to enable TLS)
  - `TLS_CERT_FILE` / `--cert` and `TLS_KEY_FILE` / `--key` (required when TLS enabled)
  - `API_RATE`, `API_BURST` (rate-limiting)

When changing configuration defaults or adding flags, update `LoadConfig()` in `cmd/web/config/config.go`.

## Local environment scripts (environments/local)

- This repo includes helper scripts and Ansible playbooks to build and run a local VM-based deployment under `environments/local`.
- Key scripts:
  - `environments/local/build_docker_images.sh` — builds `advent2024.web` and `advent2024.cli` Docker images (used by `make localci`).
  - `environments/local/build_local.sh` — boots a local VM using Terraform (downloads a minimal Ubuntu cloud image, runs `terraform init/plan/apply`). Pass `destroy` to tear it down. Called by `make localcd`.
  - `environments/local/configure_run_local.sh` — packages the `advent2024.web` Docker image to a tarball and runs two Ansible playbooks to deploy the web container and a Caddy load-balancer. Pass `cleanup` to remove the tarball.
- Ansible playbooks of interest:
  - `environments/local/ansible/deploy.web.yml` — copies the Docker tarball to the VM, loads it and runs the container exposing port 8080.
  - `environments/local/ansible/deploy.lb.yml` — deploys a Caddyfile and ensures the Caddy service runs as the load-balancer.

Use the Makefile targets to orchestrate these scripts:

```bash
make localci     # build docker images
make localrun    # build and run docker image (docker run)
make localcd     # build local and configure/run via terraform+ansible
```

When editing environment automation, prefer minimal, idempotent changes and update the Makefile if you add new orchestration steps.

## Cloud / AWS environments

- This repo includes Terraform and bootstrap flows for AWS under `environments/aws` and a bootstrap helper under `environments/bootstrap/aws/terraform` (see `docs/README.md` for full details).
- Makefile targets wired to these flows:
  - `make bootstrap` — runs Terraform in the bootstrap folder to create shared resources (S3 backend, DynamoDB locks, GH Action role).
  - `make init` / `make apply` / `make destroy` — run Terraform in the chosen environment. The Makefile reads the `ENVIRONMENT` environment variable (defaults include `prod` and `dev`).
- Prerequisites for cloud workflows:
  - Terraform installed locally
  - AWS CLI installed and configured with a user that has permissions described in `docs/README.md`
  - A domain (some bootstrap steps expect a domain)

Usage notes:

```bash
# Bootstrap shared infra
make bootstrap

# Create or update a specific environment (dev/prod)
make init
make apply

# Tear down
make destroy
```

When changing cloud automation, prefer small, reversible changes and ensure Terraform backend config (`environments/aws/backend.json`) is kept synchronized.

## Adding a new solver (concrete steps)
1. Create `pkg/dX` module matched with existing pattern (see `pkg/d1`): implement `PuzzleStruct`, `Init(io.Reader)` and `Solve(part int)`.
2. In the package's `init()` call `solver.Register("dX", func() solver.PuzzleSolver { return &PuzzleStruct{} })`.
3. Add a blank import to the runner(s) that should load it: `cmd/cli/main.go` and/or `cmd/web/main.go` (they already import many `pkg/d*` packages as blanks).
4. Run `go work sync` and `make test`.

## Routing & API notes (pay attention when editing)
- API endpoints live under `/api/` (see `cmd/web/main.go` where `globalMux` strips `/api`). Example handlers:
  - `GET /api/solvers` — list registered days (`cmd/web/api/api.go`).
  - `POST /api/solvers/{day}/{part}` — run solver on posted base64 input.
- Handlers read path params using `r.PathValue("name")`. In this repo `PathValue` is provided by the standard library in the Go toolchain used here (see `go.work` → `go 1.25.0`). When touching routing, check `cmd/web/main.go` to see the method+path `HandleFunc` patterns and add tests if you change parameter extraction.

## Code quality and style to preserve
- Small, idiomatic Go code. Keep the public solver API stable (Init/Solve). Preserve existing logging prefixes and context value keys (see middleware constants).
- Tests are per `pkg/d*`. When adding new behavior try to add a small unit test in the corresponding package.

## Quick search hints for maintainers
- Find where solvers are used: search `solver.New(` and `ListRegister()`.
- Inspect request-level logging and context keys: `cmd/web/middleware/middleware.go`.
- Swagger/OpenAPI: `cmd/web/docs` contains generated swagger files; regenerate with `swag` if you update annotations.

---
If anything above is unclear or you want me to expand any section (routing internals, adding CI steps, or a short checklist for PRs), tell me which part to iterate on.
