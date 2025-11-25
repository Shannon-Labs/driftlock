package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	
	"github.com/google/uuid"
)

// Add these methods to the store

func (s *store) tenantIDByEmail(ctx context.Context, email string) (uuid.UUID, error) {
	var id uuid.UUID
	err := s.pool.QueryRow(ctx, "SELECT id FROM tenants WHERE email = $1", email).Scan(&id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return uuid.Nil, fmt.Errorf("tenant not found")
		}
		return uuid.Nil, err
	}
	return id, nil
}

func (s *store) primaryAPIKey(ctx context.Context, tenantID uuid.UUID) (apiKeyRecord, error) {
	var k apiKeyRecord
	// Prefer 'default' key, or just the first active one
	query := `
		SELECT id, tenant_id, key_hash, stream_id, name, created_at, expires_at 
		FROM api_keys 
		WHERE tenant_id = $1 AND revoked_at IS NULL 
		ORDER BY created_at ASC 
		LIMIT 1`
		
	err := s.pool.QueryRow(ctx, query, tenantID).Scan(
		&k.ID, &k.TenantID, &k.KeyHash, &k.StreamID, &k.Name, &k.CreatedAt, &k.ExpiresAt,
	)
	if err != nil {
		return apiKeyRecord{}, err
	}
	return k, nil
}

func (s *store) listAPIKeys(ctx context.Context, tenantID uuid.UUID) ([]apiKeyRecord, error) {
	rows, err := s.pool.Query(ctx, `
		SELECT id, tenant_id, key_hash, stream_id, name, created_at, expires_at, revoked_at
		FROM api_keys
		WHERE tenant_id = $1
		ORDER BY created_at DESC`, tenantID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var keys []apiKeyRecord
	for rows.Next() {
		var k apiKeyRecord
		if err := rows.Scan(&k.ID, &k.TenantID, &k.KeyHash, &k.StreamID, &k.Name, &k.CreatedAt, &k.ExpiresAt, &k.RevokedAt); err != nil {
			return nil, err
		}
		keys = append(keys, k)
	}
	return keys, nil
}

// regenerateAPIKey revokes all existing keys and creates a new one
// Returns the plaintext key (only shown once)
func (s *store) regenerateAPIKey(ctx context.Context, tenantID uuid.UUID) (string, uuid.UUID, error) {
	// Get a stream ID for the new key (use default stream or first available)
	var streamID uuid.UUID
	err := s.pool.QueryRow(ctx, `SELECT id FROM streams WHERE tenant_id = $1 AND slug = 'default'`, tenantID).Scan(&streamID)
	if err != nil {
		// Fallback to any stream
		err = s.pool.QueryRow(ctx, `SELECT id FROM streams WHERE tenant_id = $1 LIMIT 1`, tenantID).Scan(&streamID)
		if err != nil {
			return "", uuid.Nil, fmt.Errorf("tenant has no streams")
		}
	}

	// Generate new key
	key, keyID, err := generateAPIKey()
	if err != nil {
		return "", uuid.Nil, err
	}
	hashedKey, err := hashAPIKey(key)
	if err != nil {
		return "", uuid.Nil, err
	}

	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return "", uuid.Nil, err
	}
	defer tx.Rollback(ctx)

	// Revoke all existing keys (soft delete)
	_, err = tx.Exec(ctx, `UPDATE api_keys SET revoked_at = NOW() WHERE tenant_id = $1 AND revoked_at IS NULL`, tenantID)
	if err != nil {
		return "", uuid.Nil, err
	}

	// Create new key
	_, err = tx.Exec(ctx, `INSERT INTO api_keys (id, tenant_id, name, role, key_hash, stream_id, rate_limit_rps)
		VALUES ($1, $2, $3, $4, $5, $6, $7)`,
		keyID, tenantID, "regenerated-key", "admin", hashedKey, streamID, 60)
	if err != nil {
		return "", uuid.Nil, err
	}

	if err := tx.Commit(ctx); err != nil {
		return "", uuid.Nil, err
	}

	// Refresh cache
	_ = s.loadCache(ctx)

	return key, keyID, nil
}

// createAPIKey creates a new API key for a tenant without revoking existing keys
// Returns the plaintext key (only shown once)
func (s *store) createAPIKey(ctx context.Context, tenantID uuid.UUID, name, role string) (string, uuid.UUID, error) {
	// Validate role
	if role != "admin" && role != "stream" {
		role = "admin"
	}

	// Get a stream ID for the new key (use default stream or first available)
	var streamID uuid.UUID
	err := s.pool.QueryRow(ctx, `SELECT id FROM streams WHERE tenant_id = $1 AND slug = 'default'`, tenantID).Scan(&streamID)
	if err != nil {
		// Fallback to any stream
		err = s.pool.QueryRow(ctx, `SELECT id FROM streams WHERE tenant_id = $1 LIMIT 1`, tenantID).Scan(&streamID)
		if err != nil {
			return "", uuid.Nil, fmt.Errorf("tenant has no streams")
		}
	}

	// Generate new key
	key, keyID, err := generateAPIKey()
	if err != nil {
		return "", uuid.Nil, err
	}
	hashedKey, err := hashAPIKey(key)
	if err != nil {
		return "", uuid.Nil, err
	}

	// Create new key
	_, err = s.pool.Exec(ctx, `INSERT INTO api_keys (id, tenant_id, name, role, key_hash, stream_id, rate_limit_rps)
		VALUES ($1, $2, $3, $4, $5, $6, $7)`,
		keyID, tenantID, name, role, hashedKey, streamID, 60)
	if err != nil {
		return "", uuid.Nil, err
	}

	// Refresh cache
	_ = s.loadCache(ctx)

	return key, keyID, nil
}

// softRevokeKey marks a key as revoked without deleting it
func (s *store) softRevokeKey(ctx context.Context, keyID, tenantID uuid.UUID) error {
	result, err := s.pool.Exec(ctx, `UPDATE api_keys SET revoked_at = NOW() WHERE id = $1 AND tenant_id = $2 AND revoked_at IS NULL`, keyID, tenantID)
	if err != nil {
		return err
	}
	if result.RowsAffected() == 0 {
		return fmt.Errorf("key not found or already revoked")
	}
	return nil
}





