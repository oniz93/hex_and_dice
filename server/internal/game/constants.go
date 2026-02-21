package game

import (
	"fmt"

	"github.com/teomiscia/hexbattle/internal/config"
	"github.com/teomiscia/hexbattle/internal/dice"
	"github.com/teomiscia/hexbattle/internal/hex"
	"github.com/teomiscia/hexbattle/internal/model"
)

// Balance holds the loaded game balance constants.
// Set once at startup via LoadBalance().
var Balance *config.BalanceData

// LoadBalance initializes the global balance data from a config.
func LoadBalance(b *config.BalanceData) {
	Balance = b
}

// TroopCost returns the coin cost for a troop type.
func TroopCost(t model.TroopType) int {
	if Balance == nil {
		return 0
	}
	tc, ok := Balance.Troops[string(t)]
	if !ok {
		return 0
	}
	return tc.Cost
}

// NewTroopFromBalance creates a Troop with stats from balance data.
func NewTroopFromBalance(id string, troopType model.TroopType, ownerID string, pos hex.Coord) (*model.Troop, error) {
	if Balance == nil {
		return nil, fmt.Errorf("balance data not loaded")
	}
	tc, ok := Balance.Troops[string(troopType)]
	if !ok {
		return nil, fmt.Errorf("unknown troop type: %s", troopType)
	}
	return &model.Troop{
		ID:                id,
		Type:              troopType,
		OwnerID:           ownerID,
		Hex:               pos,
		CurrentHP:         tc.HP,
		MaxHP:             tc.HP,
		ATK:               tc.ATK,
		DEF:               tc.DEF,
		Mobility:          tc.Mobility,
		Range:             tc.Range,
		Damage:            tc.Damage,
		IsReady:           false, // cannot act on purchase turn
		HasMoved:          false,
		HasAttacked:       false,
		WasInCombat:       false,
		RemainingMobility: tc.Mobility,
	}, nil
}

// NewStructureFromBalance creates a Structure with stats from balance data.
func NewStructureFromBalance(id string, sType model.StructureType, ownerID string, pos hex.Coord) (*model.Structure, error) {
	if Balance == nil {
		return nil, fmt.Errorf("balance data not loaded")
	}
	sc, ok := Balance.Structures[string(sType)]
	if !ok {
		return nil, fmt.Errorf("unknown structure type: %s", sType)
	}
	return &model.Structure{
		ID:        id,
		Type:      sType,
		OwnerID:   ownerID,
		Hex:       pos,
		CurrentHP: sc.HP,
		MaxHP:     sc.HP,
		ATK:       sc.ATK,
		DEF:       sc.DEF,
		Range:     sc.Range,
		Damage:    sc.Damage,
		Income:    sc.Income,
		CanSpawn:  sc.Spawn,
	}, nil
}

// TroopDamageDice parses and returns the damage dice notation for a troop.
func TroopDamageDice(t *model.Troop) (dice.DiceNotation, error) {
	return dice.ParseDiceNotation(t.Damage)
}

// StructureDamageDice parses and returns the damage dice notation for a structure.
func StructureDamageDice(s *model.Structure) (dice.DiceNotation, error) {
	return dice.ParseDiceNotation(s.Damage)
}

// AntiStructureMultiplier returns the damage multiplier a troop gets vs structures.
func AntiStructureMultiplier(t model.TroopType) int {
	if Balance == nil {
		return 1
	}
	tc, ok := Balance.Troops[string(t)]
	if !ok {
		return 1
	}
	if tc.AntiStructureMultiplier > 0 {
		return tc.AntiStructureMultiplier
	}
	return 1
}

// PassiveIncome returns the base passive income per turn.
func PassiveIncome() int {
	if Balance == nil {
		return 100
	}
	return Balance.Economy.PassiveIncome
}

// StructureIncome returns the income bonus per owned structure.
func StructureIncome() int {
	if Balance == nil {
		return 50
	}
	return Balance.Economy.StructureIncome
}

// StartingCoins returns the starting coin amount.
func StartingCoins() int {
	if Balance == nil {
		return 1000
	}
	return Balance.Economy.StartingCoins
}

// HealingRate returns passive healing HP per turn.
func HealingRate() int {
	if Balance == nil {
		return 2
	}
	return Balance.Healing.PassiveRate
}

// SuddenDeathThreshold returns the turn threshold for sudden death activation.
func SuddenDeathThreshold(size model.MapSize) int {
	if Balance == nil {
		switch size {
		case model.MapSizeSmall:
			return 20
		case model.MapSizeMedium:
			return 30
		case model.MapSizeLarge:
			return 40
		}
		return 30
	}
	v, ok := Balance.SuddenDeath.TurnThresholds[string(size)]
	if !ok {
		return 30
	}
	return v
}

// DominanceTurnsRequired returns how many consecutive turns of structure majority are needed to win.
func DominanceTurnsRequired() int {
	if Balance == nil {
		return 3
	}
	return Balance.WinCond.DominanceTurnsRequired
}
