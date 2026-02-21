package game

import (
	"github.com/teomiscia/hexbattle/internal/hex"
	"github.com/teomiscia/hexbattle/internal/model"
	"github.com/teomiscia/hexbattle/internal/ws"
)

// CheckSuddenDeath checks if sudden death should activate or progress.
// Returns true if sudden death is active (either just activated or already was).
func CheckSuddenDeath(gs *GameState) bool {
	threshold := SuddenDeathThreshold(gs.MapSize)

	if gs.TurnNumber > threshold {
		if !gs.SuddenDeathActive {
			gs.SuddenDeathActive = true
			gs.SuddenDeathTurn = 0
		}
		gs.SuddenDeathTurn = gs.TurnNumber - threshold
		return true
	}
	return gs.SuddenDeathActive
}

// ShrinkSafeZone reduces the safe zone radius by the configured shrink rate.
// Returns the new safe zone radius.
func ShrinkSafeZone(gs *GameState) int {
	shrinkRate := 1
	if Balance != nil {
		shrinkRate = Balance.SuddenDeath.ShrinkRate
	}

	gs.SafeZoneRadius -= shrinkRate
	if gs.SafeZoneRadius < 1 {
		gs.SafeZoneRadius = 1
	}
	return gs.SafeZoneRadius
}

// RelocateHQs moves player HQs inside the safe zone if they are outside.
// Returns a list of relocation events (for broadcasting).
func RelocateHQs(gs *GameState) []HQRelocation {
	var relocations []HQRelocation

	for i := 0; i < 2; i++ {
		hq := gs.PlayerHQ(gs.Players[i].ID)
		if hq == nil {
			continue
		}

		if hq.Hex.DistanceToOrigin() > gs.SafeZoneRadius {
			oldHex := hq.Hex
			newHex := closestPassableHexInZone(gs, hq.Hex, gs.SafeZoneRadius)
			hq.Hex = newHex
			relocations = append(relocations, HQRelocation{
				PlayerID: gs.Players[i].ID,
				FromHex:  oldHex,
				ToHex:    newHex,
			})
		}
	}

	return relocations
}

// HQRelocation records an HQ relocation event.
type HQRelocation struct {
	PlayerID string
	FromHex  hex.Coord
	ToHex    hex.Coord
}

// ApplyStormDamage applies escalating damage to all troops outside the safe zone.
// Returns damage records for each affected troop.
func ApplyStormDamage(gs *GameState) []ws.SuddenDeathDamage {
	var damages []ws.SuddenDeathDamage
	damage := gs.SuddenDeathTurn // escalating: 1, 2, 3, ...

	// Collect troop IDs first to avoid modifying map during iteration
	var troopIDs []string
	for id := range gs.Troops {
		troopIDs = append(troopIDs, id)
	}

	for _, id := range troopIDs {
		troop := gs.Troops[id]
		if troop == nil || !troop.IsAlive() {
			continue
		}

		if troop.Hex.DistanceToOrigin() > gs.SafeZoneRadius {
			killed := troop.TakeDamage(damage)
			damages = append(damages, ws.SuddenDeathDamage{
				UnitID:  troop.ID,
				Damage:  damage,
				HPAfter: troop.CurrentHP,
				Killed:  killed,
			})

			if killed {
				gs.RemoveTroop(troop.ID)
			}
		}
	}

	return damages
}

// closestPassableHexInZone finds the closest passable hex within the safe zone
// to the given position.
func closestPassableHexInZone(gs *GameState, from hex.Coord, radius int) hex.Coord {
	bestDist := -1
	bestHex := hex.Origin()

	// Search all hexes within the safe zone radius
	for _, c := range hex.Origin().Spiral(radius) {
		if !gs.Grid.Contains(c) {
			continue
		}
		if !model.IsPassable(gs.GetTerrainAt(c)) {
			continue
		}
		// Don't place on an occupied hex (by troop or another structure)
		if gs.TroopAtHex(c) != nil {
			continue
		}
		// Check for other structures at this hex
		otherStruct := gs.StructureAtHex(c)
		if otherStruct != nil {
			continue
		}

		dist := from.Distance(c)
		if bestDist < 0 || dist < bestDist {
			bestDist = dist
			bestHex = c
		}
	}

	return bestHex
}

// RunSuddenDeathPhase executes the full sudden death phase for a turn.
// Returns storm damage records and HQ relocations.
func RunSuddenDeathPhase(gs *GameState) ([]ws.SuddenDeathDamage, []HQRelocation) {
	if !CheckSuddenDeath(gs) {
		return nil, nil
	}

	// Shrink the zone
	ShrinkSafeZone(gs)

	// Relocate HQs
	relocations := RelocateHQs(gs)

	// Apply storm damage
	damages := ApplyStormDamage(gs)

	return damages, relocations
}
