package ws

import (
	"context"
	"log/slog"
	"net/http"
	"time"

	"nhooyr.io/websocket"

	"github.com/teomiscia/hexbattle/internal/player"
)

// Handler manages WebSocket upgrade requests and connection lifecycle.
type Handler struct {
	registry     *player.Registry
	pingInterval time.Duration
	pongTimeout  time.Duration

	// OnConnect is called when a new authenticated WebSocket connection is established.
	// The callback receives the connection and should set up OnMessage/OnDisconnect handlers.
	OnConnect func(conn *Connection)
}

// NewHandler creates a new WebSocket upgrade handler.
func NewHandler(registry *player.Registry, pingInterval, pongTimeout time.Duration) *Handler {
	return &Handler{
		registry:     registry,
		pingInterval: pingInterval,
		pongTimeout:  pongTimeout,
	}
}

// ServeHTTP handles the WebSocket upgrade request.
// It validates the auth token from the query parameter and upgrades the connection.
func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Extract token from query parameter
	token := r.URL.Query().Get("token")
	if token == "" {
		http.Error(w, "missing token", http.StatusUnauthorized)
		return
	}

	// Authenticate
	session, err := h.registry.Authenticate(token)
	if err != nil {
		http.Error(w, "invalid token", http.StatusUnauthorized)
		return
	}

	// Upgrade to WebSocket
	conn, err := websocket.Accept(w, r, &websocket.AcceptOptions{
		OriginPatterns: []string{"*"}, // CORS handled at Nginx level in production
	})
	if err != nil {
		slog.Error("websocket upgrade failed",
			"player_id", session.ID,
			"error", err,
		)
		return
	}

	slog.Info("websocket connected",
		"player_id", session.ID,
		"nickname", session.Nickname,
	)

	// Create managed connection
	wsConn := NewConnection(r.Context(), conn, session.ID)

	// Set up the connection callback
	if h.OnConnect != nil {
		h.OnConnect(wsConn)
	}

	// Start read/write goroutines
	wsConn.Start()

	// Start ping/pong heartbeat
	go h.heartbeat(wsConn)

	// Block until the connection is closed
	<-wsConn.ctx.Done()
}

// heartbeat sends periodic pings and monitors for pong responses.
func (h *Handler) heartbeat(conn *Connection) {
	ticker := time.NewTicker(h.pingInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			// Send ping message
			pingData, err := NewEnvelope(MsgPing, struct{}{})
			if err != nil {
				continue
			}
			if !conn.Send(pingData) {
				return
			}

			// Also use the WebSocket-level ping for connection health
			ctx, cancel := context.WithTimeout(conn.ctx, h.pongTimeout)
			err = conn.Conn.Ping(ctx)
			cancel()
			if err != nil {
				slog.Debug("ping failed, closing connection",
					"player_id", conn.PlayerID,
					"error", err,
				)
				conn.Close()
				return
			}

		case <-conn.ctx.Done():
			return
		}
	}
}
