package api

import (
	"net/http"

	"github.com/teomiscia/hexbattle/internal/lobby"
)

// MatchmakingHandler handles matchmaking queue operations.
type MatchmakingHandler struct {
	Queue *lobby.MatchmakingQueue
}

// HandleJoin handles POST /api/v1/matchmaking/join.
func (h *MatchmakingHandler) HandleJoin(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		respondError(w, http.StatusMethodNotAllowed, "METHOD_NOT_ALLOWED", "only POST is allowed")
		return
	}

	session := SessionFromContext(r.Context())
	if session == nil {
		respondError(w, http.StatusUnauthorized, "UNAUTHORIZED", "missing session")
		return
	}

	result, err := h.Queue.Join(session.ID, session.Nickname)
	if err != nil {
		respondError(w, http.StatusBadRequest, "INVALID_REQUEST", err.Error())
		return
	}

	if result != nil {
		// Match found immediately
		respondJSON(w, http.StatusOK, map[string]interface{}{
			"status":    "matched",
			"room_id":   result.RoomID,
			"room_code": result.RoomCode,
		})
		return
	}

	// Queued, waiting for opponent
	respondJSON(w, http.StatusAccepted, map[string]interface{}{
		"status": "queued",
	})
}

// HandleLeave handles DELETE /api/v1/matchmaking/leave.
func (h *MatchmakingHandler) HandleLeave(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		respondError(w, http.StatusMethodNotAllowed, "METHOD_NOT_ALLOWED", "only DELETE is allowed")
		return
	}

	session := SessionFromContext(r.Context())
	if session == nil {
		respondError(w, http.StatusUnauthorized, "UNAUTHORIZED", "missing session")
		return
	}

	removed := h.Queue.Leave(session.ID)
	if !removed {
		respondError(w, http.StatusNotFound, "NOT_FOUND", "player not in matchmaking queue")
		return
	}

	respondJSON(w, http.StatusOK, map[string]interface{}{
		"status": "removed",
	})
}

// HandleStatus handles GET /api/v1/matchmaking/status.
func (h *MatchmakingHandler) HandleStatus(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		respondError(w, http.StatusMethodNotAllowed, "METHOD_NOT_ALLOWED", "only GET is allowed")
		return
	}

	session := SessionFromContext(r.Context())
	if session == nil {
		respondError(w, http.StatusUnauthorized, "UNAUTHORIZED", "missing session")
		return
	}

	queued := h.Queue.IsQueued(session.ID)

	// Check if the player is in a room
	var matchedRoomID string
	var matchedRoomCode string
	if !queued {
		rooms := h.Queue.Manager().GetAllRooms()
		for _, room := range rooms {
			if room.HasPlayer(session.ID) && room.State == "ready" {
				matchedRoomID = room.ID
				matchedRoomCode = room.Code
				break
			}
		}
	}

	if matchedRoomID != "" {
		respondJSON(w, http.StatusOK, map[string]interface{}{
			"queued":     false,
			"queue_size": h.Queue.Size(),
			"room_id":    matchedRoomID,
			"room_code":  matchedRoomCode,
		})
		return
	}

	respondJSON(w, http.StatusOK, map[string]interface{}{
		"queued":     queued,
		"queue_size": h.Queue.Size(),
	})
}
