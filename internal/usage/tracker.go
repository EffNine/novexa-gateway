package usage

import (
	"sync"
	"time"

	"github.com/novexa/gateway/internal/database"
	"go.uber.org/zap"
)

// Record represents a usage record to be tracked
type Record struct {
	RequestID        string
	ModelID          string
	ProviderModelID  string
	Provider         string
	PromptTokens     int
	CompletionTokens int
	TotalTokens      int
	LatencyMs        int64
	StatusCode       int
	IsStream         bool
	ErrorMessage     *string
	CreatedAt        time.Time
}

// Tracker handles usage tracking
type Tracker struct {
	db     *database.Database
	logger *zap.Logger
	mu     sync.Mutex
	buffer []*Record
}

// NewTracker creates a new usage tracker
func NewTracker(db *database.Database, logger *zap.Logger) *Tracker {
	return &Tracker{
		db:     db,
		logger: logger,
		buffer: make([]*Record, 0, 100),
	}
}

// Record records a usage entry
func (t *Tracker) Record(r *Record) {
	if t == nil || t.db == nil {
		return
	}

	t.mu.Lock()
	t.buffer = append(t.buffer, r)
	t.mu.Unlock()

	// Flush immediately for low-volume usage
	t.flush()
}

// flush writes buffered records to the database
func (t *Tracker) flush() {
	t.mu.Lock()
	if len(t.buffer) == 0 {
		t.mu.Unlock()
		return
	}

	batch := make([]*Record, len(t.buffer))
	copy(batch, t.buffer)
	t.buffer = t.buffer[:0]
	t.mu.Unlock()

	// Convert to database models
	for _, r := range batch {
		record := &database.UsageRecord{
			ID:               r.RequestID,
			RequestID:        r.RequestID,
			ModelID:          r.ModelID,
			ProviderModelID:  r.ProviderModelID,
			Provider:         r.Provider,
			PromptTokens:     r.PromptTokens,
			CompletionTokens: r.CompletionTokens,
			TotalTokens:      r.TotalTokens,
			LatencyMs:        r.LatencyMs,
			EstimatedCostUSD: EstimateCost(r.ModelID, r.Provider, r.PromptTokens, r.CompletionTokens),
			StatusCode:       r.StatusCode,
			IsStream:         r.IsStream,
			ErrorMessage:     r.ErrorMessage,
			CreatedAt:        r.CreatedAt,
		}

		if err := t.db.DB.Create(record).Error; err != nil {
			t.logger.Error("failed to save usage record",
				zap.String("request_id", r.RequestID),
				zap.Error(err),
			)
		}
	}
}

// EstimateCost estimates the cost of a request
func EstimateCost(model, provider string, promptTokens, completionTokens int) float64 {
	// Default cost rates (USD per 1K tokens)
	// These should be loaded from the cost_rates table in production
	rates := getDefaultCostRates(model, provider)

	promptCost := float64(promptTokens) / 1000.0 * rates.PromptCostPer1K
	completionCost := float64(completionTokens) / 1000.0 * rates.CompletionCostPer1K

	return promptCost + completionCost
}

// costRate represents token pricing
type costRate struct {
	PromptCostPer1K     float64
	CompletionCostPer1K float64
}

// getDefaultCostRates returns default cost rates for known models
func getDefaultCostRates(model, provider string) costRate {
	rates := map[string]costRate{
		// OpenAI
		"gpt-4o":        {PromptCostPer1K: 0.0025, CompletionCostPer1K: 0.010},
		"gpt-4o-mini":   {PromptCostPer1K: 0.00015, CompletionCostPer1K: 0.0006},
		"gpt-4-turbo":   {PromptCostPer1K: 0.010, CompletionCostPer1K: 0.030},
		"gpt-3.5-turbo": {PromptCostPer1K: 0.0005, CompletionCostPer1K: 0.0015},

		// Anthropic
		"claude-sonnet-4-20250514":  {PromptCostPer1K: 0.003, CompletionCostPer1K: 0.015},
		"claude-3-5-haiku-20241022": {PromptCostPer1K: 0.0008, CompletionCostPer1K: 0.004},
		"claude-3-opus-20240229":    {PromptCostPer1K: 0.015, CompletionCostPer1K: 0.075},

		// Gemini
		"gemini-2.5-pro":   {PromptCostPer1K: 0.00125, CompletionCostPer1K: 0.005},
		"gemini-2.5-flash": {PromptCostPer1K: 0.000075, CompletionCostPer1K: 0.0003},
		"gemini-1.5-pro":   {PromptCostPer1K: 0.00125, CompletionCostPer1K: 0.005},
		"gemini-1.5-flash": {PromptCostPer1K: 0.000075, CompletionCostPer1K: 0.0003},

		// DeepSeek
		"deepseek-chat":     {PromptCostPer1K: 0.00014, CompletionCostPer1K: 0.00028},
		"deepseek-reasoner": {PromptCostPer1K: 0.00055, CompletionCostPer1K: 0.00219},

		// Groq
		"llama-3.3-70b-versatile": {PromptCostPer1K: 0.00059, CompletionCostPer1K: 0.00079},
		"llama-3.1-8b-instant":    {PromptCostPer1K: 0.00005, CompletionCostPer1K: 0.00008},
	}

	if rate, ok := rates[model]; ok {
		return rate
	}

	// Default fallback rate
	return costRate{PromptCostPer1K: 0.001, CompletionCostPer1K: 0.002}
}
