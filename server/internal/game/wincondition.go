package game

import (
	"github.com/teomiscia/hexbattle/internal/model"
	"github.com/teomiscia/hexbattle/internal/ws"
)

// CheckWinConditions evaluates all win conditions and returns a GameOverData if someone won.
// Returns nil if the game continues.
func CheckWinConditions(gs *GameState) *ws.GameOverData {
	// 1. HQ Destruction
	for i := 0; i < 2; i++ {
		hq := gs.PlayerHQ(gs.Players[i].ID)
		if hq == nil || !hq.IsOwnedBy(gs.Players[i].ID) {
			// Player lost their HQ (either destroyed/captured)
			winnerIdx := 1 - i
			return buildGameOver(gs, gs.Players[winnerIdx].ID, model.WinReasonHQDestroyed)
		}
	}

	// 2. Structure Dominance (checked per full round)
	// Dominance is checked after both players have had a turn
	for i := 0; i < 2; i++ {
		playerID := gs.Players[i].ID
		owned := gs.StructureCountOwnedBy(playerID)
		total := gs.TotalStructureCount()

		if total > 0 && owned > total/2 {
			gs.Players[i].DominanceTurnCounter++
		} else {
			gs.Players[i].DominanceTurnCounter = 0
		}

		if gs.Players[i].DominanceTurnCounter >= DominanceTurnsRequired() {
			return buildGameOver(gs, playerID, model.WinReasonStructureDominance)
		}
	}

	// 3. Sudden Death Tiebreak — zone reduced to minimum
	if gs.SuddenDeathActive && gs.SafeZoneRadius <= 1 {
		return resolveSuddenDeathTiebreak(gs)
	}

	return nil
}

// resolveSuddenDeathTiebreak determines the winner when the zone is at minimum.
func resolveSuddenDeathTiebreak(gs *GameState) *ws.GameOverData {
	p1Structs := gs.StructureCountOwnedBy(gs.Players[0].ID)
	p2Structs := gs.StructureCountOwnedBy(gs.Players[1].ID)

	// More structures wins
	if p1Structs > p2Structs {
		return buildGameOver(gs, gs.Players[0].ID, model.WinReasonSuddenDeath)
	}
	if p2Structs > p1Structs {
		return buildGameOver(gs, gs.Players[1].ID, model.WinReasonSuddenDeath)
	}

	// Tied structures — compare total troop HP
	p1HP := totalTroopHP(gs, gs.Players[0].ID)
	p2HP := totalTroopHP(gs, gs.Players[1].ID)

	if p1HP > p2HP {
		return buildGameOver(gs, gs.Players[0].ID, model.WinReasonSuddenDeath)
	}
	if p2HP > p1HP {
		return buildGameOver(gs, gs.Players[1].ID, model.WinReasonSuddenDeath)
	}

	// Still tied — draw
	return buildGameOver(gs, "", model.WinReasonDraw)
}

// totalTroopHP sums the current HP of all living troops belonging to a player.
func totalTroopHP(gs *GameState, playerID string) int {
	total := 0
	for _, t := range gs.Troops {
		if t.OwnerID == playerID && t.IsAlive() {
			total += t.CurrentHP
		}
	}
	return total
}

// buildGameOver constructs a GameOverData with stats for both players.
func buildGameOver(gs *GameState, winnerID string, reason model.WinReason) *ws.GameOverData {
	gs.Phase = model.PhaseGameOver

	// Finalize stats
	for i := 0; i < 2; i++ {
		gs.Stats[i].TurnsPlayed = gs.TurnNumber
		gs.Stats[i].StructuresHeld = gs.StructureCountOwnedBy(gs.Players[i].ID)
	}

	stats := map[string]model.GameOverStats{
		gs.Players[0].ID: gs.Stats[0],
		gs.Players[1].ID: gs.Stats[1],
	}

	return &ws.GameOverData{
		WinnerID: winnerID,
		Reason:   reason,
		Stats:    stats,
	}
}

// CheckForfeit handles a player forfeiting (disconnect timeout or surrender).
func CheckForfeit(gs *GameState, loserID string) *ws.GameOverData {
	winnerIdx := 1 - gs.PlayerIndex(loserID)
	if winnerIdx < 0 || winnerIdx > 1 {
		winnerIdx = 0
	}
	return buildGameOver(gs, gs.Players[winnerIdx].ID, model.WinReasonForfeit)
}

// CheckDisconnectForfeit handles forfeit due to reconnect timeout.
func CheckDisconnectForfeit(gs *GameState, disconnectedID string) *ws.GameOverData {
	winnerIdx := 1 - gs.PlayerIndex(disconnectedID)
	if winnerIdx < 0 || winnerIdx > 1 {
		winnerIdx = 0
	}
	return buildGameOver(gs, gs.Players[winnerIdx].ID, model.WinReasonDisconnect)
}
