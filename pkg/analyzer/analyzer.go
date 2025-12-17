package analyzer

import (
	"context"
	"fmt"
	"strings"
)

// Provider defines the interface for log analysis engines.
type Provider interface {
	Analyze(ctx context.Context, logs string, reason string) (*Result, error)
}

// Result contains the structural diagnosis.
type Result struct {
	RootCause string
	Severity  string
	Fix       string
}

func (r *Result) String() string {
	return fmt.Sprintf("⚠️ Root Cause: %s\n🔥 Severity: %s\n✅ Suggested Fix: %s", r.RootCause, r.Severity, r.Fix)
}

// LocalHeuristicProvider implements basic pattern matching without external API calls.
// It serves as a fallback or local-dev engine.
type LocalHeuristicProvider struct{}

func NewLocalProvider() *LocalHeuristicProvider {
	return &LocalHeuristicProvider{}
}

func (p *LocalHeuristicProvider) Analyze(ctx context.Context, logs string, reason string) (*Result, error) {
	// 1. Check for OOM
	if reason == "OOMKilled" || strings.Contains(logs, "OutOfMemory") || strings.Contains(logs, "java.lang.OutOfMemoryError") {
		return &Result{
			RootCause: "Memory Limit Exceeded (OOM)",
			Severity:  "Critical",
			Fix:       "Increase `resources.limits.memory` in deployment manifest or debug memory leaks.",
		}, nil
	}

	// 2. Check for Panics
	if strings.Contains(logs, "panic:") || strings.Contains(logs, "segmentation fault") {
		return &Result{
			RootCause: "Application Crash (Panic/Segfault)",
			Severity:  "High",
			Fix:       "Check application logs for null pointer dereferences or unhandled exceptions.",
		}, nil
	}

	// 3. Check for specific exit codes
	if strings.Contains(logs, "Exit Code 137") {
		return &Result{
			RootCause: "Process Terminated by SIGKILL (likely OOM)",
			Severity:  "High",
			Fix:       "Review observability metrics for memory usage spikes.",
		}, nil
	}

	// Default
	return &Result{
		RootCause: "Unknown Application Error",
		Severity:  "Medium",
		Fix:       "Review recent application logs for anomalies.",
	}, nil
}

// DefaultAnalyzer returns the configured analysis engine.
func DefaultAnalyzer() Provider {
	// TODO: Add logic to switch to OpenAI/Gemini based on env vars
	return NewLocalProvider()
}
