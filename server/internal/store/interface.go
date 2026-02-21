package store

import (
	"context"
	"time"
)

// Store defines the interface for game state persistence.
// The primary implementation uses Redis, but this interface allows
// for mock implementations in tests.
type Store interface {
	// SaveGameState persists a serialized game state snapshot.
	SaveGameState(ctx context.Context, gameID string, data []byte, ttl time.Duration) error

	// LoadGameState retrieves a game state snapshot.
	// Returns nil, nil if the key does not exist.
	LoadGameState(ctx context.Context, gameID string) ([]byte, error)

	// DeleteGameState removes a game state snapshot.
	DeleteGameState(ctx context.Context, gameID string) error

	// ListGameIDs returns all active game IDs (keys matching game:*).
	ListGameIDs(ctx context.Context) ([]string, error)

	// Ping checks the connection to the store.
	Ping(ctx context.Context) error

	// Close shuts down the store connection.
	Close() error
}
