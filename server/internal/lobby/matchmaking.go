package lobby

import (
	"fmt"
	"sync"

	"github.com/teomiscia/hexbattle/internal/model"
)

// QueueEntry represents a player waiting in the matchmaking queue.
type QueueEntry struct {
	PlayerID string
	Nickname string
}

// MatchResult is returned when two players are matched.
type MatchResult struct {
	RoomID   string
	RoomCode string
	Player1  QueueEntry
	Player2  QueueEntry
}

// MatchmakingQueue implements a simple FIFO matchmaking queue.
// When a second player joins, they are immediately paired with the first.
type MatchmakingQueue struct {
	mu      sync.Mutex
	queue   []QueueEntry
	manager *Manager
}

// NewMatchmakingQueue creates a new matchmaking queue backed by the given lobby manager.
func NewMatchmakingQueue(manager *Manager) *MatchmakingQueue {
	return &MatchmakingQueue{
		queue:   make([]QueueEntry, 0),
		manager: manager,
	}
}

// Join adds a player to the matchmaking queue.
// If another player is already waiting, they are immediately matched and a room is created.
// Returns a MatchResult if matched, nil if queued.
func (mq *MatchmakingQueue) Join(playerID, nickname string) (*MatchResult, error) {
	mq.mu.Lock()
	defer mq.mu.Unlock()

	// Check if player is already in the queue
	for _, entry := range mq.queue {
		if entry.PlayerID == playerID {
			return nil, fmt.Errorf("player already in matchmaking queue")
		}
	}

	newEntry := QueueEntry{PlayerID: playerID, Nickname: nickname}

	// If someone is already waiting, match them
	if len(mq.queue) > 0 {
		waiting := mq.queue[0]
		mq.queue = mq.queue[1:]

		// Create a room with Quick Match defaults
		settings := model.DefaultRoomSettings()
		room, err := mq.manager.CreateRoom(waiting.PlayerID, waiting.Nickname, settings)
		if err != nil {
			// Put the waiting player back and return error
			mq.queue = append([]QueueEntry{waiting}, mq.queue...)
			return nil, fmt.Errorf("failed to create match room: %w", err)
		}

		// Join the new player to the room
		_, err = mq.manager.JoinRoom(room.Code, newEntry.PlayerID, newEntry.Nickname)
		if err != nil {
			return nil, fmt.Errorf("failed to join match room: %w", err)
		}

		return &MatchResult{
			RoomID:   room.ID,
			RoomCode: room.Code,
			Player1:  waiting,
			Player2:  newEntry,
		}, nil
	}

	// No one waiting — add to queue
	mq.queue = append(mq.queue, newEntry)
	return nil, nil
}

// Leave removes a player from the matchmaking queue.
// Returns true if the player was found and removed.
func (mq *MatchmakingQueue) Leave(playerID string) bool {
	mq.mu.Lock()
	defer mq.mu.Unlock()

	for i, entry := range mq.queue {
		if entry.PlayerID == playerID {
			mq.queue = append(mq.queue[:i], mq.queue[i+1:]...)
			return true
		}
	}
	return false
}

// Size returns the number of players waiting in the queue.
func (mq *MatchmakingQueue) Size() int {
	mq.mu.Lock()
	defer mq.mu.Unlock()
	return len(mq.queue)
}

// IsQueued returns true if the given player is in the queue.
func (mq *MatchmakingQueue) IsQueued(playerID string) bool {
	mq.mu.Lock()
	defer mq.mu.Unlock()
	for _, entry := range mq.queue {
		if entry.PlayerID == playerID {
			return true
		}
	}
	return false
}

// Manager returns the lobby manager attached to the queue.
func (mq *MatchmakingQueue) Manager() *Manager {
	return mq.manager
}
