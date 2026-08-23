# Phase 04: UI Scaffolding & Initial Operations Console — Execution Plan

## 1. Scope
Initialize `web/` with Bun, build the core design system and layout shell, run Hey API codegen, and deliver interactive UI slices for Titles/Launch Calendar, Vendors, Deliveries, and Lineage DAG.

---

## 2. Tasks

### Task 1: Web Workspace & Tooling Initialization
* **Target Files**:
  - `web/package.json`
  - `web/tsconfig.json`
  - `web/vite.config.ts` (Preact preset, Vanilla Extract, TanStack Router plugin, Preact compat aliases, `/api` proxy)
  - `web/biome.json`
  - `web/openapi-ts.config.ts`
  - `web/public/favicon.svg`
  - `web/src/vite-env.d.ts`
* **Verification**: `bun install` succeeds and `bun run dev --help` works.

### Task 2: Codegen & Design System Tokens
* **Target Files**:
  - Run `bun run generate:api` against `../openapi/swagger.json` to produce `web/src/lib/api/generated/` (`sdk.gen.ts`, `types.gen.ts`, `valibot.gen.ts`).
  - `web/src/lib/api/client.ts` (configured Hey API fetch client with `/api` relative base).
  - `web/src/styles/theme.css.ts` & `web/src/styles/tokens.ts` (dark operations theme contract & variables).
  - `web/src/app.css.ts` (global reset & typography).
* **Verification**: Generated SDK exports all 20+ types and handlers; Vanilla Extract theme compiles without errors.

### Task 3: Shared UI Primitives & Providers
* **Target Files**:
  - `web/src/lib/queryClient.ts` & `web/src/lib/dbClient.ts`
  - `web/src/components/ui/button/` (`button.tsx`, `button.css.ts`, `index.ts`)
  - `web/src/components/ui/badge/` (`badge.tsx`, `badge.css.ts`, `index.ts`)
  - `web/src/components/ui/modal/` (`modal.tsx`, `modal.css.ts`, `index.ts`)
  - `web/src/components/ui/table/` (`dataTable.tsx`, `dataTable.css.ts`, `index.ts`)
  - `web/src/components/ui/input/` (`input.tsx`, `input.css.ts`, `index.ts`)
  - `web/src/components/ui/drawer/` (`drawer.tsx`, `drawer.css.ts`, `index.ts`)
  - `web/src/components/feedback/skeletonLoader/`, `errorBoundary/`, `emptyState/`
  - `web/src/main.tsx` (Preact entry point with QueryClientProvider & TanStack Router Provider).
* **Verification**: UI primitives typecheck and render in isolation.

### Task 4: Layout Shell & Navigation
* **Target Files**:
  - `web/src/features/layout/shell/` (`appShell.tsx`, `appShell.css.ts`, `index.ts`)
  - `web/src/features/layout/sidebar/` (`navigationSidebar.tsx`, `navigationSidebar.css.ts`, `index.ts`)
  - `web/src/features/layout/topbar/` (`operationsTopbar.tsx`, `operationsTopbar.css.ts`, `index.ts`)
  - `web/src/routes/__root.tsx`
* **Verification**: Root route renders top navigation, live status indicator, and sidebar links.

### Task 5: Launch Calendar Feature Slice
* **Target Files**:
  - `web/src/db/collections/titlesCollection.ts`
  - `web/src/features/calendar/grid/` (`calendarGrid.tsx`, `calendarGrid.css.ts`, `index.ts`)
  - `web/src/features/calendar/card/` (`titleCard.tsx`, `titleCard.css.ts`, `index.ts`)
  - `web/src/features/calendar/queryKeys.ts`, `queryOptions.ts`, `index.ts`
  - `web/src/routes/index.tsx` (Launch Calendar view)
  - `web/src/routes/titles/$id.tsx` (Title detail & master cut cuts view)
* **Verification**: Displays upcoming titles (*Eclipse*, *Atlas*, etc.) with live countdowns and status indicators.

### Task 6: Vendor Directory Feature Slice
* **Target Files**:
  - `web/src/db/collections/vendorsCollection.ts`
  - `web/src/features/vendors/card/` (`vendorCard.tsx`, `vendorCard.css.ts`, `index.ts`)
  - `web/src/features/vendors/chart/` (`vendorQualityChart.tsx`, `vendorQualityChart.css.ts`, `index.ts`)
  - `web/src/features/vendors/queryKeys.ts`, `queryOptions.ts`, `index.ts`
  - `web/src/routes/vendors.tsx`
* **Verification**: Displays vendor cards with quality scorecards, error rates, and specialties.

### Task 7: Territory Delivery Matrix & Lineage DAG Feature Slices
* **Target Files**:
  - `web/src/db/collections/deliveriesCollection.ts`, `packagesCollection.ts`
  - `web/src/features/deliveries/table/` (`deliveryTable.tsx`, `deliveryTable.css.ts`, `index.ts`)
  - `web/src/features/deliveries/statusPill/` (`deliveryStatusPill.tsx`, `deliveryStatusPill.css.ts`, `index.ts`)
  - `web/src/features/deliveries/overrideModal/` (`holdOverrideModal.tsx`, `holdOverrideModal.css.ts`, `index.ts`)
  - `web/src/features/deliveries/queryKeys.ts`, `queryOptions.ts`, `index.ts`
  - `web/src/routes/deliveries.tsx`
  - `web/src/features/lineage/canvas/` (`lineageCanvas.tsx`, `lineageCanvas.css.ts`, `index.ts`)
  - `web/src/features/lineage/node/` (`lineageNode.tsx`, `lineageNode.css.ts`, `index.ts`)
  - `web/src/features/lineage/hooks/useLineageLayout.ts`
  - `web/src/features/lineage/queryKeys.ts`, `queryOptions.ts`, `index.ts`
  - `web/src/routes/lineage/$id.tsx`
* **Verification**: Delivery matrix displays 40+ country statuses; Lineage DAG renders interactive node graph for a title.

### Task 8: End-to-End Build & Static Verification
* **Target Commands**:
  - `bun run typecheck` (`tsc --noEmit`)
  - `bun run lint` (`biome check src`)
  - `bun run build` (outputs to `web/dist`)
* **Verification**: 0 type errors, 0 lint errors, clean production bundle generated.
