# Novexa Gateway

> A single-user, self-hosted AI gateway that exposes one OpenAI-compatible API endpoint while routing requests to multiple AI providers transparently.

[![CI](https://github.com/novexa/gateway/actions/workflows/ci.yaml/badge.svg)](https://github.com/novexa/gateway/actions/workflows/ci.yaml)
[![Go Report Card](https://goreportcard.com/badge/github.com/novexa/gateway)](https://goreportcard.com/report/github.com/novexa/gateway)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](LICENSE)

## 🚀 One-Liner

```bash
docker run -d -p 8080:8080 -e NOVEXA_API_KEY=my-key -e OPENAI_API_KEY=sk-... novexa/gateway:latest
```

## ✨ Features

- **Single API Endpoint** — Connect any OpenAI-compatible client (VSCode, Claude Code, Open WebUI, Continue, Aider, custom apps) to one endpoint with one API key
- **9 Providers** — OpenAI, Anthropic, Gemini, DeepSeek, OpenRouter, Groq, Ollama, LM Studio, and any OpenAI-compatible endpoint
- **Streaming** — Full SSE streaming support, normalized to OpenAI format
- **Model Routing** — Configurable model→provider mappings with aliases
- **Fallback Chains** — If a provider fails, automatically try the next
- **Usage Tracking** — Track tokens, costs, and latency per request
- **Cost Estimation** — See how much each request costs across providers
- **Health Monitoring** — Automatic provider health checks
- **Dashboard API** — Usage stats, costs, logs, and provider health
- **Docker Ready** — Single container, <20MB image
- **Free Cloud Deploy** — Deploy to Railway, Fly.io, or Render for free

## 📋 Quick Start

### Prerequisites

- Docker installed
- At least one AI provider API key

### Run the Gateway

```bash
# Create a .env file
cat > .env << EOF
NOVEXA_API_KEY=your-secret-gateway-key
OPENAI_API_KEY=sk-your-openai-key
EOF

# Run with Docker
docker run -d \
  --name novexa-gateway \
  -p 8080:8080 \
  --env-file .env \
  -v novexa-data:/app/data \
  novexa/gateway:latest
```

### Test It

```bash
curl http://localhost:8080/v1/chat/completions \
  -H "Authorization: Bearer your-secret-gateway-key" \
  -H "Content-Type: application/json" \
  -d '{
    "model": "gpt-4o",
    "messages": [{"role": "user", "content": "Hello!"}]
  }'
```

## 🔧 Configuration

### Minimal Setup

Only provider API keys are required. Everything else has sensible defaults.

```bash
export NOVEXA_API_KEY=your-secret-gateway-key
export OPENAI_API_KEY=sk-your-openai-key
```

### Advanced Configuration

Create a `config.yaml` file for custom routes, aliases, fallbacks, and more:

```yaml
routes:
  "gpt-4o":
    provider: openai
  "claude-sonnet-4-20250514":
    provider: anthropic

aliases:
  "fast": "gpt-4o-mini"
  "smart": "gpt-4o"
  "coding": "deepseek-chat"

fallbacks:
  "deepseek-chat":
    - provider: deepseek
    - provider: openrouter
      model: "deepseek/deepseek-chat"
```

See [Configuration Reference](docs/configuration.md) for full details.

## 🏗️ Architecture

```
Client → API Key Check → Rate Limit → Validate → Route → Provider Adapter → Normalize → Response
```

The gateway is a middleware pipeline built with Go and Fiber. Each stage is independent and testable.

- **Auth**: Single API key via `NOVEXA_API_KEY` env var
- **Rate Limiter**: Global and per-provider limits
- **Router**: Model→provider mapping with alias resolution and fallback chains
- **Provider Adapters**: Each provider implements a common interface
- **Usage Tracker**: Async writes to SQLite for zero-latency tracking

See [Architecture](docs/architecture.md) for details.

## 📚 Documentation

- [Quick Start Guide](docs/quickstart.md) — Get running in 5 minutes
- [Configuration Reference](docs/configuration.md) — Full configuration options
- [Provider Setup](docs/providers.md) — Detailed provider configuration
- [API Reference](docs/api.md) — Complete API documentation
- [Deployment Guide](docs/deployment.md) — Deploy to Railway, Fly.io, Render
- [Contributing](docs/contributing.md) — How to contribute

## 🎯 Supported Providers

| Provider | Format | Streaming | Cost |
|----------|--------|-----------|------|
| OpenAI | OpenAI | ✅ | $0.15-$10/1M tokens |
| Anthropic | Messages API | ✅ | $0.80-$15/1M tokens |
| Google Gemini | GenerateContent | ✅ | $0.075-$5/1M tokens |
| DeepSeek | OpenAI | ✅ | $0.14-$2.19/1M tokens |
| OpenRouter | OpenAI | ✅ | Varies |
| Groq | OpenAI | ✅ | $0.05-$0.79/1M tokens |
| Ollama | OpenAI | ✅ | Free (local) |
| LM Studio | OpenAI | ✅ | Free (local) |
| Generic | OpenAI | ✅ | Varies |

## 🚢 Deployment

### Docker (Local)

```bash
docker run -d -p 8080:8080 --env-file .env -v novexa-data:/app/data novexa/gateway:latest
```

### Railway

```bash
railway login
railway init
railway variables set NOVEXA_API_KEY=your-key
railway up
```

### Fly.io

```bash
fly launch
fly secrets set NOVEXA_API_KEY=your-key
fly deploy
```

### Render

Connect your GitHub repo → Render auto-detects → Add env vars → Deploy

See [Deployment Guide](docs/deployment.md) for detailed instructions.

## 📊 Dashboard API

All endpoints use the same API key:

```bash
# Provider health
curl http://localhost:8080/api/health -H "Authorization: Bearer your-key"

# Usage stats
curl http://localhost:8080/api/usage -H "Authorization: Bearer your-key"

# Cost breakdown
curl http://localhost:8080/api/usage/costs -H "Authorization: Bearer your-key"
```

## 🛠️ Development

```bash
# Clone
git clone https://github.com/novexa/gateway.git
cd gateway

# Build
make build

# Test
make test

# Run
export NOVEXA_API_KEY=test-key
export OPENAI_API_KEY=sk-test
./bin/gateway
```

## 📄 License

MIT License — see [LICENSE](LICENSE) for details.

## 🤝 Contributing

Contributions are welcome! See [Contributing Guide](docs/contributing.md) for details.

## 🙏 Acknowledgments

- [Fiber](https://gofiber.io/) — HTTP framework
- [Viper](https://github.com/spf13/viper) — Configuration
- [Zap](https://go.uber.org/zap) — Logging
- [GORM](https://gorm.io/) — ORM
- [SQLite](https://www.sqlite.org/) — Database
