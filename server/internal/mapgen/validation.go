package mapgen

import (
	"github.com/teomiscia/hexbattle/internal/hex"
	"github.com/teomiscia/hexbattle/internal/model"
)

// ValidateMap checks all map constraints. Returns true if the map is valid.
func ValidateMap(grid *hex.Grid, terrain map[hex.Coord]model.TerrainType, hq1, hq2 hex.Coord, structurePositions []hex.Coord, minPassableRatio float64) bool {
	// 1. Connectivity: flood-fill from HQ1 reaches HQ2
	if !checkConnectivity(grid, terrain, hq1, hq2) {
		return false
	}

	// 2. Structure accessibility: every structure is reachable from both HQs
	if !checkStructureAccessibility(grid, terrain, hq1, hq2, structurePositions) {
		return false
	}

	// 3. No isolated regions: single connected component of passable terrain
	if !checkSingleComponent(grid, terrain) {
		return false
	}

	// 4. Minimum passable ratio
	if !checkPassableRatio(grid, terrain, minPassableRatio) {
		return false
	}

	// 5. HQ safety: no impassable terrain within 2 hexes of either HQ
	if !checkHQSafety(grid, terrain, hq1) || !checkHQSafety(grid, terrain, hq2) {
		return false
	}

	return true
}

// checkConnectivity does a flood-fill from start and checks if target is reachable.
func checkConnectivity(grid *hex.Grid, terrain map[hex.Coord]model.TerrainType, start, target hex.Coord) bool {
	visited := floodFill(grid, terrain, start)
	return visited[target]
}

// checkStructureAccessibility checks that every structure is reachable from both HQs.
func checkStructureAccessibility(grid *hex.Grid, terrain map[hex.Coord]model.TerrainType, hq1, hq2 hex.Coord, structures []hex.Coord) bool {
	reachableFromHQ1 := floodFill(grid, terrain, hq1)
	reachableFromHQ2 := floodFill(grid, terrain, hq2)

	for _, pos := range structures {
		if !reachableFromHQ1[pos] || !reachableFromHQ2[pos] {
			return false
		}
	}
	return true
}

// checkSingleComponent verifies there is exactly one connected component of passable terrain.
func checkSingleComponent(grid *hex.Grid, terrain map[hex.Coord]model.TerrainType) bool {
	// Find the first passable hex
	var start hex.Coord
	found := false
	for _, c := range grid.AllHexes() {
		t, ok := terrain[c]
		if !ok {
			t = model.TerrainPlains
		}
		if model.IsPassable(t) {
			start = c
			found = true
			break
		}
	}
	if !found {
		return false
	}

	visited := floodFill(grid, terrain, start)

	// Count all passable hexes
	passableCount := 0
	for _, c := range grid.AllHexes() {
		t, ok := terrain[c]
		if !ok {
			t = model.TerrainPlains
		}
		if model.IsPassable(t) {
			passableCount++
		}
	}

	return len(visited) == passableCount
}

// checkPassableRatio verifies that at least the given fraction of hexes are passable.
func checkPassableRatio(grid *hex.Grid, terrain map[hex.Coord]model.TerrainType, minRatio float64) bool {
	total := grid.HexCount()
	passable := 0
	for _, c := range grid.AllHexes() {
		t, ok := terrain[c]
		if !ok {
			t = model.TerrainPlains
		}
		if model.IsPassable(t) {
			passable++
		}
	}
	return float64(passable)/float64(total) >= minRatio
}

// checkHQSafety verifies no impassable terrain within 2 hexes of the HQ.
func checkHQSafety(grid *hex.Grid, terrain map[hex.Coord]model.TerrainType, hq hex.Coord) bool {
	hexesInRange := grid.HexesInRange(hq, 2)
	for _, c := range hexesInRange {
		t, ok := terrain[c]
		if !ok {
			t = model.TerrainPlains
		}
		if !model.IsPassable(t) {
			return false
		}
	}
	return true
}

// floodFill returns all passable hexes reachable from start.
func floodFill(grid *hex.Grid, terrain map[hex.Coord]model.TerrainType, start hex.Coord) map[hex.Coord]bool {
	visited := make(map[hex.Coord]bool)
	queue := []hex.Coord{start}
	visited[start] = true

	for len(queue) > 0 {
		current := queue[0]
		queue = queue[1:]

		for _, neighbor := range grid.Neighbors(current) {
			if visited[neighbor] {
				continue
			}
			t, ok := terrain[neighbor]
			if !ok {
				t = model.TerrainPlains
			}
			if !model.IsPassable(t) {
				continue
			}
			visited[neighbor] = true
			queue = append(queue, neighbor)
		}
	}

	return visited
}
