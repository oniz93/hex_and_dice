package model

// GamePhase represents the current state of the game FSM.
type GamePhase string

const (
	PhaseWaitingForPlayers GamePhase = "waiting_for_players"
	PhaseGeneratingMap     GamePhase = "generating_map"
	PhaseGameStarted       GamePhase = "game_started"
	PhaseTurnStart         GamePhase = "turn_start"
	PhaseStructureCombat   GamePhase = "structure_combat"
	PhasePlayerAction      GamePhase = "player_action"
	PhaseTurnTransition    GamePhase = "turn_transition"
	PhaseGameOver          GamePhase = "game_over"
)

// TurnMode determines how turns are structured.
type TurnMode string

const (
	TurnModeAlternating  TurnMode = "alternating"
	TurnModeSimultaneous TurnMode = "simultaneous" // Phase 2
)

// MapSize determines the hex grid radius and structure counts.
type MapSize string

const (
	MapSizeSmall  MapSize = "small"
	MapSizeMedium MapSize = "medium"
	MapSizeLarge  MapSize = "large"
)

// MapRadius returns the hex grid radius for this map size.
func (ms MapSize) Radius() int {
	switch ms {
	case MapSizeSmall:
		return 7
	case MapSizeMedium:
		return 10
	case MapSizeLarge:
		return 13
	default:
		return 10
	}
}

// TroopType identifies the four troop types.
type TroopType string

const (
	TroopMarine    TroopType = "marine"
	TroopSniper    TroopType = "sniper"
	TroopHoverbike TroopType = "hoverbike"
	TroopMech      TroopType = "mech"
)

// StructureType identifies the three structure types.
type StructureType string

const (
	StructureOutpost       StructureType = "outpost"
	StructureCommandCenter StructureType = "command_center"
	StructureHQ            StructureType = "hq"
)

// TerrainType identifies the five terrain types.
type TerrainType string

const (
	TerrainPlains    TerrainType = "plains"
	TerrainForest    TerrainType = "forest"
	TerrainHills     TerrainType = "hills"
	TerrainWater     TerrainType = "water"
	TerrainMountains TerrainType = "mountains"
)

// RoomState represents the lifecycle state of a room.
type RoomState string

const (
	RoomWaitingForOpponent RoomState = "waiting_for_opponent"
	RoomReady              RoomState = "ready"
	RoomGameInProgress     RoomState = "game_in_progress"
	RoomGameOver           RoomState = "game_over"
)

// WinReason describes why a game ended.
type WinReason string

const (
	WinReasonHQDestroyed        WinReason = "HQ_DESTROYED"
	WinReasonStructureDominance WinReason = "STRUCTURE_DOMINANCE"
	WinReasonSuddenDeath        WinReason = "SUDDEN_DEATH"
	WinReasonForfeit            WinReason = "FORFEIT"
	WinReasonDisconnect         WinReason = "DISCONNECT"
	WinReasonDraw               WinReason = "DRAW"
)

// ErrorCode represents machine-readable error codes.
type ErrorCode string

const (
	ErrNotYourTurn       ErrorCode = "NOT_YOUR_TURN"
	ErrInvalidMove       ErrorCode = "INVALID_MOVE"
	ErrInvalidAttack     ErrorCode = "INVALID_ATTACK"
	ErrInsufficientFunds ErrorCode = "INSUFFICIENT_FUNDS"
	ErrSpawnOccupied     ErrorCode = "SPAWN_OCCUPIED"
	ErrSpawnNotOwned     ErrorCode = "SPAWN_NOT_OWNED"
	ErrUnitAlreadyActed  ErrorCode = "UNIT_ALREADY_ACTED"
	ErrUnitNotReady      ErrorCode = "UNIT_NOT_READY"
	ErrUnitNotFound      ErrorCode = "UNIT_NOT_FOUND"
	ErrGameNotFound      ErrorCode = "GAME_NOT_FOUND"
	ErrRoomNotFound      ErrorCode = "ROOM_NOT_FOUND"
	ErrRoomFull          ErrorCode = "ROOM_FULL"
	ErrRoomExpired       ErrorCode = "ROOM_EXPIRED"
	ErrInvalidMessage    ErrorCode = "INVALID_MESSAGE"
	ErrRateLimited       ErrorCode = "RATE_LIMITED"
)
