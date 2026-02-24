package game

import (
	"github.com/teomiscia/hexbattle/internal/dice"
	"github.com/teomiscia/hexbattle/internal/hex"
	"github.com/teomiscia/hexbattle/internal/model"
	"github.com/teomiscia/hexbattle/internal/player"
	"github.com/teomiscia/hexbattle/internal/ws"
)

// ActionResult holds the output of processing a player action.
type ActionResult struct {
	Ack        bool             // true if the action was accepted
	Error      *ws.ErrorData    // set if the action was rejected
	Deltas     []interface{}    // delta messages to broadcast
	DeltaTypes []string         // message types corresponding to each delta
	GameOver   *ws.GameOverData // set if the game ended as a result
}

// ExecuteMove processes a move action.
func ExecuteMove(gs *GameState, playerID string, unitID string, target hex.Coord) *ActionResult {
	// Validate
	errCode, errMsg := ValidateMove(gs, playerID, unitID, target)
	if errCode != "" {
		return &ActionResult{
			Ack:   false,
			Error: &ws.ErrorData{Code: errCode, Message: errMsg},
		}
	}

	troop := gs.GetTroop(unitID)
	from := troop.Hex

	// Calculate cost
	cost := CanReach(gs, troop, target)

	// Execute
	troop.Hex = target
	troop.RemainingMobility -= cost
	troop.HasMoved = true

	delta := &ws.TroopMovedData{
		UnitID:            unitID,
		FromQ:             from.Q,
		FromR:             from.R,
		FromS:             from.S,
		ToQ:               target.Q,
		ToR:               target.R,
		ToS:               target.S,
		RemainingMobility: troop.RemainingMobility,
	}

	return &ActionResult{
		Ack:        true,
		Deltas:     []interface{}{delta},
		DeltaTypes: []string{ws.MsgTroopMoved},
	}
}

// ExecuteAttack processes an attack action.
func ExecuteAttack(gs *GameState, roller *dice.Roller, playerID string, unitID string, target hex.Coord) *ActionResult {
	// Validate
	errCode, errMsg := ValidateAttack(gs, playerID, unitID, target)
	if errCode != "" {
		return &ActionResult{
			Ack:   false,
			Error: &ws.ErrorData{Code: errCode, Message: errMsg},
		}
	}

	attacker := gs.GetTroop(unitID)
	attacker.HasAttacked = true
	attacker.HasMoved = true // Attacking ends unit's move phase
	var deltas []interface{}
	var deltaTypes []string

	// Determine target type: troop or structure
	enemyTroop := gs.TroopAtHex(target)
	if enemyTroop != nil {
		// Troop vs Troop combat
		combatResult, destroyed := ResolveTroopCombat(gs, roller, attacker, enemyTroop)
		deltas = append(deltas, combatResult)
		deltaTypes = append(deltaTypes, ws.MsgCombatResult)

		// Track stats
		attackerIdx := gs.PlayerIndex(attacker.OwnerID)
		defenderIdx := gs.PlayerIndex(enemyTroop.OwnerID)

		if combatResult.Hit {
			if attackerIdx >= 0 {
				gs.Stats[attackerIdx].TotalDamageDealt += combatResult.Damage
			}
		}

		for _, d := range destroyed {
			deltas = append(deltas, d)
			deltaTypes = append(deltaTypes, ws.MsgTroopDestroyed)

			// Update stats
			if d.UnitID == enemyTroop.ID {
				if attackerIdx >= 0 {
					gs.Stats[attackerIdx].TroopsKilled++
				}
				if defenderIdx >= 0 {
					gs.Stats[defenderIdx].TroopsLost++
				}
				gs.RemoveTroop(d.UnitID)
			} else if d.UnitID == attacker.ID {
				if defenderIdx >= 0 {
					gs.Stats[defenderIdx].TroopsKilled++
				}
				if attackerIdx >= 0 {
					gs.Stats[attackerIdx].TroopsLost++
				}
				gs.RemoveTroop(d.UnitID)
			}
		}
	} else {
		// Troop vs Structure combat
		structure := gs.StructureAtHex(target)
		if structure != nil {
			structResult, destroyed := ResolveStructureAttack(gs, roller, attacker, structure)
			deltas = append(deltas, structResult)
			deltaTypes = append(deltaTypes, ws.MsgStructureAttacked)

			attackerIdx := gs.PlayerIndex(attacker.OwnerID)
			if structResult.Damage > 0 && attackerIdx >= 0 {
				gs.Stats[attackerIdx].TotalDamageDealt += structResult.Damage
			}

			for _, d := range destroyed {
				deltas = append(deltas, d)
				deltaTypes = append(deltaTypes, ws.MsgTroopDestroyed)
				gs.RemoveTroop(d.UnitID)
			}
		}
	}

	// Check win conditions after attack
	gameOver := CheckWinConditions(gs, false)
	return &ActionResult{
		Ack:        true,
		Deltas:     deltas,
		DeltaTypes: deltaTypes,
		GameOver:   gameOver,
	}
}

// ExecuteBuy processes a buy (purchase troop) action.
func ExecuteBuy(gs *GameState, playerID string, troopType model.TroopType, structureID string) *ActionResult {
	// Validate
	errCode, errMsg := ValidateBuy(gs, playerID, troopType, structureID)
	if errCode != "" {
		return &ActionResult{
			Ack:   false,
			Error: &ws.ErrorData{Code: errCode, Message: errMsg},
		}
	}

	structure := gs.GetStructure(structureID)

	// Find nearest empty hex for spawn
	spawnHex := structure.Hex
	found := false
	for r := 0; r <= 3; r++ {
		for _, h := range structure.Hex.Ring(r) {
			if gs.IsHexPassable(h) && gs.TroopAtHex(h) == nil && gs.StructureAtHex(h) == nil {
				spawnHex = h
				found = true
				break
			}
		}
		if found {
			break
		}
	}

	// Deduct cost
	DeductCost(gs, playerID, troopType)

	// Create troop
	unitID := player.GenerateUnitID()
	troop, err := NewTroopFromBalance(unitID, troopType, playerID, spawnHex)
	if err != nil {
		return &ActionResult{
			Ack:   false,
			Error: &ws.ErrorData{Code: model.ErrInvalidMessage, Message: err.Error()},
		}
	}

	// Troop cannot act on purchase turn
	troop.IsReady = false

	gs.AddTroop(troop)

	idx := gs.PlayerIndex(playerID)
	coinsRemaining := 0
	if idx >= 0 {
		coinsRemaining = gs.Players[idx].Coins
	}

	delta := &ws.TroopPurchasedData{
		UnitID:         unitID,
		UnitType:       troopType,
		HexQ:           spawnHex.Q,
		HexR:           spawnHex.R,
		HexS:           spawnHex.S,
		Owner:          playerID,
		CoinsRemaining: coinsRemaining,
	}

	return &ActionResult{
		Ack:        true,
		Deltas:     []interface{}{delta},
		DeltaTypes: []string{ws.MsgTroopPurchased},
	}
}

// ExecuteEndTurn processes an end_turn action and transitions to the next turn.
func ExecuteEndTurn(gs *GameState, roller *dice.Roller, playerID string) *ActionResult {
	// Validate
	errCode, errMsg := ValidateEndTurn(gs, playerID)
	if errCode != "" {
		return &ActionResult{
			Ack:   false,
			Error: &ws.ErrorData{Code: errCode, Message: errMsg},
		}
	}

	// Transition
	gs.Phase = model.PhaseTurnTransition

	// Check win conditions before switching (end of turn for current player)
	gameOver := CheckWinConditions(gs, true)
	if gameOver != nil {
		return &ActionResult{
			Ack:      true,
			GameOver: gameOver,
		}
	}

	// Switch active player
	gs.SwitchActivePlayer()
	gs.TurnNumber++

	// Run turn start pipeline
	turnStartData := RunTurnStart(gs, roller)

	// Check win conditions after turn start (sudden death may kill things)
	gameOver = CheckWinConditions(gs, false)

	return &ActionResult{
		Ack:        true,
		Deltas:     []interface{}{turnStartData},
		DeltaTypes: []string{ws.MsgTurnStart},
		GameOver:   gameOver,
	}
}
