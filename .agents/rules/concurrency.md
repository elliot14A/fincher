# Concurrency & Async Execution Rules

These rules govern concurrent workflows and async operations:

1. **Parallel Multi-Agent Execution**:
   - The Historian and Dependency sub-agents must execute concurrently via Go goroutines / errgroup.
   - Bound goroutine execution and ensure all routines respect parent `context.Context` cancellation and timeouts.

2. **Database Concurrency**:
   - SQLite connections must use WAL mode (`_journal_mode=WAL`) and single-writer concurrency (`SetMaxOpenConns(1)`).
   - Use transactions for multi-step mutations to maintain atomicity.

3. **Real-time SSE Streaming**:
   - Event broadcaster channels must be buffered and non-blocking to prevent slow HTTP clients from stalling the core pipeline.
