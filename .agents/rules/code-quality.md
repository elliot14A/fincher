# Code Quality & Standards

These rules enforce software craftsmanship standards across Fincher:

1. **Strict Idiomatic Go**:
   - Keep abstractions simple, flat, and standard library driven.
   - Always wrap errors with context: `fmt.Errorf("evaluating policy %s: %w", policyID, err)`.
   - Propagate `context.Context` across all handlers, agent loops, MCP calls, and DB transactions.

2. **Structured Typed Outputs**:
   - Enforce strict JSON schema contracts on all LLM responses via the Google GenAI / ADK Go SDK.

3. **No Placeholders in UI or Code**:
   - Use concrete, typed fields and realistic media production datasets.
