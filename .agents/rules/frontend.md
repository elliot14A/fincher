# Frontend Engineering & Code Quality Rules (`web/`)

These rules define strict architectural invariants and coding standards for all frontend development in `web/`:

---

## 1. Naming & Case Invariants

1. **Strict `camelCase` Across the Entire Codebase**:
   * All file names, directory names, and exported symbols must strictly use **`camelCase`** (e.g. `calendarGrid.tsx`, `calendarGrid.css.ts`, `queryKeys.ts`, `queryOptions.ts`, `holdOverrideModal.tsx`, `paginationControls.tsx`).
   * **Strictly NO kebab-case** (`calendar-grid.tsx` ❌) and **NO PascalCase files** (`CalendarGrid.tsx` ❌).
   * **The Only Exceptions**:
     - TanStack Router route parameter files dictated by framework conventions (`src/routes/$id.tsx`, `__root.tsx`, `index.tsx`).
     - Machine-generated Hey API output files in `src/lib/api/generated/` (`sdk.gen.ts`, `types.gen.ts`, `valibot.gen.ts`).

---

## 2. Deep Component Co-Location Rule

1. **Never Dump Components in Flat Folders**:
   * Every UI primitive, layout piece, and feature sub-component lives in its own dedicated directory alongside its styling and index barrel:
     ```text
     src/components/ui/button/
     ├── button.tsx
     ├── button.css.ts
     └── index.ts
     ```
2. **Feature-Local Sub-Components & Hooks**:
   * Feature sub-components live inside the feature's subdirectories:
     ```text
     src/features/calendar/grid/
     ├── calendarGrid.tsx
     ├── calendarGrid.css.ts
     └── index.ts
     ```
   * Feature-local hooks live inside that feature's `hooks/` directory (e.g. `src/features/lineage/hooks/useLineageLayout.ts` + `index.ts`).
   * Shared cross-feature hooks live in `src/lib/hooks/` + `index.ts`.
3. **Index Barrel Exports**:
   * Every directory must have an `index.ts` barrel acting as the sole public export surface. External files must import from the directory path, not deep-import internal files.

---

## 3. Query Infrastructure & Universal Pagination

1. **Separated Query Keys & Query Options**:
   * Every feature must maintain:
     * `queryKeys.ts`: Deterministic query key factory object (e.g. `titlesKeys`, `deliveriesKeys`).
     * `queryOptions.ts`: TanStack Query options consuming the auto-generated Hey API client (`src/lib/api/generated/`).
   * These files sit directly at the feature root and are re-exported through the feature's top-level `index.ts`.
2. **Unified Pagination Invariant**:
   * All paginated lists, grids, and catalogs must use `<PaginationControls />` from `#/components/ui/pagination`.
   * Standard pagination envelope: `items`, `total_items`, `page`, `limit`, `total_pages`, `has_next_page`, `has_prev_page`.
   * No ad-hoc pagination UI or divergent prop naming.

---

## 4. Zero Hardcoded CSS Invariant (Vanilla Extract)

1. **Strictly NO Hardcoded CSS**:
   * **Never** hardcode raw hex colors (`#ffffff`, `#000000`, `#17181c` ❌), arbitrary `rgba()` strings, hardcoded pixel dimensions, or non-tokenized border styles directly in `*.css.ts` or component JSX.
   * **Strictly NO inline `style={{ ... }}` in component JSX**: All styling must live in co-located `*.css.ts` using `@vanilla-extract/css` or `@vanilla-extract/recipes`.
   * **Always Consume Design Tokens**: Every color, margin, padding, border radius, font size, and line height must be referenced strictly from `src/styles/theme.css.ts` (`vars.color.*`, `vars.space.*`, `vars.radii.*`, `vars.fontSize.*`, `vars.lineHeight.*`, `fonts.*`).

---

## 5. Runtime & Dev Tooling Invariants

1. **Preact Compatibility Aliases**:
   * `vite.config.ts` must maintain `react` → `preact/compat`, `react-dom` → `preact/compat`, `react/jsx-runtime` → `preact/jsx-runtime`, and `react-dom/test-utils` → `preact/test-utils`.
2. **Native Preact Icons (`lucide-preact`)**:
   * Use `lucide-preact` directly to eliminate React SVG wrapper overhead.
3. **Dev API Proxying**:
   * `vite.config.ts` dev server must proxy `/api` requests to the local Go server (`http://localhost:8080`).
4. **Contract-First Code Generation**:
   * The frontend never hand-writes backend API types. Run `bun run generate:api` to synchronize TypeScript types, Valibot schemas, and the Hey API client in `src/lib/api/generated/` directly from `openapi/swagger.json`.
5. **Biome Linting & Formatting**:
   * All code must pass `bun run lint` (`biome check src`) and `bun run typecheck` (`tsc --noEmit`) with 0 warnings before committing.
