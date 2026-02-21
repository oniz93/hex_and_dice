package mapgen

import (
	"github.com/teomiscia/hexbattle/internal/hex"
	"github.com/teomiscia/hexbattle/internal/model"
)

// ApplySymmetry takes terrain assignments for one half of the grid and
// mirrors them via 180° rotational symmetry.
// The center hex (0,0,0) is always set to Plains.
func ApplySymmetry(terrain map[hex.Coord]model.TerrainType, grid *hex.Grid) {
	// Ensure center is plains
	terrain[hex.Origin()] = model.TerrainPlains

	// For each hex in one half, set the rotated counterpart to the same terrain
	half := grid.HalfGrid()
	for _, c := range half {
		t, ok := terrain[c]
		if !ok {
			continue
		}
		rotated := c.Rotate180()
		if grid.Contains(rotated) {
			terrain[rotated] = t
		}
	}
}

// ApplyStructureSymmetry ensures neutral structures are placed in symmetric pairs.
// Given a list of structure positions for one half, returns the full set including
// their 180° rotated counterparts.
func ApplyStructureSymmetry(positions []hex.Coord) []hex.Coord {
	result := make([]hex.Coord, 0, len(positions)*2)
	for _, pos := range positions {
		result = append(result, pos)
		rotated := pos.Rotate180()
		if rotated != pos { // avoid duplicating center hex
			result = append(result, rotated)
		}
	}
	return result
}
