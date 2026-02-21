package lobby

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"sync"
	"time"

	"github.com/teomiscia/hexbattle/internal/model"
)

// Room represents a game room where two players meet before a game starts.
type Room struct {
	Code          string             `json:"code"`
	ID            string             `json:"id"`
	HostPlayerID  string             `json:"host_player_id"`
	HostNickname  string             `json:"host_nickname"`
	GuestPlayerID string             `json:"guest_player_id,omitempty"`
	GuestNickname string             `json:"guest_nickname,omitempty"`
	Settings      model.RoomSettings `json:"settings"`
	State         model.RoomState    `json:"state"`
	GameID        string             `json:"game_id,omitempty"`
	CreatedAt     time.Time          `json:"created_at"`
}

// IsFull returns true if both players are in the room.
func (r *Room) IsFull() bool {
	return r.GuestPlayerID != ""
}

// HasPlayer returns true if the given player ID is in this room.
func (r *Room) HasPlayer(playerID string) bool {
	return r.HostPlayerID == playerID || r.GuestPlayerID == playerID
}

// OpponentID returns the other player's ID, given one player's ID.
func (r *Room) OpponentID(playerID string) string {
	if r.HostPlayerID == playerID {
		return r.GuestPlayerID
	}
	return r.HostPlayerID
}

// generateRoomCode creates a 6-digit numeric string using crypto/rand.
func generateRoomCode() (string, error) {
	n, err := rand.Int(rand.Reader, big.NewInt(1000000))
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%06d", n.Int64()), nil
}

// generateRoomID creates a unique room ID.
func generateRoomID() (string, error) {
	var uuid [16]byte
	if _, err := rand.Read(uuid[:]); err != nil {
		return "", err
	}
	uuid[6] = (uuid[6] & 0x0f) | 0x40
	uuid[8] = (uuid[8] & 0x3f) | 0x80
	return fmt.Sprintf("%08x-%04x-%04x-%04x-%012x",
		uuid[0:4], uuid[4:6], uuid[6:8], uuid[8:10], uuid[10:16]), nil
}

// CreateRoomResponse is returned when a room is created.
type CreateRoomResponse struct {
	RoomCode string             `json:"room_code"`
	RoomID   string             `json:"room_id"`
	Settings model.RoomSettings `json:"settings"`
}

// JoinRoomResponse is returned when a player joins a room.
type JoinRoomResponse struct {
	RoomID       string             `json:"room_id"`
	Settings     model.RoomSettings `json:"settings"`
	HostNickname string             `json:"host_nickname"`
}

// RoomStatusResponse is returned when polling room status.
type RoomStatusResponse struct {
	Code          string             `json:"code"`
	State         model.RoomState    `json:"state"`
	Settings      model.RoomSettings `json:"settings"`
	HostNickname  string             `json:"host_nickname"`
	GuestNickname string             `json:"guest_nickname,omitempty"`
	GameID        string             `json:"game_id,omitempty"`
}

// Manager handles room creation, joining, and lifecycle.
// All room state is in-memory only.
type Manager struct {
	mu       sync.RWMutex
	byCode   map[string]*Room // code -> room
	byID     map[string]*Room // room_id -> room
	roomTTL  time.Duration
	stopChan chan struct{}
}

// NewManager creates a new lobby manager.
func NewManager(roomTTL time.Duration) *Manager {
	m := &Manager{
		byCode:   make(map[string]*Room),
		byID:     make(map[string]*Room),
		roomTTL:  roomTTL,
		stopChan: make(chan struct{}),
	}
	go m.cleanupLoop()
	return m
}

// CreateRoom creates a new room with the given settings.
func (m *Manager) CreateRoom(hostPlayerID, hostNickname string, settings model.RoomSettings) (*Room, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	// Generate a unique room code (retry on collision)
	var code string
	var err error
	for attempts := 0; attempts < 10; attempts++ {
		code, err = generateRoomCode()
		if err != nil {
			return nil, fmt.Errorf("failed to generate room code: %w", err)
		}
		if _, exists := m.byCode[code]; !exists {
			break
		}
		if attempts == 9 {
			return nil, fmt.Errorf("failed to generate unique room code after 10 attempts")
		}
	}

	roomID, err := generateRoomID()
	if err != nil {
		return nil, fmt.Errorf("failed to generate room ID: %w", err)
	}

	room := &Room{
		Code:         code,
		ID:           roomID,
		HostPlayerID: hostPlayerID,
		HostNickname: hostNickname,
		Settings:     settings,
		State:        model.RoomWaitingForOpponent,
		CreatedAt:    time.Now(),
	}

	m.byCode[code] = room
	m.byID[roomID] = room

	return room, nil
}

// JoinRoom adds a guest player to a room by code.
func (m *Manager) JoinRoom(code, guestPlayerID, guestNickname string) (*Room, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	room, ok := m.byCode[code]
	if !ok {
		return nil, fmt.Errorf("room not found")
	}

	if room.State != model.RoomWaitingForOpponent {
		if room.IsFull() {
			return nil, fmt.Errorf("room is full")
		}
		return nil, fmt.Errorf("room is not accepting players")
	}

	if room.HostPlayerID == guestPlayerID {
		return nil, fmt.Errorf("cannot join your own room")
	}

	room.GuestPlayerID = guestPlayerID
	room.GuestNickname = guestNickname
	room.State = model.RoomReady

	return room, nil
}

// GetByCode returns a room by its code.
func (m *Manager) GetByCode(code string) *Room {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.byCode[code]
}

// GetByID returns a room by its ID.
func (m *Manager) GetByID(id string) *Room {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.byID[id]
}

// GetAllRooms returns all rooms.
func (m *Manager) GetAllRooms() []*Room {
	m.mu.RLock()
	defer m.mu.RUnlock()
	rooms := make([]*Room, 0, len(m.byID))
	for _, room := range m.byID {
		rooms = append(rooms, room)
	}
	return rooms
}

// SetGameInProgress marks a room as having an active game.
func (m *Manager) SetGameInProgress(roomID, gameID string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	if room, ok := m.byID[roomID]; ok {
		room.State = model.RoomGameInProgress
		room.GameID = gameID
	}
}

// SetGameOver marks a room's game as finished.
func (m *Manager) SetGameOver(roomID string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	if room, ok := m.byID[roomID]; ok {
		room.State = model.RoomGameOver
	}
}

// RemoveRoom removes a room from the manager.
func (m *Manager) RemoveRoom(code string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	if room, ok := m.byCode[code]; ok {
		delete(m.byCode, code)
		delete(m.byID, room.ID)
	}
}

// WaitingRoomCount returns the number of rooms waiting for an opponent.
func (m *Manager) WaitingRoomCount() int {
	m.mu.RLock()
	defer m.mu.RUnlock()
	count := 0
	for _, room := range m.byCode {
		if room.State == model.RoomWaitingForOpponent {
			count++
		}
	}
	return count
}

// TotalRoomCount returns the total number of active rooms.
func (m *Manager) TotalRoomCount() int {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return len(m.byCode)
}

// cleanupLoop periodically removes expired rooms.
func (m *Manager) cleanupLoop() {
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			m.cleanupExpired()
		case <-m.stopChan:
			return
		}
	}
}

// cleanupExpired removes rooms that have exceeded their TTL while still waiting.
func (m *Manager) cleanupExpired() {
	m.mu.Lock()
	defer m.mu.Unlock()

	now := time.Now()
	for code, room := range m.byCode {
		switch room.State {
		case model.RoomWaitingForOpponent:
			// Expire rooms waiting too long for an opponent
			if now.Sub(room.CreatedAt) > m.roomTTL {
				delete(m.byCode, code)
				delete(m.byID, room.ID)
			}
		case model.RoomGameOver:
			// Clean up finished game rooms after 60 seconds
			if now.Sub(room.CreatedAt) > m.roomTTL+60*time.Second {
				delete(m.byCode, code)
				delete(m.byID, room.ID)
			}
		}
	}
}

// Stop shuts down the cleanup goroutine.
func (m *Manager) Stop() {
	close(m.stopChan)
}
