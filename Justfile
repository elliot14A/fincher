# Fincher — Autonomous Media Delivery Operations Justfile
set shell := ["bash", "-uc"]

# Show available recipes
default:
    @just --list

# Run both backend and frontend concurrently (Ctrl+C terminates both)
dev:
    #!/usr/bin/env bash
    trap 'kill 0' EXIT INT TERM
    echo "==> Starting Fincher backend (:8080) and frontend (:5173)..."
    go run ./cmd/fincher &
    (cd web && bun run dev) &
    wait

# Run backend only (Fincher Go server on :8080)
dev-be:
    go run ./cmd/fincher

# Run backend with air live-reload
dev-be-air:
    air

# Run frontend only (Vite dev server on :5173)
dev-fe:
    cd web && bun run dev

# Aliases for fast developer access
alias fe := dev-fe
alias be := dev-be
alias run := dev

# Build both frontend assets and backend Go binary
build: build-fe build-be

# Build frontend production bundle into pkg/web/dist
build-fe:
    cd web && bun run build

# Build backend Go binary into bin/fincher
build-be:
    go build -o bin/fincher ./cmd/fincher

# Build seed CLI binary into bin/seed
build-seed:
    go build -o bin/seed ./cmd/seed

# Seed SQLite catalog and ClickHouse historical QC events
seed:
    go run ./cmd/seed

# Reset and re-seed clean dataset (clears SQLite + ClickHouse tables)
seed-reset:
    go run ./cmd/seed --reset

# Run all backend and frontend tests
test: test-be test-fe

# Run Go unit and integration tests
test-be:
    go test -count=1 ./...

# Run frontend unit tests
test-fe:
    cd web && bun test

# Run code style, lint, and formatting checks
lint: lint-be lint-fe

# Run Go vet checks
lint-be:
    go vet ./...

# Run Biome lint checks on frontend
lint-fe:
    cd web && bun run lint

# Auto-format and fix frontend lint issues
lint-fe-fix:
    cd web && bun run lint:fix

# Run TypeScript type check on frontend
typecheck:
    cd web && bun run typecheck

# Regenerate frontend Hey API SDK from backend openapi/swagger.json
gen-api:
    cd web && bun run generate:api

# Start local ClickHouse and ClickHouse MCP containers
docker-up:
    docker compose up -d

# Stop local ClickHouse containers
docker-down:
    docker compose down

# Follow ClickHouse container logs
docker-logs:
    docker compose logs -f
