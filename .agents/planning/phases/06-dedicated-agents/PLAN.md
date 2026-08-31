# Phase 06 Plan: Dedicated Agents & Hackathon Workflow Engine

> **Provenance**: Validated by reference implementation in `/tmp/opencode/fincher_timemodel/` (46 tests, 0 failures), documented in `.agents/reviews/15-dedicated-agents/20260830T102605Z-timemodel-validation.md` (review session 15 corresponding to feature phase 06) and `hackathon_domain_and_workflow_spec.md`.  
> **Core Principle**: Actionable, event-driven, closed-loop. Real executed dispatches (no passive drafts), single front door ingestion, compressed-time elapsing, and title self-healing.

---

## Work Units

### Unit 1: Enrich Executor Dispatches & Add Event Constants
* **Target Files**:
  * `pkg/domain/models/event.go`
  * `internal/agent/runner.go`
  * `internal/agent/runner_test.go`
* **Details**:
  1. Add missing event type constants to `pkg/domain/models/event.go`:
     ```go
     TypeVendorAssigned       = "fincher.vendor.assigned"
     TypeVendorEmailed        = "fincher.vendor.emailed"
     TypeStakeholdersNotified = "fincher.stakeholders.notified"
     TypeSocialPosted         = "fincher.social.posted"
     TypeDeliveryHeld         = "fincher.delivery.held"
     TypeDeliveryReleased     = "fincher.delivery.released"
     ```
  2. In `internal/agent/runner.go:111`, replace string literal `"fincher.vendor.assigned"` with `models.TypeVendorAssigned`.
  3. Update `(*Event).Classify()` in `event.go` to explicitly categorize `TypeVendorAssigned`, `TypeVendorEmailed`, `TypeStakeholdersNotified`, `TypeSocialPosted`, `TypeDeliveryHeld`, `TypeDeliveryReleased` as `CategoryRoutineOutcome` (ensuring downstream dispatches never trigger an incident re-trigger loop).
  4. Rescope `RunActionPlan` in `internal/agent/runner.go` (which already implements the switch for communication actions) to enrich the existing stub handling:
     * `ActionEmailVendor`: populate `action.Payload` with receipt `{ "dispatch_id": "msg-email-" + uuid[:8], "status": "DELIVERED", "dispatched_at": now }` and append `TypeVendorEmailed` event.
     * `ActionNotifyStakeholders`: populate `action.Payload` with receipt `{ "dispatch_id": "msg-slack-" + uuid[:8], "channel": "#ops-war-room", "status": "DELIVERED", "dispatched_at": now }` and append `TypeStakeholdersNotified` event.
     * `ActionPostSocialUpdate`: populate `action.Payload` with receipt `{ "post_id": "post-x-" + uuid[:8], "platform": "x/twitter", "status": "PUBLISHED", "dispatched_at": now }` and append `TypeSocialPosted` event.
  5. Ensure all communication actions are preserved in `ExecutedActions` with their receipt payloads.
* **Verification**:
  * `go test -v ./internal/agent/... -run TestRunActionPlan` verifies dispatches contain simulated receipt IDs and emit `TypeVendorEmailed`, `TypeStakeholdersNotified`, `TypeSocialPosted` to ClickHouse.

---

### Unit 2: Closed-Loop Resolution Workflow & Title Self-Healing
* **Target Files**:
  * `internal/agent/graph/resolution.go`
  * `internal/agent/graph/resolution_test.go`
  * `internal/turso/packages/update.go`
  * `internal/turso/deliveries/update.go`
  * `internal/turso/titles/update.go`
* **Details**:
  1. Implement `ExecuteResolution(ctx context.Context, deps ResolutionDeps, input ResolutionInput) (*models.Run, error)`:
     * Step 1: Query target `MediaPackage` via `tursopackages.Get`.
       - Update package status to `models.PackageStatusValid` via `tursopackages.Update`.
     * Step 2: Multi-Package Territory Resolution (1 Delivery : N Packages):
       - Query all deliveries for the title via `tursodeliveries.List`.
       - For each delivery currently on `models.DeliveryStatusHold`:
         - Query dependent packages matching this territory/delivery (or required package IDs).
         - If ALL required packages for that delivery are `models.PackageStatusValid` $\rightarrow$ update delivery status to `models.DeliveryStatusReadyToShip` via `tursodeliveries.Update`.
     * Step 3: Title Self-Healing & Status Resolution:
       - Query all deliveries and packages for `titleSlug`.
       - Evaluate title status rules:
         * If ANY delivery is `models.DeliveryStatusHold` or any package is `models.PackageStatusInvalidated` / `models.PackageStatusReQCPending` $\rightarrow$ Title status remains `models.StatusHold`.
         * If NO delivery is `models.DeliveryStatusHold` and all packages are `models.PackageStatusValid` $\rightarrow$ update `Title.OverallStatus` to `models.StatusOnTrack` ("Ready to Ship").
         * If NO delivery is `models.DeliveryStatusHold` but some packages remain in initial `models.PackageStatusPending` (in-flight preparation without active hold) $\rightarrow$ update `Title.OverallStatus` to `models.StatusProcessing` (or `models.StatusOnTrack` per launch schedule).
       - Persist via `tursotitles.Update(ctx, client, titleID, &models.UpdateTitleInput{OverallStatus: &targetStatus})`.
     * Step 4: Downstream Event Emission:
       - Emit `TypeDeliveryReleased` (`fincher.delivery.released`) CloudEvent to ClickHouse.
     * Step 5: Persistence:
       - Create `Run` (status `COMPLETED`), `Step` (`resolution_executor`), and `WfResult` in Turso for live SSE streaming.
* **Verification**:
  * `go test -v ./internal/agent/graph/... -run TestExecuteResolution` asserting that a clean QC return unholds the package, unholds dependent deliveries, and transitions Title `OverallStatus` from `HOLD` to `ON_TRACK`.

---

### Unit 3: Single Front Door Ingestion Auto-Router
* **Target Files**:
  * `internal/api/events/create.go`
  * `internal/api/events/events_test.go`
  * `internal/api/runs/create.go`
* **Details**:
  1. Reconcile dual-endpoint architecture:
     * `POST /api/events` is the **primary front door** for all organic lifecycle events (telemetry, anomalies, clean QC, allocation requests).
     * `POST /api/runs` remains as the explicit operator manual trigger (`CategoryOperatorForced`), sharing the same background execution functions.
  2. In `internal/api/events/create.go`:
     * Validate and write incoming CloudEvents batch to ClickHouse via `chEvents.InsertBatch`.
     * For each event in the batch, inspect routing criteria:
       - Call `e.Classify()`:
         * `models.CategoryIncident`: dispatch `graph.ExecuteIncident` asynchronously.
         * `models.CategoryAllocation`: dispatch `graph.ExecuteAllocation` asynchronously.
         * `models.CategoryRoutineOutcome`:
           - Check exact type and payload: if `e.Type == models.TypeQCInspectionCompleted` and `e.Data["status"] == "PASSED"`, dispatch `graph.ExecuteResolution` asynchronously.
           - All other routine outcomes (`TypeDeliveryHeld`, `TypeDeliveryReleased`, `TypeVendorAssigned`, `TypeVendorEmailed`, `TypeStakeholdersNotified`, `TypeSocialPosted`) are logged-only to ClickHouse without spawning runs (re-trigger guard).
  3. Response payload: return `201 Created` with ingested count and list of spawned run IDs (if any).
* **Verification**:
  * Integration test in `events_test.go`:
    - Posting an anomaly event spawns an incident run.
    - Posting a clean QC event (`status: "PASSED"`) spawns a resolution run.
    - Posting routine telemetry/dispatches logs to ClickHouse without spawning runs.

---

### Unit 4: Compressed-Time Scheduler & Title Projection Subsystem
* **Target Files**:
  * `internal/config/config.go`
  * `internal/agent/scheduler/task.go`
  * `internal/agent/scheduler/scheduler.go`
  * `internal/agent/scheduler/scheduler_test.go`
  * `internal/agent/tools/projection.go`
  * `internal/agent/tools/projection_test.go`
* **Details**:
  1. Configuration:
     * Add `FINCHER_TIME_SCALE` (type `time.Duration`, default `time.Second`, kong tag `env:"FINCHER_TIME_SCALE",default="1s"`).
     * 1s real wall-clock = 1 domain hour (a 12h turnaround runs in 12 real seconds for demo watchability).
  2. Sub-Task 4A: In-Memory Task Registry & Lifecycle:
     * `Task` struct tracking: `ID`, `Kind` (`"package"` | `"master"`), `TargetID`, `VendorID`, `TurnaroundHours` (domain hours for story/history), `ActualHours` (honest duration), `FinishReal` (wall-clock deadline), `Status` (`SCHEDULED`, `RUNNING`, `COMPLETED`, `CANCELLED`).
     * `CancelTasksForPackage(packageID)`: cancels in-flight timers if a package experiences a re-incident or reassignment mid-repair (idempotency guard preventing double-completions).
  3. Sub-Task 4B: Who Calls `ScheduleTask` and When:
     * **Call Site 1**: In `internal/agent/runner.go` when `ActionReassignVendor` executes:
       - Looks up `newVendor.TurnaroundHours`.
       - Calls `scheduler.ScheduleTask("package", pkgID, vendorID, turnaroundHours, onComplete)`.
     * **Call Site 2**: In `internal/agent/graph/allocation.go` when a vendor is assigned to a newly required package.
     * **Call Site 3**: When an upstream master bump is conformed (`Kind: "master"`).
  4. Sub-Task 4C: Elapsing & Front-Door Event Emission:
     * Computes real duration: `turnaroundHours * timeScale`.
     * Schedules `time.AfterFunc(realDuration, ...)`:
       - On completion: constructs a clean QC CloudEvent (`models.TypeQCInspectionCompleted`, `status: "PASSED"`, `data.turnaround_hours: turnaroundHours`, `data.package_id: targetID`, `data.vendor_id: vendorID`).
       - Feeds this event directly through the front-door ingestion router (`POST /api/events` / `events.IngestAndRoute`), which organically fires `ExecuteResolution` (Unit 2).
  5. Sub-Task 4D: Sequential-Edge Scheduling (§3 bug fix from Python validation):
     * When a package task depends on a master reconform task, do NOT use a coarse polling tick.
     * Schedule the dependent package timer *directly from the master task's completion callback*, anchoring start time to the master's finish time (guaranteeing 6h master + 12h dub sums to 18h, rather than inflating to 24.5h).
  6. Sub-Task 4E: Title Projection Tool (`internal/agent/tools/projection.go`):
     * Implement `GetTitleReadyProjection(ctx, client, titleSlug) (*TitleProjection, error)`:
       - Computes `hours_until_premiere = time.Until(title.PremiereDate).Hours()`.
       - Computes `critical_remaining_hours` (summing sequential master+dub dependencies and maxing across parallel sibling packages).
       - Computes `buffer_hours = hours_until_premiere - critical_remaining_hours`.
       - Classifies `risk_band`: `buffer < 0 => "BREACH"`, `< 6 => "TIGHT"`, `< 24 => "WATCH"`, else `"SAFE"`.
       - Returns structured payload. AI prompt instructs Gemini to tool-call this function and narrate pre-computed numbers without doing raw date arithmetic.
* **Verification**:
  * Unit test in `scheduler_test.go` verifying task lifecycle, cancellation on re-incident, sequential duration addition (6h + 12h = 18h), and front-door resolution event emission.
  * Unit test in `projection_test.go` verifying projection math and risk bands.

---

### Unit 5: Non-Dominated Vendors & Realistic Demo Seeder
* **Target Files**:
  * `cmd/seed/main.go`
  * `data/seed/events.json`
* **Details**:
  1. Data Directory:
     * Create `data/seed/` and copy the validated 195-event dataset from `/tmp/opencode/fincher_timemodel/events.json`.
  2. SQLite Vendor Seeding:
     * Note on accuracy: SQLite `Vendor` schema strictly stores rate cards (`hourly_rate_usd`, `turnaround_hours`, `specialty`). Historical accuracy is NOT a static SQLite field; it is computed by ClickHouse from historical QC events.
     * Seed non-dominated vendor candidates in Turso:
       - *Deluxe Audio QC*: $200/hr, 12h turnaround (Rush option)
       - *Iyuno QC*: $70/hr, 36h turnaround (Economy option)
       - *Testronic Audio QC*: $120/hr, 24h turnaround (Borderline/failing partner)
       - *Pixelogic Subtitles*: $80/hr, 8h turnaround
  3. SQLite Title & Deliveries Seeding:
     * Seed Title: *Avatar: Fire & Ash* (`slug: avatar-fire-ash`, `overall_status: ON_TRACK`, `premiere_date = now + 72h`).
     * Seed Master Cut V01.
     * Seed 6 market storefront deliveries: US, DE, FR, JP, ES, BR.
     * Seed 11 localized media packages (Audio Dub + Subtitle per territory) with dependencies linking each delivery to its required audio and subtitle packages.
  4. ClickHouse Historical Ingestion:
     * Parse `data/seed/events.json` and bulk insert into ClickHouse `fincher.events` table.
     * This establishes the historical track records (*Deluxe* 99% accuracy, *Iyuno* 93%, *Testronic* 85% with prior sync drift incidents) for the Historian and Vendor Selector agents to query via MCP.
* **Verification**:
  * `go run ./cmd/seed` populates both Turso and ClickHouse cleanly with zero constraint violations.

---

### Unit 6: Frontend UI Revamp (Showcase Landing Page, Packages Route, Live Runs Console)
* **Target Files**:
  * `web/src/routes/index.tsx` (Product Showcase Landing Page)
  * `web/src/styles/routes/index.css.ts`
  * `web/src/features/layout/shell/appShell.tsx` (Conditional shell rendering for `/` vs console)
  * `web/src/routes/packages.tsx` (new dedicated route for existing Packages catalog)
  * `web/src/features/layout/sidebar/navigationSidebar.tsx` (navigation bar update)
  * `web/src/features/runs/` (queryKeys, queryOptions, components, hooks)
  * `web/src/routes/runs.tsx` (overhauled to genuine live Runs Operations Console)
* **Details**:
  1. Sub-Task 6A: Product Showcase Landing Page (`/`):
     * Rebuild `web/src/routes/index.tsx` as a full-width showcase explaining Fincher's value proposition:
       - Header: Fincher Logo, Live System Status Pill (`ClickHouse MCP Connected • Autonomous Engine Active`), and **"Open Fincher"** CTA at top right (navigating to `/titles`).
       - Hero section: Tagline (*"The Autonomous Media Supply Chain & Quality Control Engine"*), problem description (preventing global release-day disasters across streaming storefronts).
       - Interactive Architecture Strip: Visualizes CloudEvents $\rightarrow$ ClickHouse $\rightarrow$ Gemini Judges $\rightarrow$ Policy Verification $\rightarrow$ Dispatches & Self-Healing.
       - 4 Core Capability Cards: (1) Sync Drift & QC Artifact Detection, (2) Multi-Market Synchronized Releases, (3) Autonomous Remediation & Policy Gating, (4) Closed-Loop Title Self-Healing.
       - Tech Stack Badges: ClickHouse MCP, Google Gemini 2.5, Turso / SQLite, Preact, TanStack.
     * Update `web/src/features/layout/shell/appShell.tsx` to conditionally hide the navigation sidebar when on `/` (full-width showcase view), and display it on all console routes (`/titles`, `/deliveries`, `/packages`, `/vendors`, `/runs`).
  2. Sub-Task 6B: Dedicated Packages Catalog (`/packages`):
     * Create `web/src/routes/packages.tsx` and move the existing Media Packages table from `runs.tsx` into it.
     * Add `Packages` link with an icon to `navigationSidebar.tsx`.
  3. Sub-Task 6C: Live Runs Operations Console & Hero Simulator (`/runs`):
     * `queryKeys.ts` & `queryOptions.ts` consuming `getRuns` and `getRunsById`.
     * `useSSEStream.ts` hook subscribing to `GET /api/runs/:id/stream`.
     * `runRow.tsx` / `runRow.css.ts`: displays Run ID, title slug, trigger (`incident`, `allocation`, `resolution`), status badge, and stage counter.
     * `runDetailDrawer.tsx`: renders the 4-stage progression:
       - Stage 1: Triage Judge (severity, actionable verdict, rationale)
       - Stage 2: Context & Blast Radius (affected territory deliveries)
       - Stage 3: Remediation / Selection Loop (action plan attempts, verifier verdicts)
       - Stage 4: Execution & Outbound Dispatches Card (displays sent email, Slack alert, and Twitter post with receipt badges: `DELIVERED`, `PUBLISHED`).
     * Hero Simulator Controls:
       - Button 1 ("Simulate Sync Drift on German Dub"): Posts `fincher.audio.sync_drift` $\rightarrow$ Title flips to `HOLD` $\rightarrow$ Deluxe reassigned $\rightarrow$ dispatches fire.
       - Button 2 ("Simulate Deluxe Clean Return"): Posts `fincher.qc.completed` (`status: "PASSED"`) $\rightarrow$ Title self-heals back to `ON_TRACK`.
* **Verification**:
  * `bun run typecheck` passes with 0 errors.
  * `bun run lint` passes with 0 errors.
  * Browser test: Landing page renders cleanly on `/` with "Open Fincher" button; `/packages` lists packages; `/runs` streams stages and executes hero buttons live.

---

## Definition of Done
1. All 6 units implemented and verified against unit/integration tests.
2. `go test -race ./...` passes 100%.
3. Frontend `bun run typecheck` and `bun run lint` pass cleanly.
4. End-to-end rehearsal from `POST /api/events` through triage $\rightarrow$ hold $\rightarrow$ reassign $\rightarrow$ dispatches $\rightarrow$ clean return $\rightarrow$ title self-healing displays live on the browser console.
