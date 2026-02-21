# Hex & Dice — Frontend High-Level Design

## 1. Overview

This document describes the high-level architecture of the Hex & Dice frontend client. The client is built with **Flutter + Flame engine**, targeting Web, Android, and iOS from a single codebase. The client renders the game state received from the server, provides local validation for responsiveness, and sends player actions over WebSocket.

**Framework:** Flutter (stable channel, latest release)
**Game engine:** Flame (embedded in Flutter widget tree)
**State management:** Riverpod
**Navigation:** GoRouter
**Architecture principle:** Flame handles the game canvas (hex grid, sprites, animations). Flutter handles all UI chrome (menus, HUD overlays, popups, forms).

---

## 2. Project Structure

```
client/
├── lib/
│   ├── main.dart                           # App entrypoint, ProviderScope, GoRouter setup
│   ├── app.dart                            # MaterialApp.router configuration, theme
│   │
│   ├── config/
│   │   ├── theme.dart                      # ThemeData, ColorScheme, dark theme definition
│   │   ├── routes.dart                     # GoRouter route definitions
│   │   ├── constants.dart                  # API URLs, timing constants, emote list
│   │   └── environment.dart                # Dev/staging/prod environment config
│   │
│   ├── models/
│   │   ├── game_state.dart                 # GameState: full mirror of server state
│   │   ├── troop.dart                      # Troop model: id, type, owner, hex, hp, flags
│   │   ├── structure.dart                  # Structure model: id, type, owner, hex, hp
│   │   ├── hex_tile.dart                   # HexTile: coordinates, terrain type
│   │   ├── player_state.dart              # PlayerState: coins, structure list, troop list
│   │   ├── room.dart                       # Room model: code, settings, players, state
│   │   ├── combat_result.dart              # CombatResult: all dice rolls, damage, outcomes
│   │   └── enums.dart                      # TroopType, StructureType, Terrain, GamePhase, etc.
│   │
│   ├── game/
│   │   ├── hex_game.dart                   # FlameGame subclass: game world root component
│   │   ├── components/
│   │   │   ├── hex_map_component.dart      # Renders the full hex grid (terrain tiles)
│   │   │   ├── hex_tile_component.dart     # Single hex tile: terrain sprite + highlight overlay
│   │   │   ├── troop_component.dart        # Troop sprite: animation state machine, position
│   │   │   ├── structure_component.dart    # Structure sprite: type, owner color, HP bar
│   │   │   ├── highlight_component.dart    # Colored overlay for move/attack highlights
│   │   │   ├── projectile_component.dart   # Projectile/slash effect for combat animation
│   │   │   ├── damage_text_component.dart  # Floating damage number on hit
│   │   │   ├── dice_component.dart         # Inline D20 / damage dice animation
│   │   │   ├── storm_zone_component.dart   # Sudden death zone visual (darkened outer hexes)
│   │   │   └── emote_bubble_component.dart # Speech bubble for emotes above HQ
│   │   ├── systems/
│   │   │   ├── input_system.dart           # Tap/click → world coord → hex coord conversion
│   │   │   ├── selection_system.dart       # Selection state machine (Idle, TroopSelected, etc.)
│   │   │   ├── animation_queue.dart        # Sequential delta animation queue
│   │   │   └── camera_controller.dart      # Pan/zoom input handling, bounds clamping
│   │   ├── hex/
│   │   │   ├── cube_coord.dart             # Cube coordinate (q, r, s), arithmetic, distance
│   │   │   ├── hex_layout.dart             # Hex-to-pixel and pixel-to-hex conversion (pointy-top)
│   │   │   ├── hex_utils.dart              # Neighbors, rings, line drawing, distance
│   │   │   └── pathfinding.dart            # BFS reachable hexes (client-side validation)
│   │   └── data/
│   │       └── balance.dart                # Game balance constants (mirrored from server YAML)
│   │
│   ├── providers/
│   │   ├── session_provider.dart           # Player session: token, player_id, nickname
│   │   ├── connection_provider.dart        # WebSocket connection state, reconnect logic
│   │   ├── game_state_provider.dart        # GameState notifier: applies deltas, exposes state
│   │   ├── room_provider.dart              # Room state: creation, joining, status
│   │   ├── matchmaking_provider.dart       # Matchmaking queue state
│   │   ├── selection_provider.dart         # Current selection state (for HUD reactivity)
│   │   ├── audio_provider.dart             # Audio settings: volumes, mute state
│   │   └── settings_provider.dart          # Persisted local settings (nickname, server URL)
│   │
│   ├── services/
│   │   ├── api_service.dart                # REST API client (Dio): guest, rooms, matchmaking
│   │   ├── ws_service.dart                 # WebSocket client: connect, send, typed message stream
│   │   ├── message_parser.dart             # JSON → typed ServerMessage deserialization
│   │   ├── audio_service.dart              # FlameAudio wrapper: play BGM, play SFX, volume
│   │   └── storage_service.dart            # SharedPreferences wrapper: settings, reconnect data
│   │
│   ├── screens/
│   │   ├── loading_screen.dart             # Asset preloading with themed progress bar
│   │   ├── title_screen.dart               # Game logo, Play button, How to Play, Settings
│   │   ├── play_screen.dart                # Create Room / Join Room / Quick Match options
│   │   ├── room_screen.dart                # Waiting room: room code display, settings, cancel
│   │   ├── matchmaking_screen.dart         # "Searching..." with elapsed timer, cancel
│   │   ├── game_screen.dart                # Flame GameWidget + HUD overlays
│   │   ├── game_over_screen.dart           # Winner, stats, Play Again / Return to Menu
│   │   └── how_to_play_screen.dart         # Rules text with diagrams
│   │
│   ├── widgets/
│   │   ├── hud/
│   │   │   ├── top_bar.dart                # Turn number, player indicator, timer countdown
│   │   │   ├── bottom_bar.dart             # Coin count, End Turn button, emote button
│   │   │   ├── troop_popup.dart            # Positioned stat popup near tapped troop
│   │   │   ├── shop_panel.dart             # Bottom sheet with troop purchase cards
│   │   │   ├── emote_bar.dart              # Expandable emote selector
│   │   │   ├── reconnect_banner.dart       # "Reconnecting..." persistent banner
│   │   │   └── attack_confirm.dart         # Attack confirmation overlay
│   │   ├── common/
│   │   │   ├── pixel_button.dart           # Styled button matching pixel art aesthetic
│   │   │   ├── pixel_text_field.dart       # Styled text input for nickname, room code
│   │   │   ├── hp_bar.dart                 # Horizontal HP bar with color gradient
│   │   │   ├── coin_display.dart           # Coin icon + amount
│   │   │   ├── stat_row.dart               # Single stat label + value row (for popups)
│   │   │   └── settings_dialog.dart        # Audio settings popup (sliders, mute)
│   │   └── troop_card.dart                 # Troop card for shop: sprite preview, name, cost, stats
│   │
│   └── validation/
│       ├── move_validator.dart             # Client-side movement validation (BFS, range, blocking)
│       ├── attack_validator.dart           # Client-side attack range + target validation
│       ├── buy_validator.dart              # Client-side purchase validation (funds, spawn, occupied)
│       └── turn_validator.dart             # Is it my turn, has this unit acted, is unit ready
│
├── assets/
│   ├── images/
│   │   ├── terrain_atlas.png              # All terrain hex tile sprites (plains, forest, hills, water, mountains)
│   │   ├── troops_atlas.png               # All troop sprites: 4 types × (idle, walk, attack, death) × 2 team colors
│   │   ├── structures_atlas.png           # All structure sprites: 3 types × (neutral, red, blue)
│   │   ├── effects_atlas.png             # Projectile sprites, slash effects, hit/miss indicators
│   │   ├── dice_atlas.png                # D20 spin frames, D6/D8/D4 result faces
│   │   ├── ui_atlas.png                  # UI elements: icons, emote icons, HP bar, coin icon
│   │   └── logo.png                      # Game logo for title screen and loading
│   ├── audio/
│   │   ├── music/
│   │   │   ├── menu_theme.ogg            # Menu/lobby background music
│   │   │   └── battle_theme.ogg          # In-game background music
│   │   └── sfx/
│   │       ├── attack_hit.wav
│   │       ├── attack_miss.wav
│   │       ├── troop_move.wav
│   │       ├── troop_death.wav
│   │       ├── dice_roll.wav
│   │       ├── turn_start.wav
│   │       ├── structure_capture.wav
│   │       ├── coin_gain.wav
│   │       ├── purchase.wav
│   │       └── emote_pop.wav
│   └── fonts/
│       └── pixel_font.ttf                # Pixel art bitmap font (e.g., Press Start 2P)
│
├── pubspec.yaml
├── analysis_options.yaml
└── test/                                  # Placeholder for future client tests
```

### Key directories and their responsibilities

| Directory | Responsibility |
|---|---|
| `lib/config/` | App-wide configuration: theme, routes, constants, environment |
| `lib/models/` | Data classes mirroring server state structures |
| `lib/game/` | Flame engine: game world, components, rendering, hex math, animation |
| `lib/game/components/` | Individual Flame components (sprites, effects, overlays on the canvas) |
| `lib/game/systems/` | Game systems: input processing, selection FSM, animation queue, camera |
| `lib/game/hex/` | Hex grid math (cube coords, layout, pathfinding) — shared with validation |
| `lib/providers/` | Riverpod state providers for all app-wide state |
| `lib/services/` | External integrations: REST API, WebSocket, audio, local storage |
| `lib/screens/` | Top-level Flutter screens (one per route) |
| `lib/widgets/` | Reusable Flutter UI components (HUD elements, common widgets) |
| `lib/validation/` | Client-side game rule validation (mirrors server logic) |
| `assets/` | Sprite sheets, audio files, fonts |

---

## 3. External Dependencies

| Package | Purpose | pub.dev |
|---|---|---|
| **flame** | 2D game engine: sprites, animation, camera, game loop | `flame` |
| **flutter_riverpod** | State management | `flutter_riverpod` |
| **riverpod_annotation** | Code-gen for Riverpod providers | `riverpod_annotation` |
| **go_router** | Declarative navigation/routing | `go_router` |
| **dio** | HTTP client for REST API | `dio` |
| **web_socket_channel** | Cross-platform WebSocket client | `web_socket_channel` |
| **shared_preferences** | Local key-value persistence | `shared_preferences` |
| **uuid** | UUIDv4 generation (for client-side temp IDs) | `uuid` |
| **json_annotation** | JSON serialization annotations | `json_annotation` |
| **json_serializable** | JSON serialization code generation | `json_serializable` (dev) |
| **build_runner** | Code generation runner | `build_runner` (dev) |
| **flutter_lints** | Lint rules | `flutter_lints` (dev) |

Flame includes `flame_audio` (FlameAudio) as part of the core package — no separate audio dependency needed.

---

## 4. App Lifecycle

### 4.1 Startup Sequence

```
1. main() → runApp(ProviderScope(child: App()))
2. App widget initializes GoRouter, applies ThemeData
3. Initial route: /loading (LoadingScreen)
4. LoadingScreen:
   a. Initialize SharedPreferences (load saved settings, nickname, reconnect data)
   b. Preload all sprite sheet atlases into Flame image cache
   c. Preload all audio files (BGM + SFX) into FlameAudio cache
   d. Load pixel font
   e. Show themed progress bar during loading
   f. Check for reconnect data: if game_id + token exist, attempt reconnect flow
   g. On completion → navigate to /title
5. TitleScreen: start menu BGM loop
```

### 4.2 Reconnect on Startup

If reconnect data exists in SharedPreferences (game_id, player_token):

```
1. Attempt REST call to verify game still exists (GET /api/v1/rooms/{code} or similar)
2. If game active → open WebSocket → send reconnect message → navigate to /game
3. If game not found → clear reconnect data → navigate to /title normally
4. If network error → clear reconnect data → navigate to /title with error snackbar
```

### 4.3 App Backgrounding / Foregrounding

```
On background (WidgetsBindingObserver.didChangeAppLifecycleState):
  - Pause BGM
  - WebSocket stays open (OS manages TCP keepalive)
  - Save reconnect data (game_id, token) to SharedPreferences

On foreground:
  - Resume BGM
  - Check WebSocket connection health
  - If disconnected → trigger auto-reconnect flow
```

---

## 5. Navigation & Routes

### 5.1 Route Definitions

| Route | Screen | Description |
|---|---|---|
| `/loading` | LoadingScreen | Asset preloading, initial setup |
| `/title` | TitleScreen | Main menu: Play, How to Play, Settings |
| `/play` | PlayScreen | Create Room / Join Room / Quick Match |
| `/room/:code` | RoomScreen | Waiting for opponent, room code display |
| `/matchmaking` | MatchmakingScreen | Queue waiting screen |
| `/game/:id` | GameScreen | Active game: Flame canvas + HUD overlays |
| `/game-over/:id` | GameOverScreen | Winner, stats, Play Again / Menu |
| `/how-to-play` | HowToPlayScreen | Rules text and diagrams |

### 5.2 Navigation Flow

```
/loading → /title → /play → /room/:code    → /game/:id → /game-over/:id → /title
                           → /matchmaking   → /game/:id → /game-over/:id → /title
                   → /how-to-play → /title
```

### 5.3 Route Guards

- `/game/:id` requires an active WebSocket connection and valid game state. If not available, redirect to `/title` with error.
- `/room/:code` validates room exists before completing navigation. Redirects to `/play` with error if room not found.
- Back navigation from `/game/:id` shows a "Leave game? You will forfeit" confirmation dialog.

---

## 6. State Management (Riverpod)

### 6.1 Provider Architecture

```
                    ┌─────────────────────┐
                    │   SettingsProvider   │  ← SharedPreferences (persisted)
                    │  (nickname, audio,   │
                    │   reconnect data)    │
                    └─────────┬───────────┘
                              │
                    ┌─────────▼───────────┐
                    │   SessionProvider    │  ← REST: POST /api/v1/guest
                    │  (token, player_id,  │
                    │   nickname)          │
                    └─────────┬───────────┘
                              │
              ┌───────────────┼───────────────┐
              │               │               │
    ┌─────────▼─────┐ ┌──────▼──────┐ ┌──────▼──────────┐
    │ RoomProvider   │ │ Matchmaking │ │ ConnectionProv.  │
    │ (room state,   │ │ Provider    │ │ (WS state,       │
    │  settings)     │ │ (queue)     │ │  reconnect)      │
    └───────┬────────┘ └──────┬──────┘ └──────┬───────────┘
            │                 │               │
            └─────────────────┼───────────────┘
                              │
                    ┌─────────▼───────────┐
                    │ GameStateProvider    │  ← WebSocket deltas
                    │ (full game state,    │
                    │  delta application)  │
                    └─────────┬───────────┘
                              │
                    ┌─────────▼───────────┐
                    │ SelectionProvider    │  ← User interaction
                    │ (selection FSM,      │
                    │  highlights)         │
                    └─────────────────────┘
```

### 6.2 Provider Descriptions

| Provider | Type | State | Responsibility |
|---|---|---|---|
| **SettingsProvider** | StateNotifier | `{nickname, musicVol, sfxVol, muted, reconnectData}` | Load/save local settings from SharedPreferences |
| **SessionProvider** | AsyncNotifier | `{token, playerId, nickname}` | Guest registration via REST API. Holds auth credentials for the session |
| **RoomProvider** | StateNotifier | `{room, status, error}` | Room creation, joining, polling status |
| **MatchmakingProvider** | StateNotifier | `{inQueue, elapsedTime, matchedRoomId}` | Queue join/leave, match found detection |
| **ConnectionProvider** | StateNotifier | `{wsState, isConnected, reconnecting, reconnectAttempts}` | WebSocket lifecycle, auto-reconnect with backoff |
| **GameStateProvider** | StateNotifier | `GameState` | Authoritative client game state. Applies deltas from server. Exposes computed properties (my troops, opponent troops, reachable hexes) |
| **SelectionProvider** | StateNotifier | `{state, selectedUnitId, highlightedHexes, attackableHexes}` | Selection FSM state. Drives hex highlighting and available actions |
| **AudioProvider** | StateNotifier | `{musicVol, sfxVol, muted}` | Audio volume control, syncs with SettingsProvider for persistence |

### 6.3 GameState Notifier Detail

The `GameStateProvider` is the central piece. It:

1. Receives the full `game_state` snapshot on join/reconnect → replaces entire local state
2. Receives per-action deltas → applies them to the local state in order
3. Enqueues deltas into the animation queue (for the Flame engine to animate sequentially)
4. Exposes computed properties used by the HUD and Flame components:
   - `myTroops` / `opponentTroops`
   - `myStructures` / `opponentStructures` / `neutralStructures`
   - `isMyTurn`
   - `canAfford(TroopType)`
   - `myCoins`
   - `turnNumber`, `timerRemaining`
   - `gamePhase`
   - `winner` (null until game ends)

### 6.4 Delta Application

When a delta arrives from the WebSocket:

```
1. WsService emits typed ServerMessage on its stream
2. GameStateProvider listens to the stream
3. For each delta:
   a. Apply state mutation (update troop position, HP, coins, etc.)
   b. Push delta event onto the AnimationQueue
   c. Notify listeners (Riverpod rebuild triggers)
4. AnimationQueue processes events sequentially:
   a. Pop next delta
   b. Trigger corresponding Flame animation (tween, sprite anim, effect)
   c. Wait for animation completion callback
   d. Pop next delta
```

---

## 7. Networking Layer

### 7.1 API Service (REST)

Thin wrapper around Dio, exposing typed methods:

```dart
class ApiService {
  // Guest registration
  Future<Session> registerGuest(String nickname);

  // Room management
  Future<Room> createRoom(RoomSettings settings);
  Future<Room> joinRoom(String code);
  Future<Room> getRoomStatus(String code);

  // Matchmaking
  Future<void> joinMatchmaking();
  Future<void> leaveMatchmaking();
  Future<MatchmakingStatus> getMatchmakingStatus();
}
```

**Base URL:** Configured per environment (dev: `http://localhost:8080`, prod: `https://yourdomain.com`).

**Auth:** All requests include `Authorization: Bearer <token>` header (set via Dio interceptor after guest registration).

**Error handling:** Dio interceptor catches HTTP errors, maps to typed error classes with server error codes (`ROOM_NOT_FOUND`, `ROOM_FULL`, etc.). Surfaced to UI via provider error states.

### 7.2 WebSocket Service

Singleton service managing the WebSocket connection lifecycle:

```dart
class WsService {
  // Connection
  Future<void> connect(String token);
  void disconnect();

  // Sending (typed, auto-serialized, auto-sequenced)
  void sendJoinGame(String roomId);
  void sendReconnect(String gameId, String token);
  void sendMove(String unitId, CubeCoord target);
  void sendAttack(String unitId, CubeCoord target);
  void sendBuy(TroopType type, String structureId);
  void sendEndTurn();
  void sendEmote(String emoteId);
  void sendPong();

  // Receiving
  Stream<ServerMessage> get messages;        // typed parsed messages
  Stream<ConnectionState> get connectionState; // connected/disconnected/reconnecting

  // Internal
  int _seqCounter = 0;                       // auto-incrementing sequence number
  Timer? _pongTimer;                         // pong response timeout
}
```

### 7.3 Message Parser

Deserializes incoming JSON WebSocket messages into typed Dart classes:

```dart
ServerMessage parseMessage(String json) {
  final map = jsonDecode(json);
  switch (map['type']) {
    case 'game_state':    return GameStateMessage.fromJson(map['data']);
    case 'ack':           return AckMessage.fromJson(map['data']);
    case 'nack':          return NackMessage.fromJson(map['data']);
    case 'troop_moved':   return TroopMovedDelta.fromJson(map['data']);
    case 'combat_result': return CombatResultDelta.fromJson(map['data']);
    case 'troop_purchased': return TroopPurchasedDelta.fromJson(map['data']);
    // ... all server message types
  }
}
```

All message classes are generated with `json_serializable` for type-safe deserialization.

### 7.4 Auto-Reconnect with Backoff

```
On WebSocket disconnect detected (stream closes or pong timeout):

1. Set connectionState = reconnecting
2. Show reconnect banner in HUD
3. Save reconnect data to SharedPreferences (game_id, token)
4. Attempt reconnect:
   - Delay: 1s, 2s, 4s, 8s, 10s, 10s, 10s... (exponential with max 10s cap)
   - Each attempt: open new WebSocket → send reconnect message
   - On success: connectionState = connected, hide banner, receive full state snapshot
   - On failure: increment attempt counter, schedule next retry
5. After 60s total (matching server's forfeit timeout): stop retrying, show fatal dialog
```

### 7.5 Sequence Number Tracking

- Client maintains a `_seqCounter` integer, incremented for each outgoing message
- Each outgoing message includes the current `seq` value
- ACK/NACK messages from the server include the `seq` they're responding to
- Client can track pending actions (sent but not yet ACK'd) via a `Map<int, PendingAction>`
- If a NACK arrives, the pending action is removed and an error is surfaced to the UI

---

## 8. Flame Game Engine

### 8.1 HexGame Class

The root `FlameGame` subclass that manages the game world:

```dart
class HexGame extends FlameGame with HasTappables, PanDetector, ScaleDetector {
  // Components
  late HexMapComponent hexMap;          // terrain tile grid
  late CameraController cameraCtrl;     // pan/zoom handling

  // State references (set from Flutter via Riverpod)
  GameState? gameState;
  SelectionState? selectionState;

  // Systems
  late AnimationQueue animationQueue;
  late InputSystem inputSystem;
  late SelectionSystem selectionSystem;

  @override
  Future<void> onLoad() async {
    // Build hex map from game state
    // Spawn troop and structure components
    // Configure camera bounds
    // Start idle animations
  }
}
```

### 8.2 Component Hierarchy

```
HexGame (FlameGame)
└── World
    ├── HexMapComponent
    │   └── HexTileComponent × N         (terrain sprites + highlight overlays)
    ├── StructureComponent × M            (structure sprites on their hexes)
    ├── TroopComponent × K                (troop sprites with animation state machines)
    └── EffectsLayer
        ├── ProjectileComponent           (temporary, during combat)
        ├── DamageTextComponent           (temporary, floating numbers)
        ├── DiceComponent                 (temporary, during combat)
        ├── StormZoneComponent            (persistent during sudden death)
        └── EmoteBubbleComponent          (temporary, 3s display)
```

All game entities are children of the World component. The CameraComponent looks at the World, providing pan/zoom.

### 8.3 Flutter ↔ Flame Communication

Communication between the Flutter widget layer (Riverpod) and the Flame game engine:

**Flutter → Flame (pushing state and events):**

```dart
// In GameScreen widget
ref.listen(gameStateProvider, (prev, next) {
  gameRef.updateGameState(next);
});

ref.listen(selectionProvider, (prev, next) {
  gameRef.updateSelection(next);
});
```

The GameWidget's `game` instance is accessed via a key or direct reference. When Riverpod state changes, the listener calls methods on the HexGame instance that update components.

**Flame → Flutter (user interactions):**

```dart
// In HexGame, when a hex is tapped
void onHexTapped(CubeCoord hex) {
  // Call back to Flutter/Riverpod via a callback or stream
  onHexTapCallback?.call(hex);
}
```

The GameScreen widget registers callbacks on the HexGame instance. Taps on the game canvas are translated to hex coordinates, then passed up to Flutter where Riverpod providers process the interaction (selection FSM, sending actions to server).

---

## 9. Hex Grid Rendering

### 9.1 Coordinate System

Cube coordinates `(q, r, s)` where `q + r + s = 0`. Identical to the server's coordinate system.

### 9.2 Hex Layout (Pointy-Top)

Hex-to-pixel conversion for pointy-top hexes:

```dart
class HexLayout {
  final double hexSize;  // distance from center to vertex

  // Hex center in world pixel coordinates
  Offset hexToPixel(CubeCoord hex) {
    final x = hexSize * (sqrt(3) * hex.q + sqrt(3) / 2 * hex.r);
    final y = hexSize * (3.0 / 2 * hex.r);
    return Offset(x, y);
  }

  // World pixel coordinates to nearest hex (for tap detection)
  CubeCoord pixelToHex(Offset point) {
    final q = (sqrt(3) / 3 * point.dx - 1.0 / 3 * point.dy) / hexSize;
    final r = (2.0 / 3 * point.dy) / hexSize;
    return cubeRound(q, r, -q - r);  // round to nearest hex
  }
}
```

### 9.3 Hex Size

The `hexSize` determines how large hexes appear in world coordinates. With 32px sprite tiles:

- `hexSize = 18` (in world units) — each hex tile is ~36px wide, matching the 32px sprite with 2px margin per side
- The hex sprite is drawn centered on the hex's pixel position

### 9.4 Terrain Tile Rendering

Each `HexTileComponent` is a `SpriteComponent` that:

1. Looks up its terrain type from the game state
2. Selects the corresponding sprite region from `terrain_atlas.png`
3. Positions itself at `hexToPixel(coord)`
4. Optionally renders a highlight overlay on top (blue/red/yellow, semi-transparent)

The entire hex grid is built once when the game starts and doesn't change (terrain is static). Only highlight overlays are dynamically added/removed.

### 9.5 Render Order (Z-Ordering)

Components are rendered in this order (back to front):

```
1. Terrain tiles (lowest)
2. Hex highlights (move/attack overlays)
3. Storm zone darkening (sudden death)
4. Structures
5. Troops
6. Effects (projectiles, damage numbers, dice)
7. Emote bubbles (highest)
```

Z-priority is set via the `priority` property on each Flame component.

---

## 10. Sprite & Animation System

### 10.1 Sprite Atlases

| Atlas | Contents | Approximate Size |
|---|---|---|
| `terrain_atlas.png` | 5 terrain types × pointy-top hex tile variant. ~36×32px per tile. | ~256×64px |
| `troops_atlas.png` | 4 troop types × 4 animation states (idle, walk, attack, death) × 2-4 frames × 2 team colors (red/blue). 32×32px per frame. | ~512×256px |
| `structures_atlas.png` | 3 structure types × 3 ownership states (neutral, red, blue). 32×32px or 32×48px per sprite. | ~256×128px |
| `effects_atlas.png` | Projectile sprites, slash effects, hit flash, miss indicator. Various sizes. | ~256×128px |
| `dice_atlas.png` | D20 spin frames (8-12 frames), D6/D8/D4 result faces. 32×32px per frame. | ~256×128px |
| `ui_atlas.png` | Emote icons (~6), coin icon, stat icons (sword, shield, boot, crosshair), HP bar segments. | ~256×64px |

All atlases are loaded into Flame's image cache during the loading screen. Individual sprite regions are extracted using `Sprite.fromImage()` with `srcPosition` and `srcSize`.

### 10.2 Troop Animation State Machine

Each `TroopComponent` runs an internal animation state machine:

```
States:
  - Idle:    2-3 frame loop, ~400ms per frame. Default state.
  - Walking: 3-4 frame loop, ~150ms per frame. During move tween.
  - Attack:  3-4 frames, plays once, ~100ms per frame. During attack action.
  - Death:   2-3 frames, plays once, ~200ms per frame. On destruction.

Transitions:
  Idle → Walking     (on move delta received)
  Walking → Idle     (on move tween complete)
  Idle → Attack      (on attack delta received)
  Attack → Idle      (on attack animation complete, if alive)
  Attack → Death     (on destruction delta received)
  Idle → Death       (on destruction delta received)
  Death → [removed]  (component removed from game world after death anim)
```

Team color is handled via separate sprite rows in the atlas (red row, blue row). The correct row is selected based on the troop's owner.

### 10.3 Movement Animation

When a `troop_moved` delta is processed:

```
1. AnimationQueue pops the delta
2. TroopComponent switches to Walking animation state
3. EffectComponent: tween TroopComponent.position from hexToPixel(from) to hexToPixel(to) over 300ms (ease-in-out curve)
4. On tween complete: switch to Idle state
5. AnimationQueue signals completion, processes next delta
```

### 10.4 Combat Animation Sequence

When a `combat_result` delta is processed:

```
1. AnimationQueue pops the delta
2. Attacker switches to Attack animation state (plays once, ~300ms)
3. At mid-point of attack animation:
   a. Spawn ProjectileComponent at attacker position
   b. Tween projectile to defender position (~200ms)
   c. Spawn DiceComponent near combat area:
      - D20 spin animation (~800ms)
      - Land on the natural roll value
      - Display hit/miss text and modifier calculation
4. On projectile arrival at defender:
   a. If hit: defender flashes white (2 frames, ~100ms). Spawn DamageTextComponent floating up with damage number. Play "attack_hit" SFX.
   b. If miss: spawn "MISS" text floating up. Play "attack_miss" SFX.
   c. If crit: extra visual flash (golden), larger damage text with "CRIT!" label
   d. If fumble: "FUMBLE!" text above attacker
5. If damage dice exist (hit confirmed):
   a. DiceComponent shows damage dice rolling briefly (~400ms)
   b. Display final damage number
6. If counterattack occurs:
   a. Brief pause (~200ms)
   b. Repeat steps 2-5 with attacker/defender swapped, half-damage label
7. If either unit destroyed: play Death animation, then remove component
8. Remove DiceComponent and ProjectileComponent
9. AnimationQueue signals completion
```

Total combat animation time: ~2-3 seconds per engagement.

### 10.5 Dice Animation

The `DiceComponent` renders inline near the combat area (offset to not overlap troops):

- D20: 8-12 frame spin animation. Final frame shows the result face. Color-coded: green border on hit, red on miss, gold on crit, purple on fumble.
- Damage dice (if hit): smaller dice shown briefly next to the D20 result. Show the sum.
- All dice results also rendered as text below the dice sprite: `"D20: 15 + 3 ATK = 18 vs 14 DEF → HIT"` (simplified for screen space: just show the key number).

---

## 11. Camera System

### 11.1 Configuration

```dart
CameraComponent camera:
  - viewfinder.zoom: range [minZoom, maxZoom]
  - minZoom: calculated to fit entire map on screen
  - maxZoom: ~3.0 (close enough to see ~5 hex radius)
  - initial zoom: minZoom (show entire map on game start)
  - world bounds: rectangular bounding box of all hex pixel positions + padding
```

### 11.2 Input Handling

| Gesture | Platform | Action |
|---|---|---|
| Drag/pan | Mobile (single finger drag on empty area) | Pan camera |
| Pinch | Mobile (two-finger pinch) | Zoom camera |
| Scroll wheel | Web/desktop (mouse wheel) | Zoom camera |
| Click-drag | Web/desktop (click + drag on empty area) | Pan camera |
| Arrow keys | Web/desktop (keyboard) | Pan camera |

Pan and zoom have smooth interpolation (lerp) for fluid camera movement. Camera clamps to world bounds (cannot pan past the edge of the map).

### 11.3 Camera Animations

The camera auto-pans to relevant events:

- **Turn start:** Smooth pan to center on the active player's HQ (~500ms tween)
- **Combat:** If the combat is off-screen, auto-pan to center on it before playing the animation
- **Sudden death zone shrink:** Brief zoom-out to show the zone boundary, then zoom back

Auto-pan is interruptible: if the player touches/drags during an auto-pan, the auto-pan cancels immediately.

---

## 12. Selection System (Interaction FSM)

### 12.1 State Machine

```
                        ┌─────────────┐
                        │    Idle     │
                        └──────┬──────┘
                               │ tap own troop
                               ▼
                     ┌─────────────────────┐
              ┌──────│   TroopSelected     │──────┐
              │      │  (show highlights)  │      │
              │      └──────┬──────────────┘      │
              │             │                     │
      tap blue hex     tap red hex          tap empty/same troop
              │             │                     │
              ▼             ▼                     ▼
     ┌────────────┐  ┌──────────────┐     ┌──────────┐
     │  Moving    │  │ ConfirmAttack│     │   Idle   │
     │ (send move)│  │  (show dlg)  │     └──────────┘
     └─────┬──────┘  └──────┬───────┘
           │                │
      wait for ACK    confirm → send attack
           │                │
           ▼                ▼
     ┌──────────┐    ┌──────────┐
     │   Idle   │    │   Idle   │
     └──────────┘    └──────────┘
```

Additional transitions:
- **Tap another own troop** from TroopSelected → switch selection to new troop (stay in TroopSelected)
- **Tap own structure** from Idle → show stat popup + "Buy Troops" button
- **Tap enemy troop** from Idle → show stat popup (read-only)
- **Tap enemy troop** from TroopSelected (and in attack range) → show stat popup with "Attack" button → ConfirmAttack
- **Tap structure** from TroopSelected (and in attack range) → same as enemy troop attack flow

### 12.2 Highlight Computation

When entering `TroopSelected` state:

```
1. Run client-side BFS pathfinding from selected troop's position
   - Uses the troop's remaining mobility and terrain costs
   - Excludes enemy-occupied hexes and impassable terrain
   - Excludes friendly-occupied hexes (can pass through but not stop on)
   → Set of reachable hexes (blue highlights)

2. Compute attackable hexes:
   - All hexes within the troop's attack range from its current position
   - That contain an enemy troop or enemy/neutral structure
   → Set of attackable hexes (red highlights)

3. If the troop has already moved: reachable hexes = empty set
4. If the troop has already attacked: attackable hexes = empty set
5. Publish highlight sets to SelectionProvider → Flame renders overlays
```

### 12.3 Action Dispatch

When the player confirms an action (tap blue hex for move, confirm dialog for attack, tap shop card for buy):

```
1. Client-side validation (MoveValidator / AttackValidator / BuyValidator)
2. If invalid: show error snackbar (e.g., "Not enough coins"), stay in current state
3. If valid:
   a. Send action to server via WsService (move/attack/buy with seq number)
   b. Add to pending actions map: seq → PendingAction
   c. Show brief "pending" visual indicator on the acting unit (subtle spinner or opacity change)
   d. Transition selection to Idle
4. On ACK: remove from pending map, normal delta processing handles visuals
5. On NACK: remove from pending map, show error snackbar with server error message
```

---

## 13. HUD (Flutter Overlays)

### 13.1 Layout Architecture

The `GameScreen` uses a `Stack` widget:

```
Stack(
  children: [
    // Layer 1: Flame game canvas (fills entire screen)
    GameWidget(game: hexGame),

    // Layer 2: Top bar (positioned at top)
    Positioned(top: 0, left: 0, right: 0, child: TopBar()),

    // Layer 3: Bottom bar (positioned at bottom)
    Positioned(bottom: 0, left: 0, right: 0, child: BottomBar()),

    // Layer 4: Dynamic overlays (conditional)
    if (showTroopPopup) Positioned(..., child: TroopPopup()),
    if (showShopPanel) BottomSheet(child: ShopPanel()),
    if (showAttackConfirm) AttackConfirmOverlay(),
    if (showReconnectBanner) ReconnectBanner(),
    if (showEmoteBar) EmoteBar(),
  ]
)
```

### 13.2 Top Bar

```
┌──────────────────────────────────────────────────┐
│  Turn 12    ⚔ Player 1's Turn    ⏱ 0:47         │
└──────────────────────────────────────────────────┘
```

- **Turn number:** "Turn X" label
- **Active player indicator:** Player name/nickname with team color dot (red/blue). "Your Turn" / "Opponent's Turn"
- **Timer countdown:** Digital countdown (MM:SS). Color transitions: white (> 30s) → yellow (10-30s) → red (< 10s) with subtle pulse under 10s
- **Background:** Semi-transparent dark overlay for readability over the game canvas

### 13.3 Bottom Bar

```
┌──────────────────────────────────────────────────┐
│  🪙 750     [😀]                    [End Turn]    │
└──────────────────────────────────────────────────┘
```

- **Coin display:** Coin icon + current balance. Briefly animates (bounce + green text) when income is gained
- **Emote button:** Small button that expands the emote bar upward when tapped
- **End Turn button:** Large, prominent button. Disabled when it's not the player's turn. Pixel-styled button with "END TURN" text
- **Background:** Semi-transparent dark overlay

### 13.4 Troop Stat Popup

Appears as a floating panel near the tapped troop, offset to not cover the unit:

```
┌─────────────────────┐
│ Marine          ❤10/10│
│ ─────────────────── │
│ ATK  +3    DEF  14  │
│ MOB   3    RNG   1  │
│ DMG  1D6+1          │
│ Status: Ready       │
│ [Attack] (if applicable)│
└─────────────────────┘
```

- Positioned using the troop's world-to-screen coordinate (via camera projection)
- Clamped to screen bounds so it doesn't go off-screen
- Dismissed by tapping elsewhere
- Shows an "Attack" button only when viewing an enemy while a friendly troop is selected and in range

### 13.5 Shop Panel (Bottom Sheet)

Slides up from the bottom when tapping an owned spawn structure:

```
┌──────────────────────────────────────────────────┐
│  Buy Troops at [Outpost]                    [✕]  │
├──────────────────────────────────────────────────┤
│ ┌──────────┐ ┌──────────┐ ┌──────────┐ ┌──────────┐ │
│ │ [sprite] │ │ [sprite] │ │ [sprite] │ │ [sprite] │ │
│ │ Marine   │ │ Sniper   │ │ Hoverbike│ │  Mech    │ │
│ │  100🪙   │ │  150🪙   │ │  200🪙   │ │  350🪙   │ │
│ │ HP:10    │ │ HP:6     │ │ HP:8     │ │ HP:12    │ │
│ │ ATK:+3   │ │ ATK:+4   │ │ ATK:+4   │ │ ATK:+5   │ │
│ │ [Buy]    │ │ [Buy]    │ │ [Buy]    │ │ [Buy]    │ │
│ └──────────┘ └──────────┘ └──────────┘ └──────────┘ │
└──────────────────────────────────────────────────┘
```

- Horizontal scrollable row of troop cards
- Each card shows: pixel art sprite, name, cost, key stats
- "Buy" button grayed out if unaffordable (cost > current coins)
- "Buy" button grayed out if spawn hex is occupied
- Tap "Buy" → sends buy action to server, closes shop panel
- Close via [X] button or swipe down

### 13.6 Attack Confirmation

A small overlay near the attacker and target:

```
┌────────────────────────────────┐
│ Attack Sniper with Marine?     │
│ [Cancel]           [Confirm]   │
└────────────────────────────────┘
```

- Shows attacker name → target name
- Positioned between the two units on screen
- "Confirm" sends the attack action
- "Cancel" returns to TroopSelected state

### 13.7 Emote Bar

Expands upward from the emote button:

```
         ┌─────────────────────┐
         │ GG │ 👍 │ 😮 │ 😬 │ ⚔ │
         └─────────────────────┘
         [😀]  ← emote button
```

- 5-6 predefined emote icons in a horizontal row
- Tap an emote → sends emote action → closes bar
- Opponent's emote appears as a speech bubble above their HQ for 3 seconds (`EmoteBubbleComponent` in Flame)

### 13.8 Reconnect Banner

Persistent banner at the top of the screen (below the top bar):

```
┌──────────────────────────────────────────────────┐
│ ⚠ Connection lost. Reconnecting... (attempt 3)   │
└──────────────────────────────────────────────────┘
```

- Yellow/orange background
- Shows reconnect attempt count
- Auto-dismisses on successful reconnect
- Replaces with "Opponent disconnected. Waiting for reconnect... (45s)" if opponent disconnects

---

## 14. Audio System

### 14.1 FlameAudio Integration

```dart
class AudioService {
  // BGM
  Future<void> playMenuMusic();     // loops menu_theme.ogg
  Future<void> playBattleMusic();   // loops battle_theme.ogg
  Future<void> stopMusic();
  Future<void> crossfadeTo(String track, {Duration duration = 500ms});

  // SFX
  void playAttackHit();
  void playAttackMiss();
  void playTroopMove();
  void playTroopDeath();
  void playDiceRoll();
  void playTurnStart();
  void playStructureCapture();
  void playCoinGain();
  void playPurchase();
  void playEmotePop();

  // Volume control
  void setMusicVolume(double vol);  // 0.0 - 1.0
  void setSfxVolume(double vol);    // 0.0 - 1.0
  void setMuted(bool muted);
}
```

### 14.2 BGM Behavior

| Screen | Track | Behavior |
|---|---|---|
| Loading | None | Silent |
| Title, Play, Room, Matchmaking | `menu_theme.ogg` | Loop, starts on title, crossfade on entering game |
| Game | `battle_theme.ogg` | Loop, crossfade from menu music (~500ms) |
| Game Over | `menu_theme.ogg` | Crossfade back from battle music |

Music pauses on app background, resumes on foreground.

### 14.3 SFX Triggers

| Event | Sound | Triggered By |
|---|---|---|
| `troop_moved` delta | `troop_move.wav` | Movement animation start |
| `combat_result` (hit) | `dice_roll.wav` → `attack_hit.wav` | Dice animation → hit flash |
| `combat_result` (miss) | `dice_roll.wav` → `attack_miss.wav` | Dice animation → miss text |
| `troop_destroyed` delta | `troop_death.wav` | Death animation start |
| `troop_purchased` delta | `purchase.wav` | Shop panel closes |
| `turn_start` delta | `turn_start.wav` | Turn transition |
| `structure_attacked` (captured) | `structure_capture.wav` | Structure ownership change |
| `turn_start` (income > 0) | `coin_gain.wav` | Income credited |
| `emote` received | `emote_pop.wav` | Emote bubble appears |

### 14.4 Volume Persistence

Audio settings (music volume, SFX volume, mute state) are:
- Stored in `SharedPreferences` via `SettingsProvider`
- Loaded on app startup during the loading screen
- Applied to `AudioService` immediately
- Updated in real-time when the player adjusts sliders in the settings dialog

---

## 15. Client-Side Validation

### 15.1 Validation Module Structure

The validation layer mirrors the server's validation logic, implemented in Dart. It operates on the local `GameState` and returns either success or a typed error.

```dart
// Move validation
MoveValidation validateMove(GameState state, String unitId, CubeCoord target) {
  // 1. Unit exists and belongs to current player
  // 2. Unit has not already moved
  // 3. Unit is ready (not purchased this turn)
  // 4. Target is within map bounds
  // 5. Target is passable terrain
  // 6. Target is not occupied by enemy
  // 7. BFS pathfinding confirms target is reachable within mobility
  // 8. Path does not pass through enemy hexes
  return MoveValidation.valid() or MoveValidation.invalid(reason)
}

// Attack validation
AttackValidation validateAttack(GameState state, String unitId, CubeCoord target) {
  // 1. Unit exists and belongs to current player
  // 2. Unit has not already attacked
  // 3. Unit is ready
  // 4. Target within attack range (hex distance)
  // 5. Target contains enemy troop or enemy/neutral structure
  return AttackValidation.valid() or AttackValidation.invalid(reason)
}

// Buy validation
BuyValidation validateBuy(GameState state, TroopType type, String structureId) {
  // 1. Player has enough coins
  // 2. Structure belongs to player
  // 3. Structure hex is not occupied
  return BuyValidation.valid() or BuyValidation.invalid(reason)
}

// Turn validation
TurnValidation validateEndTurn(GameState state) {
  // 1. It is the current player's turn
  return TurnValidation.valid() or TurnValidation.invalid(reason)
}
```

### 15.2 Pathfinding (Client-Side BFS)

Identical algorithm to the server (see Backend HLD section 11.1):

```dart
Set<CubeCoord> reachableHexes(CubeCoord start, int mobility, GameState state) {
  // BFS with terrain movement costs
  // Skip impassable terrain
  // Skip enemy-occupied hexes
  // Allow passing through friendly-occupied hexes (but not stopping on them)
  // Return set of hexes the unit can move to
}
```

This is used both for:
- Highlight computation (showing blue hexes when a troop is selected)
- Pre-validation before sending a move action to the server

### 15.3 Balance Data Mirror

The file `lib/game/data/balance.dart` contains hardcoded Dart constants matching the server's `balance.yaml`:

```dart
const troopStats = {
  TroopType.marine: TroopStatBlock(cost: 100, hp: 10, atk: 3, def: 14, mobility: 3, range: 1, damage: '1D6+1'),
  TroopType.sniper: TroopStatBlock(cost: 150, hp: 6, atk: 4, def: 11, mobility: 2, range: 3, damage: '1D8'),
  // ...
};

const terrainModifiers = {
  Terrain.plains: TerrainMod(moveCost: 1, atkMod: 0, defMod: 0),
  Terrain.forest: TerrainMod(moveCost: 2, atkMod: 0, defMod: 2),
  // ...
};

const passiveIncome = 100;
const structureIncome = 50;
const startingCoins = 1000;
const healingRate = 2;
```

These constants must be kept manually in sync with the server's `balance.yaml`. Any balance change requires updating both files.

---

## 16. Responsive Layout

### 16.1 Design Principles

- **Mobile:** Portrait-only orientation. HUD elements at top and bottom edges. Game canvas fills the middle. Touch-first interaction.
- **Web/desktop:** No forced orientation. Layout adapts to window size. Same HUD placement (top/bottom bars scale horizontally). Mouse click = tap, scroll wheel = zoom.
- **Single layout** that flexes based on available space — no separate mobile/desktop layout variants.

### 16.2 Sizing Strategy

| Element | Sizing Approach |
|---|---|
| Top bar | Full width, fixed height (~48px mobile, ~56px desktop) |
| Bottom bar | Full width, fixed height (~56px mobile, ~64px desktop) |
| Game canvas | Fills remaining vertical space between top and bottom bars |
| Troop popup | Fixed width (~200px), positioned relative to tap point, clamped to screen |
| Shop panel | Full width bottom sheet, fixed height (~220px) |
| Buttons | Min tap target 48×48px (mobile accessibility guideline) |
| Font size | Pixel font at native pixel multiples (8px, 16px, 24px) for crisp rendering |

### 16.3 Safe Areas

On mobile devices with notches/rounded corners:
- Top bar respects `MediaQuery.of(context).padding.top` (SafeArea)
- Bottom bar respects `MediaQuery.of(context).padding.bottom`
- Game canvas extends under the safe areas for immersive feel, but HUD content stays within safe bounds

### 16.4 Input Handling Per Platform

| Input | Mobile | Web/Desktop |
|---|---|---|
| Select troop/hex | Tap | Left click |
| Pan camera | Single-finger drag on empty space | Click-drag on empty space / Arrow keys |
| Zoom camera | Two-finger pinch | Scroll wheel |
| Deselect | Tap empty space | Click empty space / Escape key |
| End turn | Tap "End Turn" button | Click button / Enter key |
| Open settings | Tap gear icon | Click gear icon |

Keyboard shortcuts (web/desktop only, optional):
- `Enter` — End turn
- `Escape` — Deselect / close popup
- `Arrow keys` — Pan camera
- `+` / `-` — Zoom in/out

---

## 17. Theming & Visual Design

### 17.1 Dark Theme

Single dark theme, defined via Flutter's `ThemeData`:

```dart
ThemeData(
  brightness: Brightness.dark,
  colorScheme: ColorScheme.dark(
    primary: Color(0xFF00D4FF),      // sci-fi cyan (buttons, active states)
    secondary: Color(0xFFFF8C00),     // accent orange (highlights, warnings)
    surface: Color(0xFF1A1A2E),       // dark navy (card backgrounds, panels)
    background: Color(0xFF0F0F1A),    // near-black (screen backgrounds)
    error: Color(0xFFFF4444),          // red (errors, enemy indicators)
    onPrimary: Color(0xFF000000),
    onSurface: Color(0xFFE0E0E0),     // light gray (body text)
    onBackground: Color(0xFFFFFFFF),  // white (headings)
  ),
  fontFamily: 'PixelFont',
)
```

### 17.2 Team Colors

| Team | Primary Color | Used For |
|---|---|---|
| Player 1 (Red) | `Color(0xFFE63946)` | Troop tint, structure tint, HUD indicator, name color |
| Player 2 (Blue) | `Color(0xFF457BF7)` | Troop tint, structure tint, HUD indicator, name color |
| Neutral | `Color(0xFF808080)` | Neutral structure tint |

### 17.3 Highlight Colors

| Highlight | Color | Opacity |
|---|---|---|
| Reachable hex (move) | Blue `(0xFF457BF7)` | 30% |
| Attackable hex | Red `(0xFFE63946)` | 30% |
| Selected unit hex | Yellow `(0xFFFFD700)` | 40% |
| Storm zone (sudden death) | Purple-black `(0xFF2D0040)` | 50% |

### 17.4 Pixel Font

- **Font:** Press Start 2P (Google Fonts, OFL license) or equivalent pixel bitmap font
- **Sizes:** Multiples of the base pixel grid for crisp rendering:
  - Headings: 16px
  - Body / HUD labels: 8px or 10px
  - Damage numbers: 12px (bold)
  - Small labels: 8px
- **Anti-aliasing:** Disabled (nearest-neighbor scaling) to maintain sharp pixel edges

---

## 18. Loading Screen

### 18.1 Sequence

```
1. Show game logo (centered, pixel art)
2. Show pixel art loading bar (empty, centered below logo)
3. Begin asset loading:
   a. SharedPreferences initialization (~instant)
   b. Sprite sheet atlas loading (terrain, troops, structures, effects, dice, UI)
   c. Audio preloading (BGM + all SFX)
   d. Font loading
4. Update progress bar proportionally as each asset group completes
5. Optional: animated pixel art element (spinning D20 or marching troop) above the progress bar
6. On all assets loaded:
   a. Check for reconnect data → attempt reconnect OR navigate to /title
   b. Brief pause (~300ms) for visual polish
   c. Fade transition to next screen
```

### 18.2 Progress Calculation

| Asset Group | Weight | % of Bar |
|---|---|---|
| Settings (SharedPreferences) | 5% | 0-5% |
| Sprite sheets (6 atlases) | 50% | 5-55% |
| Audio (2 BGM + 10 SFX) | 35% | 55-90% |
| Fonts | 5% | 90-95% |
| Finalization | 5% | 95-100% |

---

## 19. Screens Detail

### 19.1 Title Screen

```
┌──────────────────────────────────┐
│                                  │
│          ╔═══════════╗           │
│          ║ HEX & DICE║           │
│          ╚═══════════╝           │
│            [pixel logo]          │
│                                  │
│          ┌────────────┐          │
│          │    PLAY    │          │
│          └────────────┘          │
│          ┌────────────┐          │
│          │ HOW TO PLAY│          │
│          └────────────┘          │
│          ┌────────────┐          │
│          │  SETTINGS  │          │
│          └────────────┘          │
│                                  │
│              v0.1.0              │
└──────────────────────────────────┘
```

- Game logo (pixel art, animated idle or subtle glow)
- Three pixel-styled buttons, vertically stacked
- Settings opens a dialog overlay (audio sliders, mute toggle)
- Menu BGM starts/continues playing
- Version number at the bottom

### 19.2 Play Screen

```
┌──────────────────────────────────┐
│  [←]     Choose Mode             │
│                                  │
│  Nickname: [_______________]     │
│                                  │
│  ┌─────────────────────────┐     │
│  │     CREATE ROOM         │     │
│  └─────────────────────────┘     │
│  ┌─────────────────────────┐     │
│  │     JOIN ROOM           │     │
│  │  Code: [______]         │     │
│  └─────────────────────────┘     │
│  ┌─────────────────────────┐     │
│  │     QUICK MATCH         │     │
│  └─────────────────────────┘     │
│                                  │
└──────────────────────────────────┘
```

- **Nickname field:** Pre-filled with last used nickname (from SharedPreferences). Validated (3-16 chars, alphanumeric + underscore). Saved on change.
- **Create Room:** Opens room creation sub-panel with settings (map size picker, turn timer picker). Then navigates to `/room/:code`.
- **Join Room:** 6-digit code input field. Validates format client-side. Calls `POST /api/v1/rooms/join`. On success navigates to `/game/:id`.
- **Quick Match:** Immediately calls guest registration (if not already registered) → joins matchmaking queue → navigates to `/matchmaking`.

Guest registration (`POST /api/v1/guest`) is called lazily on the first action that requires a session (Create Room, Join Room, or Quick Match). Token is cached in `SessionProvider`.

### 19.3 Room Screen

```
┌──────────────────────────────────┐
│  [←]     Waiting Room            │
│                                  │
│    Share this code:              │
│                                  │
│    ┌────────────────────┐        │
│    │      482917        │        │
│    └────────────────────┘        │
│    (tap to copy)                 │
│                                  │
│    Map: Medium                   │
│    Timer: 90s                    │
│    Mode: Alternating             │
│                                  │
│    Waiting for opponent...       │
│    [animated dots/spinner]       │
│                                  │
│    ┌─────────────────────┐       │
│    │       CANCEL        │       │
│    └─────────────────────┘       │
│                                  │
└──────────────────────────────────┘
```

- Room code displayed prominently, tap-to-copy to clipboard
- Room settings summary shown below
- Polls room status via REST or receives `match_found` via WebSocket
- When opponent joins (room state = Ready): auto-navigate to `/game/:id`
- Cancel button destroys the room and returns to `/play`

### 19.4 Matchmaking Screen

```
┌──────────────────────────────────┐
│  [←]     Quick Match             │
│                                  │
│                                  │
│    Searching for opponent...     │
│    [animated pixel spinner]      │
│                                  │
│    Elapsed: 0:23                 │
│                                  │
│    ┌─────────────────────┐       │
│    │       CANCEL        │       │
│    └─────────────────────┘       │
│                                  │
└──────────────────────────────────┘
```

- Elapsed time counter (updated every second)
- When matched: auto-navigate to `/game/:id`
- Cancel removes from queue and returns to `/play`

### 19.5 Game Screen

(Described in detail in sections 8, 9, 12, 13 — Flame canvas + Flutter HUD overlays)

### 19.6 Game Over Screen

```
┌──────────────────────────────────┐
│                                  │
│           VICTORY!               │
│      (or DEFEAT / DRAW)          │
│                                  │
│    Reason: HQ Destroyed          │
│                                  │
│    ┌─────────────────────┐       │
│    │ Turns Played:   15  │       │
│    │ Troops Killed:   8  │       │
│    │ Troops Lost:     5  │       │
│    │ Structures Held: 4  │       │
│    └─────────────────────┘       │
│                                  │
│    ┌─────────────────────┐       │
│    │     PLAY AGAIN      │       │
│    └─────────────────────┘       │
│    ┌─────────────────────┐       │
│    │    RETURN TO MENU   │       │
│    └─────────────────────┘       │
│                                  │
└──────────────────────────────────┘
```

- Victory/Defeat/Draw announcement with appropriate color (green/red/yellow)
- Win reason text
- Basic game stats from the `game_over` delta
- "Play Again" navigates back to `/room/:code` (same room, new game if opponent also clicks)
- "Return to Menu" navigates to `/title`
- Clear reconnect data from SharedPreferences on this screen

### 19.7 How to Play Screen

- Scrollable page with rules text organized in sections
- Sections: Overview, Troops (with stat tables), Structures, Combat (with dice example), Terrain, Economy, Win Conditions, Sudden Death
- Inline pixel art diagrams showing:
  - Hex grid with troop movement example
  - Attack range visualization
  - D20 + ATK vs DEF formula
- Back button returns to `/title`

---

## 20. Error Handling & Feedback

### 20.1 Error Display Strategy

| Error Type | Display Method | Examples |
|---|---|---|
| REST API error | Snackbar (bottom, auto-dismiss 4s) | "Could not create room", "Room not found" |
| WebSocket disconnect | Persistent banner (top, below top bar) | "Reconnecting... (attempt 3)" |
| WS NACK (action rejected) | Snackbar (bottom, auto-dismiss 3s) | "Invalid move", "Not enough coins" |
| Fatal error | Full-screen dialog with "Return to Menu" | "Game not found", "Server unreachable after 60s" |
| Opponent disconnected | Persistent banner (top) | "Opponent disconnected. Waiting... (45s)" |
| Opponent reconnected | Brief snackbar (bottom, auto-dismiss 2s) | "Opponent reconnected" |

### 20.2 Loading / Pending States

| State | Visual Indicator |
|---|---|
| Asset loading | Full-screen loading bar (LoadingScreen) |
| Guest registration in progress | Button shows "Connecting..." with spinner |
| Room creation in progress | Button shows "Creating..." with spinner |
| Joining room in progress | Button shows "Joining..." with spinner |
| Waiting for opponent in room | Animated dots + "Waiting for opponent..." |
| Matchmaking queue | Spinner + elapsed time counter |
| Waiting for ACK after action | Subtle opacity pulse on the acting unit (~200ms) |
| Opponent's turn | "Opponent's Turn" indicator in top bar, map remains interactive for viewing |

---

## 21. Web-Specific Considerations

### 21.1 Renderer

- **CanvasKit renderer** (default): ensures pixel-perfect rendering matching mobile. Uses WebAssembly + Skia.
- Initial WASM download: ~2MB, cached by browser after first load
- Build command: `flutter build web --web-renderer canvaskit`

### 21.2 Browser Compatibility

- Targets modern evergreen browsers: Chrome, Firefox, Safari, Edge (latest 2 versions)
- No IE11 support
- WebSocket support required (available in all target browsers)

### 21.3 URL Handling

GoRouter integrates with browser URL bar:
- `/game/abc123` is a valid deep link (though reconnect requires a valid token)
- Browser back button is handled by GoRouter (with leave-game confirmation on game screen)
- Shareable room URLs: `https://yourdomain.com/#/room/482917` (Flutter web uses hash routing by default)

### 21.4 Keyboard Focus

- Game canvas captures keyboard focus when the game screen is active
- Tab / Shift+Tab navigation works for menu screens (accessibility)
- Escape key closes popups and deselects (game screen)

---

## 22. Local Storage Schema

### 22.1 SharedPreferences Keys

| Key | Type | Description |
|---|---|---|
| `nickname` | String | Last used player nickname |
| `music_volume` | double | Music volume 0.0-1.0 |
| `sfx_volume` | double | SFX volume 0.0-1.0 |
| `audio_muted` | bool | Master mute toggle |
| `reconnect_game_id` | String? | Active game ID for reconnect |
| `reconnect_token` | String? | Player token for reconnect |
| `reconnect_room_code` | String? | Room code for reconnect |
| `server_url` | String? | Override server URL (dev/staging) |

All keys are cleared (except audio settings and nickname) when the player navigates to the title screen after a game ends.

---

## 23. Build & Environment Configuration

### 23.1 Environment Modes

| Environment | Server URL | Log Level | Features |
|---|---|---|---|
| **dev** | `http://localhost:8080` | Verbose (all WS messages) | Debug overlays available |
| **staging** | `https://staging.yourdomain.com` | Info | Production-like |
| **prod** | `https://yourdomain.com` | Error only | Release mode |

Configured via `--dart-define=ENV=dev|staging|prod` at build time.

### 23.2 Build Commands

```bash
# Development (web)
flutter run -d chrome --dart-define=ENV=dev

# Development (Android)
flutter run --dart-define=ENV=dev

# Production web build
flutter build web --release --web-renderer canvaskit --dart-define=ENV=prod

# Production Android APK
flutter build apk --release --dart-define=ENV=prod

# Production iOS
flutter build ios --release --dart-define=ENV=prod
```

### 23.3 Debug Overlays (Dev Only)

When `ENV=dev`, a debug panel is available (toggled via a hidden button or keyboard shortcut):
- Current game state JSON dump
- WebSocket message log (last 50 messages)
- FPS counter
- Hex coordinate under cursor
- Current selection FSM state

---

## 24. Future Considerations (Phase 2+)

### 24.1 Simultaneous Turn Mode (Client)

- New selection flow: player queues orders for all units during planning phase
- Queued orders shown as arrows/indicators on the map (move arrows, attack target lines)
- "Submit Orders" button replaces "End Turn"
- Resolution phase: animations play sequentially (movement phase, then attack phase) as the server sends resolution deltas
- New animation types: simultaneous movement (multiple tweens at once), conflict resolution visuals

### 24.2 Fog of War (Client)

- Hex tiles outside vision range rendered with a dark overlay (50-70% opacity black)
- Enemy troops/structures outside vision range are hidden (components removed or hidden)
- Previously seen hexes show terrain (remembered) but no enemy units
- Vision range indicators: subtle circle around each own troop showing its vision range (optional toggle)
- Game state from server is already filtered — client renders exactly what it receives

### 24.3 Interactive Tutorial

- A single-player scenario with no opponent
- Scripted sequence: tooltip popups guide the player through selecting a troop, moving, attacking, buying, ending turn
- Pre-built small map with fixed positions
- Could reuse the game engine with a mock WebSocket service that returns scripted responses

### 24.4 Accounts & Persistence

- Login screen (email/password or OAuth) replacing/supplementing guest registration
- Profile screen showing win/loss stats
- Token stored persistently (not session-scoped) in secure storage (flutter_secure_storage)
- Leaderboard screen consuming new REST endpoints

### 24.5 Performance Optimization

If performance issues arise on low-end devices:
- Reduce animation frames (2 instead of 3-4)
- Disable dice animation (show numbers only)
- Use sprite batching for hex tiles (Flame's SpriteBatch)
- Reduce camera interpolation smoothness
- LOD: at far zoom levels, replace individual hex sprites with a simplified colored grid
