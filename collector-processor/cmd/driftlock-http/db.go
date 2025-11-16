package main

import (
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
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
	"golang.org/x/crypto/argon2"
)

var base32Enc = base32.StdEncoding.WithPadding(base32.NoPadding)
var errNotFound = errors.New("record not found")

type store struct {
	pool  *pgxpool.Pool
	cache *configCache
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
	if cfg.MaxConns == 0 {
		cfg.MaxConns = int32(envInt("DB_MAX_CONNS", 10))
	}
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
        SELECT t.id, t.slug, t.name, t.default_compressor, t.rate_limit_rps,
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
		if err := rows.Scan(&tenant.ID, &tenant.Slug, &tenant.Name, &tenant.DefaultCompressor, &tenant.RateLimitRPS,
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

	_, err = tx.Exec(ctx, `INSERT INTO tenants (id, name, slug, plan, default_compressor, rate_limit_rps) VALUES ($1,$2,$3,$4,$5,$6)`,
		tenantID, params.Name, slug, params.Plan, params.DefaultCompressor, params.TenantRateLimit)
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

	return &tenantCreateResult{
		TenantID:   tenantID,
		TenantSlug: slug,
		StreamID:   streamID,
		StreamSlug: streamSlug,
		APIKey:     key,
		APIKeyID:   keyID,
	}, nil
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
	tag, err := s.pool.Exec(ctx, `DELETE FROM api_keys WHERE id=$1`, keyID)
	if err != nil {
		return 0, err
	}
	return tag.RowsAffected(), nil
}

func (s *store) resolveAPIKey(ctx context.Context, id uuid.UUID) (*apiKeyRecord, error) {
	row := s.pool.QueryRow(ctx, `SELECT ak.id, ak.tenant_id, ak.role, ak.key_hash, ak.stream_id, ak.name, ak.rate_limit_rps,
        t.name, t.slug, t.default_compressor, t.rate_limit_rps
        FROM api_keys ak JOIN tenants t ON ak.tenant_id = t.id WHERE ak.id=$1`, id)
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
	ID                uuid.UUID
	Name              string
	Slug              string
	DefaultCompressor string
	RateLimitRPS      int
}

type streamRecord struct {
	ID            uuid.UUID
	TenantID      uuid.UUID
	Slug          string
	Type          string
	Seed          int64
	Compressor    string
	RetentionDays int
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
	PValueThreshold     float64
	PermutationCount    int
	DefaultCompressor   string
	Seed                int64
}

type tenantCreateResult struct {
	TenantID   uuid.UUID
	TenantSlug string
	StreamID   uuid.UUID
	StreamSlug string
	APIKeyID   uuid.UUID
	APIKey     string
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
