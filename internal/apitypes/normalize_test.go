package apitypes_test

import (
	"testing"

	"github.com/novexa/gateway/internal/apitypes"
)

func TestMessageNormalizePromotesReasoning(t *testing.T) {
	m := &apitypes.Message{Role: "assistant", Reasoning: "Hello from reasoning"}
	m.Normalize()
	if m.Content != "Hello from reasoning" {
		t.Fatalf("content = %q, want promoted reasoning", m.Content)
	}
	if m.Reasoning != "Hello from reasoning" {
		t.Fatalf("reasoning should be preserved, got %q", m.Reasoning)
	}
}

func TestMessageNormalizePrefersReasoningContent(t *testing.T) {
	m := &apitypes.Message{
		Role:             "assistant",
		Reasoning:        "ignored",
		ReasoningContent: "from deepseek",
	}
	m.Normalize()
	if m.Content != "from deepseek" {
		t.Fatalf("content = %q, want reasoning_content", m.Content)
	}
}

func TestMessageNormalizeKeepsExistingContent(t *testing.T) {
	m := &apitypes.Message{
		Role:      "assistant",
		Content:   "visible",
		Reasoning: "hidden thinking",
	}
	m.Normalize()
	if m.Content != "visible" {
		t.Fatalf("content = %q, want visible", m.Content)
	}
}

func TestNormalizeChoicesUpdatesMessageAndDelta(t *testing.T) {
	choices := []apitypes.Choice{
		{Message: &apitypes.Message{Reasoning: "msg"}},
		{Delta: &apitypes.Message{ReasoningContent: "delta"}},
	}
	apitypes.NormalizeChoices(choices)
	if choices[0].Message.Content != "msg" {
		t.Fatalf("message content = %q", choices[0].Message.Content)
	}
	if choices[1].Delta.Content != "delta" {
		t.Fatalf("delta content = %q", choices[1].Delta.Content)
	}
}
