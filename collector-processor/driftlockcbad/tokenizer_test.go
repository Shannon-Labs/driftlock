package driftlockcbad

import (
	"testing"
)

func TestTokenizeUUID(t *testing.T) {
	tokenizer := NewTokenizer(TokenizerConfig{EnableUUID: true})

	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "single UUID",
			input:    `{"id": "550e8400-e29b-41d4-a716-446655440000"}`,
			expected: `{"id": "<UUID>"}`,
		},
		{
			name:     "multiple UUIDs",
			input:    `{"a": "550e8400-e29b-41d4-a716-446655440000", "b": "6ba7b810-9dad-11d1-80b4-00c04fd430c8"}`,
			expected: `{"a": "<UUID>", "b": "<UUID>"}`,
		},
		{
			name:     "uppercase UUID",
			input:    `{"id": "550E8400-E29B-41D4-A716-446655440000"}`,
			expected: `{"id": "<UUID>"}`,
		},
		{
			name:     "no UUID",
			input:    `{"name": "test"}`,
			expected: `{"name": "test"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tokenizer.Tokenize([]byte(tt.input))
			if string(result) != tt.expected {
				t.Errorf("got %s, want %s", result, tt.expected)
			}
		})
	}
}

func TestTokenizeHash(t *testing.T) {
	tokenizer := NewTokenizer(TokenizerConfig{EnableHash: true})

	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "MD5 hash (32 chars)",
			input:    `{"hash": "d41d8cd98f00b204e9800998ecf8427e"}`,
			expected: `{"hash": "<HASH>"}`,
		},
		{
			name:     "SHA-1 hash (40 chars)",
			input:    `{"hash": "da39a3ee5e6b4b0d3255bfef95601890afd80709"}`,
			expected: `{"hash": "<HASH>"}`,
		},
		{
			name:     "SHA-256 hash (64 chars)",
			input:    `{"hash": "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855"}`,
			expected: `{"hash": "<HASH>"}`,
		},
		{
			name:     "too short (31 chars)",
			input:    `{"hash": "d41d8cd98f00b204e9800998ecf842"}`,
			expected: `{"hash": "d41d8cd98f00b204e9800998ecf842"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tokenizer.Tokenize([]byte(tt.input))
			if string(result) != tt.expected {
				t.Errorf("got %s, want %s", result, tt.expected)
			}
		})
	}
}

func TestTokenizeBase64(t *testing.T) {
	tokenizer := NewTokenizer(TokenizerConfig{EnableBase64: true})

	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "base64 with padding",
			input:    `{"data": "SGVsbG8gV29ybGQhIFRoaXM="}`,
			expected: `{"data": "<B64>"}`,
		},
		{
			name:     "base64 without padding",
			input:    `{"data": "SGVsbG8gV29ybGQhIFRoaXMgaXMgYSB0ZXN0"}`,
			expected: `{"data": "<B64>"}`,
		},
		{
			name:     "too short (19 chars)",
			input:    `{"data": "SGVsbG8gV29ybGQhI"}`,
			expected: `{"data": "SGVsbG8gV29ybGQhI"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tokenizer.Tokenize([]byte(tt.input))
			if string(result) != tt.expected {
				t.Errorf("got %s, want %s", result, tt.expected)
			}
		})
	}
}

func TestTokenizeJWT(t *testing.T) {
	tokenizer := NewTokenizer(TokenizerConfig{EnableJWT: true})

	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "valid JWT",
			input:    `{"token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c"}`,
			expected: `{"token": "<JWT>"}`,
		},
		{
			name:     "not a JWT (no second segment)",
			input:    `{"token": "eyJhbGciOiJIUzI1NiJ9"}`,
			expected: `{"token": "eyJhbGciOiJIUzI1NiJ9"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tokenizer.Tokenize([]byte(tt.input))
			if string(result) != tt.expected {
				t.Errorf("got %s, want %s", result, tt.expected)
			}
		})
	}
}

func TestPatternPrecedence(t *testing.T) {
	// JWT should be matched before Base64
	tokenizer := NewTokenizer(TokenizerConfig{EnableJWT: true, EnableBase64: true})

	jwt := `eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIn0.dozjgNryP4J3jVmNHl0w5N_XgL0n3I9PlFUP0THsR8U`
	input := `{"token": "` + jwt + `"}`
	result := tokenizer.Tokenize([]byte(input))

	// Should be <JWT>, not <B64>
	expected := `{"token": "<JWT>"}`
	if string(result) != expected {
		t.Errorf("JWT precedence failed: got %s, want %s", result, expected)
	}
}

func TestDisabledPatterns(t *testing.T) {
	tokenizer := NewTokenizer(TokenizerConfig{
		EnableUUID:   false,
		EnableHash:   false,
		EnableBase64: false,
		EnableJWT:    false,
	})

	input := `{"id": "550e8400-e29b-41d4-a716-446655440000", "hash": "d41d8cd98f00b204e9800998ecf8427e"}`
	result := tokenizer.Tokenize([]byte(input))

	// Nothing should be replaced
	if string(result) != input {
		t.Errorf("Disabled patterns should not replace: got %s, want %s", result, input)
	}
}

func TestTokenizerStats(t *testing.T) {
	tokenizer := NewTokenizer(DefaultTokenizerConfig())

	input := `{
		"uuid": "550e8400-e29b-41d4-a716-446655440000",
		"hash": "d41d8cd98f00b204e9800998ecf8427e",
		"jwt": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIn0.dozjgNryP4J3jVmNHl0w5N_XgL0n3I9PlFUP0THsR8U"
	}`

	tokenizer.Tokenize([]byte(input))
	stats := tokenizer.Stats()

	if stats.UUIDCount != 1 {
		t.Errorf("UUIDCount: got %d, want 1", stats.UUIDCount)
	}
	if stats.HashCount != 1 {
		t.Errorf("HashCount: got %d, want 1", stats.HashCount)
	}
	if stats.JWTCount != 1 {
		t.Errorf("JWTCount: got %d, want 1", stats.JWTCount)
	}
	if stats.BytesSaved <= 0 {
		t.Errorf("BytesSaved should be positive, got %d", stats.BytesSaved)
	}
}

func TestGetTokenizer(t *testing.T) {
	cfg := TokenizerConfig{EnableUUID: true}

	t1 := GetTokenizer(cfg)
	t2 := GetTokenizer(cfg)

	// Should return the same instance
	if t1 != t2 {
		t.Error("GetTokenizer should return the same instance for the same config")
	}
}

func BenchmarkTokenize(b *testing.B) {
	tokenizer := NewTokenizer(DefaultTokenizerConfig())

	// Realistic log entry with UUIDs and hashes
	input := []byte(`{
		"timestamp": "2025-01-15T10:30:00Z",
		"level": "INFO",
		"service": "api-gateway",
		"request_id": "550e8400-e29b-41d4-a716-446655440000",
		"trace_id": "6ba7b810-9dad-11d1-80b4-00c04fd430c8",
		"user_id": "usr_d41d8cd98f00b204e9800998ecf8427e",
		"session_hash": "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855",
		"message": "Request processed successfully",
		"duration_ms": 42
	}`)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		tokenizer.Tokenize(input)
	}
}

func BenchmarkTokenizeLargePayload(b *testing.B) {
	tokenizer := NewTokenizer(DefaultTokenizerConfig())

	// Build a larger payload (~10KB)
	var payload []byte
	payload = append(payload, []byte(`{"events": [`)...)
	for i := 0; i < 100; i++ {
		if i > 0 {
			payload = append(payload, ',')
		}
		payload = append(payload, []byte(`{
			"id": "550e8400-e29b-41d4-a716-446655440000",
			"hash": "d41d8cd98f00b204e9800998ecf8427e",
			"data": "SGVsbG8gV29ybGQhIFRoaXMgaXMgYSB0ZXN0"
		}`)...)
	}
	payload = append(payload, []byte(`]}`)...)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		tokenizer.Tokenize(payload)
	}
}
