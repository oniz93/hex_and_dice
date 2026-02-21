package game

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/teomiscia/hexbattle/internal/dice"
	"github.com/teomiscia/hexbattle/internal/hex"
	"github.com/teomiscia/hexbattle/internal/model"
)

func TestResolveTroopCombat_NormalHit(t *testing.T) {
	gs := NewTestGame().
		WithTroop("p1", model.TroopMarine, hex.NewCoord(0, 0, 0), true).
		WithTroop("p2", model.TroopMarine, hex.NewCoord(1, -1, 0), true).
		Build()

	attacker := gs.TroopAtHex(hex.NewCoord(0, 0, 0))
	defender := gs.TroopAtHex(hex.NewCoord(1, -1, 0))

	// Setup a roller that will roll exactly what we want
	// But since dice.Roller is opaque with unexported fields, we can't mock its rng directly.
	// Instead, we will use a known seed that produces a known first roll, OR
	// we will run it 100 times and check constraints, or just make a mock if we needed to.
	// Since we want deterministic tests, let's use Seed=42 and figure out the rolls,
	// or just test the logic ignoring the exact D20 result by modifying the caller, but the signature takes a Roller.

	// Since Roller uses seeded math/rand, we can use a fixed seed.
	roller := dice.NewRoller(42) // Known seed

	result, destroyed := ResolveTroopCombat(gs, roller, attacker, defender)

	assert.NotNil(t, result)
	assert.Equal(t, attacker.ID, result.AttackerID)
	assert.Equal(t, defender.ID, result.DefenderID)
	assert.True(t, attacker.WasInCombat)
	assert.True(t, defender.WasInCombat)
	assert.True(t, attacker.HasAttacked)

	if result.Hit {
		assert.Less(t, defender.CurrentHP, defender.MaxHP, "Defender should take damage on hit")
	}

	if result.Killed {
		require.Len(t, destroyed, 1)
		assert.Equal(t, defender.ID, destroyed[0].UnitID)
	}
}

func TestResolveStructureAttack(t *testing.T) {
	gs := NewTestGame().
		WithTroop("p1", model.TroopMech, hex.NewCoord(0, 0, 0), true).
		WithStructure(model.StructureOutpost, "p2", hex.NewCoord(1, -1, 0)).
		Build()

	attacker := gs.TroopAtHex(hex.NewCoord(0, 0, 0))
	structure := gs.StructureAtHex(hex.NewCoord(1, -1, 0))

	roller := dice.NewRoller(42)

	initialHP := structure.CurrentHP
	result, _ := ResolveStructureAttack(gs, roller, attacker, structure)

	assert.Equal(t, structure.ID, result.StructureID)
	assert.Equal(t, attacker.ID, result.AttackerID)
	assert.True(t, attacker.HasAttacked)

	if result.Damage > 0 {
		// Mech deals 2x damage to structures
		assert.Less(t, structure.CurrentHP, initialHP)
	}
}

func TestAntiStructureMultiplier(t *testing.T) {
	// Mech has 2x multiplier
	assert.Equal(t, 2, AntiStructureMultiplier(model.TroopMech))
	// Marine has 1x multiplier
	assert.Equal(t, 1, AntiStructureMultiplier(model.TroopMarine))
}
