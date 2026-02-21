package game

import (
	"github.com/teomiscia/hexbattle/internal/model"
)

// CalculateIncome computes the total income for the active player this turn.
// Returns (passiveIncome, structureIncome, totalIncome).
func CalculateIncome(gs *GameState, playerID string) (int, int, int) {
	passive := PassiveIncome()
	structIncome := 0

	for _, s := range gs.Structures {
		if s.OwnerID == playerID {
			if s.Type == model.StructureOutpost || s.Type == model.StructureCommandCenter {
				structIncome += StructureIncome()
			}
		}
	}

	return passive, structIncome, passive + structIncome
}

// CreditIncome adds the calculated income to the player's coin balance.
// Returns the total income credited.
func CreditIncome(gs *GameState, playerID string) int {
	_, _, total := CalculateIncome(gs, playerID)
	idx := gs.PlayerIndex(playerID)
	if idx >= 0 {
		gs.Players[idx].Coins += total
	}
	return total
}

// ValidatePurchase checks if a player can buy a troop at a structure.
func ValidatePurchase(gs *GameState, playerID string, troopType model.TroopType, structureID string) (model.ErrorCode, string) {
	cost := TroopCost(troopType)
	if cost == 0 {
		return model.ErrInvalidMessage, "unknown troop type"
	}

	idx := gs.PlayerIndex(playerID)
	if idx < 0 {
		return model.ErrInvalidMessage, "player not found"
	}

	if gs.Players[idx].Coins < cost {
		return model.ErrInsufficientFunds, "not enough coins"
	}

	structure := gs.GetStructure(structureID)
	if structure == nil {
		return model.ErrInvalidMessage, "structure not found"
	}

	if !structure.IsOwnedBy(playerID) {
		return model.ErrSpawnNotOwned, "structure not owned by player"
	}

	if !structure.CanSpawn {
		return model.ErrSpawnNotOwned, "structure cannot spawn troops"
	}

	if gs.TroopAtHex(structure.Hex) != nil {
		return model.ErrSpawnOccupied, "spawn hex is occupied"
	}

	return "", ""
}

// DeductCost subtracts the troop cost from the player's coins.
func DeductCost(gs *GameState, playerID string, troopType model.TroopType) {
	cost := TroopCost(troopType)
	idx := gs.PlayerIndex(playerID)
	if idx >= 0 {
		gs.Players[idx].Coins -= cost
	}
}
