package mapgen

import (
	"fmt"
	"math"
	"math/rand"

	"github.com/teomiscia/hexbattle/internal/config"
	"github.com/teomiscia/hexbattle/internal/hex"
	"github.com/teomiscia/hexbattle/internal/model"
)

// MapResult holds the output of map generation.
type MapResult struct {
	Terrain    map[hex.Coord]model.TerrainType
	HQ1        hex.Coord
	HQ2        hex.Coord
	Structures []StructurePlacement
}

// StructurePlacement holds a structure's position and type.
type StructurePlacement struct {
	Position hex.Coord
	Type     model.StructureType
	OwnerID  string // empty for neutral, player ID for HQs
}

// Generate creates a procedural hex map with the given parameters.
// It retries up to maxRetries times if validation fails.
func Generate(mapSize model.MapSize, seed int64, balance *config.BalanceData) (*MapResult, error) {
	radius := mapSize.Radius()
	grid := hex.NewGrid(radius)

	// Get config values
	maxRetries := 10
	minPassableRatio := 0.60
	structureCount := 5

	if balance != nil {
		maxRetries = balance.MapGen.MaxRetries
		minPassableRatio = balance.MapGen.MinPassableRatio
		if count, ok := balance.MapGen.StructureCounts[string(mapSize)]; ok {
			structureCount = count
		}
	}

	rng := rand.New(rand.NewSource(seed))

	for attempt := 0; attempt < maxRetries; attempt++ {
		attemptSeed := seed + int64(attempt)
		terrain := generateTerrain(grid, attemptSeed, balance)

		// Place HQs
		hq1, hq2 := PlaceHQs(grid)

		// Ensure HQ hexes and surroundings are passable
		ensurePassable(terrain, grid, hq1, 2)
		ensurePassable(terrain, grid, hq2, 2)

		// Apply symmetry
		ApplySymmetry(terrain, grid)

		// Place neutral structures
		structPositions := PlaceNeutralStructures(grid, terrain, hq1, hq2, structureCount, rng)

		// Validate
		if !ValidateMap(grid, terrain, hq1, hq2, structPositions, minPassableRatio) {
			continue
		}

		// Assign structure types
		structTypes := AssignStructureTypes(structPositions)

		// Build result
		result := &MapResult{
			Terrain: terrain,
			HQ1:     hq1,
			HQ2:     hq2,
		}

		// Add HQs as structures
		result.Structures = append(result.Structures, StructurePlacement{
			Position: hq1,
			Type:     model.StructureHQ,
			OwnerID:  "", // Will be set by caller with actual player IDs
		})
		result.Structures = append(result.Structures, StructurePlacement{
			Position: hq2,
			Type:     model.StructureHQ,
			OwnerID:  "", // Will be set by caller
		})

		// Add neutral structures
		for _, pos := range structPositions {
			sType := structTypes[pos]
			result.Structures = append(result.Structures, StructurePlacement{
				Position: pos,
				Type:     sType,
				OwnerID:  "", // neutral
			})
		}

		return result, nil
	}

	return nil, fmt.Errorf("mapgen: failed to generate valid map after %d attempts", maxRetries)
}

// generateTerrain creates terrain assignments using multi-octave noise.
func generateTerrain(grid *hex.Grid, seed int64, balance *config.BalanceData) map[hex.Coord]model.TerrainType {
	terrain := make(map[hex.Coord]model.TerrainType)
	noise := NewNoiseGenerator(seed, 0.15)

	// Get noise thresholds
	waterThreshold := 0.15
	plainsThreshold := 0.55
	forestThreshold := 0.75
	hillsThreshold := 0.88

	if balance != nil {
		if v, ok := balance.MapGen.NoiseThresholds["water"]; ok {
			waterThreshold = v
		}
		if v, ok := balance.MapGen.NoiseThresholds["plains"]; ok {
			plainsThreshold = v
		}
		if v, ok := balance.MapGen.NoiseThresholds["forest"]; ok {
			forestThreshold = v
		}
		if v, ok := balance.MapGen.NoiseThresholds["hills"]; ok {
			hillsThreshold = v
		}
	}

	// Only generate for one half; symmetry will fill the other
	half := grid.HalfGrid()
	for _, c := range half {
		// Convert hex coords to cartesian for noise sampling
		x, y := hexToCartesian(c)
		value := noise.MultiOctave(x, y)

		t := noiseToTerrain(value, waterThreshold, plainsThreshold, forestThreshold, hillsThreshold)
		terrain[c] = t
	}

	// Center is always plains
	terrain[hex.Origin()] = model.TerrainPlains

	return terrain
}

// noiseToTerrain maps a noise value [0,1] to a terrain type.
func noiseToTerrain(value, water, plains, forest, hills float64) model.TerrainType {
	switch {
	case value < water:
		return model.TerrainWater
	case value < plains:
		return model.TerrainPlains
	case value < forest:
		return model.TerrainForest
	case value < hills:
		return model.TerrainHills
	default:
		return model.TerrainMountains
	}
}

// hexToCartesian converts cube coordinates to approximate cartesian for noise sampling.
func hexToCartesian(c hex.Coord) (float64, float64) {
	x := math.Sqrt(3)*float64(c.Q) + math.Sqrt(3)/2*float64(c.R)
	y := 3.0 / 2 * float64(c.R)
	return x, y
}

// ensurePassable forces all hexes within the given radius of center to be passable.
func ensurePassable(terrain map[hex.Coord]model.TerrainType, grid *hex.Grid, center hex.Coord, radius int) {
	hexes := grid.HexesInRange(center, radius)
	for _, c := range hexes {
		t, ok := terrain[c]
		if !ok || !model.IsPassable(t) {
			terrain[c] = model.TerrainPlains
		}
	}
}
