package game

import (
	"context"
	"fmt"
	"log/slog"
	"sync"
	"time"

	"github.com/teomiscia/hexbattle/internal/store"
	"github.com/teomiscia/hexbattle/internal/ws"
)

// Manager tracks and coordinates active game engines.
type Manager struct {
	mu      sync.RWMutex
	engines map[string]*Engine
	store   store.Store
}

// NewManager creates a new Game Manager.
func NewManager(st store.Store) *Manager {
	return &Manager{
		engines: make(map[string]*Engine),
		store:   st,
	}
}

// AddEngine registers a new engine and starts it.
func (m *Manager) AddEngine(engine *Engine) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.engines[engine.State.ID] = engine
	go engine.Run()
}

// GetEngine retrieves an active engine by game ID.
func (m *Manager) GetEngine(gameID string) *Engine {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.engines[gameID]
}

// RemoveEngine stops and removes an engine from the manager.
func (m *Manager) RemoveEngine(gameID string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	if engine, ok := m.engines[gameID]; ok {
		engine.Stop()
		delete(m.engines, gameID)
	}
}

// RestoreActiveGames loads all game snapshots from Redis and resumes them.
// Games that are already in GameOver state are not resumed.
func (m *Manager) RestoreActiveGames(ctx context.Context) error {
	if m.store == nil {
		return nil
	}

	gameIDs, err := m.store.ListGameIDs(ctx)
	if err != nil {
		return fmt.Errorf("failed to list game IDs: %w", err)
	}

	restoredCount := 0
	for _, id := range gameIDs {
		data, err := m.store.LoadGameState(ctx, id)
		if err != nil {
			slog.Warn("failed to load game state for restore", "game_id", id, "error", err)
			continue
		}
		if data == nil {
			continue
		}

		state, err := DeserializeGameState(data)
		if err != nil {
			slog.Warn("failed to deserialize game state for restore", "game_id", id, "error", err)
			continue
		}

		// Don't resume games that are already over
		if string(state.Phase) == "game_over" {
			continue
		}

		// Create engine
		hub := ws.NewHub()
		engine := NewEngine(context.Background(), state, hub, m.store)

		// Mark all players as disconnected initially
		for i := 0; i < 2; i++ {
			state.Players[i].IsDisconnected = true
		}

		m.AddEngine(engine)

		// Immediately start reconnect timers for both players
		for i := 0; i < 2; i++ {
			engine.NotifyDisconnect(state.Players[i].ID)
		}

		restoredCount++
	}

	slog.Info("restored active games from redis", "count", restoredCount)
	return nil
}

// StopAll stops all running game engines and forces a final snapshot.
func (m *Manager) StopAll() {
	m.mu.Lock()
	defer m.mu.Unlock()

	var wg sync.WaitGroup
	for _, engine := range m.engines {
		wg.Add(1)
		go func(e *Engine) {
			defer wg.Done()
			e.Stop() // this will trigger snapshotState in the engine's defer
		}(engine)
	}

	// Wait with a timeout
	done := make(chan struct{})
	go func() {
		wg.Wait()
		close(done)
	}()

	select {
	case <-done:
		slog.Info("all game engines stopped gracefully")
	case <-time.After(5 * time.Second):
		slog.Warn("timeout waiting for game engines to stop")
	}
}
