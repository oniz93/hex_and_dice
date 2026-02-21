package hex

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewCoord(t *testing.T) {
	t.Run("valid coordinates", func(t *testing.T) {
		assert.NotPanics(t, func() {
			c := NewCoord(1, -1, 0)
			assert.Equal(t, 1, c.Q)
			assert.Equal(t, -1, c.R)
			assert.Equal(t, 0, c.S)
		})
	})

	t.Run("invalid coordinates panic", func(t *testing.T) {
		assert.Panics(t, func() {
			NewCoord(1, 1, 1) // q+r+s != 0
		})
	})
}

func TestNewCoordQR(t *testing.T) {
	c := NewCoordQR(2, 3)
	assert.Equal(t, 2, c.Q)
	assert.Equal(t, 3, c.R)
	assert.Equal(t, -5, c.S)
	assert.Equal(t, 0, c.Q+c.R+c.S)
}

func TestAddSubScale(t *testing.T) {
	c1 := NewCoord(1, -1, 0)
	c2 := NewCoord(2, 0, -2)

	// Add
	sum := c1.Add(c2)
	assert.Equal(t, NewCoord(3, -1, -2), sum)

	// Sub
	diff := c1.Sub(c2)
	assert.Equal(t, NewCoord(-1, -1, 2), diff)

	// Scale
	scaled := c1.Scale(3)
	assert.Equal(t, NewCoord(3, -3, 0), scaled)
}

func TestRotate180(t *testing.T) {
	c := NewCoord(2, -5, 3)
	rotated := c.Rotate180()
	assert.Equal(t, NewCoord(-2, 5, -3), rotated)
}

func TestDistance(t *testing.T) {
	c1 := NewCoord(1, -1, 0)
	c2 := NewCoord(3, -5, 2)

	dist := c1.Distance(c2)
	assert.Equal(t, 4, dist)

	// Distance to origin
	assert.Equal(t, 1, c1.DistanceToOrigin())
	assert.Equal(t, 5, c2.DistanceToOrigin())
}

func TestNeighbors(t *testing.T) {
	c := NewCoord(0, 0, 0)
	neighbors := c.Neighbors()
	require.Len(t, neighbors, 6)

	expected := [6]Coord{
		NewCoord(1, 0, -1),
		NewCoord(1, -1, 0),
		NewCoord(0, -1, 1),
		NewCoord(-1, 0, 1),
		NewCoord(-1, 1, 0),
		NewCoord(0, 1, -1),
	}
	assert.ElementsMatch(t, expected, neighbors)
}

func TestRing(t *testing.T) {
	center := NewCoord(0, 0, 0)

	// Radius 0 should be just the center
	r0 := center.Ring(0)
	require.Len(t, r0, 1)
	assert.Equal(t, center, r0[0])

	// Radius 1 should be the 6 neighbors
	r1 := center.Ring(1)
	require.Len(t, r1, 6)
	assert.ElementsMatch(t, center.Neighbors(), r1)

	// Radius 2 should be 12 hexes
	r2 := center.Ring(2)
	require.Len(t, r2, 12)
	for _, hex := range r2 {
		assert.Equal(t, 2, hex.DistanceToOrigin(), "all hexes in ring 2 should be distance 2 from center")
	}
}

func TestSpiral(t *testing.T) {
	center := NewCoord(0, 0, 0)

	// Spiral radius 2 should contain 1 + 6 + 12 = 19 hexes
	s2 := center.Spiral(2)
	require.Len(t, s2, 19)

	distCounts := make(map[int]int)
	for _, hex := range s2 {
		distCounts[hex.DistanceToOrigin()]++
	}

	assert.Equal(t, 1, distCounts[0])
	assert.Equal(t, 6, distCounts[1])
	assert.Equal(t, 12, distCounts[2])
}
