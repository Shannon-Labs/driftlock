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





