package api

import (
	"net/http"

	"github.com/teomiscia/hexbattle/internal/player"
)

// GuestHandler handles guest registration.
type GuestHandler struct {
	Registry *player.Registry
}

// GuestRequest is the request body for guest registration.
type GuestRequest struct {
	Nickname string `json:"nickname"`
}

// ServeHTTP handles POST /api/v1/guest.
func (h *GuestHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		respondError(w, http.StatusMethodNotAllowed, "METHOD_NOT_ALLOWED", "only POST is allowed")
		return
	}

	var req GuestRequest
	if err := decodeJSON(r, &req); err != nil {
		respondError(w, http.StatusBadRequest, "INVALID_REQUEST", "invalid request body")
		return
	}

	session, err := player.NewSession(req.Nickname)
	if err != nil {
		respondError(w, http.StatusBadRequest, "INVALID_NICKNAME", err.Error())
		return
	}

	h.Registry.Register(session)

	respondJSON(w, http.StatusCreated, session.ToResponse())
}
