# Configuration Reference

Novexa Gateway uses a hierarchical configuration system with environment variables taking precedence over YAML files.

## Configuration Priority

1. **Environment variables** (highest priority)
2. **YAML config file** (if provided)
3. **Default values** (lowest priority)

## Environment Variables

All configuration values can be set via environment variables using the pattern:

```
NOVEXA_<SECTION>_<KEY>
```

### Core Settings

| Variable | Description | Default | Required |
|----------|-------------|---------|----------|
| `NOVEXA_API_KEY` | Gateway API key for client authentication | - | **Yes** |
| `NOVEXA_SERVER_PORT` | HTTP server port | `8080` | No |
| `NOVEXA_SERVER_HOST` | HTTP server host | `0.0.0.0` | No |
| `NOVEXA_SERVER_READ_TIMEOUT` | Request read timeout | `30s` | No |
| `NOVEXA_SERVER_WRITE_TIMEOUT` | Response write timeout | `120s` | No |
| `NOVEXA_SERVER_MAX_REQUEST_SIZE` | Maximum request body size | `10MB` | No |

### Provider API Keys

| Variable | Description |
|----------|-------------|
| `OPENAI_API_KEY` | OpenAI API key |
| `ANTHROPIC_API_KEY` | Anthropic API key |
| `GEMINI_API_KEY` | Google Gemini API key |
| `DEEPSEEK_API_KEY` | DeepSeek API key |
| `OPENROUTER_API_KEY` | OpenRouter API key |
| `GROQ_API_KEY` | Groq API key |

### Provider Base URLs (Optional)

| Variable | Description | Default |
|----------|-------------|---------|
| `OPENAI_BASE_URL` | OpenAI API base URL | `https://api.openai.com/v1` |
| `ANTHROPIC_BASE_URL` | Anthropic API base URL | `https://api.anthropic.com` |
| `GEMINI_BASE_URL` | Gemini API base URL | `https://generativelanguage.googleapis.com/v1beta` |
| `DEEPSEEK_BASE_URL` | DeepSeek API base URL | `https://api.deepseek.com/v1` |
| `OPENROUTER_BASE_URL` | OpenRouter API base URL | `https://openrouter.ai/api/v1` |
| `GROQ_BASE_URL` | Groq API base URL | `https://api.groq.com/openai/v1` |
| `OLLAMA_BASE_URL` | Ollama API base URL | `http://localhost:11434` |
| `LMSTUDIO_BASE_URL` | LM Studio API base URL | `http://localhost:1234/v1` |

### Database Settings

| Variable | Description | Default |
|----------|-------------|---------|
| `NOVEXA_DATABASE_DRIVER` | Database driver (`sqlite` or `postgres`) | `sqlite` |
| `NOVEXA_DATABASE_DSN` | Database connection string | `./data/novexa.db` |

### Logging Settings

| Variable | Description | Default |
|----------|-------------|---------|
| `NOVEXA_LOGGING_LEVEL` | Log level (`debug`, `info`, `warn`, `error`) | `info` |
| `NOVEXA_LOGGING_FORMAT` | Log format (`json` or `console`) | `json` |
| `NOVEXA_LOGGING_LOG_PROMPTS` | Log request prompts (security risk) | `false` |
| `NOVEXA_LOGGING_LOG_RESPONSES` | Log response bodies | `false` |

### Rate Limiting

| Variable | Description | Default |
|----------|-------------|---------|
| `NOVEXA_RATELIMIT_GLOBAL_RPM` | Global requests per minute | `1000` |
| `NOVEXA_RATELIMIT_PERPROVIDER_RPM` | Per-provider requests per minute | `100` |

### Health Monitoring

| Variable | Description | Default |
|----------|-------------|---------|
| `NOVEXA_HEALTH_CHECK_INTERVAL` | Health check interval | `60s` |
| `NOVEXA_HEALTH_CHECK_TIMEOUT` | Health check timeout | `10s` |

## YAML Configuration File

Create a `config.yaml` file for advanced configuration:

```yaml
# Server configuration
server:
  host: "0.0.0.0"
  port: 8080
  read_timeout: 30s
  write_timeout: 120s
  max_request_size: 10MB
  cors:
    enabled: true
    origins: ["*"]
    methods: ["GET", "POST", "OPTIONS"]
    headers: ["Authorization", "Content-Type"]

# Gateway API key
api_key: "${NOVEXA_API_KEY}"

# Provider configuration
providers:
  openai:
    enabled: true
    api_key: "${OPENAI_API_KEY}"
    base_url: "https://api.openai.com/v1"
    timeout: 60s
    max_retries: 3
  
  anthropic:
    enabled: true
    api_key: "${ANTHROPIC_API_KEY}"
    base_url: "https://api.anthropic.com"
    timeout: 60s
    max_retries: 3
  
  gemini:
    enabled: true
    api_key: "${GEMINI_API_KEY}"
    base_url: "https://generativelanguage.googleapis.com/v1beta"
    timeout: 60s
    max_retries: 3
  
  deepseek:
    enabled: true
    api_key: "${DEEPSEEK_API_KEY}"
    base_url: "https://api.deepseek.com/v1"
    timeout: 60s
    max_retries: 3
  
  openrouter:
    enabled: true
    api_key: "${OPENROUTER_API_KEY}"
    base_url: "https://openrouter.ai/api/v1"
    timeout: 60s
    max_retries: 3
  
  groq:
    enabled: true
    api_key: "${GROQ_API_KEY}"
    base_url: "https://api.groq.com/openai/v1"
    timeout: 30s
    max_retries: 3
  
  ollama:
    enabled: false
    base_url: "http://localhost:11434"
    timeout: 120s
    max_retries: 1
  
  lmstudio:
    enabled: false
    base_url: "http://localhost:1234/v1"
    timeout: 120s
    max_retries: 1

# Model routing (optional - auto-detected if not specified)
routes:
  "gpt-4o":
    provider: openai
  "gpt-4o-mini":
    provider: openai
  "claude-sonnet-4-20250514":
    provider: anthropic
  "claude-3-5-haiku-20241022":
    provider: anthropic
  "gemini-2.5-pro":
    provider: gemini
  "gemini-2.5-flash":
    provider: gemini
  "deepseek-chat":
    provider: deepseek
  "deepseek-reasoner":
    provider: deepseek

# Model aliases (optional)
aliases:
  "fast": "gpt-4o-mini"
  "smart": "gpt-4o"
  "coding": "deepseek-chat"
  "cheap": "gpt-4o-mini"

# Fallback chains (optional)
fallbacks:
  "deepseek-chat":
    - provider: deepseek
    - provider: openrouter
      model: "deepseek/deepseek-chat"
    - provider: groq
      model: "deepseek-r1-distill-llama-70b"

# Retry configuration
retry:
  max_retries: 3
  initial_backoff: 100ms
  max_backoff: 5s
  backoff_multiplier: 2.0
  retryable_status_codes: [429, 500, 502, 503]

# Database configuration
database:
  driver: "sqlite"
  dsn: "./data/novexa.db"
  max_open_conns: 10
  max_idle_conns: 5

# Logging configuration
logging:
  level: "info"
  format: "json"
  log_prompts: false
  log_responses: false

# Rate limiting
rate_limit:
  enabled: true
  global:
    requests_per_minute: 1000
  per_provider:
    requests_per_minute: 100

# Health monitoring
health:
  check_interval: 60s
  timeout: 10s
  unhealthy_threshold: 3

# Usage tracking
usage:
  enabled: true

# Cost tracking
cost:
  enabled: true
  currency: "USD"
```

## Minimal Configuration Example

For the simplest setup, only provider API keys are required:

```bash
export NOVEXA_API_KEY=your-secret-gateway-key
export OPENAI_API_KEY=sk-your-openai-key
```

The gateway will:
- Start on port 8080
- Enable OpenAI provider
- Auto-discover available models
- Use SQLite for storage
- Use sensible defaults for all other settings

## Configuration Validation

The gateway validates configuration on startup and will exit with an error if:
- `NOVEXA_API_KEY` is not set
- No providers are configured
- Invalid timeout values
- Invalid log level

## Hot Reload

Configuration can be reloaded without restarting the gateway:

```bash
curl -X PUT http://localhost:8080/api/config/reload \
  -H "Authorization: Bearer your-secret-gateway-key"
```

Note: Provider API key changes require a restart.
