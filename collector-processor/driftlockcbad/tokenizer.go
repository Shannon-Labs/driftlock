// Package driftlockcbad provides compression-based anomaly detection.
// This file implements SHA-141: High-Entropy Tokenizer for preprocessing.

package driftlockcbad

import (
	"regexp"
	"sync"
	"sync/atomic"
)

// TokenizerConfig controls which high-entropy patterns to replace
type TokenizerConfig struct {
	EnableUUID   bool `json:"enable_uuid"`   // Replace UUIDs with <UUID>
	EnableHash   bool `json:"enable_hash"`   // Replace hex hashes with <HASH>
	EnableBase64 bool `json:"enable_base64"` // Replace Base64 with <B64>
	EnableJWT    bool `json:"enable_jwt"`    // Replace JWTs with <JWT>
}

// DefaultTokenizerConfig returns a configuration with all patterns enabled
func DefaultTokenizerConfig() TokenizerConfig {
	return TokenizerConfig{
		EnableUUID:   true,
		EnableHash:   true,
		EnableBase64: true,
		EnableJWT:    true,
	}
}

// Pre-compiled regex patterns (order matters - most specific first)
var (
	// JWT: Three Base64URL segments separated by dots
	// Must come before Base64 to avoid partial matches
	jwtPattern = regexp.MustCompile(`eyJ[A-Za-z0-9_-]+\.eyJ[A-Za-z0-9_-]+\.[A-Za-z0-9_-]+`)

	// UUID: 8-4-4-4-12 hex pattern (case-insensitive)
	// Must come before Hash to match full UUIDs
	uuidPattern = regexp.MustCompile(`(?i)[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}`)

	// Hex hashes: 32-64 hex characters (MD5, SHA-1, SHA-256, SHA-512)
	// Word boundaries to avoid matching substrings
	hashPattern = regexp.MustCompile(`(?i)\b[0-9a-f]{32,64}\b`)

	// Base64: 20+ chars with optional padding
	// Most general pattern, comes last
	base64Pattern = regexp.MustCompile(`[A-Za-z0-9+/]{20,}={0,2}`)
)

// Token replacements
var (
	jwtToken    = []byte("<JWT>")
	uuidToken   = []byte("<UUID>")
	hashToken   = []byte("<HASH>")
	base64Token = []byte("<B64>")
)

// TokenizerStats tracks tokenization statistics for observability
type TokenizerStats struct {
	JWTCount    int64
	UUIDCount   int64
	HashCount   int64
	Base64Count int64
	BytesSaved  int64
}

// Tokenizer preprocesses data to replace high-entropy fields with tokens
type Tokenizer struct {
	config TokenizerConfig
	stats  TokenizerStats
	mu     sync.RWMutex
}

// NewTokenizer creates a tokenizer with the given configuration
func NewTokenizer(cfg TokenizerConfig) *Tokenizer {
	return &Tokenizer{
		config: cfg,
	}
}

// hasAnyEnabled returns true if any pattern is enabled
func (t *Tokenizer) hasAnyEnabled() bool {
	return t.config.EnableUUID || t.config.EnableHash ||
		t.config.EnableBase64 || t.config.EnableJWT
}

// Tokenize replaces high-entropy fields with tokens
// Returns the tokenized data (new allocation, original unchanged)
func (t *Tokenizer) Tokenize(data []byte) []byte {
	// Fast path: if nothing enabled, return original
	if !t.hasAnyEnabled() {
		return data
	}

	// Make a copy to avoid modifying original
	result := make([]byte, len(data))
	copy(result, data)
	originalLen := len(result)

	// Apply patterns in order (most specific first)
	// JWT must come before Base64 (JWTs contain Base64-like segments)
	if t.config.EnableJWT {
		matches := jwtPattern.FindAllIndex(result, -1)
		if len(matches) > 0 {
			result = jwtPattern.ReplaceAll(result, jwtToken)
			atomic.AddInt64(&t.stats.JWTCount, int64(len(matches)))
		}
	}

	// UUID must come before Hash (UUIDs are special hex patterns)
	if t.config.EnableUUID {
		matches := uuidPattern.FindAllIndex(result, -1)
		if len(matches) > 0 {
			result = uuidPattern.ReplaceAll(result, uuidToken)
			atomic.AddInt64(&t.stats.UUIDCount, int64(len(matches)))
		}
	}

	// Hash patterns (after UUID to avoid conflicts)
	if t.config.EnableHash {
		matches := hashPattern.FindAllIndex(result, -1)
		if len(matches) > 0 {
			result = hashPattern.ReplaceAll(result, hashToken)
			atomic.AddInt64(&t.stats.HashCount, int64(len(matches)))
		}
	}

	// Base64 last (most general pattern)
	if t.config.EnableBase64 {
		matches := base64Pattern.FindAllIndex(result, -1)
		if len(matches) > 0 {
			result = base64Pattern.ReplaceAll(result, base64Token)
			atomic.AddInt64(&t.stats.Base64Count, int64(len(matches)))
		}
	}

	// Track bytes saved
	bytesSaved := int64(originalLen - len(result))
	if bytesSaved > 0 {
		atomic.AddInt64(&t.stats.BytesSaved, bytesSaved)
	}

	return result
}

// Stats returns a copy of the current tokenization statistics
func (t *Tokenizer) Stats() TokenizerStats {
	return TokenizerStats{
		JWTCount:    atomic.LoadInt64(&t.stats.JWTCount),
		UUIDCount:   atomic.LoadInt64(&t.stats.UUIDCount),
		HashCount:   atomic.LoadInt64(&t.stats.HashCount),
		Base64Count: atomic.LoadInt64(&t.stats.Base64Count),
		BytesSaved:  atomic.LoadInt64(&t.stats.BytesSaved),
	}
}

// ResetStats resets all statistics to zero
func (t *Tokenizer) ResetStats() {
	atomic.StoreInt64(&t.stats.JWTCount, 0)
	atomic.StoreInt64(&t.stats.UUIDCount, 0)
	atomic.StoreInt64(&t.stats.HashCount, 0)
	atomic.StoreInt64(&t.stats.Base64Count, 0)
	atomic.StoreInt64(&t.stats.BytesSaved, 0)
}

// Config returns the current tokenizer configuration
func (t *Tokenizer) Config() TokenizerConfig {
	t.mu.RLock()
	defer t.mu.RUnlock()
	return t.config
}

// UpdateConfig updates the tokenizer configuration
func (t *Tokenizer) UpdateConfig(cfg TokenizerConfig) {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.config = cfg
}

// ─────────────────────────────────────────────────────────────────────────────
// Tokenizer Pool - Singleton pattern for configuration reuse
// ─────────────────────────────────────────────────────────────────────────────

var (
	tokenizerPoolMu sync.RWMutex
	tokenizerPool   = make(map[TokenizerConfig]*Tokenizer)
)

// GetTokenizer returns a shared tokenizer for the given configuration.
// This avoids repeated regex compilation overhead.
func GetTokenizer(cfg TokenizerConfig) *Tokenizer {
	tokenizerPoolMu.RLock()
	if t, ok := tokenizerPool[cfg]; ok {
		tokenizerPoolMu.RUnlock()
		return t
	}
	tokenizerPoolMu.RUnlock()

	// Double-check lock pattern
	tokenizerPoolMu.Lock()
	defer tokenizerPoolMu.Unlock()
	if t, ok := tokenizerPool[cfg]; ok {
		return t
	}

	tokenizer := NewTokenizer(cfg)
	tokenizerPool[cfg] = tokenizer
	return tokenizer
}
