package hex

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewGrid(t *testing.T) {
	// Radius 0: 1 hex (center)
	g0 := NewGrid(0)
	assert.Equal(t, 1, g0.HexCount())
	assert.True(t, g0.Contains(Origin()))

	// Radius 1: 7 hexes (1 + 6)
	g1 := NewGrid(1)
	assert.Equal(t, 7, g1.HexCount())
	assert.True(t, g1.Contains(NewCoord(1, -1, 0)))
	assert.False(t, g1.Contains(NewCoord(2, -2, 0)))

	// Radius 2: 19 hexes (1 + 6 + 12)
	g2 := NewGrid(2)
	assert.Equal(t, 19, g2.HexCount())

	allHexes := g2.AllHexes()
	require.Len(t, allHexes, 19)
}

func TestNeighborsInGrid(t *testing.T) {
	g := NewGrid(1) // Center + 6 neighbors
	center := Origin()

	// Center should have 6 neighbors
	centerNeighbors := g.Neighbors(center)
	require.Len(t, centerNeighbors, 6)

	// An edge hex should have fewer neighbors within the grid
	edge := NewCoord(1, 0, -1)
	edgeNeighbors := g.Neighbors(edge)
	require.Len(t, edgeNeighbors, 3, "Edge hex in radius 1 grid should have 3 valid neighbors")

	// Verify the actual neighbors of the edge hex
	expectedEdgeNeighbors := []Coord{
		Origin(),
		NewCoord(1, -1, 0),
		NewCoord(0, 1, -1),
	}
	assert.ElementsMatch(t, expectedEdgeNeighbors, edgeNeighbors)
}

func TestHexesInRange(t *testing.T) {
	g := NewGrid(2)

	// All hexes within range 1 of center (should be 7)
	inRange1 := g.HexesInRange(Origin(), 1)
	assert.Len(t, inRange1, 7)

	// Hexes within range 2 of an edge hex
	// Edge hex is at distance 2. Range 2 from it covers some of the grid, but not all.
	edge := NewCoord(2, 0, -2)
	inRangeEdge := g.HexesInRange(edge, 2)
	assert.Less(t, len(inRangeEdge), 19, "Range should be constrained by grid bounds")
}

func TestHalfGrid(t *testing.T) {
	g := NewGrid(2) // 19 total hexes

	half := g.HalfGrid()
	// Total 19: 1 center, 18 others. Half should be 9.
	require.Len(t, half, 9)

	for _, c := range half {
		// Verify all are in the positive half
		assert.True(t, c.Q > 0 || (c.Q == 0 && c.R > 0))
		assert.NotEqual(t, Origin(), c) // Center not included
	}
}

func TestEdgeHexes(t *testing.T) {
	g := NewGrid(3)

	edge := g.EdgeHexes()
	// Radius 3 outer ring should have 6 * 3 = 18 hexes
	require.Len(t, edge, 18)

	for _, c := range edge {
		assert.Equal(t, 3, c.DistanceToOrigin(), "All edge hexes should be at max radius")
	}
}
