# Provider Setup Guide

Novexa Gateway supports 9 AI providers out of the box. This guide covers setup for each provider.

## Provider Overview

| Provider | Format | Translation | Streaming | Notes |
|----------|--------|-------------|-----------|-------|
| OpenAI | OpenAI | Passthrough | ✅ | Most popular |
| Anthropic | Messages API | Full translation | ✅ | Claude models |
| Google Gemini | GenerateContent | Full translation | ✅ | Gemini models |
| DeepSeek | OpenAI-compatible | Passthrough | ✅ | Cost-effective |
| OpenRouter | OpenAI-compatible | Passthrough | ✅ | Multi-model gateway |
| Groq | OpenAI-compatible | Passthrough | ✅ | Ultra-fast inference |
| Ollama | OpenAI-compatible | Passthrough | ✅ | Local models |
| LM Studio | OpenAI-compatible | Passthrough | ✅ | Local models |
| Generic | OpenAI-compatible | Passthrough | ✅ | Any compatible endpoint |

## OpenAI

### Setup

1. Get API key from [platform.openai.com](https://platform.openai.com/api-keys)
2. Set environment variable:

```bash
export OPENAI_API_KEY=sk-your-key-here
```

### Configuration

```yaml
providers:
  openai:
    enabled: true
    api_key: "${OPENAI_API_KEY}"
    base_url: "https://api.openai.com/v1"  # Default
    timeout: 60s
    max_retries: 3
```

### Available Models

- `gpt-4o`
- `gpt-4o-mini`
- `gpt-4-turbo`
- `gpt-3.5-turbo`
- `o1-preview`
- `o1-mini`

### Cost

- GPT-4o: $2.50 / 1M input tokens, $10.00 / 1M output tokens
- GPT-4o-mini: $0.15 / 1M input tokens, $0.60 / 1M output tokens

---

## Anthropic

### Setup

1. Get API key from [console.anthropic.com](https://console.anthropic.com/)
2. Set environment variable:

```bash
export ANTHROPIC_API_KEY=sk-ant-your-key-here
```

### Configuration

```yaml
providers:
  anthropic:
    enabled: true
    api_key: "${ANTHROPIC_API_KEY}"
    base_url: "https://api.anthropic.com"  # Default
    timeout: 60s
    max_retries: 3
```

### Available Models

- `claude-sonnet-4-20250514`
- `claude-3-5-haiku-20241022`
- `claude-3-opus-20240229`

### Cost

- Claude Sonnet 4: $3.00 / 1M input tokens, $15.00 / 1M output tokens
- Claude 3.5 Haiku: $0.80 / 1M input tokens, $4.00 / 1M output tokens

### Notes

Anthropic uses a different API format (Messages API). The gateway automatically translates between OpenAI and Anthropic formats.

---

## Google Gemini

### Setup

1. Get API key from [aistudio.google.com](https://aistudio.google.com/app/apikey)
2. Set environment variable:

```bash
export GEMINI_API_KEY=your-key-here
```

### Configuration

```yaml
providers:
  gemini:
    enabled: true
    api_key: "${GEMINI_API_KEY}"
    base_url: "https://generativelanguage.googleapis.com/v1beta"  # Default
    timeout: 60s
    max_retries: 3
```

### Available Models

- `gemini-2.5-pro`
- `gemini-2.5-flash`
- `gemini-1.5-pro`
- `gemini-1.5-flash`

### Cost

- Gemini 2.5 Pro: $1.25 / 1M input tokens, $5.00 / 1M output tokens
- Gemini 2.5 Flash: $0.075 / 1M input tokens, $0.30 / 1M output tokens

### Notes

Gemini uses a different API format (GenerateContent). The gateway automatically translates between OpenAI and Gemini formats.

---

## DeepSeek

### Setup

1. Get API key from [platform.deepseek.com](https://platform.deepseek.com/)
2. Set environment variable:

```bash
export DEEPSEEK_API_KEY=your-key-here
```

### Configuration

```yaml
providers:
  deepseek:
    enabled: true
    api_key: "${DEEPSEEK_API_KEY}"
    base_url: "https://api.deepseek.com/v1"  # Default
    timeout: 60s
    max_retries: 3
```

### Available Models

- `deepseek-chat`
- `deepseek-reasoner`

### Cost

- DeepSeek Chat: $0.14 / 1M input tokens, $0.28 / 1M output tokens
- DeepSeek Reasoner: $0.55 / 1M input tokens, $2.19 / 1M output tokens

### Notes

DeepSeek is OpenAI-compatible, so no translation is needed. Very cost-effective for coding tasks.

---

## OpenRouter

### Setup

1. Get API key from [openrouter.ai](https://openrouter.ai/keys)
2. Set environment variable:

```bash
export OPENROUTER_API_KEY=your-key-here
```

### Configuration

```yaml
providers:
  openrouter:
    enabled: true
    api_key: "${OPENROUTER_API_KEY}"
    base_url: "https://openrouter.ai/api/v1"  # Default
    timeout: 60s
    max_retries: 3
```

### Available Models

OpenRouter provides access to hundreds of models. Popular ones:

- `openai/gpt-4o`
- `anthropic/claude-sonnet-4`
- `google/gemini-2.5-pro`
- `deepseek/deepseek-chat`
- `meta-llama/llama-3.3-70b-instruct`

### Cost

Varies by model. OpenRouter adds a small markup on top of provider costs.

### Notes

OpenRouter is a meta-provider that routes to many other providers. Useful as a fallback when direct providers are unavailable.

---

## Groq

### Setup

1. Get API key from [console.groq.com](https://console.groq.com/)
2. Set environment variable:

```bash
export GROQ_API_KEY=your-key-here
```

### Configuration

```yaml
providers:
  groq:
    enabled: true
    api_key: "${GROQ_API_KEY}"
    base_url: "https://api.groq.com/openai/v1"  # Default
    timeout: 30s
    max_retries: 3
```

### Available Models

- `llama-3.3-70b-versatile`
- `llama-3.1-8b-instant`
- `mixtral-8x7b-32768`
- `gemma2-9b-it`

### Cost

- Llama 3.3 70B: $0.59 / 1M input tokens, $0.79 / 1M output tokens
- Llama 3.1 8B: $0.05 / 1M input tokens, $0.08 / 1M output tokens

### Notes

Groq is extremely fast (ultra-low latency) but has smaller context windows. Great for quick responses.

---

## Ollama (Local)

### Setup

1. Install Ollama from [ollama.com](https://ollama.com/)
2. Pull a model:

```bash
ollama pull llama3.1:8b
```

3. Ollama runs locally on `http://localhost:11434` by default

### Configuration

```yaml
providers:
  ollama:
    enabled: true
    base_url: "http://localhost:11434"  # Default
    timeout: 120s
    max_retries: 1
```

Note: No API key needed for local Ollama.

### Available Models

Any model you've pulled with Ollama:

- `llama3.1:8b`
- `llama3.1:70b`
- `mistral:7b`
- `codellama:13b`
- `phi3:mini`

### Cost

Free (runs on your hardware)

### Notes

- Requires Ollama to be running on the same machine or accessible network
- Performance depends on your hardware (GPU recommended)
- No API key required
- Longer timeouts recommended (120s+)

---

## LM Studio (Local)

### Setup

1. Download LM Studio from [lmstudio.ai](https://lmstudio.ai/)
2. Load a model in LM Studio
3. Start the local server (default: `http://localhost:1234`)

### Configuration

```yaml
providers:
  lmstudio:
    enabled: true
    base_url: "http://localhost:1234/v1"  # Default
    timeout: 120s
    max_retries: 1
```

Note: No API key needed for local LM Studio.

### Available Models

Any model loaded in LM Studio. The model name is whatever you loaded.

### Cost

Free (runs on your hardware)

### Notes

- Requires LM Studio to be running with a model loaded
- Performance depends on your hardware
- No API key required
- Longer timeouts recommended (120s+)

---

## Generic (Any OpenAI-Compatible Endpoint)

### Setup

Use this for any OpenAI-compatible API endpoint not covered above.

### Configuration

```yaml
providers:
  my-custom-provider:
    enabled: true
    api_key: "${CUSTOM_API_KEY}"
    base_url: "https://my-custom-api.com/v1"
    timeout: 60s
    max_retries: 3
```

### Notes

- Must be fully OpenAI-compatible
- Supports streaming
- Useful for self-hosted models (vLLM, text-generation-inference, etc.)

---

## Model Routing

### Auto-Detection

If you don't specify routes in config, the gateway auto-detects available models from each provider at startup.

### Manual Routing

For explicit control, define routes:

```yaml
routes:
  "gpt-4o":
    provider: openai
  "claude-sonnet-4-20250514":
    provider: anthropic
  "my-custom-model":
    provider: my-custom-provider
```

### Aliases

Create shortcuts for frequently used models:

```yaml
aliases:
  "fast": "gpt-4o-mini"
  "smart": "gpt-4o"
  "coding": "deepseek-chat"
  "local": "llama3.1:8b"
```

Then use the alias in requests:

```bash
curl http://localhost:8080/v1/chat/completions \
  -H "Authorization: Bearer your-key" \
  -d '{"model": "fast", "messages": [...]}'
```

### Fallback Chains

If a provider fails, automatically try the next:

```yaml
fallbacks:
  "deepseek-chat":
    - provider: deepseek
    - provider: openrouter
      model: "deepseek/deepseek-chat"
    - provider: groq
      model: "deepseek-r1-distill-llama-70b"
```

---

## Provider Health Monitoring

The gateway automatically monitors provider health every 60 seconds (configurable).

### Check Health

```bash
curl http://localhost:8080/api/health \
  -H "Authorization: Bearer your-key"
```

### Response Example

```json
{
  "providers": [
    {
      "name": "openai",
      "healthy": true,
      "latency_ms": 245,
      "last_error": null,
      "checked_at": "2026-07-19T10:30:00Z"
    },
    {
      "name": "anthropic",
      "healthy": false,
      "latency_ms": 0,
      "last_error": "connection timeout",
      "checked_at": "2026-07-19T10:30:00Z"
    }
  ]
}
```

### Unhealthy Providers

If a provider is marked unhealthy:
- Requests to that provider will fail fast
- Fallback chains will skip unhealthy providers
- Health checks continue to run
- Provider is automatically re-enabled when healthy again

---

## Troubleshooting

### Provider Returns 401 Unauthorized

- Check API key is correct
- Ensure environment variable is set
- Verify no typos in config

### Provider Returns 429 Too Many Requests

- You've hit the provider's rate limit
- Wait and retry
- Consider using a different provider or fallback chain

### Local Provider (Ollama/LM Studio) Not Responding

- Ensure the service is running
- Check the base_url is correct
- Verify network connectivity (if not on localhost)
- Increase timeout in config

### Streaming Not Working

- Ensure provider supports streaming
- Check client supports Server-Sent Events (SSE)
- Test with curl using `--no-buffer` flag
