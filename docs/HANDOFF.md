# Handoff: Novexa Gateway — Unified API Key + Merged Model Picker

## Goal

Build a single-operator, self-hosted AI API gateway that exposes one OpenAI-compatible API key and routes requests across multiple upstream AI provider subscriptions. The VS Code/coding-CLI client sees a merged model picker with provider-qualified duplicates and the operator can check usage/cost across all providers.

## What has been decided

All domain decisions are captured in [CONTEXT.md](../CONTEXT.md):

- Single operator, many provider keys/plans.
- Canonical terms: **Model ID** (user-facing), **Provider Model ID** (upstream slug), **Alias**, **Route**, **Fallback**, **Model Catalog**, **Static Model List**, **Provider Key**, **Usage**, **Cost Rate**.
- `/v1/models` queries provider catalogs; duplicates get provider-prefixed IDs (e.g. `groq/llama3-8b`); prefix is stripped on routing.
- Aliases are config-only shortcuts, not advertised in the model list.
- Routing is explicit (aliases → routes → fallbacks); no auto-provider-selection.
- Usage is token-centric with optional extra counters for non-token providers; costs estimated in USD.
- Cost rates come from public pricing APIs/lists plus per-request cost when available, with manual fallback.
- Dashboard MVP: models, usage, health, logs.

Implementation plan is in [docs/PLAN.md](PLAN.md) with six vertical slices.

## What has been implemented

### Slice 1: Domain cleanup — COMPLETE
- Renamed `internal/model` → `internal/apitypes`.
- Split overloaded `model` field:
  - `RouteConfig.Model` → `RouteConfig.ModelID`
  - `FallbackConfig.Model` → `FallbackConfig.ModelID`
  - `ResolvedRoute.Model` → `ResolvedRoute.ProviderModelID`
  - Added `ResolvedRoute.ModelID`
  - `usage.Record` and `database.UsageRecord` now store both `ModelID` and `ProviderModelID`.
- Updated `config/config.example.yaml` to use `model_id`.
- `go build ./...` passes.

### Slice 2: Provider interface expansion — COMPLETE
- Added to `Provider` interface:
  - `ListModels(ctx) ([]ModelInfo, error)`
  - `GetPricing(ctx) (map[string]PricingInfo, error)`
- Added domain types:
  - `UnitType` (`UnitToken`, `UnitRequest`, `UnitMinute`, `UnitCharacter`)
  - `ModelInfo`, `PricingInfo`
  - `ErrNotImplemented` sentinel.
- Implemented real `ListModels` and `GetPricing` for OpenAI provider.
- Added stub implementations for anthropic, deepseek, gemini, generic, groq, lmstudio, ollama, openrouter.
- Updated `cmd/gateway/main.go` to register all stub providers when enabled.
- `go build ./...` passes.

### Slice 3: Merged `/v1/models` — COMPLETE
- Added `internal/catalog` aggregator: merges provider `ListModels`, qualifies duplicate base IDs with `provider/` prefixes, falls back to static `providers.*.models` when listing fails.
- `/v1/models` uses the catalog; aliases are not advertised.
- Router strips a registered provider prefix before route lookup and rejects prefix/provider mismatches.
- Config: `ProviderConfig.Models` for Static Model List; example YAML updated for ollama/lmstudio.
- Tests cover merge, prefixing, static fallback, alias exclusion, and prefix routing.
- `go test ./...` and `go build ./...` pass.

## Remaining work

Slices 4–6 from [docs/PLAN.md](PLAN.md):

4. **Usage/cost enhancements** — extra usage counters, fetch pricing, per-request cost reconciliation, manual fallback.
5. **Dashboard API endpoints** — `/api/models`, `/api/usage`, `/api/health`, `/api/logs`.
6. **Documentation reconciliation** — rewrite README/docs to match actual capabilities.

## Important notes for the next agent

- Always use the vocabulary from [CONTEXT.md](../CONTEXT.md); challenge any re-introduction of the ambiguous term `model`.
- Router still has legacy auto-detect via `SupportsModel` when no route exists; CONTEXT says routing should be explicit-only — consider removing auto-detect in a later slice.
- Slice 4 should extend usage/cost without changing the catalog prefix rules.
- Provider prefixes in the model list are a serialization concern only; route config still uses bare Model IDs.

## Suggested skills for the next session

- `tdd` — for test-first implementation of slices 4–5.
- `diagnose` — if routing or model-list behavior does not match expectations.
- `review` — before merging any slice.
- `to-issues` — if the remaining slices need to be tracked as separate issues.

## Artifacts

- Domain glossary: [CONTEXT.md](../CONTEXT.md)
- Implementation plan: [docs/PLAN.md](PLAN.md)
- This handoff: [docs/HANDOFF.md](HANDOFF.md)
