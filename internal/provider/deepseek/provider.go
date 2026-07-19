package deepseek

import (
	"context"
	"net/http"
	"time"

	"github.com/novexa/gateway/internal/apitypes"
	"github.com/novexa/gateway/internal/provider"
)

// Provider implements the provider.Provider interface for DeepSeek.
type Provider struct {
	name    string
	apiKey  string
	baseURL string
	client  *http.Client
}

// NewProvider creates a new DeepSeek provider.
func NewProvider(apiKey, baseURL string, timeout time.Duration) *Provider {
	return &Provider{
		name:    "deepseek",
		apiKey:  apiKey,
		baseURL: baseURL,
		client: &http.Client{
			Timeout: timeout,
		},
	}
}

// Name returns the provider name.
func (p *Provider) Name() string { return p.name }

// ChatCompletion is not yet implemented.
func (p *Provider) ChatCompletion(ctx context.Context, req *apitypes.ChatCompletionRequest) (*apitypes.ChatCompletionResponse, error) {
	return nil, provider.ErrNotImplemented
}

// ChatCompletionStream is not yet implemented.
func (p *Provider) ChatCompletionStream(ctx context.Context, req *apitypes.ChatCompletionRequest) (<-chan apitypes.StreamChunk, error) {
	return nil, provider.ErrNotImplemented
}

// Embeddings is not yet implemented.
func (p *Provider) Embeddings(ctx context.Context, req *apitypes.EmbeddingRequest) (*apitypes.EmbeddingResponse, error) {
	return nil, provider.ErrNotImplemented
}

// ListModels is not yet implemented.
func (p *Provider) ListModels(ctx context.Context) ([]provider.ModelInfo, error) {
	return nil, provider.ErrNotImplemented
}

// GetPricing is not yet implemented.
func (p *Provider) GetPricing(ctx context.Context) (map[string]provider.PricingInfo, error) {
	return nil, provider.ErrNotImplemented
}

// HealthCheck is not yet implemented.
func (p *Provider) HealthCheck(ctx context.Context) (*provider.HealthStatus, error) {
	return nil, provider.ErrNotImplemented
}

// SupportsModel is not yet implemented.
func (p *Provider) SupportsModel(modelID string) bool { return false }
