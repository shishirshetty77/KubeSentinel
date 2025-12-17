package analyzer

import (
	"context"
	"fmt"
	"strings"
)

// AIProvider defines the interface for different LLMs (Gemini, OpenAI, etc.)
type AIProvider interface {
	AnalyzeLog(ctx context.Context, logs string, errorReason string) (string, error)
}

// MockProvider is a placeholder that returns static analysis
type MockProvider struct{}

func (m *MockProvider) AnalyzeLog(ctx context.Context, logs string, errorReason string) (string, error) {
	// In a real implementation, you would make an HTTP POST to OpenAI/Gemini here.

	// Simple heuristic for demo purposes
	if strings.Contains(logs, "panic") {
		return "⚠️ **Root Cause:** Application Panic.\n✅ **Fix:** Check code for nil pointer dereference or unhandled errors. See logs for stack trace.", nil
	}
	if strings.Contains(logs, "OutOfMemory") || errorReason == "OOMKilled" {
		return "⚠️ **Root Cause:** Memory Limit Exceeded.\n✅ **Fix:** Increase `resources.limits.memory` in deployment YAML.", nil
	}

	return "⚠️ **Root Cause:** Unknown Application Error.\n✅ **Fix:** Check logs for more details.", nil
}

// Analyze triggers the AI analysis
func Analyze(logs string, reason string) string {
	provider := &MockProvider{}
	analysis, err := provider.AnalyzeLog(context.Background(), logs, reason)
	if err != nil {
		return fmt.Sprintf("Failed to analyze: %v", err)
	}
	return analysis
}
