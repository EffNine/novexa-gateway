# Novexa Gateway Architecture

## Overview

Novexa Gateway is a single-user, self-hosted AI gateway that exposes an OpenAI-compatible API while routing requests to multiple AI providers transparently.

## System Architecture

```
Client (VSCode, Claude Code, mobile apps, etc.)
    │
    ▼
┌─────────────────────────────────────────────┐
│         Novexa Gateway (Go/Fiber)            │
│                                              │
│  API Key Check → Rate Limit → Validate       │
│       → Route → Provider Adapter             │
│       → Normalize → Usage Track → Response   │
│                                              │
│  ┌─────────────────────────────────────────┐ │
│  │  Providers: OpenAI, Anthropic, Gemini,  │ │
│  │  DeepSeek, OpenRouter, Groq, Ollama,    │ │
│  │  LM Studio, Generic                     │ │
│  └─────────────────────────────────────────┘ │
│                                              │
│  SQLite (usage, costs, health, logs)         │
└─────────────────────────────────────────────┘
```

## Key Design Decisions

### Single-User Model
- Single API key via `NOVEXA_API_KEY` environment variable
- No user management or authentication complexity
- Operator is the only user

### Provider Abstraction
- All providers implement a common `Provider` interface
- Each provider is a separate package with translation logic
- Adding a new provider requires only implementing the interface
- No provider-specific code outside adapter packages

### Configuration-Driven Routing
- All model→provider mappings defined in YAML config
- Environment variable overrides for all settings
- Zero hardcoded routing logic
- Sensible defaults for minimal configuration

### Streaming via SSE Proxy
- Real-time streaming without buffering full responses
- Provider-specific SSE formats normalized to OpenAI format
- Goroutine per stream with context cancellation

### SQLite-First Storage
- SQLite for usage tracking, costs, health, and logs
- Perfect for single-user local deployment
- PostgreSQL available as opt-in for advanced users
- WAL mode for concurrent read/write

### Async Usage Tracking
- Usage records written directly to SQLite (no batching needed at <1K req/day)
- No latency impact on request path
- Cost estimation calculated on write

## Request Lifecycle

1. Client sends `POST /v1/chat/completions` with `Authorization: Bearer <key>`
2. API key middleware validates against `NOVEXA_API_KEY`
3. Rate limiter checks global and per-provider limits
4. Request validator checks payload structure and size
5. Router reads `model` field → resolves alias → resolves route → selects provider
6. If fallback chain configured, router prepares ordered provider list
7. Provider adapter translates OpenAI request → provider-native format (if needed)
8. Provider adapter sends HTTP request to provider
9. Response is normalized back to OpenAI format
10. Usage tracker records tokens, latency, cost
11. Normalized response returned to client

## Streaming Flow

Same as above through step 7, then:
8. Provider adapter opens streaming connection
9. Gateway reads SSE chunks from provider in a goroutine
10. Each chunk is normalized to OpenAI `data: {...}\n\n` format
11. Chunks are forwarded to client via Fiber's streaming response
12. On `data: [DONE]`, connection closes
13. Final token counts are estimated and tracked

## Security Model

- Provider API keys stored in environment variables only
- Keys never logged (Zap field filters)
- Keys redacted in `/api/config` endpoint
- Request/response bodies not logged by default
- CORS configurable for cross-origin clients
- Payload size limits enforced

## Deployment Model

- Single binary, single process
- SQLite database file (default: `./data/novexa.db`)
- No external dependencies required
- Docker image <20MB
- Deployable on Railway, Fly.io, Render, or local machine
