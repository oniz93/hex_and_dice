package api

import (
	"net/http"

	"github.com/teomiscia/hexbattle/internal/lobby"
	"github.com/teomiscia/hexbattle/internal/player"
	"github.com/teomiscia/hexbattle/internal/store"
	"github.com/teomiscia/hexbattle/internal/ws"
	"time"
)

// Router sets up all HTTP routes for the server.
type Router struct {
	Mux       *http.ServeMux
	Registry  *player.Registry
	Lobby     *lobby.Manager
	Queue     *lobby.MatchmakingQueue
	Store     store.Store
	WSHandler *ws.Handler
}

// RouterConfig holds the dependencies needed to create the router.
type RouterConfig struct {
	Registry    *player.Registry
	Lobby       *lobby.Manager
	Queue       *lobby.MatchmakingQueue
	Store       store.Store
	WSHandler   *ws.Handler
	CORSOrigins []string
	StartTime   time.Time
}

// NewRouter creates and configures the HTTP router with all routes.
func NewRouter(cfg RouterConfig) *Router {
	mux := http.NewServeMux()

	r := &Router{
		Mux:       mux,
		Registry:  cfg.Registry,
		Lobby:     cfg.Lobby,
		Queue:     cfg.Queue,
		Store:     cfg.Store,
		WSHandler: cfg.WSHandler,
	}

	// Create handlers
	guestHandler := &GuestHandler{Registry: cfg.Registry}
	roomsHandler := &RoomsHandler{Lobby: cfg.Lobby}
	matchmakingHandler := &MatchmakingHandler{Queue: cfg.Queue}
	healthHandler := &HealthHandler{
		Registry:  cfg.Registry,
		Lobby:     cfg.Lobby,
		Queue:     cfg.Queue,
		Store:     cfg.Store,
		StartTime: cfg.StartTime,
	}

	// Auth middleware wrapper
	authMW := AuthMiddleware(cfg.Registry)

	// --- Public routes ---
	mux.Handle("POST /api/v1/guest", guestHandler)
	mux.Handle("GET /health", healthHandler)

	// --- Protected routes ---
	mux.Handle("POST /api/v1/rooms", authMW(http.HandlerFunc(roomsHandler.HandleCreate)))
	mux.Handle("POST /api/v1/rooms/join", authMW(http.HandlerFunc(roomsHandler.HandleJoin)))
	mux.Handle("POST /api/v1/rooms/bot", authMW(http.HandlerFunc(roomsHandler.HandleCreateBotGame)))
	mux.Handle("GET /api/v1/rooms/", authMW(http.HandlerFunc(roomsHandler.HandleGetStatus)))

	mux.Handle("POST /api/v1/matchmaking/join", authMW(http.HandlerFunc(matchmakingHandler.HandleJoin)))
	mux.Handle("DELETE /api/v1/matchmaking/leave", authMW(http.HandlerFunc(matchmakingHandler.HandleLeave)))
	mux.Handle("GET /api/v1/matchmaking/status", authMW(http.HandlerFunc(matchmakingHandler.HandleStatus)))

	// --- WebSocket ---
	mux.Handle("GET /ws", cfg.WSHandler)

	return r
}

// Handler returns the top-level HTTP handler with global middleware applied.
func (r *Router) Handler(corsOrigins []string) http.Handler {
	var handler http.Handler = r.Mux
	handler = CORSMiddleware(corsOrigins)(handler)
	handler = LoggingMiddleware(handler)
	return handler
}
