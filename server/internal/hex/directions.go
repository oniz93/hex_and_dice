package hex

// Direction constants for pointy-top hexagonal grids.
// Directions are numbered 0-5, starting from east and going counter-clockwise.
//
//	Direction 0: East        (+1, 0, -1)
//	Direction 1: NorthEast   (+1, -1, 0)
//	Direction 2: NorthWest   (0, -1, +1)
//	Direction 3: West        (-1, 0, +1)
//	Direction 4: SouthWest   (-1, +1, 0)
//	Direction 5: SouthEast   (0, +1, -1)

var directions = [6]Coord{
	{Q: 1, R: 0, S: -1}, // East
	{Q: 1, R: -1, S: 0}, // NorthEast
	{Q: 0, R: -1, S: 1}, // NorthWest
	{Q: -1, R: 0, S: 1}, // West
	{Q: -1, R: 1, S: 0}, // SouthWest
	{Q: 0, R: 1, S: -1}, // SouthEast
}

const (
	DirEast      = 0
	DirNorthEast = 1
	DirNorthWest = 2
	DirWest      = 3
	DirSouthWest = 4
	DirSouthEast = 5
)

// Direction returns the coordinate offset for the given direction (0-5).
// It wraps around using modulo, so dir=6 is equivalent to dir=0.
func Direction(dir int) Coord {
	d := dir % 6
	if d < 0 {
		d += 6
	}
	return directions[d]
}

// AllDirections returns all 6 direction offsets.
func AllDirections() [6]Coord {
	return directions
}

// Scale multiplies a direction offset by k.
func (c Coord) ScaleDir(k int) Coord {
	return Coord{Q: c.Q * k, R: c.R * k, S: c.S * k}
}
