package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/teomiscia/hexbattle/internal/api"
	"github.com/teomiscia/hexbattle/internal/config"
	"github.com/teomiscia/hexbattle/internal/game"
	"github.com/teomiscia/hexbattle/internal/lobby"
	"github.com/teomiscia/hexbattle/internal/player"
	"github.com/teomiscia/hexbattle/internal/store"
	"github.com/teomiscia/hexbattle/internal/ws"
)

func main() {
	// 1. Parse configuration
	cfg := config.Load()

	// Configure logger
	var logLevel slog.Level
	switch cfg.LogLevel {
	case "DEBUG":
		logLevel = slog.LevelDebug
	case "WARN":
		logLevel = slog.LevelWarn
	case "ERROR":
		logLevel = slog.LevelError
	default:
		logLevel = slog.LevelInfo
	}
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: logLevel}))
	slog.SetDefault(logger)

	slog.Info("starting server",
		"port", cfg.Port,
		"log_level", cfg.LogLevel,
		"redis_url", cfg.RedisURL,
		"balance_file", cfg.BalanceFile,
	)

	// 2. Load balance data
	balance, err := config.LoadBalance(cfg.BalanceFile)
	if err != nil {
		slog.Error("failed to load balance data", "error", err)
		os.Exit(1)
	}
	game.LoadBalance(balance)
	slog.Info("balance data loaded", "file", cfg.BalanceFile)

	// 3. Connect to Redis
	redisStore, err := store.NewRedisStore(cfg.RedisURL)
	if err != nil {
		slog.Warn("failed to connect to redis, running without persistence", "error", err)
	} else {
		slog.Info("redis connected")
	}

	// 4. Initialize player registry
	registry := player.NewRegistry()

	// 5. Initialize lobby manager
	lobbyManager := lobby.NewManager(cfg.RoomTTL)
	defer lobbyManager.Stop()

	// 6. Initialize matchmaking queue
	matchQueue := lobby.NewMatchmakingQueue(lobbyManager)

	// 7. Set up WebSocket handler
	wsHandler := ws.NewHandler(registry, cfg.WSPingInterval, cfg.WSPongTimeout)
	wsHandler.OnConnect = func(conn *ws.Connection) {
		slog.Debug("new websocket connection",
			"player_id", conn.PlayerID,
		)
		// The connection's OnMessage and OnDisconnect will be set up
		// when the player joins a game via the join_game message.
		// For now, set a default message handler that routes to game engines.
		conn.OnMessage = func(playerID string, env ws.Envelope) {
			// Route messages based on type
			// For join_game, we need to look up the game engine and register
			// For other messages, route to the associated game engine
			slog.Debug("ws message received",
				"player_id", playerID,
				"type", env.Type,
			)
		}
		conn.OnDisconnect = func(playerID string) {
			slog.Info("player disconnected",
				"player_id", playerID,
			)
		}
	}

	// 8. Set up HTTP router
	startTime := time.Now()
	var st store.Store
	if redisStore != nil {
		st = redisStore
	}

	router := api.NewRouter(api.RouterConfig{
		Registry:    registry,
		Lobby:       lobbyManager,
		Queue:       matchQueue,
		Store:       st,
		WSHandler:   wsHandler,
		CORSOrigins: cfg.CORSOrigins,
		StartTime:   startTime,
	})

	// 9. Create HTTP server
	handler := router.Handler(cfg.CORSOrigins)
	server := &http.Server{
		Addr:         ":" + cfg.Port,
		Handler:      handler,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// 10. Start server in a goroutine
	go func() {
		slog.Info("server listening", "addr", server.Addr)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			slog.Error("server error", "error", err)
			os.Exit(1)
		}
	}()

	// 11. Wait for shutdown signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	sig := <-quit

	slog.Info("shutdown signal received", "signal", sig.String())

	// 12. Graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), cfg.ShutdownDrainTimeout)
	defer cancel()

	// Stop accepting new connections
	if err := server.Shutdown(ctx); err != nil {
		slog.Error("server shutdown error", "error", err)
	}

	// Close Redis connection
	if redisStore != nil {
		if err := redisStore.Close(); err != nil {
			slog.Error("redis close error", "error", err)
		}
	}

	slog.Info("server stopped")
}
