package dice

import (
	"fmt"
	"math/rand"
	"regexp"
	"strconv"
)

// Roller provides seeded dice rolling for a single game instance.
// Each game gets its own Roller with a unique seed for reproducibility.
type Roller struct {
	rng *rand.Rand
}

// NewRoller creates a new Roller with the given seed.
func NewRoller(seed int64) *Roller {
	return &Roller{
		rng: rand.New(rand.NewSource(seed)),
	}
}

// Roll returns a random integer in [1, sides].
func (r *Roller) Roll(sides int) int {
	if sides < 1 {
		return 1
	}
	return r.rng.Intn(sides) + 1
}

// D20 rolls a single 20-sided die.
func (r *Roller) D20() int {
	return r.Roll(20)
}

// D6 rolls a single 6-sided die.
func (r *Roller) D6() int {
	return r.Roll(6)
}

// D8 rolls a single 8-sided die.
func (r *Roller) D8() int {
	return r.Roll(8)
}

// D4 rolls a single 4-sided die.
func (r *Roller) D4() int {
	return r.Roll(4)
}

// diceNotationRe matches strings like "2D6+2", "1D8", "1D4+1".
var diceNotationRe = regexp.MustCompile(`^(\d+)[Dd](\d+)(?:\+(\d+))?$`)

// DiceNotation represents a parsed dice expression like "2D6+2".
type DiceNotation struct {
	Count    int // number of dice
	Sides    int // sides per die
	Modifier int // flat bonus added after rolling
}

// ParseDiceNotation parses a string like "2D6+2" into a DiceNotation.
func ParseDiceNotation(s string) (DiceNotation, error) {
	matches := diceNotationRe.FindStringSubmatch(s)
	if matches == nil {
		return DiceNotation{}, fmt.Errorf("dice: invalid notation %q", s)
	}

	count, _ := strconv.Atoi(matches[1])
	sides, _ := strconv.Atoi(matches[2])
	modifier := 0
	if matches[3] != "" {
		modifier, _ = strconv.Atoi(matches[3])
	}

	return DiceNotation{Count: count, Sides: sides, Modifier: modifier}, nil
}

// String returns the notation in standard form, e.g. "2D6+2".
func (dn DiceNotation) String() string {
	if dn.Modifier > 0 {
		return fmt.Sprintf("%dD%d+%d", dn.Count, dn.Sides, dn.Modifier)
	}
	return fmt.Sprintf("%dD%d", dn.Count, dn.Sides)
}

// StepDown reduces the die size by one step: D8->D6, D6->D4, D4->D4 (minimum).
func (dn DiceNotation) StepDown() DiceNotation {
	newSides := dn.Sides
	switch dn.Sides {
	case 8:
		newSides = 6
	case 6:
		newSides = 4
		// D4 is the minimum; no further reduction
	}
	return DiceNotation{Count: dn.Count, Sides: newSides, Modifier: dn.Modifier}
}

// RollDamage rolls damage using the given notation and returns the total.
func (r *Roller) RollDamage(dn DiceNotation) int {
	total := 0
	for i := 0; i < dn.Count; i++ {
		total += r.Roll(dn.Sides)
	}
	total += dn.Modifier
	return total
}

// RollHalfDamage rolls damage and halves it (rounded down, minimum 1).
// Used for counterattack damage.
func (r *Roller) RollHalfDamage(dn DiceNotation) int {
	total := r.RollDamage(dn)
	half := total / 2
	if half < 1 {
		half = 1
	}
	return half
}
