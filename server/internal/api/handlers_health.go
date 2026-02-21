package api

import (
	"context"
	"net/http"
	"time"

	"github.com/teomiscia/hexbattle/internal/lobby"
	"github.com/teomiscia/hexbattle/internal/player"
	"github.com/teomiscia/hexbattle/internal/store"
)

// HealthHandler handles the health check endpoint.
type HealthHandler struct {
	Registry  *player.Registry
	Lobby     *lobby.Manager
	Queue     *lobby.MatchmakingQueue
	Store     store.Store
	StartTime time.Time
}

// ServeHTTP handles GET /health.
func (h *HealthHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		respondError(w, http.StatusMethodNotAllowed, "METHOD_NOT_ALLOWED", "only GET is allowed")
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 3*time.Second)
	defer cancel()

	redisConnected := false
	var redisLatency int64

	if h.Store != nil {
		if rs, ok := h.Store.(*store.RedisStore); ok {
			latency, err := rs.PingLatency(ctx)
			if err == nil {
				redisConnected = true
				redisLatency = latency
			}
		}
	}

	status := http.StatusOK
	statusText := "healthy"
	if !redisConnected {
		status = http.StatusServiceUnavailable
		statusText = "degraded"
	}

	uptime := time.Since(h.StartTime).Seconds()

	respondJSON(w, status, map[string]interface{}{
		"status":                 statusText,
		"uptime_seconds":         int(uptime),
		"connected_players":      h.Registry.Count(),
		"waiting_rooms":          h.Lobby.WaitingRoomCount(),
		"matchmaking_queue_size": h.Queue.Size(),
		"redis": map[string]interface{}{
			"connected":  redisConnected,
			"latency_ms": redisLatency,
		},
	})
}
