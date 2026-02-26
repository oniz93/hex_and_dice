# Bot System Implementation Summary

This document summarizes all changes made to implement a PvE bot system in hexbattle.

## Overview

Added a server-side bot that plays like a human player: buying troops, moving, attacking, and ending turns. The bot resides in the `server` component and can be configured with difficulty levels (easy/medium/hard).

---

## Server Changes

### New Files

#### 1. `server/internal/game/bot.go`
Defines the bot interface and action types:
- `BotActionType` - enum for `move`, `attack`, `buy` actions
- `BotAction` - struct containing action details (unit ID, target, troop type, structure ID)
- `BotPlayer` interface:
  - `NextAction(gs *GameState) *BotAction` - called each turn to get the next action; returns `nil` when done
  - `PlayerID() string` - returns the bot's player ID

#### 2. `server/internal/bot/bot.go`
Implements the bot AI:
- `Difficulty` type - `easy`, `medium`, `hard`
- `Bot` struct with:
  - `NextAction(gs *GameState)` - main AI loop implementing turn phases: Buy → Attack → Move → End
- **Buy Phase:**
  - Finds spawnable structures owned by the bot
  - Chooses troop type based on difficulty:
    - Easy: only marines
    - Medium: marines + snipers + hoverbikes  
    - Hard: all troop types including mechs
  - Spawns near the enemy HQ
- **Attack Phase:**
  - Iterates through all bot troops
  - Finds in-range targets (enemy troops first, then structures)
  - Prioritizes low-HP enemies
- **Move Phase:**
  - Uses `game.ReachableHexes` (Dijkstra-based pathfinding)
  - Moves toward nearest enemy troop (within 5 hexes), uncaptured structures, or enemy HQ
- Designed for easy extensibility - adding difficulty tiers just requires adjusting `chooseTroopType`, `chooseObjective`, etc.

---

### Modified Files

#### 3. `server/internal/game/engine.go`
Added bot integration to the game engine:

- **New fields:**
  - `Bot BotPlayer` - holds the bot instance (nil for PvP games)
  - `botTimer *time.Timer` - schedules bot turns

- **New methods:**
  - `botTimerChan() <-chan time.Time` - returns bot timer channel
  - `IsBotGame() bool` - checks if engine has a bot
  - `triggerBotIfNeeded()` - schedules bot turn with 800ms delay if it's the bot's turn
  - `playBotTurn()` - executes bot actions one at a time with delays between them

- **Modified methods:**
  - `Run()` - added bot timer channel handling and calls `triggerBotIfNeeded()` on resume
  - `handleJoinGame()` - for bot games, only 1 human connection needed to start (instead of 2)
  - `startGame()` - logs "bot game starting" when applicable, triggers bot if bot goes first
  - `handleEndTurn()` - triggers bot after human ends turn
  - `handleTurnTimeout()` - triggers bot after timeout ends turn
  - `handleDisconnect()` - ignores disconnects from bot player (no reconnect timer)

#### 4. `server/internal/lobby/room.go`
Added bot game support to room management:

- **Room struct - new fields:**
  - `IsBotGame bool` - indicates if this is a PvE game
  - `BotDifficulty string` - difficulty level ("easy", "medium", "hard")

- **New method:**
  - `Manager.CreateBotRoom()` - creates a room pre-filled with a bot as the guest player, immediately ready to start

#### 5. `server/internal/api/handlers_rooms.go`
Added bot game creation endpoint:

- **New request type:** `CreateBotGameRequest`
  - `MapSize string` - optional map size
  - `TurnTimer int` - optional turn timer
  - `Difficulty string` - "easy", "medium", "hard" (default: easy)

- **New handler:** `HandleCreateBotGame` - `POST /api/v1/rooms/bot`

- **New response type:** `BotGameResult` (in api_service.dart)
  - `roomId`, `roomCode`, `settings`, `botPlayerId`, `botDifficulty`

#### 6. `server/internal/api/router.go`
Added bot route:
- `POST /api/v1/rooms/bot` → `roomsHandler.HandleCreateBotGame`

#### 7. `server/cmd/server/main.go`
Wired up bot in game creation flow:

- Added import: `"github.com/teomiscia/hexbattle/internal/bot"`
- Modified engine creation (around line 175):
  - After creating `game.Engine`, checks if `room.IsBotGame`
  - If true, creates a `bot.Bot` with the appropriate difficulty
  - Assigns to `engine.Bot`

---

## Client Changes

### Modified Files

#### 8. `client/lib/services/api_service.dart`
Added bot game support:

- **New response class:** `BotGameResult`
  - `roomId`, `roomCode`, `settings`, `botPlayerId`, `botDifficulty`

- **New method:** `createBotGame()`
  - Calls `POST /api/v1/rooms/bot`
  - Parameters: `mapSize`, `turnTimer`, `difficulty`
  - Returns `BotGameResult`

#### 9. `client/lib/screens/play_screen.dart`
Added bot game buttons to the play screen:

- **Play vs Bot (Easy)** - Green button, calls `createBotGame(difficulty: 'easy')`
- **Play vs Bot (Hard)** - Red button, calls `createBotGame(difficulty: 'hard')`

Both buttons navigate to `/game/{roomId}` after successfully creating the bot game.

---

## API Flow

1. Client calls `POST /api/v1/rooms/bot` with optional `map_size`, `turn_timer`, `difficulty`
2. Server creates a room with bot as guest, state = `ready`
3. Client connects via WebSocket and sends `join_game`
4. Engine starts with only 1 human connection (bot is server-side)
5. On bot's turn:
   - `triggerBotIfNeeded()` schedules 800ms delay
   - `playBotTurn()` calls `bot.NextAction()` in a loop
   - Each action executed via `ExecuteBuy/Move/Attack`, deltas broadcast to human
   - When `NextAction` returns nil, bot ends turn

---

## Testing

- Server compiles: `GO111MODULE=on go build ./cmd/server/...` ✓
- All tests pass: `GO111MODULE=on go test ./...` ✓
- Client analyzes: `flutter analyze` (warnings are pre-existing, no errors) ✓

---

## Future Improvements

- Add medium difficulty option to client UI
- Add map size / turn timer selection for bot games
- More sophisticated AI (e.g., evaluate board state, prioritize objectives)
- Bot persistence across reconnects (currently bot state is in-memory)
- Per-difficulty personality traits (aggressive vs defensive)
