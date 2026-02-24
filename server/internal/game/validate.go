package game

import (
	"github.com/teomiscia/hexbattle/internal/hex"
	"github.com/teomiscia/hexbattle/internal/model"
)

// ValidateMove checks if a move action is legal.
// Returns (errorCode, errorMessage) — empty strings if valid.
func ValidateMove(gs *GameState, playerID string, unitID string, target hex.Coord) (model.ErrorCode, string) {
	// Check it's this player's turn
	if !gs.IsActivePlayer(playerID) {
		return model.ErrNotYourTurn, "it is not your turn"
	}

	// Check game phase
	if gs.Phase != model.PhasePlayerAction {
		return model.ErrNotYourTurn, "not in player action phase"
	}

	// Check unit exists
	troop := gs.GetTroop(unitID)
	if troop == nil {
		return model.ErrUnitNotFound, "unit not found"
	}

	// Check unit belongs to player
	if troop.OwnerID != playerID {
		return model.ErrUnitNotFound, "unit does not belong to you"
	}

	// Check unit is ready
	if !troop.IsReady {
		return model.ErrUnitNotReady, "unit was purchased this turn and cannot act"
	}

	// Check unit hasn't already moved
	if troop.HasMoved {
		return model.ErrUnitAlreadyActed, "unit has already moved this turn"
	}

	// Check target is in bounds
	if !gs.Grid.Contains(target) {
		return model.ErrInvalidMove, "target hex is out of bounds"
	}

	// Check target is passable
	if !gs.IsHexPassable(target) {
		return model.ErrInvalidMove, "target hex is impassable terrain"
	}

	// Check target is not occupied by enemy
	if gs.IsHexOccupiedByEnemy(target, playerID) {
		return model.ErrInvalidMove, "target hex is occupied by an enemy unit"
	}

	// Check target is not occupied by friendly troop
	if occupant := gs.TroopAtHex(target); occupant != nil {
		return model.ErrInvalidMove, "target hex is occupied by a unit"
	}

	// Check target is not occupied by a structure
	if occupant := gs.StructureAtHex(target); occupant != nil {
		return model.ErrInvalidMove, "target hex is occupied by a structure"
	}

	// Check BFS reachability within mobility budget
	cost := CanReach(gs, troop, target)
	if cost < 0 {
		return model.ErrInvalidMove, "target hex is not reachable within mobility range"
	}

	return "", ""
}

// ValidateAttack checks if an attack action is legal.
func ValidateAttack(gs *GameState, playerID string, unitID string, target hex.Coord) (model.ErrorCode, string) {
	// Check it's this player's turn
	if !gs.IsActivePlayer(playerID) {
		return model.ErrNotYourTurn, "it is not your turn"
	}

	// Check game phase
	if gs.Phase != model.PhasePlayerAction {
		return model.ErrNotYourTurn, "not in player action phase"
	}

	// First turn restriction: player 0 cannot attack on turn 1
	if gs.FirstTurnRestriction && gs.TurnNumber == 1 && gs.ActivePlayer == 0 {
		return model.ErrInvalidAttack, "cannot attack on the first turn"
	}

	// Check unit exists
	troop := gs.GetTroop(unitID)
	if troop == nil {
		return model.ErrUnitNotFound, "unit not found"
	}

	// Check unit belongs to player
	if troop.OwnerID != playerID {
		return model.ErrUnitNotFound, "unit does not belong to you"
	}

	// Check unit is ready
	if !troop.IsReady {
		return model.ErrUnitNotReady, "unit was purchased this turn and cannot act"
	}

	// Check unit hasn't already attacked
	if troop.HasAttacked {
		return model.ErrUnitAlreadyActed, "unit has already attacked this turn"
	}

	// Check target is in attack range
	if !CanAttackTarget(troop, target) {
		return model.ErrInvalidAttack, "target is out of attack range"
	}

	// Check there is a valid target at that hex
	enemyTroop := gs.TroopAtHex(target)
	enemyStructure := gs.StructureAtHex(target)

	if enemyTroop == nil && enemyStructure == nil {
		return model.ErrInvalidAttack, "no target at hex"
	}

	// If attacking a troop, it must be an enemy
	if enemyTroop != nil && enemyTroop.OwnerID == playerID {
		return model.ErrInvalidAttack, "cannot attack your own unit"
	}

	// If attacking a structure, it must be enemy or neutral
	if enemyTroop == nil && enemyStructure != nil {
		if enemyStructure.OwnerID == playerID {
			return model.ErrInvalidAttack, "cannot attack your own structure"
		}
	}

	return "", ""
}

// ValidateBuy checks if a buy action is legal.
func ValidateBuy(gs *GameState, playerID string, troopType model.TroopType, structureID string) (model.ErrorCode, string) {
	// Check it's this player's turn
	if !gs.IsActivePlayer(playerID) {
		return model.ErrNotYourTurn, "it is not your turn"
	}

	// Check game phase
	if gs.Phase != model.PhasePlayerAction {
		return model.ErrNotYourTurn, "not in player action phase"
	}

	return ValidatePurchase(gs, playerID, troopType, structureID)
}

// ValidateEndTurn checks if ending the turn is legal.
func ValidateEndTurn(gs *GameState, playerID string) (model.ErrorCode, string) {
	if !gs.IsActivePlayer(playerID) {
		return model.ErrNotYourTurn, "it is not your turn"
	}

	if gs.Phase != model.PhasePlayerAction {
		return model.ErrNotYourTurn, "not in player action phase"
	}

	return "", ""
}
