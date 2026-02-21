package game

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/teomiscia/hexbattle/internal/hex"
	"github.com/teomiscia/hexbattle/internal/model"
)

func TestCheckWinConditions_HQDestroyed(t *testing.T) {
	// P1 has HQ, P2 HQ is missing (destroyed)
	gs := NewTestGame().
		WithStructure(model.StructureHQ, "p1", hex.NewCoord(0, -5, 5)).
		Build()

	gameOver := CheckWinConditions(gs)
	assert.NotNil(t, gameOver)
	assert.Equal(t, "p1", gameOver.WinnerID)
	assert.Equal(t, model.WinReasonHQDestroyed, gameOver.Reason)
}

func TestCheckWinConditions_NoWinYet(t *testing.T) {
	// Both have HQs, no dominance, no sudden death tie
	gs := NewTestGame().
		WithStructure(model.StructureHQ, "p1", hex.NewCoord(0, -5, 5)).
		WithStructure(model.StructureHQ, "p2", hex.NewCoord(0, 5, -5)).
		WithStructure(model.StructureOutpost, "p1", hex.NewCoord(1, -1, 0)).
		Build()

	gameOver := CheckWinConditions(gs)
	assert.Nil(t, gameOver, "Game should continue")

	// P1 has 2 structures, P2 has 1. Total = 3.
	// P1 has 2 > 3/2 (1.5). P1 has dominance, but not for 3 turns yet.
	assert.Equal(t, 1, gs.Players[0].DominanceTurnCounter)
}

func TestCheckWinConditions_StructureDominance(t *testing.T) {
	gs := NewTestGame().
		WithStructure(model.StructureHQ, "p1", hex.NewCoord(0, -5, 5)).
		WithStructure(model.StructureHQ, "p2", hex.NewCoord(0, 5, -5)).
		WithStructure(model.StructureOutpost, "p1", hex.NewCoord(1, -1, 0)).
		Build()

	// Simulate 3 turns of dominance for P1
	gs.Players[0].DominanceTurnCounter = 3

	gameOver := CheckWinConditions(gs)
	assert.NotNil(t, gameOver)
	assert.Equal(t, "p1", gameOver.WinnerID)
	assert.Equal(t, model.WinReasonStructureDominance, gameOver.Reason)
}

func TestCheckWinConditions_SuddenDeathTiebreak_Structures(t *testing.T) {
	gs := NewTestGame().
		WithStructure(model.StructureHQ, "p1", hex.NewCoord(0, 0, 0)).
		WithStructure(model.StructureHQ, "p2", hex.NewCoord(1, -1, 0)).
		WithStructure(model.StructureOutpost, "p2", hex.NewCoord(-1, 1, 0)).
		Build()

	// Trigger sudden death tiebreak manually
	gs.SuddenDeathActive = true
	gs.SafeZoneRadius = 1

	gameOver := CheckWinConditions(gs)
	assert.NotNil(t, gameOver)
	// P2 has 2 structures, P1 has 1
	assert.Equal(t, "p2", gameOver.WinnerID)
	assert.Equal(t, model.WinReasonSuddenDeath, gameOver.Reason)
}

func TestCheckWinConditions_SuddenDeathTiebreak_HP(t *testing.T) {
	gs := NewTestGame().
		WithStructure(model.StructureHQ, "p1", hex.NewCoord(0, 0, 0)).
		WithStructure(model.StructureHQ, "p2", hex.NewCoord(1, -1, 0)).
		WithTroop("p1", model.TroopMech, hex.NewCoord(0, 1, -1), true).   // 12 HP
		WithTroop("p2", model.TroopMarine, hex.NewCoord(1, 0, -1), true). // 10 HP
		Build()

	// Both have 1 structure. P1 has more HP (12 vs 10).
	gs.SuddenDeathActive = true
	gs.SafeZoneRadius = 1

	gameOver := CheckWinConditions(gs)
	assert.NotNil(t, gameOver)
	assert.Equal(t, "p1", gameOver.WinnerID)
	assert.Equal(t, model.WinReasonSuddenDeath, gameOver.Reason)
}
