package ws

import (
	"log/slog"
	"sync"
)

// Hub manages WebSocket connections for a single game instance.
// It handles broadcasting messages to both players and direct sends.
type Hub struct {
	mu    sync.RWMutex
	conns map[string]*Connection // player_id -> connection
}

// NewHub creates a new per-game message hub.
func NewHub() *Hub {
	return &Hub{
		conns: make(map[string]*Connection),
	}
}

// Register adds a connection to the hub.
func (h *Hub) Register(conn *Connection) {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.conns[conn.PlayerID] = conn
}

// Unregister removes a connection from the hub.
func (h *Hub) Unregister(playerID string) {
	h.mu.Lock()
	defer h.mu.Unlock()
	delete(h.conns, playerID)
}

// GetConnection returns the connection for a player, or nil if not connected.
func (h *Hub) GetConnection(playerID string) *Connection {
	h.mu.RLock()
	defer h.mu.RUnlock()
	return h.conns[playerID]
}

// IsConnected returns true if the player has an active connection.
func (h *Hub) IsConnected(playerID string) bool {
	h.mu.RLock()
	defer h.mu.RUnlock()
	_, ok := h.conns[playerID]
	return ok
}

// Broadcast sends a message to all connected players in the hub.
func (h *Hub) Broadcast(data []byte) {
	h.mu.RLock()
	defer h.mu.RUnlock()
	for _, conn := range h.conns {
		if !conn.Send(data) {
			slog.Warn("failed to broadcast to player",
				"player_id", conn.PlayerID,
			)
		}
	}
}

// BroadcastMessage marshals and broadcasts a typed message to all players.
func (h *Hub) BroadcastMessage(msgType string, data interface{}) error {
	bytes, err := NewEnvelope(msgType, data)
	if err != nil {
		return err
	}
	h.Broadcast(bytes)
	return nil
}

// SendTo sends a message to a specific player.
func (h *Hub) SendTo(playerID string, data []byte) bool {
	h.mu.RLock()
	defer h.mu.RUnlock()
	conn, ok := h.conns[playerID]
	if !ok {
		return false
	}
	return conn.Send(data)
}

// SendMessageTo marshals and sends a typed message to a specific player.
func (h *Hub) SendMessageTo(playerID string, msgType string, data interface{}) error {
	bytes, err := NewEnvelope(msgType, data)
	if err != nil {
		return err
	}
	if !h.SendTo(playerID, bytes) {
		return nil // player not connected, not an error
	}
	return nil
}

// ConnectedCount returns the number of connected players.
func (h *Hub) ConnectedCount() int {
	h.mu.RLock()
	defer h.mu.RUnlock()
	return len(h.conns)
}

// CloseAll closes all connections in the hub.
func (h *Hub) CloseAll() {
	h.mu.Lock()
	defer h.mu.Unlock()
	for _, conn := range h.conns {
		conn.Close()
	}
	h.conns = make(map[string]*Connection)
}
