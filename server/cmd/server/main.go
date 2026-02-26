package main

import (
	"context"
	"encoding/json"
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
	"github.com/teomiscia/hexbattle/internal/mapgen"
	"github.com/teomiscia/hexbattle/internal/model"
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
	var st store.Store
	if err != nil {
		slog.Warn("failed to connect to redis, running without persistence", "error", err)
	} else {
		slog.Info("redis connected")
		st = redisStore
	}

	// 4. Initialize registries and managers
	registry := player.NewRegistry()
	lobbyManager := lobby.NewManager(cfg.RoomTTL)
	defer lobbyManager.Stop()
	matchQueue := lobby.NewMatchmakingQueue(lobbyManager)
	gameManager := game.NewManager(st)

	// Restore active games if persistence is enabled
	if st != nil {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		if err := gameManager.RestoreActiveGames(ctx); err != nil {
			slog.Error("failed to restore active games", "error", err)
		}
		cancel()
	}

	// 5. Set up WebSocket handler
	wsHandler := ws.NewHandler(registry, cfg.WSPingInterval, cfg.WSPongTimeout, cfg.CORSOrigins)
	wsHandler.OnConnect = func(conn *ws.Connection) {
		slog.Debug("new websocket connection", "player_id", conn.PlayerID)

		conn.OnMessage = func(playerID string, env ws.Envelope) {
			slog.Debug("ws message received", "player_id", playerID, "type", env.Type)

			// Intercept join_game to route or create the game engine
			if env.Type == ws.MsgJoinGame {
				var data ws.JoinGameData
				if err := json.Unmarshal(env.Data, &data); err != nil {
					conn.SendNack(env.Seq, env.Type, string(model.ErrInvalidMessage), "invalid join_game data")
					return
				}

				room := lobbyManager.GetByID(data.RoomID)
				if room == nil {
					conn.SendNack(env.Seq, env.Type, string(model.ErrRoomNotFound), "room not found")
					return
				}

				if !room.HasPlayer(playerID) {
					conn.SendNack(env.Seq, env.Type, string(model.ErrNotYourTurn), "you are not in this room")
					return
				}

				var engine *game.Engine

				// If the game doesn't exist yet, we might need to create it
				if room.GameID == "" {
					// Only start if the room is full
					if !room.IsFull() {
						conn.SendNack(env.Seq, env.Type, string(model.ErrInvalidMessage), "waiting for opponent")
						return
					}

					// Lock to prevent race condition if both players send join_game simultaneously
					// Since lobby Room is mutable, we need to handle this carefully.
					// For simplicity, generate ID, and use GetEngine to see if it was just created.
					newGameID, _ := player.GenerateGameID()
					room.GameID = newGameID
					lobbyManager.SetGameInProgress(room.ID, newGameID)

					// Build players
					p1Session := registry.GetByID(room.HostPlayerID)
					p2Session := registry.GetByID(room.GuestPlayerID)

					// Fallbacks if sessions are somehow missing (shouldn't happen)
					p1Name := room.HostNickname
					if p1Session != nil {
						p1Name = p1Session.Nickname
					}
					p2Name := room.GuestNickname
					if p2Session != nil {
						p2Name = p2Session.Nickname
					}

					p1 := model.PlayerState{ID: room.HostPlayerID, Nickname: p1Name}
					p2 := model.PlayerState{ID: room.GuestPlayerID, Nickname: p2Name}

					// Create Game State
					seed := time.Now().UnixNano()
					state := game.NewGameState(newGameID, room.Settings, p1, p2, seed)

					// Generate map
					mapResult, err := mapgen.Generate(state.MapSize, seed, balance)
					if err != nil {
						conn.SendNack(env.Seq, env.Type, string(model.ErrInvalidMessage), "failed to generate map")
						return
					}

					// Apply map to state
					state.Terrain = mapResult.Terrain
					// Apply structures (HQs and Neutral)
					for _, sp := range mapResult.Structures {
						owner := sp.OwnerID
						if sp.Type == model.StructureHQ {
							// First HQ goes to P1, second to P2
							if state.PlayerHQ(p1.ID) == nil {
								owner = p1.ID
							} else {
								owner = p2.ID
							}
						}
						id := player.GenerateStructureID()
						s, _ := game.NewStructureFromBalance(id, sp.Type, owner, sp.Position)
						state.AddStructure(s)
					}

					hub := ws.NewHub()
					engine = game.NewEngine(context.Background(), state, hub, st)
					gameManager.AddEngine(engine)
				} else {
					engine = gameManager.GetEngine(room.GameID)
					if engine == nil {
						conn.SendNack(env.Seq, env.Type, string(model.ErrGameNotFound), "game engine not found")
						return
					}
				}

				// Register connection with the engine's hub
				conn.GameID = engine.State.ID
				engine.Hub.Register(conn)

				// Submit the join_game action to the engine
				engine.SubmitAction(game.PlayerAction{
					PlayerID: playerID,
					Seq:      env.Seq,
					Type:     env.Type,
					Data:     env.Data,
					Conn:     conn,
				})

				return
			} else if env.Type == ws.MsgReconnect {
				// Handle reconnect explicitly
				var data ws.ReconnectData
				if err := json.Unmarshal(env.Data, &data); err != nil {
					conn.SendNack(env.Seq, env.Type, string(model.ErrInvalidMessage), "invalid reconnect data")
					return
				}

				engine := gameManager.GetEngine(data.GameID)
				if engine == nil {
					conn.SendNack(env.Seq, env.Type, string(model.ErrGameNotFound), "game not found")
					return
				}

				// The player's token is already validated by the ws.Handler
				// We just need to verify they are part of the game
				if engine.State.PlayerIndex(playerID) < 0 {
					conn.SendNack(env.Seq, env.Type, string(model.ErrNotYourTurn), "you are not in this game")
					return
				}

				conn.GameID = engine.State.ID
				conn.SendAck(env.Seq, env.Type)

				// Notify engine to resume connection
				engine.NotifyReconnect(game.ReconnectEvent{
					PlayerID: playerID,
					Conn:     conn,
				})
				return
			} else if env.Type == ws.MsgPong {
				// Pong is a heartbeat response and should be allowed even before joining a game
				// No need to send ack for pong, just handle it silently
				return
			}

			// For all other messages, route to the associated engine
			if conn.GameID == "" {
				conn.SendNack(env.Seq, env.Type, string(model.ErrInvalidMessage), "connection not associated with a game")
				return
			}

			engine := gameManager.GetEngine(conn.GameID)
			if engine != nil {
				engine.SubmitAction(game.PlayerAction{
					PlayerID: playerID,
					Seq:      env.Seq,
					Type:     env.Type,
					Data:     env.Data,
					Conn:     conn,
				})
			}
		}

		conn.OnDisconnect = func(playerID string) {
			slog.Info("player disconnected", "player_id", playerID, "game_id", conn.GameID)
			if conn.GameID != "" {
				engine := gameManager.GetEngine(conn.GameID)
				if engine != nil {
					engine.Hub.Unregister(playerID)
					engine.NotifyDisconnect(playerID)
				}
			}
		}
	}

	// 6. Set up HTTP router
	startTime := time.Now()
	router := api.NewRouter(api.RouterConfig{
		Registry:    registry,
		Lobby:       lobbyManager,
		Queue:       matchQueue,
		Store:       st,
		WSHandler:   wsHandler,
		CORSOrigins: cfg.CORSOrigins,
		StartTime:   startTime,
	})

	// 7. Create HTTP server
	handler := router.Handler(cfg.CORSOrigins)
	server := &http.Server{
		Addr:         ":" + cfg.Port,
		Handler:      handler,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// 8. Start server in a goroutine
	go func() {
		slog.Info("server listening", "addr", server.Addr)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			slog.Error("server error", "error", err)
			os.Exit(1)
		}
	}()

	// 9. Wait for shutdown signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	sig := <-quit

	slog.Info("shutdown signal received", "signal", sig.String())

	// 10. Graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), cfg.ShutdownDrainTimeout)
	defer cancel()

	// Stop accepting new connections
	if err := server.Shutdown(ctx); err != nil {
		slog.Error("server shutdown error", "error", err)
	}

	// Stop all active games and force snapshots
	gameManager.StopAll()

	// Close Redis connection
	if st != nil {
		if err := st.Close(); err != nil {
			slog.Error("redis close error", "error", err)
		}
	}

	slog.Info("server stopped")
}
