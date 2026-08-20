#!/usr/bin/env bash
set -euo pipefail

# Read stdin to consume payload
cat > /dev/null

ACTIVE_PHASE="01-foundation"
if [ -f ".planning/STATE.md" ]; then
  ACTIVE_PHASE=$(grep -E '^\* \*\*Active Phase\*\*' .planning/STATE.md | head -n1 | sed -E 's/.*`([^`]+)`.*/\1/' || echo "01-foundation")
fi

cat << EOF
{
  "injectSteps": [
    {
      "ephemeralMessage": "[GSD] Active Phase: ${ACTIVE_PHASE} | Rules: (1) Agents read-only via MCP HTTP (FINCHER_MCP_URL), (2) DB credentials isolated in MCP container, (3) Envs use FINCHER_{SERVICE}_*, (4) Actions gated by DB-backed policies."
    }
  ]
}
EOF
