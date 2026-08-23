# Phase 04: UI Scaffolding & Initial Operations Console — Context

## 1. Goal
Scaffold the Preact + Vite frontend application in `web/` using Bun, Vanilla Extract (`.css.ts`), TanStack Router, TanStack DB, Hey API codegen, and implement initial operations console views for the completed backend entities (Launch Calendar, Title detail, Vendor scorecards, Territory Delivery Matrix, and Lineage DAG).

## 2. Invariants & Rules
1. **Strict camelCase Naming**: All files and directories in `web/src/` must be `camelCase` (except TanStack Router `$id.tsx`, `__root.tsx`, `index.tsx` and generated SDK in `src/lib/api/generated/`).
2. **Deep Component Co-location**: No flat dumping grounds. Every component primitive or feature sub-component has its own folder with `*.tsx`, `*.css.ts`, and `index.ts`.
3. **Zero-Runtime Styling**: All styles authored in `*.css.ts` using `@vanilla-extract/css` and `@vanilla-extract/recipes` bound to `src/styles/theme.css.ts`.
4. **Preact Compatibility**: `vite.config.ts` must alias `react` → `preact/compat`, etc.
5. **Contract-First Codegen**: SDK generated directly from backend `openapi/swagger.json` via `@hey-api/openapi-ts`.

## 3. Reference
* Architecture specification in `.agents/reviews/06-frontend-architecture/20260822T081500Z.md`.
