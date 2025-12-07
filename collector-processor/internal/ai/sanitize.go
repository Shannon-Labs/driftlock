package ai

import (
	"encoding/json"
	"regexp"
	"strings"
)

// Common prompt injection patterns to filter
var injectionPatterns = []string{
	"ignore previous",
	"disregard above",
	"ignore all",
	"forget everything",
	"new instructions",
	"system prompt",
	"system:",
	"assistant:",
	"user:",
	"[system]",
	"[assistant]",
	"</s>",
	"<|im_end|>",
	"<|im_start|>",
	"<|endoftext|>",
}

// injectionRegex matches common injection attempt patterns
var injectionRegex = regexp.MustCompile(`(?i)(ignore|disregard|forget).*(previous|above|instructions|context)`)

// SanitizeEventForPrompt sanitizes user-provided event data for safe inclusion in AI prompts.
// It applies several protections:
// 1. Length limiting to prevent context stuffing
// 2. Clear data boundary markers
// 3. Pattern filtering for known injection phrases
// 4. JSON structure preservation where possible
func SanitizeEventForPrompt(event json.RawMessage, maxLen int) string {
	if len(event) == 0 {
		return "[empty event]"
	}

	// Convert to string
	s := string(event)

	// Limit length
	if maxLen <= 0 {
		maxLen = 2048
	}
	if len(s) > maxLen {
		s = s[:maxLen] + "...[truncated]"
	}

	// Filter known injection patterns (case-insensitive)
	lower := strings.ToLower(s)
	for _, pattern := range injectionPatterns {
		if strings.Contains(lower, pattern) {
			// Replace the pattern while preserving case
			idx := strings.Index(lower, pattern)
			if idx >= 0 {
				s = s[:idx] + "[FILTERED]" + s[idx+len(pattern):]
				lower = strings.ToLower(s)
			}
		}
	}

	// Apply regex filter for more complex patterns
	s = injectionRegex.ReplaceAllString(s, "[FILTERED]")

	// Wrap in clear data boundaries
	return "---USER-DATA-START---\n" + s + "\n---USER-DATA-END---"
}

// SanitizeStringForPrompt sanitizes a plain string for inclusion in AI prompts.
func SanitizeStringForPrompt(s string, maxLen int) string {
	return SanitizeEventForPrompt(json.RawMessage(s), maxLen)
}
