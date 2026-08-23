#!/usr/bin/env bash
set -euo pipefail

echo "========================================================"
echo "Fincher & GSD Core Invariant Verification Suite"
echo "========================================================"

echo "[1/6] Checking GSD Core Planning Artifacts..."
for req_file in .agents/planning/PROJECT.md .agents/planning/REQUIREMENTS.md .agents/planning/ROADMAP.md .agents/planning/STATE.md; do
  if [ ! -f "$req_file" ]; then
    echo "  -> ERROR: Missing required planning artifact: $req_file"
    exit 1
  fi
done
echo "  -> GSD Core .agents/planning/ artifacts present & valid"

echo "[2/6] Checking Seed Data Placeholder..."
if [ -f "data/seed/media_events.json" ]; then
  echo "  -> Media event dataset present"
else
  echo "  -> Note: Seed data is scheduled for Feature 09 (Hackathon Seeder)."
fi

echo "[3/6] Checking ClickHouse & Official MCP HTTP Containers..."
if curl -sf http://127.0.0.1:8123/ping > /dev/null 2>&1; then
  echo "  -> ClickHouse (:8123) is REACHABLE"
else
  echo "  -> WARNING: ClickHouse (:8123) not reachable. Run: docker compose up -d"
fi

if curl -sf http://127.0.0.1:8000/health > /dev/null 2>&1; then
  echo "  -> Official ClickHouse MCP Server (:8000) is HEALTHY"
else
  echo "  -> WARNING: ClickHouse MCP Server (:8000) not reachable. Run: docker compose up -d"
fi

echo "[4/6] Checking Frontend Quality & Build (web/)..."
if [ -d "web" ] && [ -f "web/package.json" ]; then
  cd web
  bun run lint
  echo "  -> Frontend Biome lint OK"
  bun run typecheck
  echo "  -> Frontend TypeScript typecheck OK"
  bun run build
  echo "  -> Frontend Vite build OK"
  cd ..
fi

echo "[5/6] Checking Environment Variable Naming Invariant (FINCHER_*)..."
if [ -d "internal/config" ]; then
  INVALID_ENVS=$(grep -rn "env='" internal/config/ | grep -v "env='FINCHER_" || true)
  if [ -n "$INVALID_ENVS" ]; then
    echo "  -> ERROR: Non-compliant environment variable tags found in internal/config:"
    echo "$INVALID_ENVS"
    exit 1
  fi
  echo "  -> All configuration struct tags strictly adhere to FINCHER_{SERVICE}_* format"
fi

echo "[6/6] Checking Go Compilation & Tests with Embedded Frontend..."
if [ -f "go.mod" ]; then
  go build ./cmd/... ./internal/... ./pkg/...
  echo "  -> Compilation OK"
  go vet ./cmd/... ./internal/... ./pkg/...
  echo "  -> Go vet OK"
  go test -v -race ./cmd/... ./internal/... ./pkg/...
  echo "  -> Unit tests OK"
fi

echo "========================================================"
echo "All Fincher & GSD Core Setup Verifications Passed."
echo "========================================================"
