package model

import (
	"github.com/teomiscia/hexbattle/internal/hex"
)

// Troop represents a single unit on the battlefield.
type Troop struct {
	ID                string    `json:"id"`
	Type              TroopType `json:"type"`
	OwnerID           string    `json:"owner_id"`
	Hex               hex.Coord `json:"hex"`
	CurrentHP         int       `json:"current_hp"`
	MaxHP             int       `json:"max_hp"`
	ATK               int       `json:"atk"`
	DEF               int       `json:"def"`
	Mobility          int       `json:"mobility"`
	Range             int       `json:"range"`
	Damage            string    `json:"damage"` // Dice notation, e.g. "1D6+1"
	IsReady           bool      `json:"is_ready"`
	HasMoved          bool      `json:"has_moved"`
	HasAttacked       bool      `json:"has_attacked"`
	WasInCombat       bool      `json:"was_in_combat"`
	RemainingMobility int       `json:"remaining_mobility"`
}

// IsAlive returns true if the troop has HP remaining.
func (t *Troop) IsAlive() bool {
	return t.CurrentHP > 0
}

// CanAct returns true if the troop is ready and alive.
func (t *Troop) CanAct() bool {
	return t.IsAlive() && t.IsReady
}

// CanMove returns true if the troop can still move this turn.
func (t *Troop) CanMove() bool {
	return t.CanAct() && !t.HasMoved && t.RemainingMobility > 0
}

// CanAttack returns true if the troop can still attack this turn.
func (t *Troop) CanAttack() bool {
	return t.CanAct() && !t.HasAttacked
}

// ResetForTurn resets per-turn action flags and restores mobility.
func (t *Troop) ResetForTurn() {
	t.HasMoved = false
	t.HasAttacked = false
	t.WasInCombat = false
	t.RemainingMobility = t.Mobility
}

// TakeDamage applies damage to the troop and returns true if it died.
func (t *Troop) TakeDamage(amount int) bool {
	t.CurrentHP -= amount
	if t.CurrentHP < 0 {
		t.CurrentHP = 0
	}
	return t.CurrentHP <= 0
}

// Heal restores HP up to MaxHP. Returns the amount actually healed.
func (t *Troop) Heal(amount int) int {
	before := t.CurrentHP
	t.CurrentHP += amount
	if t.CurrentHP > t.MaxHP {
		t.CurrentHP = t.MaxHP
	}
	return t.CurrentHP - before
}

// IsMelee returns true if the troop has a range of 1.
func (t *Troop) IsMelee() bool {
	return t.Range == 1
}
