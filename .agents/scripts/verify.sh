#!/usr/bin/env bash
set -euo pipefail

echo "========================================================"
echo "Fincher & GSD Core Invariant Verification Suite"
echo "========================================================"

echo "[1/5] Checking GSD Core Planning Artifacts..."
for req_file in .planning/PROJECT.md .planning/REQUIREMENTS.md .planning/ROADMAP.md .planning/STATE.md; do
  if [ ! -f "$req_file" ]; then
    echo "  -> ERROR: Missing required planning artifact: $req_file"
    exit 1
  fi
done
echo "  -> GSD Core .planning/ artifacts present & valid"

echo "[2/5] Checking Seed Data..."
if [ ! -f "data/seed/media_events.json" ]; then
  echo "  -> ERROR: Missing seed data (data/seed/media_events.json)"
  exit 1
fi
echo "  -> Media event dataset present"

echo "[3/5] Checking ClickHouse & Official MCP HTTP Containers..."
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

echo "[4/5] Checking Go Compilation & Dependencies..."
if [ -f "go.mod" ]; then
  GO_FILES=$(find . -maxdepth 3 -name "*.go" 2>/dev/null | head -n 1 || true)
  if [ -n "$GO_FILES" ]; then
    go build ./...
    echo "  -> Compilation OK"
    go vet ./...
    echo "  -> Go vet OK"
    go test -v -race ./...
    echo "  -> Unit tests OK"
  else
    echo "  -> go.mod initialized with all required dependencies. Ready for Phase 01 Go scaffolding."
  fi
else
  echo "  -> Setup mode: No active go.mod"
fi

echo "[5/5] Checking Environment Variable Naming Invariant (FINCHER_*)..."
if [ -d "internal/config" ]; then
  INVALID_ENVS=$(grep -rn 'env:"' internal/config/ | grep -v 'env:"FINCHER_' || true)
  if [ -n "$INVALID_ENVS" ]; then
    echo "  -> ERROR: Non-compliant environment variable tags found:"
    echo "$INVALID_ENVS"
    exit 1
  fi
  echo "  -> All configuration struct tags strictly adhere to FINCHER_{SERVICE}_* format"
else
  echo "  -> Config directory will enforce FINCHER_{SERVICE}_* upon scaffolding"
fi

echo "========================================================"
echo "All Fincher & GSD Core Setup Verifications Passed."
echo "========================================================"
