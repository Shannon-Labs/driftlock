package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
)

// StateManager handles distributed state for CBAD processing across multiple instances
type StateManager struct {
	client *redis.Client
	prefix string
}

// NewStateManager creates a new Redis-based state manager
func NewStateManager(addr, password, prefix string, db int) *StateManager {
	rdb := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
	})

	return &StateManager{
		client: rdb,
		prefix: prefix,
	}
}

// CBADWindowState represents the state of a CBAD window for a specific stream
type CBADWindowState struct {
	StreamType       string    `json:"stream_type"`
	WindowID         string    `json:"window_id"`
	BaselineData     []byte    `json:"baseline_data"`
	WindowData       []byte    `json:"window_data"`
	Timestamp        time.Time `json:"timestamp"`
	CompressionRatio float64   `json:"compression_ratio"`
	NCD              float64   `json:"ncd"`
	PermutationCount int       `json:"permutation_count"`
}

// SaveWindowState stores CBAD window state in Redis
func (sm *StateManager) SaveWindowState(ctx context.Context, streamID string, state *CBADWindowState) error {
	key := fmt.Sprintf("%s:cbad:window:%s", sm.prefix, streamID)
	data, err := json.Marshal(state)
	if err != nil {
		return fmt.Errorf("marshal window state: %w", err)
	}

	// Store with expiration to clean up old states automatically
	err = sm.client.SetEX(ctx, key, data, 24*time.Hour).Err()
	if err != nil {
		return fmt.Errorf("redis set: %w", err)
	}

	return nil
}

// GetWindowState retrieves CBAD window state from Redis
func (sm *StateManager) GetWindowState(ctx context.Context, streamID string) (*CBADWindowState, error) {
	key := fmt.Sprintf("%s:cbad:window:%s", sm.prefix, streamID)
	data, err := sm.client.Get(ctx, key).Result()
	if err != nil {
		if err == redis.Nil {
			return nil, nil // No state found
		}
		return nil, fmt.Errorf("redis get: %w", err)
	}

	var state CBADWindowState
	err = json.Unmarshal([]byte(data), &state)
	if err != nil {
		return nil, fmt.Errorf("unmarshal window state: %w", err)
	}

	return &state, nil
}

// DeleteWindowState removes CBAD window state from Redis
func (sm *StateManager) DeleteWindowState(ctx context.Context, streamID string) error {
	key := fmt.Sprintf("%s:cbad:window:%s", sm.prefix, streamID)
	err := sm.client.Del(ctx, key).Err()
	if err != nil {
		return fmt.Errorf("redis del: %w", err)
	}
	return nil
}

// SaveStreamConfiguration stores CBAD configuration for a specific stream
func (sm *StateManager) SaveStreamConfiguration(ctx context.Context, streamID string, config map[string]interface{}) error {
	key := fmt.Sprintf("%s:cbad:config:%s", sm.prefix, streamID)
	data, err := json.Marshal(config)
	if err != nil {
		return fmt.Errorf("marshal config: %w", err)
	}

	err = sm.client.SetEX(ctx, key, data, 7*24*time.Hour).Err()
	if err != nil {
		return fmt.Errorf("redis set config: %w", err)
	}

	return nil
}

// GetStreamConfiguration retrieves CBAD configuration for a specific stream
func (sm *StateManager) GetStreamConfiguration(ctx context.Context, streamID string) (map[string]interface{}, error) {
	key := fmt.Sprintf("%s:cbad:config:%s", sm.prefix, streamID)
	data, err := sm.client.Get(ctx, key).Result()
	if err != nil {
		if err == redis.Nil {
			return nil, nil // No config found
		}
		return nil, fmt.Errorf("redis get config: %w", err)
	}

	var config map[string]interface{}
	err = json.Unmarshal([]byte(data), &config)
	if err != nil {
		return nil, fmt.Errorf("unmarshal config: %w", err)
	}

	return config, nil
}

// Close closes the Redis connection
func (sm *StateManager) Close() error {
	return sm.client.Close()
}