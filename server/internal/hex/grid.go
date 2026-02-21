package hex

// Grid represents a hexagonal grid with a given radius.
// The grid contains all hexes where Distance(hex, origin) <= Radius.
type Grid struct {
	Radius int
	hexes  map[Coord]bool
}

// NewGrid creates a hexagonal grid with the given radius.
// The grid contains all hexes within radius distance from the origin.
func NewGrid(radius int) *Grid {
	g := &Grid{
		Radius: radius,
		hexes:  make(map[Coord]bool),
	}
	for _, c := range Origin().Spiral(radius) {
		g.hexes[c] = true
	}
	return g
}

// Contains returns true if the coordinate is within the grid bounds.
func (g *Grid) Contains(c Coord) bool {
	return g.hexes[c]
}

// HexCount returns the total number of hexes in the grid.
// For a hex grid of radius r, this is 3*r*(r+1) + 1.
func (g *Grid) HexCount() int {
	return len(g.hexes)
}

// AllHexes returns all coordinates in the grid.
func (g *Grid) AllHexes() []Coord {
	result := make([]Coord, 0, len(g.hexes))
	for c := range g.hexes {
		result = append(result, c)
	}
	return result
}

// Neighbors returns all valid neighbors of the given coordinate that are within the grid.
func (g *Grid) Neighbors(c Coord) []Coord {
	all := c.Neighbors()
	result := make([]Coord, 0, 6)
	for _, n := range all {
		if g.Contains(n) {
			result = append(result, n)
		}
	}
	return result
}

// HexesInRange returns all hexes within the given distance from center
// that are also within the grid bounds.
func (g *Grid) HexesInRange(center Coord, distance int) []Coord {
	result := make([]Coord, 0)
	for _, c := range center.Spiral(distance) {
		if g.Contains(c) {
			result = append(result, c)
		}
	}
	return result
}

// EdgeHexes returns all hexes at exactly the grid radius (the outermost ring).
func (g *Grid) EdgeHexes() []Coord {
	return Origin().Ring(g.Radius)
}

// HalfGrid returns all hexes in one "hemisphere" of the grid, used for
// rotational symmetry map generation. Includes hexes where q > 0, or
// where q == 0 and r > 0. The center hex (0,0,0) is excluded.
func (g *Grid) HalfGrid() []Coord {
	result := make([]Coord, 0, len(g.hexes)/2)
	for c := range g.hexes {
		if c.Q > 0 || (c.Q == 0 && c.R > 0) {
			result = append(result, c)
		}
	}
	return result
}
