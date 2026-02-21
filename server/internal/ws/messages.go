package ws

import (
	"encoding/json"
	"fmt"

	"github.com/teomiscia/hexbattle/internal/model"
)

// --- Message Envelope ---

// Envelope is the top-level message wrapper for all WebSocket messages.
type Envelope struct {
	Type string          `json:"type"`
	Seq  int             `json:"seq,omitempty"`
	Data json.RawMessage `json:"data"`
}

// NewEnvelope creates a new message envelope with the given type and data.
func NewEnvelope(msgType string, data interface{}) ([]byte, error) {
	dataBytes, err := json.Marshal(data)
	if err != nil {
		return nil, fmt.Errorf("ws: failed to marshal message data: %w", err)
	}
	env := Envelope{
		Type: msgType,
		Data: dataBytes,
	}
	return json.Marshal(env)
}

// NewEnvelopeWithSeq creates a new message envelope with a sequence number.
func NewEnvelopeWithSeq(msgType string, seq int, data interface{}) ([]byte, error) {
	dataBytes, err := json.Marshal(data)
	if err != nil {
		return nil, fmt.Errorf("ws: failed to marshal message data: %w", err)
	}
	env := Envelope{
		Type: msgType,
		Seq:  seq,
		Data: dataBytes,
	}
	return json.Marshal(env)
}

// --- Client → Server Message Types ---

const (
	MsgJoinGame  = "join_game"
	MsgReconnect = "reconnect"
	MsgMove      = "move"
	MsgAttack    = "attack"
	MsgBuy       = "buy"
	MsgEndTurn   = "end_turn"
	MsgEmote     = "emote"
	MsgPong      = "pong"
)

// JoinGameData is sent by the client to associate with a game room.
type JoinGameData struct {
	RoomID string `json:"room_id"`
}

// ReconnectData is sent by the client to reconnect to an active game.
type ReconnectData struct {
	GameID      string `json:"game_id"`
	PlayerToken string `json:"player_token"`
}

// MoveData is sent by the client to move a troop.
type MoveData struct {
	UnitID  string `json:"unit_id"`
	TargetQ int    `json:"target_q"`
	TargetR int    `json:"target_r"`
	TargetS int    `json:"target_s"`
}

// AttackData is sent by the client to attack a target.
type AttackData struct {
	UnitID  string `json:"unit_id"`
	TargetQ int    `json:"target_q"`
	TargetR int    `json:"target_r"`
	TargetS int    `json:"target_s"`
}

// BuyData is sent by the client to purchase a troop.
type BuyData struct {
	UnitType    model.TroopType `json:"unit_type"`
	StructureID string          `json:"structure_id"`
}

// EmoteData is sent/received for emote messages.
type EmoteData struct {
	PlayerID string `json:"player_id,omitempty"`
	EmoteID  string `json:"emote_id"`
}

// --- Server → Client Message Types ---

const (
	MsgGameState          = "game_state"
	MsgAck                = "ack"
	MsgNack               = "nack"
	MsgTroopMoved         = "troop_moved"
	MsgCombatResult       = "combat_result"
	MsgTroopPurchased     = "troop_purchased"
	MsgTroopDestroyed     = "troop_destroyed"
	MsgStructureAttacked  = "structure_attacked"
	MsgStructureFires     = "structure_fires"
	MsgTurnStart          = "turn_start"
	MsgGameOver           = "game_over"
	MsgPlayerDisconnected = "player_disconnected"
	MsgPlayerReconnected  = "player_reconnected"
	MsgPing               = "ping"
	MsgMatchFound         = "match_found"
	MsgError              = "error"
)

// AckData acknowledges a client action.
type AckData struct {
	Seq        int    `json:"seq"`
	ActionType string `json:"action_type"`
}

// NackData rejects a client action.
type NackData struct {
	Seq        int       `json:"seq"`
	ActionType string    `json:"action_type"`
	Error      ErrorData `json:"error"`
}

// ErrorData holds a machine-readable code and human-readable message.
type ErrorData struct {
	Code    model.ErrorCode `json:"code"`
	Message string          `json:"message"`
}

// TroopMovedData is broadcast when a troop moves.
type TroopMovedData struct {
	UnitID            string `json:"unit_id"`
	FromQ             int    `json:"from_q"`
	FromR             int    `json:"from_r"`
	FromS             int    `json:"from_s"`
	ToQ               int    `json:"to_q"`
	ToR               int    `json:"to_r"`
	ToS               int    `json:"to_s"`
	RemainingMobility int    `json:"remaining_mobility"`
}

// CombatResultData is broadcast when combat is resolved.
type CombatResultData struct {
	AttackerID  string `json:"attacker_id"`
	DefenderID  string `json:"defender_id"`
	HitRoll     int    `json:"hit_roll"`
	NaturalRoll int    `json:"natural_roll"`
	Hit         bool   `json:"hit"`
	DamageRoll  int    `json:"damage_roll"`
	Damage      int    `json:"damage"`
	DefenderHP  int    `json:"defender_hp"`
	Killed      bool   `json:"killed"`
	Crit        bool   `json:"crit"`
	Fumble      bool   `json:"fumble"`
	// Counterattack fields (zero values if no counter)
	HasCounter     bool `json:"has_counter"`
	CounterHitRoll int  `json:"counter_hit_roll,omitempty"`
	CounterNatural int  `json:"counter_natural_roll,omitempty"`
	CounterHit     bool `json:"counter_hit,omitempty"`
	CounterDamage  int  `json:"counter_damage,omitempty"`
	AttackerHP     int  `json:"attacker_hp"`
	AttackerKilled bool `json:"attacker_killed"`
}

// TroopPurchasedData is broadcast when a troop is purchased.
type TroopPurchasedData struct {
	UnitID         string          `json:"unit_id"`
	UnitType       model.TroopType `json:"unit_type"`
	HexQ           int             `json:"hex_q"`
	HexR           int             `json:"hex_r"`
	HexS           int             `json:"hex_s"`
	Owner          string          `json:"owner"`
	CoinsRemaining int             `json:"coins_remaining"`
}

// TroopDestroyedData is broadcast when a troop is destroyed.
type TroopDestroyedData struct {
	UnitID string `json:"unit_id"`
	HexQ   int    `json:"hex_q"`
	HexR   int    `json:"hex_r"`
	HexS   int    `json:"hex_s"`
	Cause  string `json:"cause"` // "combat", "sudden_death", "structure_fire"
}

// StructureAttackedData is broadcast when a structure takes damage.
type StructureAttackedData struct {
	StructureID string `json:"structure_id"`
	AttackerID  string `json:"attacker_id"`
	HitRoll     int    `json:"hit_roll"`
	Damage      int    `json:"damage"`
	StructureHP int    `json:"structure_hp"`
	Captured    bool   `json:"captured"`
	NewOwner    string `json:"new_owner,omitempty"`
}

// StructureFiresData is broadcast when a structure attacks a troop.
type StructureFiresData struct {
	StructureID string `json:"structure_id"`
	TargetID    string `json:"target_id"`
	HitRoll     int    `json:"hit_roll"`
	Damage      int    `json:"damage"`
	TargetHP    int    `json:"target_hp"`
	Killed      bool   `json:"killed"`
}

// HealedUnit records a unit that was healed at turn start.
type HealedUnit struct {
	UnitID   string `json:"unit_id"`
	HPBefore int    `json:"hp_before"`
	HPAfter  int    `json:"hp_after"`
}

// StructureRegen records a structure that regenerated HP.
type StructureRegen struct {
	StructureID string `json:"structure_id"`
	HPBefore    int    `json:"hp_before"`
	HPAfter     int    `json:"hp_after"`
}

// SuddenDeathDamage records storm damage to a troop.
type SuddenDeathDamage struct {
	UnitID  string `json:"unit_id"`
	Damage  int    `json:"damage"`
	HPAfter int    `json:"hp_after"`
	Killed  bool   `json:"killed"`
}

// TurnStartData is broadcast when a new turn begins.
type TurnStartData struct {
	TurnNumber         int                 `json:"turn_number"`
	ActivePlayerID     string              `json:"active_player_id"`
	TimerSeconds       int                 `json:"timer_seconds"`
	IncomeGained       int                 `json:"income_gained"`
	StructureIncome    int                 `json:"structure_income"`
	TotalCoins         int                 `json:"total_coins"`
	HealedUnits        []HealedUnit        `json:"healed_units"`
	StructureRegens    []StructureRegen    `json:"structure_regens"`
	SuddenDeathDamages []SuddenDeathDamage `json:"sudden_death_damage"`
}

// GameOverData is broadcast when the game ends.
type GameOverData struct {
	WinnerID string                         `json:"winner_id"`
	Reason   model.WinReason                `json:"reason"`
	Stats    map[string]model.GameOverStats `json:"stats"` // player_id -> stats
}

// PlayerDisconnectedData is broadcast when a player disconnects.
type PlayerDisconnectedData struct {
	PlayerID string `json:"player_id"`
}

// PlayerReconnectedData is broadcast when a player reconnects.
type PlayerReconnectedData struct {
	PlayerID string `json:"player_id"`
}

// MatchFoundData is sent to a waiting player when a match is found.
type MatchFoundData struct {
	RoomID string `json:"room_id"`
}
