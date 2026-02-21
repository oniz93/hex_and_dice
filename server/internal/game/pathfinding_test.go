package game

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/teomiscia/hexbattle/internal/hex"
	"github.com/teomiscia/hexbattle/internal/model"
)

func TestReachableHexes(t *testing.T) {
	// Setup map where Marine is at center, surrounded by plains
	gs := NewTestGame().
		WithMapSize(model.MapSizeSmall).
		WithTroop("p1", model.TroopMarine, hex.NewCoord(0, 0, 0), true).
		Build()

	troop := gs.TroopAtHex(hex.NewCoord(0, 0, 0))

	// Marine mobility is 3 on plains.
	// Max distance 3 means it can reach up to ring 3.
	reachable := ReachableHexes(gs, troop)

	// Range 3 on all plains = 1 + 6 + 12 + 18 = 37 hexes, minus the center itself = 36 reachable destinations.
	assert.Len(t, reachable, 36)

	// Check specific hexes
	target1 := hex.NewCoord(1, 0, -1) // Dist 1
	cost, ok := reachable[target1]
	assert.True(t, ok)
	assert.Equal(t, 1, cost)

	target3 := hex.NewCoord(3, 0, -3) // Dist 3
	cost3, ok := reachable[target3]
	assert.True(t, ok)
	assert.Equal(t, 3, cost3)

	target4 := hex.NewCoord(4, 0, -4) // Dist 4
	_, ok4 := reachable[target4]
	assert.False(t, ok4, "Should not reach distance 4")
}

func TestReachableHexes_TerrainCostsAndBlocking(t *testing.T) {
	// Setup: Marine (Mobility 3) at center
	// Forest (Cost 2) at (1, 0, -1)
	// Mountain (Impassable) at (-1, 0, 1)
	// Enemy Troop blocking at (0, 1, -1)
	gs := NewTestGame().
		WithMapSize(model.MapSizeSmall).
		WithTroop("p1", model.TroopMarine, hex.NewCoord(0, 0, 0), true).
		WithTerrain(hex.NewCoord(1, 0, -1), model.TerrainForest).
		WithTerrain(hex.NewCoord(-1, 0, 1), model.TerrainMountains).
		WithTroop("p2", model.TroopMarine, hex.NewCoord(0, 1, -1), true). // enemy block
		WithTroop("p1", model.TroopSniper, hex.NewCoord(0, -1, 1), true). // friendly block (can't stop, can pass)
		Build()

	troop := gs.TroopAtHex(hex.NewCoord(0, 0, 0))
	reachable := ReachableHexes(gs, troop)

	// Forest costs 2, moving to it costs 2. Moving through it to next plains costs 2+1=3 (ok).
	assert.Equal(t, 2, reachable[hex.NewCoord(1, 0, -1)])
	assert.Equal(t, 3, reachable[hex.NewCoord(2, 0, -2)])

	// Mountain is impassable
	_, ok := reachable[hex.NewCoord(-1, 0, 1)]
	assert.False(t, ok)

	// Enemy hex is not reachable
	_, ok = reachable[hex.NewCoord(0, 1, -1)]
	assert.False(t, ok)

	// Friendly hex is passable but not a valid stopping point
	_, ok = reachable[hex.NewCoord(0, -1, 1)]
	assert.False(t, ok, "Cannot stop on friendly troop")
	// But can move THROUGH friendly troop to the next hex (cost 2)
	assert.Equal(t, 2, reachable[hex.NewCoord(0, -2, 2)])
}

func TestCanAttackTarget(t *testing.T) {
	gs := NewTestGame().
		WithTroop("p1", model.TroopSniper, hex.NewCoord(0, 0, 0), true). // Range 3
		Build()

	troop := gs.TroopAtHex(hex.NewCoord(0, 0, 0))

	assert.True(t, CanAttackTarget(troop, hex.NewCoord(1, 0, -1)))  // Dist 1
	assert.True(t, CanAttackTarget(troop, hex.NewCoord(3, 0, -3)))  // Dist 3
	assert.False(t, CanAttackTarget(troop, hex.NewCoord(4, 0, -4))) // Dist 4
	assert.False(t, CanAttackTarget(troop, hex.NewCoord(0, 0, 0)))  // Self (Dist 0) is not attackable
}
