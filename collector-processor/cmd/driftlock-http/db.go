package main

import (
	"bytes"
	"context"
	"crypto/rand"
	"crypto/sha256"
	"database/sql"
	"encoding/base32"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/go-redis/redis/v8"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
	"golang.org/x/crypto/argon2"
)

var base32Enc = base32.StdEncoding.WithPadding(base32.NoPadding)
var errNotFound = errors.New("record not found")

type store struct {
	pool          *pgxpool.Pool
	cache         *configCache
	baselineTTL   time.Duration
	baselineRedis *redis.Client
}

func newStore(pool *pgxpool.Pool) *store {
	return &store{
		pool:  pool,
		cache: newConfigCache(),
	}
}

func connectDB(ctx context.Context) (*pgxpool.Pool, error) {
	dbURL := env("DATABASE_URL", "")
	if dbURL == "" {
		return nil, fmt.Errorf("DATABASE_URL not set")
	}
	cfg, err := pgxpool.ParseConfig(dbURL)
	if err != nil {
		return nil, err
	}
	// Pool sizing: defaults increased for production load
	// DB_MAX_CONNS: Maximum connections (default 25, was 10)
	// DB_MIN_CONNS: Minimum idle connections to keep warm (default 5)
	if cfg.MaxConns == 0 {
		cfg.MaxConns = int32(envInt("DB_MAX_CONNS", 25))
	}
	cfg.MinConns = int32(envInt("DB_MIN_CONNS", 5))
	pool, err := pgxpool.NewWithConfig(ctx, cfg)
	if err != nil {
		return nil, err
	}
	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, err
	}
	return pool, nil
}

func (s *store) Close() {
	if s.pool != nil {
		s.pool.Close()
	}
	if s.baselineRedis != nil {
		_ = s.baselineRedis.Close()
	}
}

func (s *store) checkTenantFirebaseUID(ctx context.Context, firebaseUID string) (bool, error) {
	var exists bool
	err := s.pool.QueryRow(ctx, "SELECT EXISTS(SELECT 1 FROM tenants WHERE firebase_uid = $1)", firebaseUID).Scan(&exists)
	if err != nil {
		return false, err
	}
	return exists, nil
}

func (s *store) checkTenantEmail(ctx context.Context, email string) (bool, error) {
	var exists bool
	err := s.pool.QueryRow(ctx, `SELECT EXISTS(SELECT 1 FROM tenants WHERE email = $1)`, email).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("failed to check email: %w", err)
	}
	return exists, nil
}

func runMigrations(ctx context.Context, action string) error {
	dbURL := env("DATABASE_URL", "")
	if dbURL == "" {
		return fmt.Errorf("DATABASE_URL not set")
	}
	goose.SetBaseFS(nil)
	goose.SetLogger(goose.NopLogger())
	db, err := sql.Open("pgx", dbURL)
	if err != nil {
		return err
	}
	defer db.Close()
	path := migrationsPath()

	switch action {
	case "up":
		return goose.UpContext(ctx, db, path)
	case "status":
		return goose.Status(db, path)
	default:
		return fmt.Errorf("unknown migrate command %q", action)
	}
}

func (s *store) loadCache(ctx context.Context) error {
	rows, err := s.pool.Query(ctx, `
        WITH latest_config AS (
            SELECT DISTINCT ON (stream_id) stream_id, version, config
            FROM stream_configs
            ORDER BY stream_id, version DESC
        )
               SELECT t.id, t.slug, t.name, t.plan, t.default_compressor, t.rate_limit_rps,
               t.stripe_customer_id, t.stripe_subscription_id, t.stripe_status,
               s.id, s.slug, s.type, s.seed, s.compressor, s.retention_days,
               COALESCE(lc.config, '{}'::jsonb)
        FROM tenants t
        JOIN streams s ON s.tenant_id = t.id
        LEFT JOIN latest_config lc ON lc.stream_id = s.id`)
	if err != nil {
		return err
	}
	defer rows.Close()

	entries := make([]struct {
		tenant tenantRecord
		stream streamRecord
		cfg    streamSettings
	}, 0)
	for rows.Next() {
		var (
			tenant tenantRecord
			stream streamRecord
			cfgRaw []byte
		)
		if err := rows.Scan(&tenant.ID, &tenant.Slug, &tenant.Name, &tenant.Plan, &tenant.DefaultCompressor, &tenant.RateLimitRPS,
			&tenant.StripeCustomerID, &tenant.StripeSubscriptionID, &tenant.StripeStatus,
			&stream.ID, &stream.Slug, &stream.Type, &stream.Seed, &stream.Compressor, &stream.RetentionDays, &cfgRaw); err != nil {
			return err
		}
		stream.TenantID = tenant.ID
		var cfg streamSettings
		if len(cfgRaw) > 0 {
			_ = json.Unmarshal(cfgRaw, &cfg)
		}
		cfg.applyDefaults()
		entries = append(entries, struct {
			tenant tenantRecord
			stream streamRecord
			cfg    streamSettings
		}{tenant: tenant, stream: stream, cfg: cfg})
	}
	if rows.Err() != nil {
		return rows.Err()
	}

	s.cache.replace(entries)
	return nil
}

// baselineKey returns the Redis key for a stream's persisted baseline buffer.
func baselineKey(streamID uuid.UUID) string {
	return fmt.Sprintf("cbad:baseline:%s", streamID.String())
}

// loadBaselineSnapshot fetches a persisted baseline from Redis as a slice of events.
// Returns (nil, nil) when Redis is disabled or no snapshot exists.
func (s *store) loadBaselineSnapshot(ctx context.Context, streamID uuid.UUID) ([]json.RawMessage, error) {
	if s == nil || s.baselineRedis == nil {
		return nil, nil
	}
	data, err := s.baselineRedis.Get(ctx, baselineKey(streamID)).Bytes()
	if err == redis.Nil {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	lines := bytes.Split(data, []byte{'\n'})
	events := make([]json.RawMessage, 0, len(lines))
	for _, line := range lines {
		line = bytes.TrimSpace(line)
		if len(line) == 0 {
			continue
		}
		// Skip invalid JSON to avoid poisoning the detector
		if !json.Valid(line) {
			continue
		}
		copyLine := append([]byte(nil), line...)
		events = append(events, json.RawMessage(copyLine))
	}
	return events, nil
}

// saveBaselineSnapshot stores the trailing "limit" events in Redis with TTL for reuse across requests.
func (s *store) saveBaselineSnapshot(ctx context.Context, streamID uuid.UUID, events []json.RawMessage, limit int) error {
	if s == nil || s.baselineRedis == nil || limit <= 0 || len(events) == 0 {
		return nil
	}
	if len(events) > limit {
		events = events[len(events)-limit:]
	}

	var buf bytes.Buffer
	for i, ev := range events {
		compact, err := compactEvent(ev, nil) // Store compacted JSON to keep payload small
		if err != nil {
			// Skip invalid entries instead of failing persistence
			continue
		}
		buf.Write(compact)
		if i < len(events)-1 {
			buf.WriteByte('\n')
		}
	}

	if buf.Len() == 0 {
		return nil
	}

	ttl := s.baselineTTL
	if ttl <= 0 {
		ttl = 24 * time.Hour
	}
	return s.baselineRedis.Set(ctx, baselineKey(streamID), buf.Bytes(), ttl).Err()
}

func (s *store) tenantByID(id uuid.UUID) (tenantRecord, bool) {
	return s.cache.tenant(id)
}

func (s *store) streamBySlugOrID(tenantID uuid.UUID, value string) (streamRecord, streamSettings, bool) {
	if id, err := uuid.Parse(value); err == nil {
		return s.cache.streamByID(id)
	}
	return s.cache.streamBySlug(tenantID, value)
}

func (s *store) streamByID(id uuid.UUID) (streamRecord, streamSettings, bool) {
	return s.cache.streamByID(id)
}

func (s *store) defaultStream(tenantID uuid.UUID) (streamRecord, streamSettings, bool) {
	return s.cache.defaultStream(tenantID)
}

func (s *store) createTenantWithKey(ctx context.Context, params tenantCreateParams) (*tenantCreateResult, error) {
	if params.Name == "" {
		return nil, fmt.Errorf("tenant name required")
	}
	role := strings.ToLower(params.KeyRole)
	if role != "admin" && role != "stream" {
		return nil, fmt.Errorf("invalid key role %q", params.KeyRole)
	}

	slug := params.Slug
	if slug == "" {
		slug = slugify(params.Name)
	}
	streamSlug := params.StreamSlug
	if streamSlug == "" {
		streamSlug = "default"
	}

	tenantID := uuid.New()
	streamID := uuid.New()
	cfg := streamSettings{
		BaselineSize:     params.DefaultBaseline,
		WindowSize:       params.DefaultWindow,
		HopSize:          params.DefaultHop,
		NCDThreshold:     params.NCDThreshold,
		PValueThreshold:  params.PValueThreshold,
		PermutationCount: params.PermutationCount,
		Compressor:       params.DefaultCompressor,
	}

	key, keyID, err := generateAPIKey()
	if err != nil {
		return nil, err
	}
	hashedKey, err := hashAPIKey(key)
	if err != nil {
		return nil, err
	}

	tx, err := s.pool.BeginTx(ctx, pgx.TxOptions{IsoLevel: pgx.Serializable})
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)

	now := time.Now()
	status := params.Status
	if status == "" {
		status = "active"
	}
	_, err = tx.Exec(ctx, `INSERT INTO tenants (id, name, slug, plan, default_compressor, rate_limit_rps, email, signup_ip, signup_source, created_at, verification_token, status) 
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12)`,
		tenantID, params.Name, slug, params.Plan, params.DefaultCompressor, params.TenantRateLimit, params.Email, params.SignupIP, params.Source, now, params.VerificationToken, status)
	if err != nil {
		return nil, err
	}

	_, err = tx.Exec(ctx, `INSERT INTO streams (id, tenant_id, slug, type, description, seed, compressor, queue_mode, retention_days)
        VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9)`,
		streamID, tenantID, streamSlug, params.StreamType, params.StreamDescription, params.Seed, params.DefaultCompressor, "memory", params.StreamRetentionDays)
	if err != nil {
		return nil, err
	}

	cfgJSON, _ := json.Marshal(cfg)
	_, err = tx.Exec(ctx, `INSERT INTO stream_configs (id, stream_id, version, config, created_by) VALUES ($1,$2,$3,$4,$5)`,
		uuid.New(), streamID, 1, cfgJSON, "cli")
	if err != nil {
		return nil, err
	}

	_, err = tx.Exec(ctx, `INSERT INTO api_keys (id, tenant_id, name, role, key_hash, stream_id, rate_limit_rps) VALUES ($1,$2,$3,$4,$5,$6,$7)`,
		keyID, tenantID, params.KeyName, role, hashedKey, streamID, params.KeyRateLimit)
	if err != nil {
		return nil, err
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, err
	}

	// refresh cache
	_ = s.loadCache(ctx)

	result := &tenantCreateResult{
		TenantID:   tenantID,
		TenantSlug: slug,
		StreamID:   streamID,
		StreamSlug: streamSlug,
		APIKey:     key,
		APIKeyID:   keyID,
	}
	result.Tenant.ID = tenantID
	result.Tenant.Name = params.Name
	result.Tenant.Slug = slug
	result.Tenant.Plan = params.Plan
	result.Tenant.CreatedAt = now

	return result, nil
}

func (s *store) createPendingTenant(ctx context.Context, params tenantCreateParams) (*tenantCreateResult, error) {
	if params.Name == "" {
		return nil, fmt.Errorf("tenant name required")
	}

	slug := params.Slug
	if slug == "" {
		slug = slugify(params.Name)
	}
	streamSlug := params.StreamSlug
	if streamSlug == "" {
		streamSlug = "default"
	}

	tenantID := uuid.New()
	streamID := uuid.New()
	cfg := streamSettings{
		BaselineSize:     params.DefaultBaseline,
		WindowSize:       params.DefaultWindow,
		HopSize:          params.DefaultHop,
		NCDThreshold:     params.NCDThreshold,
		PValueThreshold:  params.PValueThreshold,
		PermutationCount: params.PermutationCount,
		Compressor:       params.DefaultCompressor,
	}

	tx, err := s.pool.BeginTx(ctx, pgx.TxOptions{IsoLevel: pgx.Serializable})
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)

	now := time.Now()
	status := "pending_verification"
	// Verification tokens expire after 15 minutes for security
	tokenExpiresAt := now.Add(15 * time.Minute)

	// Insert tenant (firebase_uid is optional, use NULL if empty)
	var firebaseUID *string
	if params.FirebaseUID != "" {
		firebaseUID = &params.FirebaseUID
	}
	_, err = tx.Exec(ctx, `INSERT INTO tenants (id, name, slug, plan, default_compressor, rate_limit_rps, email, signup_ip, signup_source, created_at, verification_token, verification_token_expires_at, status, firebase_uid)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14)`,
		tenantID, params.Name, slug, params.Plan, params.DefaultCompressor, params.TenantRateLimit, params.Email, params.SignupIP, params.Source, now, params.VerificationToken, tokenExpiresAt, status, firebaseUID)
	if err != nil {
		return nil, err
	}

	_, err = tx.Exec(ctx, `INSERT INTO streams (id, tenant_id, slug, type, description, seed, compressor, queue_mode, retention_days)
        VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9)`,
		streamID, tenantID, streamSlug, params.StreamType, params.StreamDescription, params.Seed, params.DefaultCompressor, "memory", params.StreamRetentionDays)
	if err != nil {
		return nil, err
	}

	cfgJSON, _ := json.Marshal(cfg)
	_, err = tx.Exec(ctx, `INSERT INTO stream_configs (id, stream_id, version, config, created_by) VALUES ($1,$2,$3,$4,$5)`,
		uuid.New(), streamID, 1, cfgJSON, "onboarding")
	if err != nil {
		return nil, err
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, err
	}

	// refresh cache
	_ = s.loadCache(ctx)

	result := &tenantCreateResult{
		TenantID:   tenantID,
		TenantSlug: slug,
		StreamID:   streamID,
		StreamSlug: streamSlug,
		// No API Key for pending tenant
	}
	result.Tenant.ID = tenantID
	result.Tenant.Name = params.Name
	result.Tenant.Slug = slug
	result.Tenant.Plan = params.Plan
	result.Tenant.CreatedAt = now

	return result, nil
}

var errTokenExpired = errors.New("verification token has expired")

func (s *store) verifyAndActivateTenant(ctx context.Context, token string) (*tenantCreateResult, error) {
	var tenantID uuid.UUID
	var tenantName, tenantSlug, tenantPlan, tenantCompressor string
	var tenantRateLimit int
	var createdAt time.Time
	var tokenExpiresAt sql.NullTime

	// 1. Find pending tenant by token
	err := s.pool.QueryRow(ctx, `SELECT id, name, slug, plan, default_compressor, rate_limit_rps, created_at, verification_token_expires_at
		FROM tenants WHERE verification_token = $1 AND status = 'pending_verification'`, token).Scan(
		&tenantID, &tenantName, &tenantSlug, &tenantPlan, &tenantCompressor, &tenantRateLimit, &createdAt, &tokenExpiresAt)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("invalid or expired verification token")
		}
		return nil, err
	}

	// Check if token has expired (15-minute TTL)
	if tokenExpiresAt.Valid && time.Now().After(tokenExpiresAt.Time) {
		return nil, errTokenExpired
	}

	// 2. Generate API Key
	key, keyID, err := generateAPIKey()
	if err != nil {
		return nil, err
	}
	hashedKey, err := hashAPIKey(key)
	if err != nil {
		return nil, err
	}

	// 3. Find default stream ID (needed for API key)
	var streamID uuid.UUID
	var streamSlug string
	err = s.pool.QueryRow(ctx, `SELECT id, slug FROM streams WHERE tenant_id = $1 AND slug = 'default'`, tenantID).Scan(&streamID, &streamSlug)
	if err != nil {
		// Fallback to any stream if default doesn't exist (shouldn't happen)
		err = s.pool.QueryRow(ctx, `SELECT id, slug FROM streams WHERE tenant_id = $1 LIMIT 1`, tenantID).Scan(&streamID, &streamSlug)
		if err != nil {
			return nil, fmt.Errorf("tenant has no streams")
		}
	}

	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)

	// 4. Update status
	_, err = tx.Exec(ctx, `UPDATE tenants SET status = 'active', verified_at = NOW(), verification_token = NULL WHERE id = $1`, tenantID)
	if err != nil {
		return nil, err
	}

	// 5. Insert API Key
	// Defaults for onboarding key
	keyName := "onboarding-key"
	keyRole := "admin"
	// Use tenant rate limit or default
	keyRateLimit := tenantRateLimit
	if keyRateLimit == 0 {
		keyRateLimit = 60
	}

	_, err = tx.Exec(ctx, `INSERT INTO api_keys (id, tenant_id, name, role, key_hash, stream_id, rate_limit_rps) VALUES ($1,$2,$3,$4,$5,$6,$7)`,
		keyID, tenantID, keyName, keyRole, hashedKey, streamID, keyRateLimit)
	if err != nil {
		return nil, err
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, err
	}

	// Refresh cache
	_ = s.loadCache(ctx)

	result := &tenantCreateResult{
		TenantID:   tenantID,
		TenantSlug: tenantSlug,
		StreamID:   streamID,
		StreamSlug: streamSlug,
		APIKey:     key,
		APIKeyID:   keyID,
	}
	result.Tenant.ID = tenantID
	result.Tenant.Name = tenantName
	result.Tenant.Slug = tenantSlug
	result.Tenant.Plan = tenantPlan
	result.Tenant.CreatedAt = createdAt

	return result, nil
}

func (s *store) listKeys(ctx context.Context, tenantSlug string) ([]apiKeyInfo, error) {
	rows, err := s.pool.Query(ctx, `SELECT ak.id, ak.name, ak.role, ak.stream_id, ak.created_at, ak.last_used_at
        FROM api_keys ak JOIN tenants t ON ak.tenant_id = t.id WHERE t.slug=$1 ORDER BY ak.created_at ASC`, tenantSlug)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []apiKeyInfo
	for rows.Next() {
		var info apiKeyInfo
		var streamID *uuid.UUID
		var lastUsed *time.Time
		err := rows.Scan(&info.ID, &info.Name, &info.Role, &streamID, &info.CreatedAt, &lastUsed)
		if err != nil {
			return nil, err
		}
		if streamID != nil {
			tmp := *streamID
			info.StreamID = &tmp
		}
		if lastUsed != nil {
			info.LastUsedAt = lastUsed
		}
		out = append(out, info)
	}
	return out, rows.Err()
}

func (s *store) revokeKey(ctx context.Context, keyID uuid.UUID) (int64, error) {
	// Soft delete: set revoked_at instead of deleting
	// This allows for audit trail and prevents key reuse
	tag, err := s.pool.Exec(ctx, `UPDATE api_keys SET revoked_at = NOW() WHERE id = $1 AND revoked_at IS NULL`, keyID)
	if err != nil {
		return 0, err
	}
	return tag.RowsAffected(), nil
}

func (s *store) resolveAPIKey(ctx context.Context, id uuid.UUID) (*apiKeyRecord, error) {
	// Exclude revoked keys - they should not authenticate
	row := s.pool.QueryRow(ctx, `SELECT ak.id, ak.tenant_id, ak.role, ak.key_hash, ak.stream_id, ak.name, ak.rate_limit_rps,
        t.name, t.slug, t.default_compressor, t.rate_limit_rps
        FROM api_keys ak JOIN tenants t ON ak.tenant_id = t.id WHERE ak.id=$1 AND ak.revoked_at IS NULL`, id)
	var rec apiKeyRecord
	var streamID *uuid.UUID
	err := row.Scan(&rec.ID, &rec.TenantID, &rec.Role, &rec.KeyHash, &streamID, &rec.Name, &rec.KeyRateLimit,
		&rec.TenantName, &rec.TenantSlug, &rec.DefaultCompressor, &rec.TenantRateLimit)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("api key not found")
		}
		return nil, err
	}
	if streamID != nil {
		rec.StreamID = *streamID
	}
	return &rec, nil
}

func (s *store) insertBatch(ctx context.Context, tenantID, streamID uuid.UUID, batchHash string, worker string) (uuid.UUID, error) {
	batchID := uuid.New()
	_, err := s.pool.Exec(ctx, `INSERT INTO ingest_batches (id, tenant_id, stream_id, batch_hash, status, worker)
        VALUES ($1,$2,$3,$4,$5,$6)`, batchID, tenantID, streamID, batchHash, "completed", worker)
	if err != nil {
		return uuid.Nil, err
	}
	return batchID, nil
}

func (s *store) insertAnomalies(ctx context.Context, batchID uuid.UUID, tenantID uuid.UUID, streamID uuid.UUID, records []persistedAnomaly) ([]uuid.UUID, error) {
	if len(records) == 0 {
		return nil, nil
	}
	ids := make([]uuid.UUID, 0, len(records))
	batch := &pgx.Batch{}
	for _, rec := range records {
		id := uuid.New()
		ids = append(ids, id)
		uri := rec.EvidenceURI
		if uri == "" {
			uri = fmt.Sprintf("local://evidence/%s.md", id)
		}
		checksum := rec.EvidenceChecksum
		if checksum == "" && len(rec.Details) > 0 {
			sum := sha256.Sum256(rec.Details)
			checksum = fmt.Sprintf("%x", sum[:])
		}
		batch.Queue(`INSERT INTO anomalies (id, tenant_id, stream_id, ingest_batch_id, ncd, compression_ratio, entropy_change, p_value, confidence, explanation, status, detected_at, details)
            VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13)`,
			id, tenantID, streamID, batchID, rec.NCD, rec.CompressionRatio, rec.EntropyChange, rec.PValue, rec.Confidence, rec.Explanation, "new", time.Now().UTC(), rec.Details)
		batch.Queue(`INSERT INTO anomaly_evidence (id, anomaly_id, format, uri, checksum) VALUES ($1,$2,$3,$4,$5)`,
			uuid.New(), id, rec.EvidenceFormat, uri, checksum)
	}
	br := s.pool.SendBatch(ctx, batch)
	defer br.Close()
	for range records {
		if _, err := br.Exec(); err != nil {
			return nil, err
		}
		if _, err := br.Exec(); err != nil {
			return nil, err
		}
	}
	return ids, nil
}

func (s *store) incrementUsage(ctx context.Context, tenantID, streamID uuid.UUID, eventCount, requestCount, anomalyCount int) error {
	_, err := s.pool.Exec(ctx, `
		INSERT INTO usage_metrics (tenant_id, stream_id, date, event_count, api_request_count, anomaly_count)
		VALUES ($1, $2, CURRENT_DATE, $3, $4, $5)
		ON CONFLICT (tenant_id, stream_id, date)
		DO UPDATE SET
			event_count = usage_metrics.event_count + EXCLUDED.event_count,
			api_request_count = usage_metrics.api_request_count + EXCLUDED.api_request_count,
			anomaly_count = usage_metrics.anomaly_count + EXCLUDED.anomaly_count`,
		tenantID, streamID, eventCount, requestCount, anomalyCount)
	return err
}

func (s *store) getMonthlyUsage(ctx context.Context, tenantID uuid.UUID) (int64, error) {
	// Sum events for current calendar month
	var total int64
	err := s.pool.QueryRow(ctx, `
		SELECT COALESCE(SUM(event_count), 0)
		FROM usage_metrics
		WHERE tenant_id = $1 AND date_trunc('month', date) = date_trunc('month', CURRENT_DATE)`,
		tenantID).Scan(&total)
	return total, err
}

type dailyUsageRecord struct {
	Date         time.Time `json:"date"`
	EventCount   int64     `json:"event_count"`
	RequestCount int64     `json:"request_count"`
	AnomalyCount int64     `json:"anomaly_count"`
}

func (s *store) getDailyUsage(ctx context.Context, tenantID uuid.UUID, days int) ([]dailyUsageRecord, error) {
	rows, err := s.pool.Query(ctx, `
		WITH date_series AS (
			SELECT generate_series(
				CURRENT_DATE - ($2 - 1) * INTERVAL '1 day',
				CURRENT_DATE,
				'1 day'::interval
			)::date AS date
		)
		SELECT
			ds.date,
			COALESCE(SUM(um.event_count), 0) as event_count,
			COALESCE(SUM(um.api_request_count), 0) as request_count,
			COALESCE(SUM(um.anomaly_count), 0) as anomaly_count
		FROM date_series ds
		LEFT JOIN usage_metrics um ON um.date = ds.date AND um.tenant_id = $1
		GROUP BY ds.date
		ORDER BY ds.date ASC`,
		tenantID, days)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []dailyUsageRecord
	for rows.Next() {
		var rec dailyUsageRecord
		if err := rows.Scan(&rec.Date, &rec.EventCount, &rec.RequestCount, &rec.AnomalyCount); err != nil {
			return nil, err
		}
		results = append(results, rec)
	}
	return results, rows.Err()
}

type streamUsageRecord struct {
	StreamID     uuid.UUID `json:"stream_id"`
	StreamName   string    `json:"stream_name"`
	EventCount   int64     `json:"event_count"`
	RequestCount int64     `json:"request_count"`
	AnomalyCount int64     `json:"anomaly_count"`
}

func (s *store) getStreamUsage(ctx context.Context, tenantID uuid.UUID) ([]streamUsageRecord, error) {
	rows, err := s.pool.Query(ctx, `
		SELECT
			um.stream_id,
			COALESCE(st.name, 'default') as stream_name,
			COALESCE(SUM(um.event_count), 0) as event_count,
			COALESCE(SUM(um.api_request_count), 0) as request_count,
			COALESCE(SUM(um.anomaly_count), 0) as anomaly_count
		FROM usage_metrics um
		LEFT JOIN streams st ON st.id = um.stream_id
		WHERE um.tenant_id = $1
		  AND date_trunc('month', um.date) = date_trunc('month', CURRENT_DATE)
		GROUP BY um.stream_id, st.name
		ORDER BY event_count DESC`,
		tenantID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []streamUsageRecord
	for rows.Next() {
		var rec streamUsageRecord
		if err := rows.Scan(&rec.StreamID, &rec.StreamName, &rec.EventCount, &rec.RequestCount, &rec.AnomalyCount); err != nil {
			return nil, err
		}
		results = append(results, rec)
	}
	return results, rows.Err()
}

func (s *store) fetchAnomaly(ctx context.Context, tenantID uuid.UUID, anomalyID uuid.UUID) (anomalyDetailRecord, []anomalyEvidenceRecord, error) {
	var rec anomalyDetailRecord
	query := `SELECT id, stream_id, ingest_batch_id, ncd, compression_ratio, entropy_change, p_value, confidence, explanation, status, detected_at, details, baseline_snapshot, window_snapshot
        FROM anomalies WHERE tenant_id=$1 AND id=$2`
	row := s.pool.QueryRow(ctx, query, tenantID, anomalyID)
	var (
		dbID     uuid.UUID
		streamID uuid.UUID
		batchID  *uuid.UUID
		details  []byte
		baseline []byte
		window   []byte
	)
	if err := row.Scan(&dbID, &streamID, &batchID, &rec.NCD, &rec.CompressionRatio, &rec.EntropyChange, &rec.PValue, &rec.Confidence, &rec.Explanation, &rec.Status, &rec.DetectedAt, &details, &baseline, &window); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return anomalyDetailRecord{}, nil, errNotFound
		}
		return anomalyDetailRecord{}, nil, err
	}
	rec.ID = dbID.String()
	rec.StreamID = streamID.String()
	if batchID != nil {
		rec.BatchID = batchID.String()
	}
	rec.Details = details
	rec.BaselineSnapshot = baseline
	rec.WindowSnapshot = window
	rec.Explanation = strings.TrimSpace(rec.Explanation)

	rows, err := s.pool.Query(ctx, `SELECT format, uri, checksum, COALESCE(size_bytes,0), created_at FROM anomaly_evidence WHERE anomaly_id=$1 ORDER BY created_at ASC`, anomalyID)
	if err != nil {
		return anomalyDetailRecord{}, nil, err
	}
	defer rows.Close()
	var evidence []anomalyEvidenceRecord
	for rows.Next() {
		var ev anomalyEvidenceRecord
		if err := rows.Scan(&ev.Format, &ev.URI, &ev.Checksum, &ev.SizeBytes, &ev.CreatedAt); err != nil {
			return anomalyDetailRecord{}, nil, err
		}
		evidence = append(evidence, ev)
	}
	if rows.Err() != nil {
		return anomalyDetailRecord{}, nil, rows.Err()
	}
	return rec, evidence, nil
}

func (s *store) updateAnomalyExplanation(ctx context.Context, anomalyID string, explanation string) error {
	_, err := s.pool.Exec(ctx, `UPDATE anomalies SET explanation = $1, status = 'analyzed' WHERE id = $2`, explanation, anomalyID)
	return err
}

func (s *store) createExportJob(ctx context.Context, tenantID uuid.UUID, format string, filters, delivery []byte) (uuid.UUID, error) {
	if len(filters) == 0 {
		filters = []byte("{}")
	}
	if len(delivery) == 0 {
		delivery = []byte(`{"type":"inline"}`)
	}
	jobID := uuid.New()
	_, err := s.pool.Exec(ctx, `INSERT INTO export_jobs (id, tenant_id, format, filters, delivery, status)
        VALUES ($1,$2,$3,$4,$5,$6)`, jobID, tenantID, strings.ToLower(format), filters, delivery, "pending")
	if err != nil {
		return uuid.Nil, err
	}
	return jobID, nil
}

type persistedAnomaly struct {
	NCD              float64
	CompressionRatio float64
	EntropyChange    float64
	PValue           float64
	Confidence       float64
	Explanation      string
	Details          []byte
	EvidenceURI      string
	EvidenceFormat   string
	EvidenceChecksum string
}

type anomalyDetailRecord struct {
	ID               string
	StreamID         string
	BatchID          string
	Status           string
	DetectedAt       time.Time
	Explanation      string
	NCD              float64
	CompressionRatio float64
	EntropyChange    float64
	PValue           float64
	Confidence       float64
	Details          []byte
	BaselineSnapshot []byte
	WindowSnapshot   []byte
}

type anomalyEvidenceRecord struct {
	Format    string
	URI       string
	Checksum  string
	SizeBytes int64
	CreatedAt time.Time
}

// Credential helpers

func generateAPIKey() (string, uuid.UUID, error) {
	id := uuid.New()
	raw := make([]byte, 20)
	if _, err := rand.Read(raw); err != nil {
		return "", uuid.Nil, err
	}
	secret := strings.ToLower(base32Enc.EncodeToString(raw))
	key := fmt.Sprintf("dlk_%s.%s", id.String(), secret)
	return key, id, nil
}

func hashAPIKey(key string) (string, error) {
	salt := make([]byte, 16)
	if _, err := rand.Read(salt); err != nil {
		return "", err
	}
	hash := argon2.IDKey([]byte(key), salt, 1, 64*1024, 1, 32)
	return fmt.Sprintf("argon2id$v=19$m=65536,t=1,p=1$%s$%s", base64.RawStdEncoding.EncodeToString(salt), base64.RawStdEncoding.EncodeToString(hash)), nil
}

func verifyAPIKey(hash, candidate string) bool {
	parts := strings.Split(hash, "$")
	if len(parts) != 5 {
		return false
	}
	salt, err := base64.RawStdEncoding.DecodeString(parts[3])
	if err != nil {
		return false
	}
	expected, err := base64.RawStdEncoding.DecodeString(parts[4])
	if err != nil {
		return false
	}
	derived := argon2.IDKey([]byte(candidate), salt, 1, 64*1024, 1, uint32(len(expected)))
	return subtleCompare(expected, derived)
}

func subtleCompare(a, b []byte) bool {
	if len(a) != len(b) {
		return false
	}
	var res byte
	for i := range a {
		res |= a[i] ^ b[i]
	}
	return res == 0
}

// Config cache structures

type tenantRecord struct {
	ID                   uuid.UUID
	Name                 string
	Slug                 string
	Plan                 string
	DefaultCompressor    string
	RateLimitRPS         int
	StripeCustomerID     *string
	StripeSubscriptionID *string
	StripeStatus         *string
	Email                string
}

type streamRecord struct {
	ID            uuid.UUID
	TenantID      uuid.UUID
	Slug          string
	Type          string
	Seed          int64
	Compressor    string
	RetentionDays int
	// Cold start calibration fields (SHA-139)
	EventsIngested  int64
	IsCalibrated    bool
	MinBaselineSize int
}

func (sr streamRecord) cacheKey() string {
	return sr.TenantID.String() + ":" + sr.Slug
}

type streamSettings struct {
	BaselineSize     int     `json:"baseline_size"`
	WindowSize       int     `json:"window_size"`
	HopSize          int     `json:"hop_size"`
	NCDThreshold     float64 `json:"ncd_threshold"`
	PValueThreshold  float64 `json:"p_value_threshold"`
	PermutationCount int     `json:"permutation_count"`
	Compressor       string  `json:"compressor"`
	// SHA-141: Tokenizer settings
	TokenizerEnabled bool `json:"tokenizer_enabled"`
	TokenizeUUID     bool `json:"tokenize_uuid"`
	TokenizeHash     bool `json:"tokenize_hash"`
	TokenizeBase64   bool `json:"tokenize_base64"`
	TokenizeJWT      bool `json:"tokenize_jwt"`
	// SHA-143: Numeric outlier detection
	NumericOutlierEnabled bool    `json:"numeric_outlier_enabled"`
	NumericKSigma         float64 `json:"numeric_k_sigma"`
}

func (s *streamSettings) applyDefaults() {
	if s.BaselineSize <= 0 {
		s.BaselineSize = 400
	}
	if s.WindowSize <= 0 {
		s.WindowSize = 50
	}
	if s.HopSize <= 0 {
		s.HopSize = 10
	}
	if s.NCDThreshold == 0 {
		s.NCDThreshold = 0.3
	}
	if s.PValueThreshold == 0 {
		s.PValueThreshold = 0.05
	}
	if s.PermutationCount == 0 {
		s.PermutationCount = 1000
	}
	if s.Compressor == "" {
		s.Compressor = "zstd"
	}
	// SHA-141: Tokenizer defaults - disabled by default, all patterns on when enabled
	// (no defaults needed - bool zero values are false which is correct)
	// SHA-143: Numeric outlier defaults
	if s.NumericKSigma == 0 {
		s.NumericKSigma = 3.0
	}
}

// TokenizerEnabled returns true if tokenization is enabled
func (s *streamSettings) GetTokenizerEnabled() bool {
	return s.TokenizerEnabled
}

// GetTokenizerPatterns returns which tokenizer patterns are enabled
// Returns: enableUUID, enableHash, enableBase64, enableJWT
func (s *streamSettings) GetTokenizerPatterns() (bool, bool, bool, bool) {
	if !s.TokenizerEnabled {
		return false, false, false, false
	}
	// If tokenizer is enabled but no specific patterns set, enable all
	anySet := s.TokenizeUUID || s.TokenizeHash || s.TokenizeBase64 || s.TokenizeJWT
	if !anySet {
		return true, true, true, true
	}
	return s.TokenizeUUID, s.TokenizeHash, s.TokenizeBase64, s.TokenizeJWT
}

type configCache struct {
	mu           sync.RWMutex
	tenants      map[uuid.UUID]tenantRecord
	streams      map[uuid.UUID]streamRecord
	streamConfig map[uuid.UUID]streamSettings
	slugIndex    map[string]uuid.UUID
	defaults     map[uuid.UUID]uuid.UUID
}

func newConfigCache() *configCache {
	return &configCache{
		tenants:      make(map[uuid.UUID]tenantRecord),
		streams:      make(map[uuid.UUID]streamRecord),
		streamConfig: make(map[uuid.UUID]streamSettings),
		slugIndex:    make(map[string]uuid.UUID),
		defaults:     make(map[uuid.UUID]uuid.UUID),
	}
}

func (c *configCache) replace(entries []struct {
	tenant tenantRecord
	stream streamRecord
	cfg    streamSettings
}) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.tenants = make(map[uuid.UUID]tenantRecord)
	c.streams = make(map[uuid.UUID]streamRecord)
	c.streamConfig = make(map[uuid.UUID]streamSettings)
	c.slugIndex = make(map[string]uuid.UUID)
	c.defaults = make(map[uuid.UUID]uuid.UUID)
	for _, entry := range entries {
		c.tenants[entry.tenant.ID] = entry.tenant
		c.streams[entry.stream.ID] = entry.stream
		c.streamConfig[entry.stream.ID] = entry.cfg
		c.slugIndex[entry.stream.cacheKey()] = entry.stream.ID
		if _, exists := c.defaults[entry.stream.TenantID]; !exists {
			c.defaults[entry.stream.TenantID] = entry.stream.ID
		}
	}
}

func (c *configCache) tenant(id uuid.UUID) (tenantRecord, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	t, ok := c.tenants[id]
	return t, ok
}

func (c *configCache) streamByID(id uuid.UUID) (streamRecord, streamSettings, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	s, ok := c.streams[id]
	if !ok {
		return streamRecord{}, streamSettings{}, false
	}
	cfg := c.streamConfig[id]
	return s, cfg, true
}

func (c *configCache) streamBySlug(tenantID uuid.UUID, slug string) (streamRecord, streamSettings, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	id, ok := c.slugIndex[tenantID.String()+":"+slug]
	if !ok {
		return streamRecord{}, streamSettings{}, false
	}
	s := c.streams[id]
	cfg := c.streamConfig[id]
	return s, cfg, true
}

func (c *configCache) defaultStream(tenantID uuid.UUID) (streamRecord, streamSettings, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	id, ok := c.defaults[tenantID]
	if !ok {
		return streamRecord{}, streamSettings{}, false
	}
	s := c.streams[id]
	cfg := c.streamConfig[id]
	return s, cfg, true
}

// CLI structures

type tenantCreateParams struct {
	Name                string
	Slug                string
	Plan                string
	StreamSlug          string
	StreamType          string
	StreamDescription   string
	StreamRetentionDays int
	KeyRole             string
	KeyName             string
	KeyRateLimit        int
	TenantRateLimit     int
	DefaultBaseline     int
	DefaultWindow       int
	DefaultHop          int
	NCDThreshold        float64
	Email               string
	SignupIP            string
	Source              string
	PValueThreshold     float64
	PermutationCount    int
	DefaultCompressor   string
	Seed                int64
	VerificationToken   string
	Status              string
	FirebaseUID         string // Optional: Firebase Auth UID to link tenant to user account
}

type tenantCreateResult struct {
	TenantID   uuid.UUID
	TenantSlug string
	StreamID   uuid.UUID
	StreamSlug string
	APIKeyID   uuid.UUID
	APIKey     string
	Tenant     struct {
		ID        uuid.UUID
		Name      string
		Slug      string
		Plan      string
		CreatedAt time.Time
	}
}

type apiKeyInfo struct {
	ID         uuid.UUID
	Name       string
	Role       string
	StreamID   *uuid.UUID
	CreatedAt  time.Time
	LastUsedAt *time.Time
}

type apiKeyRecord struct {
	ID                uuid.UUID
	TenantID          uuid.UUID
	Role              string
	KeyHash           string
	StreamID          uuid.UUID
	Name              string
	KeyRateLimit      int
	TenantName        string
	TenantSlug        string
	DefaultCompressor string
	TenantRateLimit   int
	CreatedAt         time.Time
	ExpiresAt         *time.Time
	RevokedAt         *time.Time
}

type tenantContext struct {
	Tenant tenantRecord
	Key    apiKeyRecord
}

func batchHash(tenantID uuid.UUID, timestamp time.Time, payload []byte) string {
	sum := sha256.Sum256(append(append([]byte(tenantID.String()), byte('|')), payload...))
	return fmt.Sprintf("%x", sum[:])
}

func slugify(name string) string {
	lower := strings.ToLower(name)
	builder := strings.Builder{}
	for _, r := range lower {
		if (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') {
			builder.WriteRune(r)
			continue
		}
		if r == ' ' || r == '-' || r == '_' {
			builder.WriteRune('-')
		}
	}
	slug := builder.String()
	slug = strings.Trim(slug, "-")
	if slug == "" {
		slug = fmt.Sprintf("tenant-%d", time.Now().Unix())
	}
	return slug
}

// Billing-related database functions

func (s *store) getTenantByStripeCustomer(ctx context.Context, customerID string) (*tenantRecord, error) {
	var t tenantRecord
	err := s.pool.QueryRow(ctx, `
		SELECT id, name, slug, plan, default_compressor, rate_limit_rps,
		       stripe_customer_id, stripe_subscription_id, stripe_status, COALESCE(email, '')
		FROM tenants WHERE stripe_customer_id = $1`, customerID).Scan(
		&t.ID, &t.Name, &t.Slug, &t.Plan, &t.DefaultCompressor, &t.RateLimitRPS,
		&t.StripeCustomerID, &t.StripeSubscriptionID, &t.StripeStatus, &t.Email)
	if err != nil {
		return nil, err
	}
	return &t, nil
}

func (s *store) setGracePeriod(ctx context.Context, tenantID uuid.UUID, gracePeriodEnd time.Time) error {
	_, err := s.pool.Exec(ctx, `
		UPDATE tenants
		SET grace_period_ends_at = $2,
		    payment_failure_count = COALESCE(payment_failure_count, 0) + 1,
		    stripe_status = 'past_due',
		    updated_at = NOW()
		WHERE id = $1`, tenantID, gracePeriodEnd)
	return err
}

func (s *store) clearGracePeriod(ctx context.Context, tenantID uuid.UUID) error {
	_, err := s.pool.Exec(ctx, `
		UPDATE tenants
		SET grace_period_ends_at = NULL,
		    payment_failure_count = 0,
		    stripe_status = 'active',
		    updated_at = NOW()
		WHERE id = $1`, tenantID)
	return err
}

// BillingStatus represents the billing state for a tenant
type BillingStatus struct {
	Status              string     `json:"status"` // free, trialing, active, grace_period, past_due, canceled
	Plan                string     `json:"plan"`
	StripeStatus        *string    `json:"stripe_status,omitempty"`
	TrialEndsAt         *time.Time `json:"trial_ends_at,omitempty"`
	TrialDaysRemaining  *int       `json:"trial_days_remaining,omitempty"`
	GracePeriodEndsAt   *time.Time `json:"grace_period_ends_at,omitempty"`
	PaymentFailureCount int        `json:"payment_failure_count"`
	CurrentPeriodEnd    *time.Time `json:"current_period_end,omitempty"`
}

func (s *store) getBillingStatus(ctx context.Context, tenantID uuid.UUID) (*BillingStatus, error) {
	var bs BillingStatus
	var trialEndsAt, gracePeriodEndsAt, currentPeriodEnd sql.NullTime

	err := s.pool.QueryRow(ctx, `
		SELECT plan, stripe_status, trial_ends_at, grace_period_ends_at,
		       COALESCE(payment_failure_count, 0), current_period_end
		FROM tenants WHERE id = $1`, tenantID).Scan(
		&bs.Plan, &bs.StripeStatus, &trialEndsAt, &gracePeriodEndsAt,
		&bs.PaymentFailureCount, &currentPeriodEnd)
	if err != nil {
		return nil, err
	}

	if trialEndsAt.Valid {
		bs.TrialEndsAt = &trialEndsAt.Time
	}
	if gracePeriodEndsAt.Valid {
		bs.GracePeriodEndsAt = &gracePeriodEndsAt.Time
	}
	if currentPeriodEnd.Valid {
		bs.CurrentPeriodEnd = &currentPeriodEnd.Time
	}

	// Compute derived fields
	bs.Status = computeBillingStatus(bs)
	if bs.TrialEndsAt != nil && bs.TrialEndsAt.After(time.Now()) {
		days := int(time.Until(*bs.TrialEndsAt).Hours() / 24)
		bs.TrialDaysRemaining = &days
	}

	return &bs, nil
}

func computeBillingStatus(bs BillingStatus) string {
	now := time.Now()

	// Check grace period first
	if bs.GracePeriodEndsAt != nil {
		if bs.GracePeriodEndsAt.After(now) {
			return "grace_period"
		}
		// Grace period expired - should be downgraded
		return "expired"
	}

	// Check trial status
	if bs.TrialEndsAt != nil && bs.TrialEndsAt.After(now) {
		return "trialing"
	}

	// Map Stripe status
	if bs.StripeStatus != nil {
		switch *bs.StripeStatus {
		case "active":
			return "active"
		case "past_due":
			return "past_due"
		case "canceled", "unpaid":
			return "canceled"
		case "trialing":
			return "trialing"
		}
	}

	// Default to free tier
	return "free"
}

// regenerateVerificationToken generates a new token for a pending tenant
func (s *store) regenerateVerificationToken(ctx context.Context, email string, newToken string) (string, error) {
	now := time.Now()
	expiresAt := now.Add(15 * time.Minute)

	// Update token for pending tenant with matching email
	tag, err := s.pool.Exec(ctx, `
		UPDATE tenants
		SET verification_token = $2,
		    verification_token_expires_at = $3
		WHERE email = $1
		  AND status = 'pending_verification'`,
		email, newToken, expiresAt)
	if err != nil {
		return "", err
	}
	if tag.RowsAffected() == 0 {
		return "", fmt.Errorf("no pending account found for this email")
	}

	// Fetch company name for email template
	var companyName string
	err = s.pool.QueryRow(ctx, `SELECT name FROM tenants WHERE email = $1 AND status = 'pending_verification'`, email).Scan(&companyName)
	if err != nil {
		companyName = "there" // fallback
	}

	return companyName, nil
}

// addToWaitlist adds an email to the pre-launch waitlist
func (s *store) addToWaitlist(ctx context.Context, email, source, ipAddress string) error {
	_, err := s.pool.Exec(ctx, `
		INSERT INTO waitlist (email, source, ip_address)
		VALUES ($1, $2, $3)
		ON CONFLICT (email) DO NOTHING`,
		email, source, ipAddress)
	return err
}

// ─────────────────────────────────────────────────────────────────────────────
// SHA-139: Cold Start Guardrail - Stream Calibration Functions
// ─────────────────────────────────────────────────────────────────────────────

// CalibrationStatus holds the current calibration state for a stream
type CalibrationStatus struct {
	EventsIngested  int64
	IsCalibrated    bool
	MinBaselineSize int
	EventsNeeded    int64
	ProgressPercent int
}

// getStreamCalibrationStatus returns the current calibration state for a stream
func (s *store) getStreamCalibrationStatus(ctx context.Context, streamID uuid.UUID) (*CalibrationStatus, error) {
	var status CalibrationStatus
	err := s.pool.QueryRow(ctx, `
		SELECT events_ingested, is_calibrated, min_baseline_size
		FROM streams WHERE id = $1`,
		streamID).Scan(&status.EventsIngested, &status.IsCalibrated, &status.MinBaselineSize)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errNotFound
		}
		return nil, fmt.Errorf("failed to get stream calibration status: %w", err)
	}

	// Calculate derived fields
	status.EventsNeeded = int64(status.MinBaselineSize) - status.EventsIngested
	if status.EventsNeeded < 0 {
		status.EventsNeeded = 0
	}
	if status.MinBaselineSize > 0 {
		status.ProgressPercent = int(float64(status.EventsIngested) / float64(status.MinBaselineSize) * 100)
		if status.ProgressPercent > 100 {
			status.ProgressPercent = 100
		}
	}

	return &status, nil
}

// incrementStreamEvents atomically increments the events_ingested counter
// and updates is_calibrated if the threshold is reached.
// Returns the updated CalibrationStatus.
func (s *store) incrementStreamEvents(ctx context.Context, streamID uuid.UUID, count int) (*CalibrationStatus, error) {
	var status CalibrationStatus
	err := s.pool.QueryRow(ctx, `
		UPDATE streams
		SET events_ingested = events_ingested + $2,
			is_calibrated = CASE
				WHEN events_ingested + $2 >= min_baseline_size THEN TRUE
				ELSE is_calibrated
			END,
			calibrated_at = CASE
				WHEN NOT is_calibrated AND events_ingested + $2 >= min_baseline_size THEN NOW()
				ELSE calibrated_at
			END
		WHERE id = $1
		RETURNING events_ingested, is_calibrated, min_baseline_size`,
		streamID, count).Scan(&status.EventsIngested, &status.IsCalibrated, &status.MinBaselineSize)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errNotFound
		}
		return nil, fmt.Errorf("failed to increment stream events: %w", err)
	}

	// Calculate derived fields
	status.EventsNeeded = int64(status.MinBaselineSize) - status.EventsIngested
	if status.EventsNeeded < 0 {
		status.EventsNeeded = 0
	}
	if status.MinBaselineSize > 0 {
		status.ProgressPercent = int(float64(status.EventsIngested) / float64(status.MinBaselineSize) * 100)
		if status.ProgressPercent > 100 {
			status.ProgressPercent = 100
		}
	}

	return &status, nil
}

// ─────────────────────────────────────────────────────────────────────────────
// SHA-140: Anchor Baseline for Drift Detection
// ─────────────────────────────────────────────────────────────────────────────

// StreamAnchor represents a frozen baseline snapshot for drift detection
type StreamAnchor struct {
	ID                       uuid.UUID
	StreamID                 uuid.UUID
	AnchorData               []byte
	Compressor               string
	EventCount               int
	CalibrationCompletedAt   time.Time
	IsActive                 bool
	BaselineEntropy          *float64
	BaselineCompressionRatio *float64
	BaselineNCDSelf          *float64
	DriftNCDThreshold        float64
	CreatedAt                time.Time
}

// createStreamAnchor creates a new anchor for a stream, deactivating any existing one
func (s *store) createStreamAnchor(ctx context.Context, anchor *StreamAnchor) error {
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	// Deactivate existing active anchor if any
	_, err = tx.Exec(ctx, `
		UPDATE stream_anchors
		SET is_active = FALSE, superseded_at = NOW()
		WHERE stream_id = $1 AND is_active = TRUE`,
		anchor.StreamID)
	if err != nil {
		return fmt.Errorf("failed to deactivate existing anchor: %w", err)
	}

	// Insert new anchor
	err = tx.QueryRow(ctx, `
		INSERT INTO stream_anchors (
			stream_id, anchor_data, compressor, event_count,
			calibration_completed_at, is_active,
			baseline_entropy, baseline_compression_ratio, baseline_ncd_self,
			drift_ncd_threshold
		) VALUES ($1, $2, $3, $4, $5, TRUE, $6, $7, $8, $9)
		RETURNING id, created_at`,
		anchor.StreamID, anchor.AnchorData, anchor.Compressor, anchor.EventCount,
		anchor.CalibrationCompletedAt,
		anchor.BaselineEntropy, anchor.BaselineCompressionRatio, anchor.BaselineNCDSelf,
		anchor.DriftNCDThreshold,
	).Scan(&anchor.ID, &anchor.CreatedAt)
	if err != nil {
		return fmt.Errorf("failed to insert anchor: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("failed to commit anchor transaction: %w", err)
	}

	return nil
}

// getActiveAnchor retrieves the currently active anchor for a stream
func (s *store) getActiveAnchor(ctx context.Context, streamID uuid.UUID) (*StreamAnchor, error) {
	var anchor StreamAnchor
	err := s.pool.QueryRow(ctx, `
		SELECT id, stream_id, anchor_data, compressor, event_count,
			   calibration_completed_at, is_active,
			   baseline_entropy, baseline_compression_ratio, baseline_ncd_self,
			   drift_ncd_threshold, created_at
		FROM stream_anchors
		WHERE stream_id = $1 AND is_active = TRUE`,
		streamID).Scan(
		&anchor.ID, &anchor.StreamID, &anchor.AnchorData, &anchor.Compressor,
		&anchor.EventCount, &anchor.CalibrationCompletedAt, &anchor.IsActive,
		&anchor.BaselineEntropy, &anchor.BaselineCompressionRatio, &anchor.BaselineNCDSelf,
		&anchor.DriftNCDThreshold, &anchor.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil // No active anchor (not an error)
		}
		return nil, fmt.Errorf("failed to get active anchor: %w", err)
	}
	return &anchor, nil
}

// deactivateAnchor deactivates the current anchor for a stream
func (s *store) deactivateAnchor(ctx context.Context, streamID uuid.UUID) error {
	_, err := s.pool.Exec(ctx, `
		UPDATE stream_anchors
		SET is_active = FALSE, superseded_at = NOW()
		WHERE stream_id = $1 AND is_active = TRUE`,
		streamID)
	if err != nil {
		return fmt.Errorf("failed to deactivate anchor: %w", err)
	}
	return nil
}

// AnchorSettings contains anchor-related settings for a stream
type AnchorSettings struct {
	AnchorEnabled      bool
	DriftNCDThreshold  float64
	AnchorResetOnDrift bool
}

// getAnchorSettings retrieves anchor settings for a stream
func (s *store) getAnchorSettings(ctx context.Context, streamID uuid.UUID) (*AnchorSettings, error) {
	var settings AnchorSettings
	err := s.pool.QueryRow(ctx, `
		SELECT anchor_enabled, drift_ncd_threshold, anchor_reset_on_drift
		FROM streams WHERE id = $1`,
		streamID).Scan(&settings.AnchorEnabled, &settings.DriftNCDThreshold, &settings.AnchorResetOnDrift)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errNotFound
		}
		return nil, fmt.Errorf("failed to get anchor settings: %w", err)
	}
	return &settings, nil
}
