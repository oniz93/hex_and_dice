package game

import (
	"context"
	"encoding/json"
	"log/slog"
	"time"

	"github.com/teomiscia/hexbattle/internal/dice"
	"github.com/teomiscia/hexbattle/internal/hex"
	"github.com/teomiscia/hexbattle/internal/model"
	"github.com/teomiscia/hexbattle/internal/store"
	"github.com/teomiscia/hexbattle/internal/ws"
)

// Engine runs a single game instance in its own goroutine.
type Engine struct {
	State  *GameState
	Hub    *ws.Hub
	Roller *dice.Roller
	Store  store.Store

	actionChan     chan PlayerAction
	disconnectChan chan string
	reconnectChan  chan ReconnectEvent
	turnTimer      *time.Timer
	reconnectTimer *time.Timer
	disconnectedID string
	ctx            context.Context
	cancel         context.CancelFunc
	logger         *slog.Logger
}

// PlayerAction wraps an incoming action from a player.
type PlayerAction struct {
	PlayerID string
	Seq      int
	Type     string
	Data     json.RawMessage
	Conn     *ws.Connection
}

// ReconnectEvent signals that a player has reconnected.
type ReconnectEvent struct {
	PlayerID string
	Conn     *ws.Connection
}

// NewEngine creates a new game engine for the given state.
func NewEngine(ctx context.Context, state *GameState, hub *ws.Hub, st store.Store) *Engine {
	ctx, cancel := context.WithCancel(ctx)
	return &Engine{
		State:          state,
		Hub:            hub,
		Roller:         dice.NewRoller(state.Seed),
		Store:          st,
		actionChan:     make(chan PlayerAction, 32),
		disconnectChan: make(chan string, 2),
		reconnectChan:  make(chan ReconnectEvent, 2),
		ctx:            ctx,
		cancel:         cancel,
		logger: slog.Default().With(
			"game_id", state.ID,
		),
	}
}

// SubmitAction sends a player action to the engine's event loop.
func (e *Engine) SubmitAction(action PlayerAction) {
	select {
	case e.actionChan <- action:
	default:
		e.logger.Warn("action channel full, dropping action",
			"player_id", action.PlayerID,
			"type", action.Type,
		)
	}
}

// NotifyDisconnect signals that a player has disconnected.
func (e *Engine) NotifyDisconnect(playerID string) {
	select {
	case e.disconnectChan <- playerID:
	default:
	}
}

// NotifyReconnect signals that a player has reconnected.
func (e *Engine) NotifyReconnect(event ReconnectEvent) {
	select {
	case e.reconnectChan <- event:
	default:
	}
}

// Stop signals the engine to shut down.
func (e *Engine) Stop() {
	e.cancel()
}

// Run starts the game engine event loop. This should be called in a goroutine.
func (e *Engine) Run() {
	e.logger.Info("game engine started",
		"players", []string{e.State.Players[0].ID, e.State.Players[1].ID},
		"map_size", e.State.MapSize,
	)

	defer func() {
		if e.turnTimer != nil {
			e.turnTimer.Stop()
		}
		if e.reconnectTimer != nil {
			e.reconnectTimer.Stop()
		}
		e.logger.Info("game engine stopped")
	}()

	// If the game is already in progress (restored from snapshot), resume
	if e.State.Phase == model.PhasePlayerAction {
		e.startTurnTimer()
	}

	for {
		select {
		case action := <-e.actionChan:
			e.handleAction(action)

		case <-e.turnTimerChan():
			e.handleTurnTimeout()

		case playerID := <-e.disconnectChan:
			e.handleDisconnect(playerID)

		case event := <-e.reconnectChan:
			e.handleReconnect(event)

		case <-e.reconnectTimerChan():
			e.handleReconnectTimeout()

		case <-e.ctx.Done():
			e.snapshotState()
			return
		}
	}
}

// turnTimerChan returns the turn timer's channel, or a nil channel if no timer is active.
func (e *Engine) turnTimerChan() <-chan time.Time {
	if e.turnTimer == nil {
		return nil
	}
	return e.turnTimer.C
}

// reconnectTimerChan returns the reconnect timer's channel.
func (e *Engine) reconnectTimerChan() <-chan time.Time {
	if e.reconnectTimer == nil {
		return nil
	}
	return e.reconnectTimer.C
}

// handleAction dispatches a player action to the appropriate handler.
func (e *Engine) handleAction(action PlayerAction) {
	e.logger.Debug("processing action",
		"player_id", action.PlayerID,
		"type", action.Type,
		"seq", action.Seq,
		"phase", e.State.Phase,
	)

	switch action.Type {
	case ws.MsgJoinGame:
		e.handleJoinGame(action)
	case ws.MsgMove:
		e.handleMove(action)
	case ws.MsgAttack:
		e.handleAttack(action)
	case ws.MsgBuy:
		e.handleBuy(action)
	case ws.MsgEndTurn:
		e.handleEndTurn(action)
	case ws.MsgEmote:
		e.handleEmote(action)
	case ws.MsgPong:
		// No-op, handled at connection level
	default:
		e.sendNack(action, string(model.ErrInvalidMessage), "unknown message type")
	}
}

// handleJoinGame processes a join_game message.
func (e *Engine) handleJoinGame(action PlayerAction) {
	// Player is joining the game — send them the full state
	e.sendAck(action)
	e.sendFullState(action.PlayerID)

	// Check if both players are connected
	if e.Hub.ConnectedCount() >= 2 && e.State.Phase == model.PhaseWaitingForPlayers {
		e.startGame()
	}
}

// startGame transitions from WaitingForPlayers to the first turn.
func (e *Engine) startGame() {
	e.State.Phase = model.PhaseGeneratingMap
	e.logger.Info("both players connected, game starting")

	// Map should already be generated before engine starts
	e.State.Phase = model.PhaseGameStarted

	// Broadcast full game state to both players
	e.broadcastFullState()

	// Start first turn
	e.State.TurnNumber = 1
	turnStart := RunTurnStart(e.State, e.Roller)

	e.State.Phase = model.PhasePlayerAction
	e.Hub.BroadcastMessage(ws.MsgTurnStart, turnStart)
	e.startTurnTimer()
}

// handleMove processes a move action.
func (e *Engine) handleMove(action PlayerAction) {
	var data ws.MoveData
	if err := json.Unmarshal(action.Data, &data); err != nil {
		e.sendNack(action, string(model.ErrInvalidMessage), "invalid move data")
		return
	}

	target := hex.NewCoord(data.TargetQ, data.TargetR, data.TargetS)
	result := ExecuteMove(e.State, action.PlayerID, data.UnitID, target)

	if !result.Ack {
		e.sendNack(action, string(result.Error.Code), result.Error.Message)
		return
	}

	e.sendAck(action)
	e.broadcastDeltas(result)
}

// handleAttack processes an attack action.
func (e *Engine) handleAttack(action PlayerAction) {
	var data ws.AttackData
	if err := json.Unmarshal(action.Data, &data); err != nil {
		e.sendNack(action, string(model.ErrInvalidMessage), "invalid attack data")
		return
	}

	target := hex.NewCoord(data.TargetQ, data.TargetR, data.TargetS)
	result := ExecuteAttack(e.State, e.Roller, action.PlayerID, data.UnitID, target)

	if !result.Ack {
		e.sendNack(action, string(result.Error.Code), result.Error.Message)
		return
	}

	e.sendAck(action)
	e.broadcastDeltas(result)

	if result.GameOver != nil {
		e.endGame(result.GameOver)
	}
}

// handleBuy processes a buy action.
func (e *Engine) handleBuy(action PlayerAction) {
	var data ws.BuyData
	if err := json.Unmarshal(action.Data, &data); err != nil {
		e.sendNack(action, string(model.ErrInvalidMessage), "invalid buy data")
		return
	}

	result := ExecuteBuy(e.State, action.PlayerID, data.UnitType, data.StructureID)

	if !result.Ack {
		e.sendNack(action, string(result.Error.Code), result.Error.Message)
		return
	}

	e.sendAck(action)
	e.broadcastDeltas(result)
}

// handleEndTurn processes an end_turn action.
func (e *Engine) handleEndTurn(action PlayerAction) {
	result := ExecuteEndTurn(e.State, e.Roller, action.PlayerID)

	if !result.Ack {
		e.sendNack(action, string(result.Error.Code), result.Error.Message)
		return
	}

	e.sendAck(action)

	if e.turnTimer != nil {
		e.turnTimer.Stop()
	}

	// Run structure combat phase
	e.runStructureCombat()

	// Broadcast turn start delta
	e.broadcastDeltas(result)

	if result.GameOver != nil {
		e.endGame(result.GameOver)
		return
	}

	// Snapshot state to Redis after turn ends
	e.snapshotState()

	// Start new turn timer
	e.startTurnTimer()
}

// handleEmote forwards an emote to the opponent.
func (e *Engine) handleEmote(action PlayerAction) {
	var data ws.EmoteData
	if err := json.Unmarshal(action.Data, &data); err != nil {
		return
	}

	data.PlayerID = action.PlayerID
	opponentID := e.State.Players[1-e.State.PlayerIndex(action.PlayerID)].ID
	e.Hub.SendMessageTo(opponentID, ws.MsgEmote, data)
}

// handleTurnTimeout auto-ends the turn when the timer expires.
func (e *Engine) handleTurnTimeout() {
	if e.State.Phase != model.PhasePlayerAction {
		return
	}

	e.logger.Info("turn timer expired",
		"turn", e.State.TurnNumber,
		"player_id", e.State.ActivePlayerID(),
	)

	result := ExecuteEndTurn(e.State, e.Roller, e.State.ActivePlayerID())
	if result.Ack {
		e.runStructureCombat()
		e.broadcastDeltas(result)

		if result.GameOver != nil {
			e.endGame(result.GameOver)
			return
		}

		e.snapshotState()
		e.startTurnTimer()
	}
}

// handleDisconnect handles a player disconnection.
func (e *Engine) handleDisconnect(playerID string) {
	e.logger.Info("player disconnected",
		"player_id", playerID,
	)

	idx := e.State.PlayerIndex(playerID)
	if idx >= 0 {
		e.State.Players[idx].IsDisconnected = true
	}

	e.disconnectedID = playerID

	// Notify opponent
	e.Hub.BroadcastMessage(ws.MsgPlayerDisconnected, ws.PlayerDisconnectedData{
		PlayerID: playerID,
	})

	// Start 60-second reconnect timer
	e.reconnectTimer = time.NewTimer(60 * time.Second)
}

// handleReconnect handles a player reconnecting.
func (e *Engine) handleReconnect(event ReconnectEvent) {
	e.logger.Info("player reconnected",
		"player_id", event.PlayerID,
	)

	idx := e.State.PlayerIndex(event.PlayerID)
	if idx >= 0 {
		e.State.Players[idx].IsDisconnected = false
	}

	// Cancel reconnect timer
	if e.reconnectTimer != nil {
		e.reconnectTimer.Stop()
		e.reconnectTimer = nil
	}
	e.disconnectedID = ""

	// Register the new connection
	e.Hub.Register(event.Conn)

	// Send full state to reconnecting player
	e.sendFullState(event.PlayerID)

	// Notify opponent
	e.Hub.BroadcastMessage(ws.MsgPlayerReconnected, ws.PlayerReconnectedData{
		PlayerID: event.PlayerID,
	})
}

// handleReconnectTimeout forfeits the disconnected player.
func (e *Engine) handleReconnectTimeout() {
	if e.disconnectedID == "" {
		return
	}

	e.logger.Info("reconnect timeout expired",
		"player_id", e.disconnectedID,
	)

	gameOver := CheckDisconnectForfeit(e.State, e.disconnectedID)
	e.endGame(gameOver)
}

// runStructureCombat executes the structure auto-attack phase.
func (e *Engine) runStructureCombat() {
	e.State.Phase = model.PhaseStructureCombat

	// All structures fire (if they have a valid target).
	// Neutral structures fire every turn transition.
	// Owned structures fire every turn transition (at enemies).
	for _, structure := range e.State.Structures {
		target := FindStructureTarget(e.State, e.Roller, structure)
		if target == nil {
			continue
		}

		result := ResolveStructureFire(e.State, e.Roller, structure, target)
		e.Hub.BroadcastMessage(ws.MsgStructureFires, result)

		if result.Killed {
			destroyedData := &ws.TroopDestroyedData{
				UnitID: target.ID,
				HexQ:   target.Hex.Q,
				HexR:   target.Hex.R,
				HexS:   target.Hex.S,
				Cause:  "structure_fire",
			}
			e.Hub.BroadcastMessage(ws.MsgTroopDestroyed, destroyedData)
			e.State.RemoveTroop(target.ID)
		}
	}

	e.State.Phase = model.PhasePlayerAction
}

// startTurnTimer starts the turn countdown timer.
func (e *Engine) startTurnTimer() {
	if e.turnTimer != nil {
		e.turnTimer.Stop()
	}
	duration := time.Duration(e.State.TurnTimer) * time.Second
	e.turnTimer = time.NewTimer(duration)
	e.State.TurnStartedAt = time.Now()
}

// endGame handles game over state.
func (e *Engine) endGame(gameOver *ws.GameOverData) {
	e.State.Phase = model.PhaseGameOver

	if e.turnTimer != nil {
		e.turnTimer.Stop()
	}

	e.logger.Info("game over",
		"winner_id", gameOver.WinnerID,
		"reason", gameOver.Reason,
	)

	e.Hub.BroadcastMessage(ws.MsgGameOver, gameOver)
	e.snapshotState()
}

// sendAck sends an ACK to the acting player.
func (e *Engine) sendAck(action PlayerAction) {
	if action.Conn != nil {
		action.Conn.SendAck(action.Seq, action.Type)
	}
}

// sendNack sends a NACK to the acting player.
func (e *Engine) sendNack(action PlayerAction, code, message string) {
	if action.Conn != nil {
		action.Conn.SendNack(action.Seq, action.Type, code, message)
	}
}

// broadcastDeltas sends all delta messages from an ActionResult to both players.
func (e *Engine) broadcastDeltas(result *ActionResult) {
	for i, delta := range result.Deltas {
		e.Hub.BroadcastMessage(result.DeltaTypes[i], delta)
	}
}

// sendFullState sends the complete game state to a specific player.
func (e *Engine) sendFullState(playerID string) {
	e.Hub.SendMessageTo(playerID, ws.MsgGameState, e.State)
}

// broadcastFullState sends the complete game state to all connected players.
func (e *Engine) broadcastFullState() {
	e.Hub.BroadcastMessage(ws.MsgGameState, e.State)
}

// snapshotState persists the game state to Redis.
func (e *Engine) snapshotState() {
	if e.Store == nil {
		return
	}

	data, err := e.State.Serialize()
	if err != nil {
		e.logger.Error("failed to serialize game state",
			"error", err,
		)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	ttl := 24 * time.Hour
	if e.State.Phase == model.PhaseGameOver {
		ttl = 1 * time.Hour
	}

	if err := e.Store.SaveGameState(ctx, e.State.ID, data, ttl); err != nil {
		e.logger.Error("failed to snapshot game state to redis",
			"error", err,
		)
	}
}

// RunTurnStart executes the turn start pipeline and returns the delta data.
func RunTurnStart(gs *GameState, roller *dice.Roller) *ws.TurnStartData {
	gs.Phase = model.PhaseTurnStart
	activePlayerID := gs.ActivePlayerID()

	// 1. Advance turn counter (already done in ExecuteEndTurn for subsequent turns)
	// For the first turn, it's set in startGame

	// 2. Sudden death
	var sdDamages []ws.SuddenDeathDamage
	sdDamages, _ = RunSuddenDeathPhase(gs)

	// 3. Passive healing (+2 HP to troops not in combat last turn)
	var healed []ws.HealedUnit
	for _, troop := range gs.Troops {
		if troop.OwnerID == activePlayerID && troop.IsAlive() && !troop.WasInCombat {
			before := troop.CurrentHP
			amount := troop.Heal(HealingRate())
			if amount > 0 {
				healed = append(healed, ws.HealedUnit{
					UnitID:   troop.ID,
					HPBefore: before,
					HPAfter:  troop.CurrentHP,
				})
			}
		}
	}

	// 4. Structure passive regen
	var structRegens []ws.StructureRegen
	for _, structure := range gs.Structures {
		if structure.IsOwnedBy(activePlayerID) && structure.IsAlive() {
			before := structure.CurrentHP
			amount := structure.Heal(HealingRate())
			if amount > 0 {
				structRegens = append(structRegens, ws.StructureRegen{
					StructureID: structure.ID,
					HPBefore:    before,
					HPAfter:     structure.CurrentHP,
				})
			}
		}
	}

	// 5-6. Calculate and credit income
	passive, structIncome, totalIncome := CalculateIncome(gs, activePlayerID)
	CreditIncome(gs, activePlayerID)
	_ = passive

	// 7-8. Reset troop action flags and mark purchased troops as ready
	for _, troop := range gs.Troops {
		if troop.OwnerID == activePlayerID {
			if !troop.IsReady {
				troop.IsReady = true // troops purchased last turn become ready
			}
			troop.ResetForTurn()
		}
	}

	idx := gs.PlayerIndex(activePlayerID)
	totalCoins := 0
	if idx >= 0 {
		totalCoins = gs.Players[idx].Coins
	}

	return &ws.TurnStartData{
		TurnNumber:         gs.TurnNumber,
		ActivePlayerID:     activePlayerID,
		TimerSeconds:       gs.TurnTimer,
		IncomeGained:       totalIncome,
		StructureIncome:    structIncome,
		TotalCoins:         totalCoins,
		HealedUnits:        healed,
		StructureRegens:    structRegens,
		SuddenDeathDamages: sdDamages,
	}
}
