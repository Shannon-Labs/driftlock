package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/Shannon-Labs/driftlock/collector-processor/cmd/driftlock-http/plans"
	"github.com/google/uuid"
)

func handleCLI(args []string, cfg config) bool {
	if len(args) == 0 {
		return false
	}

	switch args[0] {
	case "migrate":
		action := "up"
		if len(args) > 1 {
			action = strings.ToLower(args[1])
		}
		if err := runMigrations(context.Background(), action); err != nil {
			log.Fatalf("migrate %s failed: %v", action, err)
		}
		fmt.Printf("migrate %s succeeded\n", action)
		return true
	case "create-tenant":
		createTenantCommand(args[1:], cfg)
		return true
	case "list-keys":
		listKeysCommand(args[1:])
		return true
	case "revoke-key":
		revokeKeyCommand(args[1:])
		return true
	default:
		return false
	}
}

func createTenantCommand(args []string, cfg config) {
	fs := flag.NewFlagSet("create-tenant", flag.ExitOnError)
	var name string
	var slug string
	var plan string
	var streamSlug string
	var streamType string
	var keyRole string
	var keyName string
	var description string
	var retention int
	var rateLimit int
	var keyRateLimit int
	var jsonOutput bool
	fs.StringVar(&name, "name", "", "Tenant name")
	fs.StringVar(&slug, "slug", "", "Tenant slug (optional)")
	fs.StringVar(&plan, "plan", plans.Pulse, "Plan tier (pulse, radar, tensor, orbit)")
	fs.StringVar(&streamSlug, "stream", "default", "Initial stream slug")
	fs.StringVar(&streamType, "stream-type", "logs", "Stream type (logs|metrics|traces|llm)")
	fs.StringVar(&description, "stream-description", "CLI created stream", "Stream description")
	fs.StringVar(&keyRole, "key-role", "admin", "Key role (admin|stream)")
	fs.StringVar(&keyName, "key-name", "cli", "API key label")
	fs.IntVar(&retention, "retention-days", 14, "Stream retention days")
	fs.IntVar(&rateLimit, "tenant-rate-limit", cfg.DefaultRateLimit(), "Tenant rate limit rps")
	fs.IntVar(&keyRateLimit, "key-rate-limit", 0, "Key-specific rate limit (optional)")
	fs.BoolVar(&jsonOutput, "json", false, "Print tenant metadata as JSON")
	_ = fs.Parse(args)
	if name == "" {
		log.Fatalf("--name is required")
	}

	// Normalize plan name
	normalizedPlan, _ := plans.NormalizePlan(plan)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	res, err := createTenant(ctx, cfg, tenantCreateParams{
		Name:                name,
		Slug:                slug,
		Plan:                normalizedPlan,
		StreamSlug:          streamSlug,
		StreamType:          streamType,
		StreamDescription:   description,
		StreamRetentionDays: retention,
		KeyRole:             keyRole,
		KeyName:             keyName,
		KeyRateLimit:        keyRateLimit,
		TenantRateLimit:     rateLimit,
		DefaultBaseline:     cfg.DefaultBaseline,
		DefaultWindow:       cfg.DefaultWindow,
		DefaultHop:          cfg.DefaultHop,
		NCDThreshold:        cfg.NCDThreshold,
		PValueThreshold:     cfg.PValueThreshold,
		PermutationCount:    cfg.PermutationCount,
		DefaultCompressor:   cfg.DefaultAlgo,
		Seed:                int64(cfg.Seed),
		SignupIP:            "127.0.0.1",
	})
	if err != nil {
		log.Fatalf("create-tenant failed: %v", err)
	}

	if jsonOutput {
		payload := map[string]any{
			"tenant_id":   res.TenantID.String(),
			"tenant_slug": res.TenantSlug,
			"stream_id":   res.StreamID.String(),
			"stream_slug": res.StreamSlug,
			"api_key":     res.APIKey,
			"api_key_id":  res.APIKeyID.String(),
			"role":        strings.ToLower(keyRole),
		}
		enc := json.NewEncoder(os.Stdout)
		enc.SetIndent("", "  ")
		_ = enc.Encode(payload)
		return
	}

	fmt.Printf("tenant %s (%s) created\n", res.TenantSlug, res.TenantID)
	fmt.Printf("stream %s (%s) created\n", res.StreamSlug, res.StreamID)
	fmt.Printf("api key (role=%s): %s\n", strings.ToLower(keyRole), res.APIKey)
}

func listKeysCommand(args []string) {
	fs := flag.NewFlagSet("list-keys", flag.ExitOnError)
	var tenantSlug string
	fs.StringVar(&tenantSlug, "tenant", "", "Tenant slug")
	_ = fs.Parse(args)
	if tenantSlug == "" {
		log.Fatalf("--tenant is required")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	pool, err := connectDB(ctx)
	if err != nil {
		log.Fatalf("db connect failed: %v", err)
	}
	defer pool.Close()

	s := newStore(pool)
	keys, err := s.listKeys(ctx, tenantSlug)
	if err != nil {
		log.Fatalf("list-keys failed: %v", err)
	}
	if len(keys) == 0 {
		fmt.Println("no keys found")
		return
	}
	for _, k := range keys {
		stream := "*"
		if k.StreamID != nil {
			stream = k.StreamID.String()
		}
		lastUsed := "never"
		if k.LastUsedAt != nil {
			lastUsed = k.LastUsedAt.Format(time.RFC3339)
		}
		fmt.Printf("%s\t%s\t%s\tstream=%s\tcreated=%s\tlast=%s\n", k.ID, k.Name, k.Role, stream, k.CreatedAt.Format(time.RFC3339), lastUsed)
	}
}

func revokeKeyCommand(args []string) {
	fs := flag.NewFlagSet("revoke-key", flag.ExitOnError)
	var keyIDStr string
	fs.StringVar(&keyIDStr, "id", "", "API key ID (UUID)")
	_ = fs.Parse(args)
	if keyIDStr == "" {
		log.Fatalf("--id is required")
	}
	keyID, err := uuid.Parse(keyIDStr)
	if err != nil {
		log.Fatalf("invalid key id: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	pool, err := connectDB(ctx)
	if err != nil {
		log.Fatalf("db connect failed: %v", err)
	}
	defer pool.Close()

	s := newStore(pool)
	n, err := s.revokeKey(ctx, keyID)
	if err != nil {
		log.Fatalf("revoke failed: %v", err)
	}
	if n == 0 {
		fmt.Println("no rows deleted")
	} else {
		fmt.Printf("revoked key %s\n", keyID)
	}
}

func createTenant(ctx context.Context, cfg config, params tenantCreateParams) (*tenantCreateResult, error) {
	pool, err := connectDB(ctx)
	if err != nil {
		return nil, err
	}
	defer pool.Close()
	s := newStore(pool)
	if err := s.loadCache(ctx); err != nil {
		return nil, err
	}
	return s.createTenantWithKey(ctx, params)
}
