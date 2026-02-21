package model

// PlayerState holds a player's in-game state (coins, owned structures, etc.).
type PlayerState struct {
	ID                   string `json:"id"`
	Nickname             string `json:"nickname"`
	Coins                int    `json:"coins"`
	DominanceTurnCounter int    `json:"dominance_turn_counter"`
	IsDisconnected       bool   `json:"is_disconnected"`
}

// RoomSettings holds the configurable options for a game room.
type RoomSettings struct {
	MapSize   MapSize  `json:"map_size"`
	TurnTimer int      `json:"turn_timer"` // seconds: 60, 90, or 120
	TurnMode  TurnMode `json:"turn_mode"`
}

// DefaultRoomSettings returns the Quick Match defaults.
func DefaultRoomSettings() RoomSettings {
	return RoomSettings{
		MapSize:   MapSizeMedium,
		TurnTimer: 90,
		TurnMode:  TurnModeAlternating,
	}
}

// GameOverStats holds the end-of-game statistics.
type GameOverStats struct {
	TurnsPlayed      int `json:"turns_played"`
	TroopsKilled     int `json:"troops_killed"`
	TroopsLost       int `json:"troops_lost"`
	StructuresHeld   int `json:"structures_held"`
	TotalDamageDealt int `json:"total_damage_dealt"`
}
