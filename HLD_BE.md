# Hex & Dice — Backend High-Level Design

## 1. Overview

This document describes the high-level architecture of the Hex & Dice backend server. The server is the **authoritative source of truth** for all game state. It handles player sessions, matchmaking, room management, map generation, game logic execution, combat resolution, and real-time state synchronization over WebSockets.

**Language:** Go
**Architecture:** Modular monolith (single binary, internal package boundaries)
**Concurrency model:** One goroutine per active game + read/write goroutines per WebSocket connection, communicating via channels

---

## 2. Project Structure

```
server/
├── cmd/
│   └── server/
│       └── main.go                  # Entrypoint: config loading, dependency wiring, server startup
├── internal/
│   ├── config/
│   │   └── config.go                # Environment variable parsing, server configuration
│   ├── game/
│   │   ├── engine.go                # Game loop goroutine, FSM transitions, turn pipeline
│   │   ├── state.go                 # GameState struct, state mutation methods
│   │   ├── actions.go               # Action types: move, attack, buy, end_turn
│   │   ├── validate.go              # Per-action validation functions
│   │   ├── combat.go                # Hit resolution, damage resolution, counterattacks
│   │   ├── pathfinding.go           # BFS reachable hexes, movement cost calculation
│   │   ├── economy.go               # Income calculation, purchase validation
│   │   ├── wincondition.go          # Win condition checks (HQ destruction, structure dominance)
│   │   ├── suddendeath.go           # Shrinking zone logic, escalating damage, HQ relocation
│   │   └── constants.go             # Re-exports of balance data loaded from YAML
│   ├── mapgen/
│   │   ├── generator.go             # Procedural map generation entry point
│   │   ├── noise.go                 # OpenSimplex noise wrapper
│   │   ├── symmetry.go              # Rotational symmetry (180°) application
│   │   ├── placement.go             # Structure placement (even spread, contested center)
│   │   └── validation.go            # Connectivity validation, constraint checking, retry logic
│   ├── hex/
│   │   ├── coords.go                # Cube coordinate type (q, r, s), arithmetic, distance
│   │   ├── grid.go                  # HexGrid type, neighbor lookup, ring/spiral iterators
│   │   └── directions.go            # 6 hex directions, coordinate offsets
│   ├── dice/
│   │   ├── roller.go                # Seeded RNG per game, D20/D6/D8/D4 roll methods
│   │   └── types.go                 # DiceResult, CriticalHit, Fumble types
│   ├── lobby/
│   │   ├── manager.go               # Room creation, join, lookup, TTL expiry
│   │   ├── room.go                  # Room struct, lifecycle states
│   │   └── matchmaking.go           # FIFO queue, match pairing logic
│   ├── player/
│   │   ├── session.go               # Player session: token, ID, nickname, connection state
│   │   └── registry.go              # Active player registry, token-to-session lookup
│   ├── ws/
│   │   ├── handler.go               # WebSocket upgrade handler, auth via query param
│   │   ├── connection.go            # Connection wrapper: read/write goroutines, channels
│   │   ├── messages.go              # Message envelope types, serialization/deserialization
│   │   └── hub.go                   # Per-game message hub: broadcast, direct send
│   ├── api/
│   │   ├── router.go                # REST route definitions (/api/v1/...)
│   │   ├── middleware.go            # Logging, CORS, rate limiting, auth middleware
│   │   ├── handlers_guest.go        # POST /api/v1/guest — guest registration
│   │   ├── handlers_rooms.go        # Room CRUD: create, join, get
│   │   ├── handlers_matchmaking.go  # POST /api/v1/matchmaking/join, DELETE .../leave
│   │   └── handlers_health.go       # GET /health — server health check
│   ├── store/
│   │   ├── redis.go                 # Redis client wrapper, game state snapshot/restore
│   │   └── interface.go             # Store interface (for testing with mocks)
│   └── model/
│       ├── troop.go                 # Troop struct, TroopType enum
│       ├── structure.go             # Structure struct, StructureType enum
│       ├── terrain.go               # Terrain enum, movement cost/combat modifier lookups
│       ├── player.go                # PlayerState: coins, owned structures, troop list
│       └── enums.go                 # Shared enums: GamePhase, TurnMode, MapSize, etc.
├── data/
│   └── balance.yaml                 # Game balance constants (troop stats, economy, terrain)
├── Dockerfile                       # Multi-stage build: Go build → minimal runtime image
├── docker-compose.yml               # Server + Redis + PostgreSQL (future)
├── go.mod
└── go.sum
```

### Key packages and their responsibilities

| Package | Responsibility |
|---|---|
| `cmd/server` | Binary entrypoint. Loads config, initializes dependencies, starts HTTP server |
| `internal/game` | Core game engine: state machine, action processing, combat, economy, win conditions |
| `internal/mapgen` | Procedural hex map generation with noise, symmetry, structure placement, validation |
| `internal/hex` | Hex grid math: cube coordinates, distance, neighbors, rings |
| `internal/dice` | Dice rolling with seeded per-game RNG |
| `internal/lobby` | Room management and matchmaking queue |
| `internal/player` | Player session and identity management |
| `internal/ws` | WebSocket connection lifecycle, message framing, per-game hub |
| `internal/api` | REST endpoint handlers and middleware |
| `internal/store` | Redis persistence layer |
| `internal/model` | Shared data structures and enums used across packages |
| `data/` | YAML files with game balance constants |

---

## 3. External Dependencies

| Dependency | Purpose | Import Path |
|---|---|---|
| **nhooyr/websocket** | WebSocket server implementation | `nhooyr.io/websocket` |
| **go-redis** | Redis client | `github.com/redis/go-redis/v9` |
| **opensimplex-go** | OpenSimplex noise for map generation | `github.com/ojrac/opensimplex-go` |
| **testify** | Test assertions and mocks | `github.com/stretchr/testify` |
| **gopkg.in/yaml.v3** | YAML parsing for balance data | `gopkg.in/yaml.v3` |

No HTTP framework — uses Go stdlib `net/http` with `http.ServeMux` (Go 1.22+).

---

## 4. Server Lifecycle

### 4.1 Startup Sequence

```
1. Parse environment variables → Config struct
2. Load balance.yaml → game constants
3. Connect to Redis (retry with backoff)
4. Initialize player registry (in-memory)
5. Initialize lobby manager (in-memory rooms + matchmaking queue)
6. Restore any active games from Redis snapshots (server restart recovery)
7. Register REST routes on ServeMux
8. Register WebSocket upgrade handler
9. Start HTTP server (listen on configured port)
10. Log "server started" with port, config summary
```

### 4.2 Graceful Shutdown

Triggered by `SIGTERM` or `SIGINT`:

```
1. Stop accepting new HTTP connections (server.Shutdown with context)
2. Stop matchmaking queue (reject new entries)
3. Signal all active game goroutines to drain (via context cancellation)
4. Wait for active games to reach a save point (max 30s timeout)
5. Snapshot all active game states to Redis
6. Close Redis connection
7. Log "server stopped", exit 0
```

If the 30-second drain timeout is exceeded, force-snapshot remaining games and exit.

---

## 5. Player Session Management

### 5.1 Guest Registration

- **Endpoint:** `POST /api/v1/guest`
- **Request body:** `{"nickname": "PlayerOne"}`
- **Server generates:**
  - `player_id`: UUIDv4
  - `token`: 32-byte hex string via `crypto/rand`
- **Response:** `{"player_id": "uuid", "token": "hex_string", "nickname": "PlayerOne"}`
- **Storage:** In-memory `map[token] → PlayerSession` in the player registry
- **Lifetime:** Session-scoped. Token is valid as long as the player holds it. No expiration. If the server restarts, all sessions are lost (players re-register — acceptable for guest-only MVP).

### 5.2 Token Usage

- **REST requests:** `Authorization: Bearer <token>` header
- **WebSocket upgrade:** Query parameter `ws://host/ws?token=<token>`
- Server validates the token against the player registry before processing any request

### 5.3 Nickname Validation

- 3-16 characters
- Alphanumeric + underscores only
- No uniqueness constraint (guests can share nicknames)
- Trimmed and sanitized server-side

---

## 6. Lobby & Matchmaking

### 6.1 Room Management

All room state is **in-memory only** (not persisted to Redis). Rooms are ephemeral pre-game constructs.

#### Room Struct

```
Room {
    Code         string       // 6-digit numeric (e.g., "482917")
    HostPlayerID uuid
    GuestPlayerID uuid | nil
    Settings     RoomSettings {
        MapSize   enum(Small, Medium, Large)
        TurnTimer enum(60, 90, 120)  // seconds
        TurnMode  enum(Alternating)  // Simultaneous in Phase 2
    }
    State        enum(WaitingForOpponent, Ready, GameInProgress, GameOver)
    CreatedAt    time.Time
}
```

#### Room Code Generation

- 6-digit numeric string (`000000` – `999999`)
- Generated via `crypto/rand`
- Checked for uniqueness against active rooms (collision retry)
- 1 million combinations — sufficient for MVP concurrent room counts

#### Room Lifecycle

```
Created → WaitingForOpponent → Ready → GameInProgress → GameOver → Cleaned up
                                                                      ↑
                                                              60s delay for
                                                              "Play Again"
```

#### Room TTL

- If no opponent joins within **5 minutes** of creation, the room auto-expires
- A background goroutine sweeps expired rooms periodically (every 30s)
- After `GameOver`, room lingers for 60 seconds to allow "Play Again" (which creates a new game in the same room). If no action, the room is destroyed.

#### REST Endpoints

| Method | Path | Description |
|---|---|---|
| `POST` | `/api/v1/rooms` | Create a room. Body: `{settings}`. Returns: `{room_code, room_id}` |
| `POST` | `/api/v1/rooms/join` | Join a room. Body: `{code}`. Returns: `{room_id, settings, host_nickname}` |
| `GET` | `/api/v1/rooms/{code}` | Get room status (for polling before WS connect) |

### 6.2 Matchmaking Queue

- **Endpoint:** `POST /api/v1/matchmaking/join` — add player to queue
- **Endpoint:** `DELETE /api/v1/matchmaking/leave` — remove player from queue
- **Algorithm:** Simple FIFO. When a second player joins, they are immediately matched.
- **Quick Match defaults:** Medium map, 90s timer, Alternating turn mode
- **On match:**
  1. Create a room with Quick Match defaults
  2. Assign both players to the room
  3. Room state = `Ready`
  4. Return the room ID to both players (via response to the joining player, and via a push notification to the waiting player over WebSocket if they have one open, or via polling)

#### Queue Wait Notification

When a player joins the matchmaking queue:
- They can either poll `GET /api/v1/matchmaking/status` or
- Open a WebSocket connection first. The server sends a `match_found` message with the room ID when matched.

---

## 7. WebSocket Protocol

### 7.1 Connection Lifecycle

```
1. Client sends HTTP upgrade request: GET /ws?token=<token>
2. Server validates token against player registry
3. If invalid → 401 Unauthorized, reject upgrade
4. If valid → complete WebSocket upgrade
5. Server spawns read goroutine and write goroutine for this connection
6. Client sends initial message: {"type": "join_game", "data": {"room_id": "..."}}
7. Server associates connection with the game instance
8. Server sends full game state snapshot to the joining player
9. Normal game communication begins
```

### 7.2 Connection Architecture

Each WebSocket connection is managed by a `Connection` struct:

```
Connection {
    PlayerID    uuid
    Conn        *websocket.Conn
    SendChan    chan []byte       // write goroutine reads from this
    GameID      string | nil      // set after joining a game
}
```

- **Read goroutine:** Reads messages from the WebSocket, deserializes, validates structure, routes to the game goroutine via channel
- **Write goroutine:** Reads from `SendChan`, writes to WebSocket. If write fails, signals disconnect
- **Backpressure:** `SendChan` has a buffer (e.g., 64 messages). If full, the connection is considered stuck and is closed

### 7.3 Message Envelope

All messages (client→server and server→client) use the same envelope:

```json
{
    "type": "action_type",
    "seq": 123,
    "data": { ... }
}
```

- `type` (string, required): Message type identifier
- `seq` (integer, required for client→server): Monotonically increasing sequence number per connection. Server echoes it in ACK/NACK responses for request-response correlation
- `data` (object, required): Type-specific payload

### 7.4 Client → Server Messages

| Type | Data | Description |
|---|---|---|
| `join_game` | `{room_id}` | Associate this connection with a game/room |
| `reconnect` | `{game_id, player_token}` | Reconnect to an active game after disconnect |
| `move` | `{unit_id, target_q, target_r, target_s}` | Move a troop to a hex |
| `attack` | `{unit_id, target_q, target_r, target_s}` | Attack a target at hex |
| `buy` | `{unit_type, structure_id}` | Purchase a troop at a spawn structure |
| `end_turn` | `{}` | End the current turn |
| `emote` | `{emote_id}` | Send a predefined emote |
| `pong` | `{}` | Response to server ping |

### 7.5 Server → Client Messages

| Type | Data | Description |
|---|---|---|
| `game_state` | `{full_state}` | Full game state snapshot (on join/reconnect) |
| `ack` | `{seq, action_type}` | Action accepted |
| `nack` | `{seq, action_type, error}` | Action rejected with error |
| `troop_moved` | `{unit_id, from, to, remaining_mobility}` | Troop movement delta |
| `combat_result` | `{attacker_id, defender_id, hit_roll, natural_roll, hit, damage_roll, damage, defender_hp, killed, counter_hit_roll, counter_natural_roll, counter_hit, counter_damage, attacker_hp, attacker_killed, crit, fumble}` | Full combat resolution delta |
| `troop_purchased` | `{unit_id, unit_type, hex, owner, coins_remaining}` | New troop purchased |
| `troop_destroyed` | `{unit_id, hex, cause}` | Troop died (combat, sudden death, structure fire) |
| `structure_attacked` | `{structure_id, attacker_id, hit_roll, damage, structure_hp, captured, new_owner}` | Structure took damage or was captured |
| `structure_fires` | `{structure_id, target_id, hit_roll, damage, target_hp, killed}` | Structure attacked a troop |
| `turn_start` | `{turn_number, active_player_id, timer_seconds, income_gained, structure_income, total_coins, healed_units[], structure_regen[], sudden_death_damage[]}` | New turn begins with all passive effects |
| `game_over` | `{winner_id, reason, stats}` | Game ended |
| `player_disconnected` | `{player_id}` | Opponent disconnected, reconnect timer started |
| `player_reconnected` | `{player_id}` | Opponent reconnected |
| `emote` | `{player_id, emote_id}` | Emote from opponent |
| `ping` | `{}` | Server heartbeat (expect pong) |
| `match_found` | `{room_id}` | Matchmaking found an opponent |
| `error` | `{code, message}` | General error (not tied to a specific action) |

### 7.6 Server Response Pattern

When a client sends an action:

1. Server validates the message structure (strict validation)
2. Server validates game logic (is this action legal in the current state?)
3. If invalid → send `nack` to the acting player with the `seq` number and error
4. If valid →
   a. Execute the action, mutate game state
   b. Send `ack` to the acting player with the `seq` number
   c. Broadcast the resulting delta(s) to **both** players

Deltas are **per-action granularity**: each action produces one or more delta messages (e.g., an attack may produce `combat_result` + `troop_destroyed` if the target dies, or `combat_result` + `combat_result` if a counterattack occurs).

### 7.7 Error Response Format

```json
{
    "error": {
        "code": "INVALID_MOVE",
        "message": "Hex (2, 3, -5) is occupied by an enemy unit"
    }
}
```

Error codes are machine-readable uppercase strings. Error messages are human-readable descriptions.

#### Error Code Catalog

| Code | Description |
|---|---|
| `NOT_YOUR_TURN` | Player attempted an action when it's not their turn |
| `INVALID_MOVE` | Move target is unreachable, occupied, or out of bounds |
| `INVALID_ATTACK` | Target out of range, no target at hex, attacking own unit |
| `INSUFFICIENT_FUNDS` | Not enough coins to purchase the unit |
| `SPAWN_OCCUPIED` | Spawn hex is already occupied by another troop |
| `SPAWN_NOT_OWNED` | Structure is not owned by the player (or neutral) |
| `UNIT_ALREADY_ACTED` | Unit has already moved/attacked this turn |
| `UNIT_NOT_READY` | Unit was purchased this turn and cannot act |
| `UNIT_NOT_FOUND` | Referenced unit_id does not exist |
| `GAME_NOT_FOUND` | Referenced game_id does not exist |
| `ROOM_NOT_FOUND` | Referenced room code does not exist |
| `ROOM_FULL` | Room already has two players |
| `ROOM_EXPIRED` | Room TTL expired |
| `INVALID_MESSAGE` | Malformed message structure |
| `RATE_LIMITED` | Too many actions in a short period |

### 7.8 Heartbeat / Keep-Alive

- Server sends `{"type": "ping"}` every **15 seconds** to each connected client
- Client must respond with `{"type": "pong"}` within **10 seconds**
- If no pong is received, the server considers the connection lost
- On connection loss:
  1. Mark the player as disconnected
  2. Notify the opponent via `player_disconnected` message
  3. Start the **60-second reconnect timer**
  4. If the player reconnects within 60s → send `player_reconnected` to opponent + full game state snapshot to the reconnecting player
  5. If 60s expires → the disconnected player forfeits, game ends

### 7.9 Reconnect Protocol

1. Client opens a new WebSocket connection: `GET /ws?token=<token>`
2. Client sends: `{"type": "reconnect", "seq": 1, "data": {"game_id": "...", "player_token": "..."}}`
3. Server validates token + game_id + confirms player was in this game
4. Server cancels the disconnect forfeit timer
5. Server sends full `game_state` snapshot (current complete state)
6. Server sends `player_reconnected` to the opponent
7. Normal gameplay resumes

---

## 8. Game Engine

### 8.1 Game State Machine

Each game instance runs as a single goroutine with an explicit finite state machine:

```
                    ┌──────────────────┐
                    │ WaitingForPlayers │
                    └────────┬─────────┘
                             │ both players connected
                             ▼
                    ┌────────────────┐
                    │ GeneratingMap  │
                    └────────┬───────┘
                             │ map generated + validated
                             ▼
                    ┌─────────────────┐
                    │ GameStarted     │
                    └────────┬────────┘
                             │ initial state broadcast
                             ▼
                ┌───────────────────────┐
          ┌────►│ TurnStart             │◄──────────────────┐
          │     │ (passive effects)     │                   │
          │     └────────┬──────────────┘                   │
          │              │ effects applied                  │
          │              ▼                                  │
          │     ┌────────────────────┐                      │
          │     │ StructureCombat    │                      │
          │     │ (structures fire)  │                      │
          │     └────────┬───────────┘                      │
          │              │ structure attacks resolved       │
          │              ▼                                  │
          │     ┌────────────────────┐                      │
          │     │ PlayerAction       │──── win condition ──►│ GameOver │
          │     │ (awaiting input)   │     met              └──────────┘
          │     └────────┬───────────┘                      │
          │              │ end_turn / timer expired         │
          │              ▼                                  │
          │     ┌────────────────────┐                      │
          │     │ TurnTransition     │──── win condition ──►│
          │     │ (switch player)    │     met              │
          │     └────────┬───────────┘                      │
          │              │                                  │
          └──────────────┘                                  │
                                                            │
          (sudden death zone shrinks during TurnStart       │
           when turn counter exceeds threshold)             │
```

### 8.2 Game Goroutine Event Loop

Each game runs in its own goroutine with a `select` loop:

```
game goroutine:
    for {
        select {
        case action := <-actionChan:       // player action from WS read goroutine
            validate(action)
            execute(action)
            broadcastDeltas()

        case <-turnTimer.C:                // turn timer expired
            autoEndTurn()
            transitionTurn()
            broadcastDeltas()

        case disconnect := <-disconnectChan: // player disconnected
            startReconnectTimer(disconnect.playerID)

        case reconnect := <-reconnectChan:   // player reconnected
            cancelReconnectTimer()
            sendFullState(reconnect.conn)

        case <-reconnectTimeout.C:          // reconnect window expired
            forfeit(disconnectedPlayer)
            endGame()

        case <-ctx.Done():                  // server shutdown
            snapshotToRedis()
            return
        }
    }
```

### 8.3 Turn Start Pipeline (Synchronous)

When a new turn begins, the following effects are applied in order:

```
1. Advance turn counter
2. Check sudden death activation/progression
   a. If sudden death active: shrink safe zone radius by 1
   b. Relocate HQs inside safe zone if necessary
   c. Apply escalating damage to troops outside safe zone
   d. Check for troop deaths from storm damage
3. Apply passive healing (+2 HP to troops that were not in combat last turn, up to max HP)
4. Apply structure passive regen (+2 HP to owned structures, up to max HP)
5. Calculate income: passive (100) + structure bonuses (50 × owned income-generating structures)
6. Credit income to the active player's coin balance
7. Mark all of active player's troops as "ready" (can act this turn)
   - Except troops purchased last turn that are now becoming active
8. Reset per-turn action flags on all troops (has_moved, has_attacked, was_in_combat)
9. Check win conditions (structure dominance counter)
10. If win condition met → transition to GameOver
11. Broadcast turn_start delta with all passive effect details
12. Start turn timer
```

### 8.4 Structure Combat Phase

After passive effects are applied, before the active player can act:

```
1. Gather all structures owned by the active player
2. For each structure, find enemy troops within range
3. Each structure attacks one enemy troop in range (closest first, random tiebreak)
4. Resolve each attack using standard D20 + ATK vs DEF formula
5. Apply damage if hit
6. Broadcast structure_fires deltas

For neutral structures (processed at the start of the round, before Player 1's first action):
1. Find all troops from any player within range
2. Use reduced neutral stats (ATK -2, damage dice one step lower)
3. Resolve attacks, broadcast deltas
```

### 8.5 Action Processing

Each action received from a player goes through:

```
1. Check game phase is PlayerAction
2. Check it is this player's turn
3. Dispatch to action-specific validator
4. If validation fails → NACK with error code
5. If validation passes → ACK + execute + broadcast delta(s)
```

#### Move Action

- **Validation:**
  1. Unit exists and belongs to the acting player
  2. Unit has not already moved this turn
  3. Unit is "ready" (not purchased this turn)
  4. Target hex is within the map bounds
  5. Target hex is passable terrain (not water/mountain)
  6. Target hex is not occupied by an enemy troop
  7. A valid path exists from unit's current hex to target hex within the unit's mobility budget (BFS with terrain costs)
  8. Path does not pass through enemy-occupied hexes (friendly hexes are passable)
- **Execution:**
  1. Update unit's hex position
  2. Deduct movement cost from remaining mobility for this turn
  3. Set `has_moved = true`
- **Delta:** `troop_moved {unit_id, from, to, remaining_mobility}`

#### Attack Action

- **Validation:**
  1. Unit exists and belongs to the acting player
  2. Unit has not already attacked this turn
  3. Unit is "ready"
  4. If unit has not moved yet, it may still move first (but this is just an attack order — valid)
  5. Target hex is within the unit's attack range (hex distance from unit's current position)
  6. Target hex contains an enemy troop OR an enemy/neutral structure
  7. Target is a valid combat target (not own troop, not own structure)
- **Execution:**
  1. Resolve hit: roll D20 + attacker ATK vs defender DEF (with terrain modifier on defender's hex)
  2. Check for critical hit (natural 20) or fumble (natural 1)
  3. If hit: roll damage dice, apply terrain ATK modifier, apply double damage on crit, apply 2x structure damage for Mech
  4. Subtract damage from defender HP
  5. If defender HP ≤ 0: destroy defender (remove from game state); if structure at 0 HP: transfer ownership to attacker at 1 HP
  6. If fumble: defender gets free counterattack (half damage)
  7. If melee vs melee and not a fumble: defender gets normal counterattack (half damage)
  8. Counterattack: roll D20 + defender ATK vs attacker DEF (with terrain modifier). On hit: roll half damage dice
  9. If attacker HP ≤ 0 from counterattack: destroy attacker
  10. Mark both units as `was_in_combat = true` (prevents healing next turn)
  11. Set attacker `has_attacked = true`
- **Deltas:** `combat_result` (primary), optionally `troop_destroyed`, `structure_attacked`

#### Buy Action

- **Validation:**
  1. Player has enough coins for the unit type
  2. Structure exists, belongs to the player, and is a valid spawn point (HQ, Outpost, or Command Center)
  3. Structure's hex is not occupied by another troop
- **Execution:**
  1. Deduct cost from player coins
  2. Create new troop on the structure's hex
  3. Mark troop as `not_ready` (cannot act until next turn)
  4. Assign a unique `unit_id` (UUIDv4)
- **Delta:** `troop_purchased {unit_id, unit_type, hex, owner, coins_remaining}`

#### End Turn Action

- **Validation:**
  1. It is this player's turn
  2. Game phase is PlayerAction
- **Execution:**
  1. Transition to TurnTransition
  2. Check win conditions
  3. Switch active player
  4. Run TurnStart pipeline for the new active player

### 8.6 Unit Action Tracking

Each troop maintains per-turn flags:

```
Troop {
    ID              uuid
    Type            TroopType
    OwnerID         uuid
    Hex             CubeCoord
    CurrentHP       int
    MaxHP           int
    IsReady         bool    // false on the turn purchased, true thereafter
    HasMoved        bool    // reset each turn
    HasAttacked     bool    // reset each turn
    WasInCombat     bool    // set if attacked or was attacked; reset each turn
    RemainingMobility int  // set to max at turn start, decremented on move
}
```

---

## 9. Combat System Implementation

### 9.1 Hit Resolution

```
attackRoll = D20()  // 1-20 uniform random
naturalRoll = attackRoll
attackRoll += attacker.ATK
attackRoll += terrain.ATKModifier(attacker.Hex)  // hills: +1

targetDEF = defender.DEF + terrain.DEFModifier(defender.Hex)  // forest: +2, hills: +1

if naturalRoll == 20:
    hit = true    // auto-hit regardless of DEF
    isCrit = true
elif naturalRoll == 1:
    hit = false   // auto-miss regardless of modifiers
    isFumble = true
else:
    hit = (attackRoll >= targetDEF)
```

### 9.2 Damage Resolution

```
if hit:
    damageRoll = rollDamageDice(attacker.DamageDice)  // e.g., 2D6+2
    if isCrit:
        damageRoll = rollDamageDice(attacker.DamageDice) * 2  // double the dice result
    if attacker.Type == Mech AND defender is Structure:
        damageRoll *= 2  // anti-structure bonus
    defender.CurrentHP -= damageRoll
```

### 9.3 Counterattack Resolution

Triggered when:
- Melee attacker (range=1) attacks a melee defender (range=1), OR
- Attacker rolled a fumble (natural 1) — defender gets a free counter regardless of range

```
counterRoll = D20()
counterNatural = counterRoll
counterRoll += defender.ATK
counterDEF = attacker.DEF + terrain.DEFModifier(attacker.Hex)

if counterNatural == 20:
    counterHit = true
elif counterNatural == 1:
    counterHit = false
else:
    counterHit = (counterRoll >= counterDEF)

if counterHit:
    counterDamage = rollDamageDice(defender.DamageDice) / 2  // half damage, rounded down
    attacker.CurrentHP -= counterDamage
```

### 9.4 Damage Dice Notation

Parsed from the balance data YAML:

| Notation | Meaning |
|---|---|
| `1D6+1` | Roll 1 six-sided die, add 1 |
| `1D8` | Roll 1 eight-sided die |
| `2D6+2` | Roll 2 six-sided dice, sum them, add 2 |
| `1D4` | Roll 1 four-sided die |

Half damage (counterattack): divide the **total roll result** by 2, round down. Minimum 1 damage on hit.

---

## 10. Map Generation Engine

### 10.1 Algorithm Overview

```
1. Create hexagonal grid of the specified radius
2. Generate OpenSimplex noise values for each hex in one half of the map
3. Map noise values to terrain types using thresholds
4. Apply 180° rotational symmetry (copy half to the other half)
5. Place player HQs at opposite poles of the hex grid
6. Place neutral structures (even spread with contested center)
7. Validate constraints:
   a. Path exists between HQs through passable terrain
   b. All structures are on passable terrain
   c. No isolated passable regions (all passable hexes connected)
   d. Minimum distance between structures
8. If validation fails → regenerate (up to 10 retries)
9. Return completed map
```

### 10.2 Noise-to-Terrain Mapping

Noise values are normalized to [0, 1]:

| Noise Range | Terrain |
|---|---|
| 0.00 – 0.15 | Water |
| 0.15 – 0.55 | Plains |
| 0.55 – 0.75 | Forest |
| 0.75 – 0.88 | Hills |
| 0.88 – 1.00 | Mountains |

Thresholds may be tuned via the balance YAML. Multiple noise octaves at different scales create natural-looking biome clusters.

### 10.3 Rotational Symmetry

- Generate terrain for one half of the hex grid (one "hemisphere" split along the center)
- For each hex `(q, r, s)` in the generated half, set the hex at `(-q, -r, -s)` to the same terrain type
- Center hex (0, 0, 0) is always Plains
- This guarantees both players face identical terrain relative to their starting position

### 10.4 Structure Placement

| Map Size | Radius | Neutral Structures | HQs |
|---|---|---|---|
| Small | 7 | 3 | 2 |
| Medium | 10 | 5 | 2 |
| Large | 13 | 7 | 2 |

**HQ placement:** At opposite poles along one axis of the hexagonal grid (maximizing distance).

**Neutral structure placement algorithm:**
1. Define candidate hexes: passable terrain, minimum 3 hexes from any HQ, minimum 2 hexes from any other structure
2. Place ~40% of neutral structures in the center ring (within radius/3 of center) — contested zone
3. Place ~30% in the mid ring — distributed evenly
4. Place ~30% near each player's side (but not too close to HQ) — early expansion targets
5. Structures respect rotational symmetry (placed in symmetric pairs, with one in center if odd count)

**Structure type distribution:**
- Outposts: ~70% of neutral structures
- Command Centers: ~30% of neutral structures
- Types assigned to maintain symmetric pairs

### 10.5 Validation Constraints

All constraints must pass for a map to be accepted:

1. **Connectivity:** Flood-fill from HQ1 reaches HQ2 through passable terrain
2. **Structure accessibility:** Every structure is reachable from both HQs
3. **No isolated regions:** Single connected component of passable terrain
4. **Minimum passable ratio:** At least 60% of hexes are passable
5. **HQ safety:** No impassable terrain within 2 hexes of either HQ

If any constraint fails, the entire map is regenerated with a different noise seed. Max 10 retries before returning an error (should be extremely rare with reasonable noise thresholds).

---

## 11. Pathfinding

### 11.1 Reachable Hexes (BFS with Cost)

Used when a player selects a troop to see where it can move.

```
function reachableHexes(start, mobility, grid, friendlyUnits, enemyUnits):
    frontier = priority queue [(start, 0)]
    reached = {start: 0}

    while frontier not empty:
        current, costSoFar = frontier.pop()
        for each neighbor of current:
            if neighbor is impassable terrain: skip
            if neighbor is occupied by enemy: skip
            moveCost = terrain.movementCost(neighbor)
            totalCost = costSoFar + moveCost
            if totalCost > mobility: skip
            if neighbor not in reached OR totalCost < reached[neighbor]:
                reached[neighbor] = totalCost
                frontier.push(neighbor, totalCost)

    // Remove hexes occupied by friendly units (can pass through but not stop on)
    for hex in reached:
        if hex is occupied by friendly unit AND hex != start:
            remove from result (can't stop here) but keep for pathfinding through

    return reached (set of hexes the unit can move to)
```

### 11.2 Path Reconstruction

When the client sends a move action, the server validates that a valid path exists by running BFS from the unit's position and checking if the target hex is in the reachable set. The server does **not** require the client to send the full path — only the destination. The path taken is irrelevant since only the destination matters (no "opportunity attacks" or path-dependent effects).

### 11.3 Attack Range Check

Simple hex distance check:

```
function canAttack(attacker, targetHex):
    distance = hexDistance(attacker.Hex, targetHex)
    return distance <= attacker.Range AND distance >= 1
```

Hex distance in cube coordinates: `max(|q1-q2|, |r1-r2|, |s1-s2|)`

---

## 12. Economy Engine

### 12.1 Income Calculation (Per Turn Start)

```
income = PASSIVE_INCOME  // 100

for each structure owned by active player:
    if structure.Type == Outpost OR structure.Type == CommandCenter:
        income += STRUCTURE_INCOME  // 50

player.Coins += income
```

### 12.2 Purchase Validation

```
function validatePurchase(player, unitType, structureID):
    cost = TROOP_COSTS[unitType]
    if player.Coins < cost:
        return error(INSUFFICIENT_FUNDS)

    structure = getStructure(structureID)
    if structure.OwnerID != player.ID:
        return error(SPAWN_NOT_OWNED)
    if hexOccupied(structure.Hex):
        return error(SPAWN_OCCUPIED)

    return ok
```

---

## 13. Win Condition Checks

Checked at the end of every action and during turn transitions:

### 13.1 HQ Destruction

```
if opponent.HQ.CurrentHP <= 0:
    winner = activePlayer
    reason = "HQ_DESTROYED"
```

### 13.2 Structure Dominance

```
totalStructures = count of all structures on map (including HQs)
playerStructures = count of structures owned by player

if playerStructures > totalStructures / 2:
    player.DominanceTurnCounter++
else:
    player.DominanceTurnCounter = 0

if player.DominanceTurnCounter >= 3:
    winner = player
    reason = "STRUCTURE_DOMINANCE"
```

The dominance counter increments once per **full round** (after both players have taken a turn).

### 13.3 Sudden Death Tiebreak

If the shrinking zone reduces to a single hex and no winner yet:

```
1. Player with more structures wins
2. If tied: player with more total remaining troop HP wins
3. If still tied: draw
```

---

## 14. Sudden Death Implementation

### 14.1 Activation

```
thresholds = {Small: 20, Medium: 30, Large: 40}

if turnCounter > thresholds[mapSize]:
    suddenDeathActive = true
    suddenDeathTurn = turnCounter - thresholds[mapSize]  // 1, 2, 3, ...
```

### 14.2 Zone Shrink

```
safeZoneRadius = mapRadius - suddenDeathTurn

// Clamp to minimum 1
safeZoneRadius = max(safeZoneRadius, 1)
```

### 14.3 HQ Relocation

```
for each player HQ:
    if hexDistance(HQ.Hex, center) > safeZoneRadius:
        // Find the closest passable hex within the safe zone
        newHex = closestPassableHexInZone(HQ.Hex, safeZoneRadius)
        HQ.Hex = newHex
        // Broadcast HQ relocation delta
```

### 14.4 Storm Damage

```
for each troop on the map:
    if hexDistance(troop.Hex, center) > safeZoneRadius:
        damage = suddenDeathTurn  // 1 on first SD turn, 2 on second, etc.
        troop.CurrentHP -= damage
        if troop.CurrentHP <= 0:
            destroyTroop(troop)
```

---

## 15. Redis Persistence

### 15.1 Snapshot Strategy

- **Primary storage:** In-memory Go structs (authoritative during gameplay)
- **Backup:** Redis snapshots after every turn ends
- **Key format:** `game:<game_id>` → JSON-serialized full game state
- **TTL:** Keys expire 24 hours after last update (auto-cleanup of stale games)

### 15.2 Snapshot Trigger Points

| Event | Action |
|---|---|
| Turn ends (either player) | Snapshot full game state to Redis |
| Game ends | Snapshot final state, set short TTL (1 hour) |
| Server shutdown | Snapshot all active games |

### 15.3 Game State Recovery

On server startup:

```
1. Scan Redis for keys matching game:*
2. For each key, deserialize the game state
3. Re-create the game goroutine with the restored state
4. Game enters a "waiting for reconnect" state
5. Both players must reconnect within 60s or the game is forfeit
```

### 15.4 Redis Key Schema

| Key Pattern | Value | TTL |
|---|---|---|
| `game:<game_id>` | JSON game state blob | 24h (refreshed on each snapshot) |

Future (Phase 2+):

| Key Pattern | Value | TTL |
|---|---|---|
| `stats:<player_id>` | JSON player stats | None |
| `leaderboard` | Sorted set (score → player_id) | None |

---

## 16. Configuration

### 16.1 Environment Variables

| Variable | Default | Description |
|---|---|---|
| `PORT` | `8080` | HTTP server listen port |
| `REDIS_URL` | `redis://localhost:6379` | Redis connection string |
| `LOG_LEVEL` | `INFO` | Minimum log level: DEBUG, INFO, WARN, ERROR |
| `CORS_ORIGINS` | `*` | Comma-separated allowed origins. `*` for dev, explicit domains for production |
| `BALANCE_FILE` | `data/balance.yaml` | Path to game balance data file |
| `WS_PING_INTERVAL` | `15s` | WebSocket ping interval |
| `WS_PONG_TIMEOUT` | `10s` | Time to wait for pong before disconnect |
| `RECONNECT_TIMEOUT` | `60s` | Time allowed for player reconnection |
| `ROOM_TTL` | `5m` | Room expiry if opponent doesn't join |
| `SHUTDOWN_DRAIN_TIMEOUT` | `30s` | Max wait time for active games during shutdown |

### 16.2 Balance Data File (`data/balance.yaml`)

```yaml
economy:
  starting_coins: 1000
  passive_income: 100
  structure_income: 50

troops:
  marine:
    cost: 100
    hp: 10
    atk: 3
    def: 14
    mobility: 3
    range: 1
    damage: "1D6+1"
  sniper:
    cost: 150
    hp: 6
    atk: 4
    def: 11
    mobility: 2
    range: 3
    damage: "1D8"
  hoverbike:
    cost: 200
    hp: 8
    atk: 4
    def: 12
    mobility: 5
    range: 1
    damage: "1D8+1"
  mech:
    cost: 350
    hp: 12
    atk: 5
    def: 10
    mobility: 1
    range: 3
    damage: "2D6+2"
    anti_structure_multiplier: 2

structures:
  outpost:
    hp: 8
    atk: 2
    def: 12
    range: 2
    damage: "1D4"
    income: 50
    spawn: true
  command_center:
    hp: 15
    atk: 4
    def: 15
    range: 3
    damage: "1D6+2"
    income: 50
    spawn: true
  hq:
    hp: 20
    atk: 3
    def: 16
    range: 2
    damage: "1D6"
    income: 0
    spawn: true

neutral_modifiers:
  atk_reduction: 2
  damage_step_down: 1  # reduce each die size by 1 step (D6→D4, D8→D6, etc.)

terrain:
  plains:
    movement_cost: 1
    atk_modifier: 0
    def_modifier: 0
  forest:
    movement_cost: 2
    atk_modifier: 0
    def_modifier: 2
  hills:
    movement_cost: 2
    atk_modifier: 1
    def_modifier: 1
  water:
    passable: false
  mountains:
    passable: false

healing:
  passive_rate: 2  # HP per turn if not in combat

sudden_death:
  turn_thresholds:
    small: 20
    medium: 30
    large: 40
  shrink_rate: 1  # hexes per turn

map_generation:
  noise_thresholds:
    water: 0.15
    plains: 0.55
    forest: 0.75
    hills: 0.88
    # > 0.88 = mountains
  structure_counts:
    small: 3
    medium: 5
    large: 7
  min_passable_ratio: 0.60
  max_retries: 10

matchmaking:
  quick_match_defaults:
    map_size: "medium"
    turn_timer: 90
    turn_mode: "alternating"

win_conditions:
  dominance_turns_required: 3
```

---

## 17. Rate Limiting

### 17.1 REST API Rate Limits

Per-IP using an in-memory token bucket:

| Endpoint Group | Rate | Burst |
|---|---|---|
| `POST /api/v1/guest` | 5/min | 2 |
| `POST /api/v1/rooms` | 10/min | 3 |
| `POST /api/v1/rooms/join` | 10/min | 3 |
| `POST /api/v1/matchmaking/*` | 10/min | 3 |
| `GET /health` | 60/min | 10 |

### 17.2 WebSocket Rate Limits

Per-connection:

| Action | Rate | Behavior on exceed |
|---|---|---|
| Game actions (move/attack/buy/end_turn) | 10/sec | NACK with `RATE_LIMITED` |
| Emotes | 2/sec | Silently dropped |
| Malformed messages | 5 total | Connection terminated |

---

## 18. Logging

### 18.1 Library & Format

- **Library:** `log/slog` (Go stdlib, 1.21+)
- **Output format:** JSON in production, text in development (configurable via `LOG_LEVEL`)
- **Structured fields:** All log entries include contextual fields

### 18.2 Log Level Guidelines

| Level | Use For | Examples |
|---|---|---|
| **DEBUG** | Per-action game details, dice rolls, pathfinding steps, state deltas | `"troop moved" unit_id=x from=(1,2,-3) to=(2,1,-3)` |
| **INFO** | Lifecycle events, server operations | `"game created" game_id=x players=[a,b]`, `"turn ended" game_id=x turn=5` |
| **WARN** | Recoverable issues, suspicious activity | `"reconnect attempt" game_id=x player_id=y`, `"invalid action rejected" code=INVALID_MOVE` |
| **ERROR** | Unexpected failures, infrastructure issues | `"redis snapshot failed" game_id=x err=...`, `"websocket write error" player_id=y err=...` |

### 18.3 Standard Log Fields

Every game-related log entry includes:

```
game_id: string
player_id: string (when applicable)
turn: int (when applicable)
phase: string (current FSM state)
```

---

## 19. Health Check

**Endpoint:** `GET /health`

**Response (200 OK):**

```json
{
    "status": "healthy",
    "uptime_seconds": 3600,
    "active_games": 12,
    "connected_players": 24,
    "waiting_rooms": 3,
    "matchmaking_queue_size": 1,
    "redis": {
        "connected": true,
        "latency_ms": 2
    }
}
```

**Response (503 Service Unavailable):** when Redis is disconnected or the server is shutting down.

---

## 20. Deployment

### 20.1 Dockerfile (Multi-Stage)

```
Stage 1: Go build
  - Base: golang:1.23-alpine
  - Copy go.mod, go.sum, download dependencies
  - Copy source, build binary with CGO_DISABLED=1

Stage 2: Runtime
  - Base: alpine:3.19 (for ca-certificates, timezone data)
  - Copy binary from stage 1
  - Copy data/balance.yaml
  - Expose port
  - Entrypoint: the compiled binary
```

### 20.2 docker-compose.yml

```yaml
services:
  server:
    build: ./server
    ports:
      - "8080:8080"
    environment:
      - REDIS_URL=redis://redis:6379
      - LOG_LEVEL=INFO
      - CORS_ORIGINS=https://yourdomain.com
    depends_on:
      - redis
    restart: unless-stopped

  redis:
    image: redis:7-alpine
    volumes:
      - redis_data:/data
    restart: unless-stopped

  # PostgreSQL reserved for Phase 2+
  # postgres:
  #   image: postgres:16-alpine
  #   environment:
  #     POSTGRES_DB: hexdice
  #     POSTGRES_USER: hexdice
  #     POSTGRES_PASSWORD: ${POSTGRES_PASSWORD}
  #   volumes:
  #     - pg_data:/var/lib/postgresql/data

volumes:
  redis_data:
  # pg_data:
```

### 20.3 Nginx Configuration

```
server {
    listen 443 ssl;
    server_name yourdomain.com;

    ssl_certificate /etc/letsencrypt/live/yourdomain.com/fullchain.pem;
    ssl_certificate_key /etc/letsencrypt/live/yourdomain.com/privkey.pem;

    # Flutter web static files
    location / {
        root /var/www/hexdice/web;
        try_files $uri $uri/ /index.html;
    }

    # REST API
    location /api/ {
        proxy_pass http://localhost:8080;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
    }

    # WebSocket
    location /ws {
        proxy_pass http://localhost:8080;
        proxy_http_version 1.1;
        proxy_set_header Upgrade $http_upgrade;
        proxy_set_header Connection "upgrade";
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_read_timeout 3600s;
    }

    # Health check
    location /health {
        proxy_pass http://localhost:8080;
    }
}
```

---

## 21. Testing Strategy

### 21.1 Scope

Focus on game engine core logic. No HTTP handler or WebSocket integration tests for MVP.

### 21.2 Test Categories

| Category | Package | What's Tested |
|---|---|---|
| **Combat math** | `internal/game` | Hit resolution (D20 + ATK vs DEF), crit/fumble, damage dice, counterattack, terrain modifiers, structure combat |
| **Movement** | `internal/game` | BFS reachable hex calculation, terrain cost, enemy blocking, friendly passthrough, mobility budget |
| **Economy** | `internal/game` | Income calculation, purchase validation, insufficient funds |
| **Win conditions** | `internal/game` | HQ destruction, structure dominance counter, tiebreak logic |
| **Sudden death** | `internal/game` | Zone shrink, escalating damage, HQ relocation, troop destruction |
| **Turn transitions** | `internal/game` | Passive healing, structure regen, income credit, action flag reset |
| **Action validation** | `internal/game` | All error cases for move/attack/buy/end_turn |
| **Hex math** | `internal/hex` | Cube coordinate distance, neighbors, rings, rotation |
| **Map generation** | `internal/mapgen` | Symmetry validation, constraint checking, structure placement rules |
| **Dice rolling** | `internal/dice` | Distribution sanity, seeded reproducibility |

### 21.3 Testing Tools

- **Framework:** Go stdlib `testing` + `github.com/stretchr/testify` (assertions, require)
- **Test data:** Builder functions that construct game states programmatically

```go
// Example builder pattern
game := NewTestGame().
    WithMapSize(Small).
    WithTroop(Player1, Marine, Hex(0, 1, -1)).
    WithTroop(Player2, Sniper, Hex(2, -1, -1)).
    WithStructure(Outpost, Neutral, Hex(3, 0, -3)).
    WithCoins(Player1, 500).
    WithTurn(5).
    Build()
```

### 21.4 Seeded Dice Tests

Since the RNG is seeded per-game, tests can use fixed seeds to get deterministic dice rolls:

```go
game := NewTestGame().WithSeed(42).Build()
// D20 roll with seed 42 will always produce the same sequence
result := game.ResolveAttack(attackerID, targetHex)
assert.Equal(t, 15, result.HitRoll)  // deterministic
```

---

## 22. Capacity Estimation

### 22.1 Per-Game Resource Usage

| Resource | Estimate |
|---|---|
| Game state in memory | ~5-20 KB (depending on map size and troop count) |
| Goroutines per game | 1 (game loop) + 4 (2 read + 2 write per player connection) = 5 |
| Goroutine stack | ~4 KB each = ~20 KB per game |
| WebSocket buffers | ~8 KB per connection × 2 = ~16 KB per game |
| Total per game | ~50-60 KB |

### 22.2 Single VPS Capacity (4 GB RAM, 2 vCPU)

- Available for game server (after OS, Redis, Nginx): ~2 GB
- Theoretical max concurrent games: ~30,000-40,000
- Practical limit (CPU for game logic, WS I/O): ~1,000-5,000 concurrent games
- MVP target: comfortably handles hundreds of concurrent games

### 22.3 WebSocket Message Throughput

- Average message size: ~200-500 bytes (JSON)
- Peak messages per game per second: ~5 (during rapid action sequences)
- At 1,000 concurrent games: ~5,000 messages/sec — well within single-server capacity

---

## 23. Future Considerations (Phase 2+)

### 23.1 Simultaneous Turn Mode

The game engine FSM will need additional states:

```
PlanningPhase → MovementResolution → AttackResolution → TurnTransition
```

- Both players submit orders during PlanningPhase (with a timer)
- Server collects and resolves movement conflicts (speed priority)
- Server resolves attacks with hex-targeting and adjacency splash penalty
- New delta types needed: `orders_submitted`, `movement_resolved`, `attacks_resolved`

### 23.2 Fog of War

- Each troop and structure gets a `vision_range` field
- Server computes visible hexes per player each turn
- Game state snapshots sent to clients are **filtered** to only include visible information
- Delta broadcasts are filtered per player (each player sees different deltas)
- New concept: `VisibilityMask` per player, recalculated on every troop move

### 23.3 PostgreSQL Integration

- Add `internal/store/postgres.go` implementing the same store interface
- Schema: `players` table, `match_history` table, `stats` table
- Migrate from guest-only to optional accounts
- Write match results at game end

### 23.4 Horizontal Scaling

If demand exceeds single-server capacity:
- Use Redis Pub/Sub for cross-server game state coordination
- Sticky sessions (route both players to the same server) via Nginx
- Matchmaking queue in Redis (shared across servers)
- Consider dedicated game server instances behind a routing layer
