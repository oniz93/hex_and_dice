package mapgen

import (
	"math/rand"
	"sort"

	"github.com/teomiscia/hexbattle/internal/hex"
	"github.com/teomiscia/hexbattle/internal/model"
)

// PlaceHQs places two HQs at opposite poles of the hex grid (maximizing distance).
// Returns (hq1Pos, hq2Pos).
func PlaceHQs(grid *hex.Grid) (hex.Coord, hex.Coord) {
	radius := grid.Radius
	// Place along the Q axis: one at north pole, one at south pole
	hq1 := hex.NewCoord(0, -radius, radius) // top
	hq2 := hex.NewCoord(0, radius, -radius) // bottom (180° rotation of hq1)
	return hq1, hq2
}

// PlaceNeutralStructures places neutral structures on the map with even distribution.
// count is the total number of neutral structures to place.
// Returns the positions for the structures.
func PlaceNeutralStructures(grid *hex.Grid, terrain map[hex.Coord]model.TerrainType, hq1, hq2 hex.Coord, count int, rng *rand.Rand) []hex.Coord {
	radius := grid.Radius

	// Gather candidate hexes
	candidates := gatherCandidates(grid, terrain, hq1, hq2)
	if len(candidates) == 0 {
		return nil
	}

	// Categorize candidates into zones
	centerRadius := radius / 3
	midRadius := radius * 2 / 3

	var centerCandidates, midCandidates, outerCandidates []hex.Coord
	for _, c := range candidates {
		dist := c.DistanceToOrigin()
		switch {
		case dist <= centerRadius:
			centerCandidates = append(centerCandidates, c)
		case dist <= midRadius:
			midCandidates = append(midCandidates, c)
		default:
			outerCandidates = append(outerCandidates, c)
		}
	}

	// For symmetric placement, we place structures in the "positive half" only
	// and mirror them. For odd counts, one goes in the center.
	halfCount := count / 2
	hasCenter := count%2 == 1

	var positions []hex.Coord

	// Place center structure (if odd count)
	if hasCenter && len(centerCandidates) > 0 {
		// Pick the candidate closest to the origin
		sort.Slice(centerCandidates, func(i, j int) bool {
			return centerCandidates[i].DistanceToOrigin() < centerCandidates[j].DistanceToOrigin()
		})
		positions = append(positions, centerCandidates[0])
	}

	// Distribute half the structures across zones
	// ~40% center, ~30% mid, ~30% outer (of the half count)
	centerCount := max(1, halfCount*40/100)
	midCount := max(1, halfCount*30/100)
	outerCount := halfCount - centerCount - midCount
	if outerCount < 0 {
		outerCount = 0
		midCount = halfCount - centerCount
	}

	// Filter to "positive half" only for symmetric placement
	filterHalf := func(coords []hex.Coord) []hex.Coord {
		var result []hex.Coord
		for _, c := range coords {
			if c.Q > 0 || (c.Q == 0 && c.R > 0) {
				result = append(result, c)
			}
		}
		return result
	}

	centerHalf := filterHalf(centerCandidates)
	midHalf := filterHalf(midCandidates)
	outerHalf := filterHalf(outerCandidates)

	// Pick from each zone
	placed := make(map[hex.Coord]bool)
	for _, pos := range positions {
		placed[pos] = true
		placed[pos.Rotate180()] = true
	}

	pickFromZone := func(zone []hex.Coord, n int) []hex.Coord {
		shuffle(zone, rng)
		var picked []hex.Coord
		for _, c := range zone {
			if len(picked) >= n {
				break
			}
			if placed[c] || placed[c.Rotate180()] {
				continue
			}
			// Check minimum distance from all placed structures
			tooClose := false
			for p := range placed {
				if c.Distance(p) < 2 {
					tooClose = true
					break
				}
			}
			if tooClose {
				continue
			}
			picked = append(picked, c)
			placed[c] = true
			placed[c.Rotate180()] = true
		}
		return picked
	}

	halfPositions := pickFromZone(centerHalf, centerCount)
	halfPositions = append(halfPositions, pickFromZone(midHalf, midCount)...)
	halfPositions = append(halfPositions, pickFromZone(outerHalf, outerCount)...)

	// If we didn't get enough, try all remaining candidates
	if len(halfPositions) < halfCount {
		allHalf := filterHalf(candidates)
		remaining := halfCount - len(halfPositions)
		halfPositions = append(halfPositions, pickFromZone(allHalf, remaining)...)
	}

	// Apply symmetry to get the full set
	for _, pos := range halfPositions {
		positions = append(positions, pos)
		rotated := pos.Rotate180()
		if rotated != pos {
			positions = append(positions, rotated)
		}
	}

	return positions
}

// AssignStructureTypes assigns types to neutral structure positions.
// ~70% Outposts, ~30% Command Centers, maintaining symmetric pairs.
func AssignStructureTypes(positions []hex.Coord) map[hex.Coord]model.StructureType {
	result := make(map[hex.Coord]model.StructureType)

	ccCount := max(1, len(positions)*30/100)
	assigned := 0

	// Assign in pairs (symmetric structures get the same type)
	paired := make(map[hex.Coord]bool)
	for _, pos := range positions {
		if paired[pos] {
			continue
		}
		rotated := pos.Rotate180()

		var sType model.StructureType
		if assigned < ccCount {
			sType = model.StructureCommandCenter
		} else {
			sType = model.StructureOutpost
		}
		assigned++

		result[pos] = sType
		paired[pos] = true
		if rotated != pos {
			result[rotated] = sType
			paired[rotated] = true
		}
	}

	return result
}

// gatherCandidates returns all hexes that are valid candidates for structure placement.
func gatherCandidates(grid *hex.Grid, terrain map[hex.Coord]model.TerrainType, hq1, hq2 hex.Coord) []hex.Coord {
	var candidates []hex.Coord
	for _, c := range grid.AllHexes() {
		t := terrain[c]
		if !model.IsPassable(t) {
			continue
		}
		// Minimum 3 hexes from any HQ
		if c.Distance(hq1) < 3 || c.Distance(hq2) < 3 {
			continue
		}
		// Don't place on HQ hexes
		if c == hq1 || c == hq2 {
			continue
		}
		candidates = append(candidates, c)
	}
	return candidates
}

func shuffle(coords []hex.Coord, rng *rand.Rand) {
	rng.Shuffle(len(coords), func(i, j int) {
		coords[i], coords[j] = coords[j], coords[i]
	})
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
