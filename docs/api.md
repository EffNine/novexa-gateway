# API Reference

Novexa Gateway exposes an OpenAI-compatible API along with dashboard endpoints for monitoring.

## Authentication

All endpoints require authentication via the `Authorization` header:

```
Authorization: Bearer <your-api-key>
```

The API key is set via the `NOVEXA_API_KEY` environment variable.

---

## OpenAI-Compatible Endpoints

### Chat Completions

**Endpoint**: `POST /v1/chat/completions`

Creates a model response for the given chat conversation.

#### Request Body

```json
{
  "model": "gpt-4o",
  "messages": [
    {"role": "system", "content": "You are a helpful assistant."},
    {"role": "user", "content": "Hello!"}
  ],
  "temperature": 0.7,
  "max_tokens": 1024,
  "stream": false
}
```

**Required Fields**:
- `model` (string): Model ID or alias to use
- `messages` (array): Array of message objects

**Optional Fields**:
- `temperature` (number): Sampling temperature (0-2). Default: 1.0
- `max_tokens` (number): Maximum tokens to generate. Default: provider default
- `stream` (boolean): Enable streaming. Default: false
- `top_p` (number): Nucleus sampling parameter. Default: 1.0
- `frequency_penalty` (number): Frequency penalty (-2 to 2). Default: 0
- `presence_penalty` (number): Presence penalty (-2 to 2). Default: 0
- `stop` (string|array): Stop sequences. Default: null

#### Non-Streaming Response

```json
{
  "id": "chatcmpl-abc123",
  "object": "chat.completion",
  "created": 1677652288,
  "model": "gpt-4o",
  "choices": [
    {
      "index": 0,
      "message": {
        "role": "assistant",
        "content": "Hello! How can I help you today?"
      },
      "finish_reason": "stop"
    }
  ],
  "usage": {
    "prompt_tokens": 12,
    "completion_tokens": 9,
    "total_tokens": 21
  }
}
```

#### Streaming Response

When `stream: true`, returns Server-Sent Events (SSE):

```
data: {"id":"chatcmpl-abc123","object":"chat.completion.chunk","created":1677652288,"model":"gpt-4o","choices":[{"index":0,"delta":{"role":"assistant","content":"Hello"},"finish_reason":null}]}

data: {"id":"chatcmpl-abc123","object":"chat.completion.chunk","created":1677652288,"model":"gpt-4o","choices":[{"index":0,"delta":{"content":"!"},"finish_reason":null}]}

data: {"id":"chatcmpl-abc123","object":"chat.completion.chunk","created":1677652288,"model":"gpt-4o","choices":[{"index":0,"delta":{},"finish_reason":"stop"}]}

data: [DONE]
```

---

### List Models

**Endpoint**: `GET /v1/models`

Lists all available models.

#### Response

```json
{
  "object": "list",
  "data": [
    {
      "id": "gpt-4o",
      "object": "model",
      "created": 1677652288,
      "owned_by": "novexa"
    },
    {
      "id": "claude-sonnet-4-20250514",
      "object": "model",
      "created": 1677652288,
      "owned_by": "novexa"
    }
  ]
}
```

**Note**: The `owned_by` field is always `"novexa"` to hide the actual provider.

---

### Embeddings

**Endpoint**: `POST /v1/embeddings`

Creates an embedding vector representing the input text.

#### Request Body

```json
{
  "model": "text-embedding-3-small",
  "input": "The food was delicious"
}
```

**Required Fields**:
- `model` (string): Model ID to use
- `input` (string|array): Input text to embed

#### Response

```json
{
  "object": "list",
  "data": [
    {
      "object": "embedding",
      "embedding": [0.0023, -0.009, 0.015, ...],
      "index": 0
    }
  ],
  "model": "text-embedding-3-small",
  "usage": {
    "prompt_tokens": 6,
    "total_tokens": 6
  }
}
```

---

## Dashboard Endpoints

### Health Check

**Endpoint**: `GET /health`

Simple health check (no authentication required).

#### Response

```json
{
  "status": "ok"
}
```

---

### Provider Health

**Endpoint**: `GET /api/health`

Returns health status of all configured providers.

#### Response

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

---

### List Providers

**Endpoint**: `GET /api/providers`

Lists all configured providers with their status.

#### Response

```json
{
  "providers": [
    {
      "name": "openai",
      "enabled": true,
      "base_url": "https://api.openai.com/v1",
      "healthy": true,
      "models_count": 12
    }
  ]
}
```

---

### Usage Statistics

**Endpoint**: `GET /api/usage`

Returns usage statistics with optional filters.

#### Query Parameters

- `from` (string): Start date (ISO 8601). Default: 24 hours ago
- `to` (string): End date (ISO 8601). Default: now
- `model` (string): Filter by model
- `provider` (string): Filter by provider

#### Response

```json
{
  "period": {
    "from": "2026-07-18T10:30:00Z",
    "to": "2026-07-19T10:30:00Z"
  },
  "summary": {
    "total_requests": 150,
    "total_tokens": 45000,
    "prompt_tokens": 30000,
    "completion_tokens": 15000,
    "total_cost_usd": 0.45,
    "avg_latency_ms": 1200
  },
  "by_model": [
    {
      "model": "gpt-4o",
      "requests": 80,
      "tokens": 24000,
      "cost_usd": 0.30
    }
  ],
  "by_provider": [
    {
      "provider": "openai",
      "requests": 80,
      "tokens": 24000,
      "cost_usd": 0.30
    }
  ]
}
```

---

### Cost Breakdown

**Endpoint**: `GET /api/usage/costs`

Returns detailed cost breakdown.

#### Query Parameters

- `period` (string): Time period (`day`, `week`, `month`). Default: `day`
- `provider` (string): Filter by provider
- `model` (string): Filter by model

#### Response

```json
{
  "period": "day",
  "total_cost_usd": 0.45,
  "currency": "USD",
  "breakdown": [
    {
      "date": "2026-07-19",
      "cost_usd": 0.45,
      "by_provider": [
        {
          "provider": "openai",
          "cost_usd": 0.30
        },
        {
          "provider": "anthropic",
          "cost_usd": 0.15
        }
      ]
    }
  ]
}
```

---

### Request Logs

**Endpoint**: `GET /api/logs`

Returns recent request logs.

#### Query Parameters

- `limit` (number): Number of logs to return. Default: 100, Max: 1000
- `offset` (number): Offset for pagination. Default: 0
- `status` (string): Filter by status (`success`, `error`)

#### Response

```json
{
  "logs": [
    {
      "id": "log-abc123",
      "timestamp": "2026-07-19T10:30:00Z",
      "method": "POST",
      "path": "/v1/chat/completions",
      "model": "gpt-4o",
      "provider": "openai",
      "status_code": 200,
      "latency_ms": 1200,
      "tokens": 150,
      "cost_usd": 0.003
    }
  ],
  "total": 150,
  "limit": 100,
  "offset": 0
}
```

---

### Current Configuration

**Endpoint**: `GET /api/config`

Returns current configuration (secrets redacted).

#### Response

```json
{
  "server": {
    "port": 8080,
    "host": "0.0.0.0"
  },
  "providers": {
    "openai": {
      "enabled": true,
      "base_url": "https://api.openai.com/v1",
      "api_key": "***"
    }
  },
  "routes": {
    "gpt-4o": {
      "provider": "openai"
    }
  },
  "aliases": {
    "fast": "gpt-4o-mini"
  }
}
```

---

### Reload Configuration

**Endpoint**: `PUT /api/config/reload`

Reloads configuration from file and environment variables.

#### Response

```json
{
  "status": "ok",
  "message": "Configuration reloaded successfully"
}
```

**Note**: Provider API key changes require a restart.

---

## Error Responses

All errors follow the OpenAI error format:

```json
{
  "error": {
    "message": "Model 'invalid-model' not found",
    "type": "invalid_request_error",
    "param": "model",
    "code": "model_not_found"
  }
}
```

### Error Types

| HTTP Status | Type | Description |
|-------------|------|-------------|
| 400 | `invalid_request_error` | Invalid request format or parameters |
| 401 | `authentication_error` | Invalid or missing API key |
| 403 | `permission_error` | API key doesn't have permission |
| 404 | `not_found_error` | Resource not found |
| 429 | `rate_limit_error` | Rate limit exceeded |
| 500 | `server_error` | Internal server error |
| 502 | `provider_error` | Provider returned an error |
| 503 | `service_unavailable` | Service temporarily unavailable |

### Common Error Codes

- `model_not_found`: Specified model doesn't exist
- `invalid_api_key`: API key is invalid
- `rate_limit_exceeded`: Too many requests
- `provider_unavailable`: Provider is down or unreachable
- `invalid_request`: Request format is invalid
- `context_length_exceeded`: Input too long for model

---

## Rate Limiting

Rate limits are enforced per API key:

- **Default**: 1000 requests per minute (global)
- **Per Provider**: 100 requests per minute

When rate limited, returns:

```json
{
  "error": {
    "message": "Rate limit exceeded",
    "type": "rate_limit_error",
    "code": "rate_limit_exceeded"
  }
}
```

HTTP Status: `429 Too Many Requests`

Headers:
- `X-RateLimit-Limit`: Maximum requests allowed
- `X-RateLimit-Remaining`: Requests remaining
- `X-RateLimit-Reset`: Time when limit resets (Unix timestamp)

---

## CORS

CORS is enabled by default for all origins. Configure in `config.yaml`:

```yaml
server:
  cors:
    enabled: true
    origins: ["*"]
    methods: ["GET", "POST", "OPTIONS"]
    headers: ["Authorization", "Content-Type"]
```

---

## Examples

### Basic Chat Completion

```bash
curl http://localhost:8080/v1/chat/completions \
  -H "Authorization: Bearer your-api-key" \
  -H "Content-Type: application/json" \
  -d '{
    "model": "gpt-4o",
    "messages": [
      {"role": "user", "content": "Hello!"}
    ]
  }'
```

### Streaming Chat Completion

```bash
curl http://localhost:8080/v1/chat/completions \
  -H "Authorization: Bearer your-api-key" \
  -H "Content-Type: application/json" \
  -d '{
    "model": "gpt-4o",
    "messages": [
      {"role": "user", "content": "Hello!"}
    ],
    "stream": true
  }' \
  --no-buffer
```

### Using Model Alias

```bash
curl http://localhost:8080/v1/chat/completions \
  -H "Authorization: Bearer your-api-key" \
  -H "Content-Type: application/json" \
  -d '{
    "model": "fast",
    "messages": [
      {"role": "user", "content": "Hello!"}
    ]
  }'
```

### Get Usage Statistics

```bash
curl "http://localhost:8080/api/usage?period=day" \
  -H "Authorization: Bearer your-api-key"
```

### Get Provider Health

```bash
curl http://localhost:8080/api/health \
  -H "Authorization: Bearer your-api-key"
```
