package hex

import "math"

// Coord represents a position in cube coordinate space.
// The invariant q + r + s = 0 must always hold.
type Coord struct {
	Q int `json:"q"`
	R int `json:"r"`
	S int `json:"s"`
}

// NewCoord creates a cube coordinate. It panics if q+r+s != 0.
func NewCoord(q, r, s int) Coord {
	if q+r+s != 0 {
		panic("hex: invalid cube coordinate: q+r+s must equal 0")
	}
	return Coord{Q: q, R: r, S: s}
}

// NewCoordQR creates a cube coordinate from axial (q, r), deriving s = -q-r.
func NewCoordQR(q, r int) Coord {
	return Coord{Q: q, R: r, S: -q - r}
}

// Origin returns the center hex (0, 0, 0).
func Origin() Coord {
	return Coord{}
}

// Add returns the sum of two coordinates.
func (c Coord) Add(other Coord) Coord {
	return Coord{Q: c.Q + other.Q, R: c.R + other.R, S: c.S + other.S}
}

// Sub returns the difference of two coordinates.
func (c Coord) Sub(other Coord) Coord {
	return Coord{Q: c.Q - other.Q, R: c.R - other.R, S: c.S - other.S}
}

// Scale multiplies each component by a scalar.
func (c Coord) Scale(k int) Coord {
	return Coord{Q: c.Q * k, R: c.R * k, S: c.S * k}
}

// Rotate180 returns the coordinate rotated 180 degrees around the origin.
// This is simply the negation of all components.
func (c Coord) Rotate180() Coord {
	return Coord{Q: -c.Q, R: -c.R, S: -c.S}
}

// Distance returns the hex distance between two coordinates.
// In cube coordinates: max(|dq|, |dr|, |ds|).
func (c Coord) Distance(other Coord) int {
	dq := abs(c.Q - other.Q)
	dr := abs(c.R - other.R)
	ds := abs(c.S - other.S)
	return max(dq, max(dr, ds))
}

// DistanceToOrigin returns the hex distance from this coordinate to (0,0,0).
func (c Coord) DistanceToOrigin() int {
	return c.Distance(Origin())
}

// Length returns the distance from the origin (same as DistanceToOrigin).
func (c Coord) Length() int {
	return c.DistanceToOrigin()
}

// Neighbor returns the adjacent hex in the given direction (0-5).
func (c Coord) Neighbor(dir int) Coord {
	return c.Add(Direction(dir))
}

// Neighbors returns all 6 adjacent hexes.
func (c Coord) Neighbors() [6]Coord {
	var result [6]Coord
	for i := 0; i < 6; i++ {
		result[i] = c.Neighbor(i)
	}
	return result
}

// Ring returns all hexes at exactly the given radius from this coordinate.
// radius must be >= 1. For radius 0, returns a slice containing only this coord.
func (c Coord) Ring(radius int) []Coord {
	if radius == 0 {
		return []Coord{c}
	}

	results := make([]Coord, 0, 6*radius)
	// Start at the hex radius steps in direction 4 (southwest)
	current := c.Add(Direction(4).Scale(radius))

	for i := 0; i < 6; i++ {
		for j := 0; j < radius; j++ {
			results = append(results, current)
			current = current.Neighbor(i)
		}
	}
	return results
}

// Spiral returns all hexes from center outward up to (and including) the given radius.
func (c Coord) Spiral(radius int) []Coord {
	results := []Coord{c}
	for r := 1; r <= radius; r++ {
		results = append(results, c.Ring(r)...)
	}
	return results
}

// PixelCenter returns the pixel position of the hex center for pointy-top layout.
// hexSize is the distance from the center to any vertex.
func (c Coord) PixelCenter(hexSize float64) (x, y float64) {
	x = hexSize * (math.Sqrt(3)*float64(c.Q) + math.Sqrt(3)/2*float64(c.R))
	y = hexSize * (3.0 / 2 * float64(c.R))
	return
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
