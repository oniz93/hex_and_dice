package game

import (
	"github.com/teomiscia/hexbattle/internal/hex"
	"github.com/teomiscia/hexbattle/internal/model"
)

// BotActionType identifies the type of bot action.
type BotActionType string

const (
	BotActionMove   BotActionType = "move"
	BotActionAttack BotActionType = "attack"
	BotActionBuy    BotActionType = "buy"
)

// BotAction represents a single action the bot wants to take.
type BotAction struct {
	Type        BotActionType
	UnitID      string
	Target      hex.Coord
	TroopType   model.TroopType
	StructureID string
}

// BotPlayer is the interface that bot implementations must satisfy.
// The engine calls NextAction repeatedly during the bot's turn with
// the current (mutated) game state. Return nil to signal the bot is
// done and the turn should end.
type BotPlayer interface {
	// NextAction returns the next action the bot wants to take.
	// Returns nil when the bot is done acting (the engine will end the turn).
	NextAction(gs *GameState) *BotAction

	// PlayerID returns the bot's player ID.
	PlayerID() string
}
