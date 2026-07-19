# Quick Start Guide

Get Novexa Gateway running in 5 minutes.

## Prerequisites

- Docker installed
- At least one AI provider API key (OpenAI, Anthropic, etc.)

## Option 1: Docker (Recommended)

### 1. Create a `.env` file

```bash
cat > .env << EOF
NOVEXA_API_KEY=your-secret-gateway-key
OPENAI_API_KEY=sk-your-openai-key
EOF
```

### 2. Run the gateway

```bash
docker run -d \
  --name novexa-gateway \
  -p 8080:8080 \
  --env-file .env \
  -v novexa-data:/app/data \
  novexa/gateway:latest
```

### 3. Test it

```bash
curl http://localhost:8080/v1/chat/completions \
  -H "Authorization: Bearer your-secret-gateway-key" \
  -H "Content-Type: application/json" \
  -d '{
    "model": "gpt-4o",
    "messages": [{"role": "user", "content": "Hello!"}]
  }'
```

## Option 2: Docker Compose

### 1. Create `docker-compose.yaml`

```yaml
version: '3.8'

services:
  gateway:
    image: novexa/gateway:latest
    ports:
      - "8080:8080"
    environment:
      - NOVEXA_API_KEY=your-secret-gateway-key
      - OPENAI_API_KEY=sk-your-openai-key
    volumes:
      - novexa-data:/app/data
    restart: unless-stopped

volumes:
  novexa-data:
```

### 2. Start the gateway

```bash
docker-compose up -d
```

## Option 3: Build from Source

### 1. Clone the repository

```bash
git clone https://github.com/novexa/gateway.git
cd gateway
```

### 2. Build the binary

```bash
make build
```

### 3. Run the gateway

```bash
export NOVEXA_API_KEY=your-secret-gateway-key
export OPENAI_API_KEY=sk-your-openai-key
./bin/gateway
```

## Configuration

### Minimal Configuration

Only provider API keys are required. Everything else has sensible defaults.

```bash
export NOVEXA_API_KEY=your-secret-gateway-key
export OPENAI_API_KEY=sk-your-openai-key
```

### Adding More Providers

```bash
export ANTHROPIC_API_KEY=sk-ant-your-key
export GEMINI_API_KEY=your-gemini-key
export DEEPSEEK_API_KEY=your-deepseek-key
```

### Custom Configuration

Create a `config.yaml` file:

```yaml
server:
  port: 8080

api_key: "${NOVEXA_API_KEY}"

providers:
  openai:
    api_key: "${OPENAI_API_KEY}"
  anthropic:
    api_key: "${ANTHROPIC_API_KEY}"

routes:
  "gpt-4o":
    provider: openai
  "claude-sonnet-4-20250514":
    provider: anthropic

aliases:
  "fast": "gpt-4o-mini"
  "smart": "gpt-4o"
```

Run with config file:

```bash
docker run -d \
  -p 8080:8080 \
  --env-file .env \
  -v $(pwd)/config.yaml:/app/config.yaml \
  -v novexa-data:/app/data \
  novexa/gateway:latest
```

## Using with Clients

### VS Code (Continue, Copilot alternatives)

Configure your OpenAI-compatible extension:

```json
{
  "apiBase": "http://localhost:8080/v1",
  "apiKey": "your-secret-gateway-key",
  "model": "gpt-4o"
}
```

### Claude Code

```bash
export ANTHROPIC_BASE_URL=http://localhost:8080/v1
export ANTHROPIC_API_KEY=your-secret-gateway-key
```

### Open WebUI

In Open WebUI settings:
- API Base URL: `http://localhost:8080/v1`
- API Key: `your-secret-gateway-key`

### cURL

```bash
curl http://localhost:8080/v1/chat/completions \
  -H "Authorization: Bearer your-secret-gateway-key" \
  -H "Content-Type: application/json" \
  -d '{
    "model": "gpt-4o",
    "messages": [{"role": "user", "content": "Hello!"}],
    "stream": true
  }'
```

## Checking Status

### Health Check

```bash
curl http://localhost:8080/health
```

### Provider Health

```bash
curl http://localhost:8080/api/health \
  -H "Authorization: Bearer your-secret-gateway-key"
```

### Usage Statistics

```bash
curl http://localhost:8080/api/usage \
  -H "Authorization: Bearer your-secret-gateway-key"
```

### Cost Breakdown

```bash
curl http://localhost:8080/api/usage/costs \
  -H "Authorization: Bearer your-secret-gateway-key"
```

## Next Steps

- [Configuration Reference](configuration.md) - Full configuration options
- [Provider Setup](providers.md) - Detailed provider configuration
- [Deployment Guide](deployment.md) - Deploy to Railway, Fly.io, Render
- [API Reference](api.md) - Complete API documentation

## Troubleshooting

### Gateway won't start

Check logs:

```bash
docker logs novexa-gateway
```

Common issues:
- Missing `NOVEXA_API_KEY` environment variable
- Port 8080 already in use (change with `SERVER_PORT`)
- Invalid provider API key format

### Provider returns errors

Check provider health:

```bash
curl http://localhost:8080/api/health \
  -H "Authorization: Bearer your-secret-gateway-key"
```

Common issues:
- Invalid provider API key
- Provider API quota exceeded
- Network connectivity issues

### Streaming not working

Ensure your client supports Server-Sent Events (SSE). Most OpenAI-compatible clients support streaming.

Test with curl:

```bash
curl http://localhost:8080/v1/chat/completions \
  -H "Authorization: Bearer your-secret-gateway-key" \
  -H "Content-Type: application/json" \
  -d '{
    "model": "gpt-4o",
    "messages": [{"role": "user", "content": "Hello!"}],
    "stream": true
  }' \
  --no-buffer
```
