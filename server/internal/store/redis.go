package store

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/redis/go-redis/v9"
)

const gameKeyPrefix = "game:"

// RedisStore implements the Store interface using Redis.
type RedisStore struct {
	client *redis.Client
}

// NewRedisStore creates a new Redis-backed store.
// The redisURL should be in the format "redis://host:port" or "redis://host:port/db".
func NewRedisStore(redisURL string) (*RedisStore, error) {
	opts, err := redis.ParseURL(redisURL)
	if err != nil {
		return nil, fmt.Errorf("store: invalid redis URL %q: %w", redisURL, err)
	}

	client := redis.NewClient(opts)

	// Verify connection with retries
	var lastErr error
	for i := 0; i < 5; i++ {
		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		if err := client.Ping(ctx).Err(); err != nil {
			lastErr = err
			cancel()
			time.Sleep(time.Duration(i+1) * time.Second)
			continue
		}
		cancel()
		return &RedisStore{client: client}, nil
	}

	return nil, fmt.Errorf("store: failed to connect to redis after 5 attempts: %w", lastErr)
}

// SaveGameState persists a serialized game state snapshot to Redis.
func (s *RedisStore) SaveGameState(ctx context.Context, gameID string, data []byte, ttl time.Duration) error {
	key := gameKeyPrefix + gameID
	return s.client.Set(ctx, key, data, ttl).Err()
}

// LoadGameState retrieves a game state snapshot from Redis.
// Returns nil, nil if the key does not exist.
func (s *RedisStore) LoadGameState(ctx context.Context, gameID string) ([]byte, error) {
	key := gameKeyPrefix + gameID
	data, err := s.client.Get(ctx, key).Bytes()
	if err == redis.Nil {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("store: failed to load game %s: %w", gameID, err)
	}
	return data, nil
}

// DeleteGameState removes a game state snapshot from Redis.
func (s *RedisStore) DeleteGameState(ctx context.Context, gameID string) error {
	key := gameKeyPrefix + gameID
	return s.client.Del(ctx, key).Err()
}

// ListGameIDs returns all active game IDs by scanning keys matching "game:*".
func (s *RedisStore) ListGameIDs(ctx context.Context) ([]string, error) {
	var gameIDs []string
	var cursor uint64

	for {
		keys, nextCursor, err := s.client.Scan(ctx, cursor, gameKeyPrefix+"*", 100).Result()
		if err != nil {
			return nil, fmt.Errorf("store: failed to scan game keys: %w", err)
		}

		for _, key := range keys {
			gameID := strings.TrimPrefix(key, gameKeyPrefix)
			gameIDs = append(gameIDs, gameID)
		}

		cursor = nextCursor
		if cursor == 0 {
			break
		}
	}

	return gameIDs, nil
}

// Ping checks the Redis connection.
func (s *RedisStore) Ping(ctx context.Context) error {
	return s.client.Ping(ctx).Err()
}

// PingLatency returns the round-trip latency to Redis in milliseconds.
func (s *RedisStore) PingLatency(ctx context.Context) (int64, error) {
	start := time.Now()
	if err := s.client.Ping(ctx).Err(); err != nil {
		return 0, err
	}
	return time.Since(start).Milliseconds(), nil
}

// Close shuts down the Redis connection.
func (s *RedisStore) Close() error {
	return s.client.Close()
}
