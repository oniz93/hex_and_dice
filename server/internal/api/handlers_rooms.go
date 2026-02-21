package api

import (
	"net/http"
	"strings"

	"github.com/teomiscia/hexbattle/internal/lobby"
	"github.com/teomiscia/hexbattle/internal/model"
)

// RoomsHandler handles room creation, joining, and status.
type RoomsHandler struct {
	Lobby *lobby.Manager
}

// CreateRoomRequest is the request body for creating a room.
type CreateRoomRequest struct {
	MapSize   string `json:"map_size"`
	TurnTimer int    `json:"turn_timer"`
	TurnMode  string `json:"turn_mode"`
}

// JoinRoomRequest is the request body for joining a room.
type JoinRoomRequest struct {
	Code string `json:"code"`
}

// HandleCreate handles POST /api/v1/rooms.
func (h *RoomsHandler) HandleCreate(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		respondError(w, http.StatusMethodNotAllowed, "METHOD_NOT_ALLOWED", "only POST is allowed")
		return
	}

	session := SessionFromContext(r.Context())
	if session == nil {
		respondError(w, http.StatusUnauthorized, "UNAUTHORIZED", "missing session")
		return
	}

	var req CreateRoomRequest
	if err := decodeJSON(r, &req); err != nil {
		respondError(w, http.StatusBadRequest, "INVALID_REQUEST", "invalid request body")
		return
	}

	// Parse settings with defaults
	settings := model.DefaultRoomSettings()
	if req.MapSize != "" {
		switch model.MapSize(req.MapSize) {
		case model.MapSizeSmall, model.MapSizeMedium, model.MapSizeLarge:
			settings.MapSize = model.MapSize(req.MapSize)
		default:
			respondError(w, http.StatusBadRequest, "INVALID_REQUEST", "invalid map size")
			return
		}
	}
	if req.TurnTimer > 0 {
		switch req.TurnTimer {
		case 60, 90, 120:
			settings.TurnTimer = req.TurnTimer
		default:
			respondError(w, http.StatusBadRequest, "INVALID_REQUEST", "turn timer must be 60, 90, or 120")
			return
		}
	}
	if req.TurnMode != "" {
		if model.TurnMode(req.TurnMode) == model.TurnModeAlternating {
			settings.TurnMode = model.TurnModeAlternating
		} else {
			respondError(w, http.StatusBadRequest, "INVALID_REQUEST", "only alternating turn mode is supported")
			return
		}
	}

	room, err := h.Lobby.CreateRoom(session.ID, session.Nickname, settings)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "INTERNAL_ERROR", err.Error())
		return
	}

	respondJSON(w, http.StatusCreated, lobby.CreateRoomResponse{
		RoomCode: room.Code,
		RoomID:   room.ID,
		Settings: room.Settings,
	})
}

// HandleJoin handles POST /api/v1/rooms/join.
func (h *RoomsHandler) HandleJoin(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		respondError(w, http.StatusMethodNotAllowed, "METHOD_NOT_ALLOWED", "only POST is allowed")
		return
	}

	session := SessionFromContext(r.Context())
	if session == nil {
		respondError(w, http.StatusUnauthorized, "UNAUTHORIZED", "missing session")
		return
	}

	var req JoinRoomRequest
	if err := decodeJSON(r, &req); err != nil {
		respondError(w, http.StatusBadRequest, "INVALID_REQUEST", "invalid request body")
		return
	}

	if req.Code == "" {
		respondError(w, http.StatusBadRequest, "INVALID_REQUEST", "room code is required")
		return
	}

	room, err := h.Lobby.JoinRoom(req.Code, session.ID, session.Nickname)
	if err != nil {
		msg := err.Error()
		switch {
		case strings.Contains(msg, "not found"):
			respondError(w, http.StatusNotFound, string(model.ErrRoomNotFound), msg)
		case strings.Contains(msg, "full"):
			respondError(w, http.StatusConflict, string(model.ErrRoomFull), msg)
		default:
			respondError(w, http.StatusBadRequest, "INVALID_REQUEST", msg)
		}
		return
	}

	respondJSON(w, http.StatusOK, lobby.JoinRoomResponse{
		RoomID:       room.ID,
		Settings:     room.Settings,
		HostNickname: room.HostNickname,
	})
}

// HandleGetStatus handles GET /api/v1/rooms/{code}.
func (h *RoomsHandler) HandleGetStatus(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		respondError(w, http.StatusMethodNotAllowed, "METHOD_NOT_ALLOWED", "only GET is allowed")
		return
	}

	// Extract room code from path
	path := r.URL.Path
	parts := strings.Split(strings.TrimSuffix(path, "/"), "/")
	if len(parts) == 0 {
		respondError(w, http.StatusBadRequest, "INVALID_REQUEST", "room code is required")
		return
	}
	code := parts[len(parts)-1]

	room := h.Lobby.GetByCode(code)
	if room == nil {
		respondError(w, http.StatusNotFound, string(model.ErrRoomNotFound), "room not found")
		return
	}

	respondJSON(w, http.StatusOK, lobby.RoomStatusResponse{
		Code:          room.Code,
		State:         room.State,
		Settings:      room.Settings,
		HostNickname:  room.HostNickname,
		GuestNickname: room.GuestNickname,
		GameID:        room.GameID,
	})
}
