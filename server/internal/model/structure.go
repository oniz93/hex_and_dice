package model

import (
	"github.com/teomiscia/hexbattle/internal/hex"
)

// Structure represents a capturable building on the map.
type Structure struct {
	ID        string        `json:"id"`
	Type      StructureType `json:"type"`
	OwnerID   string        `json:"owner_id"` // Empty string = neutral
	Hex       hex.Coord     `json:"hex"`
	CurrentHP int           `json:"current_hp"`
	MaxHP     int           `json:"max_hp"`
	ATK       int           `json:"atk"`
	DEF       int           `json:"def"`
	Range     int           `json:"range"`
	Damage    string        `json:"damage"` // Dice notation
	Income    int           `json:"income"`
	CanSpawn  bool          `json:"can_spawn"`
}

// IsNeutral returns true if the structure has no owner.
func (s *Structure) IsNeutral() bool {
	return s.OwnerID == ""
}

// IsOwnedBy returns true if the structure is owned by the given player.
func (s *Structure) IsOwnedBy(playerID string) bool {
	return s.OwnerID == playerID
}

// IsAlive returns true if the structure has HP remaining.
func (s *Structure) IsAlive() bool {
	return s.CurrentHP > 0
}

// TakeDamage applies damage to the structure and returns true if it reached 0 HP.
func (s *Structure) TakeDamage(amount int) bool {
	s.CurrentHP -= amount
	if s.CurrentHP < 0 {
		s.CurrentHP = 0
	}
	return s.CurrentHP <= 0
}

// Capture transfers ownership to the given player and sets HP to max.
func (s *Structure) Capture(newOwnerID string) {
	s.OwnerID = newOwnerID
	s.CurrentHP = s.MaxHP
}

// Heal restores HP up to MaxHP. Returns the amount actually healed.
func (s *Structure) Heal(amount int) int {
	before := s.CurrentHP
	s.CurrentHP += amount
	if s.CurrentHP > s.MaxHP {
		s.CurrentHP = s.MaxHP
	}
	return s.CurrentHP - before
}
