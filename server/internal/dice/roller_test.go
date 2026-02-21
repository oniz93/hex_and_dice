package dice

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRoller_Determinism(t *testing.T) {
	// With the same seed, two rollers should produce the same sequence
	r1 := NewRoller(42)
	r2 := NewRoller(42)

	for i := 0; i < 100; i++ {
		assert.Equal(t, r1.Roll(20), r2.Roll(20), "Rolls should be deterministic")
	}
}

func TestRoller_Rolls(t *testing.T) {
	r := NewRoller(12345)

	t.Run("D20 rolls are within [1, 20]", func(t *testing.T) {
		for i := 0; i < 100; i++ {
			roll := r.D20()
			assert.GreaterOrEqual(t, roll, 1)
			assert.LessOrEqual(t, roll, 20)
		}
	})

	t.Run("D8 rolls are within [1, 8]", func(t *testing.T) {
		for i := 0; i < 100; i++ {
			roll := r.D8()
			assert.GreaterOrEqual(t, roll, 1)
			assert.LessOrEqual(t, roll, 8)
		}
	})

	t.Run("Invalid sides returns 1", func(t *testing.T) {
		assert.Equal(t, 1, r.Roll(-5))
		assert.Equal(t, 1, r.Roll(0))
	})
}

func TestParseDiceNotation(t *testing.T) {
	tests := []struct {
		notation string
		expected DiceNotation
		err      bool
	}{
		{"1D6", DiceNotation{Count: 1, Sides: 6, Modifier: 0}, false},
		{"2D6+2", DiceNotation{Count: 2, Sides: 6, Modifier: 2}, false},
		{"1D8+1", DiceNotation{Count: 1, Sides: 8, Modifier: 1}, false},
		{"1D4", DiceNotation{Count: 1, Sides: 4, Modifier: 0}, false},
		{"10D20+100", DiceNotation{Count: 10, Sides: 20, Modifier: 100}, false},

		// Invalid cases
		{"D6", DiceNotation{}, true},
		{"1D", DiceNotation{}, true},
		{"1D6+", DiceNotation{}, true},
		{"abc", DiceNotation{}, true},
	}

	for _, tt := range tests {
		t.Run(tt.notation, func(t *testing.T) {
			result, err := ParseDiceNotation(tt.notation)
			if tt.err {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected, result)
				assert.Equal(t, tt.notation, result.String())
			}
		})
	}
}

func TestDiceNotation_StepDown(t *testing.T) {
	tests := []struct {
		input    DiceNotation
		expected DiceNotation
	}{
		{DiceNotation{Count: 1, Sides: 8, Modifier: 1}, DiceNotation{Count: 1, Sides: 6, Modifier: 1}},
		{DiceNotation{Count: 2, Sides: 6, Modifier: 2}, DiceNotation{Count: 2, Sides: 4, Modifier: 2}},
		{DiceNotation{Count: 1, Sides: 4, Modifier: 0}, DiceNotation{Count: 1, Sides: 4, Modifier: 0}}, // Minimum
	}

	for _, tt := range tests {
		t.Run(tt.input.String(), func(t *testing.T) {
			assert.Equal(t, tt.expected, tt.input.StepDown())
		})
	}
}

func TestRoller_Damage(t *testing.T) {
	r := NewRoller(42) // Fixed seed for deterministic behavior

	dn := DiceNotation{Count: 2, Sides: 6, Modifier: 2} // 2D6+2

	// Roll 1: deterministic sequence for seed 42
	// For seed 42, first two rolls (using rand.Intn)
	// Just verify bounds
	for i := 0; i < 50; i++ {
		total := r.RollDamage(dn)
		assert.GreaterOrEqual(t, total, 4) // Minimum: 2*1 + 2 = 4
		assert.LessOrEqual(t, total, 14)   // Maximum: 2*6 + 2 = 14
	}

	// Roll Half Damage
	for i := 0; i < 50; i++ {
		half := r.RollHalfDamage(dn)
		assert.GreaterOrEqual(t, half, 2) // Minimum: 4 / 2 = 2
		assert.LessOrEqual(t, half, 7)    // Maximum: 14 / 2 = 7
	}
}
