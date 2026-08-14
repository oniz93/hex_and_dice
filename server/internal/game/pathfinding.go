package game

import (
	"container/heap"

	"github.com/teomiscia/hexbattle/internal/hex"
	"github.com/teomiscia/hexbattle/internal/model"
)

// ReachableHexes computes all hexes a troop can move to given its remaining mobility.
// Uses Dijkstra's algorithm (BFS with terrain costs).
// Returns a map of reachable hex -> movement cost to reach it.
func ReachableHexes(gs *GameState, troop *model.Troop) map[hex.Coord]int {
	start := troop.Hex
	mobility := troop.RemainingMobility
	ownerID := troop.OwnerID

	reached := make(map[hex.Coord]int, 64)
	reached[start] = 0

	frontier := &costQueue{}
	heap.Init(frontier)
	heap.Push(frontier, costEntry{pos: start, cost: 0})

	for frontier.Len() > 0 {
		current := heap.Pop(frontier).(costEntry)

		if current.cost > reached[current.pos] {
			continue // stale entry
		}

		neighbors := gs.Grid.Neighbors(current.pos)
		for _, neighbor := range neighbors {
			// Skip impassable terrain
			terrain := gs.GetTerrainAt(neighbor)
			if !model.IsPassable(terrain) {
				continue
			}

			// Skip hexes occupied by enemy troops (can't even pass through)
			if gs.IsHexOccupiedByEnemy(neighbor, ownerID) {
				continue
			}

			moveCost := model.MovementCost(terrain)
			totalCost := current.cost + moveCost

			// Minimum movement rule: if adjacent to start and have mobility left,
			// always allow moving to at least one cell (unless impassable).
			if current.pos == start && mobility > 0 && totalCost > mobility {
				totalCost = mobility
			}

			if totalCost > mobility {
				continue
			}

			if prevCost, visited := reached[neighbor]; !visited || totalCost < prevCost {
				reached[neighbor] = totalCost
				heap.Push(frontier, costEntry{pos: neighbor, cost: totalCost})
			}
		}
	}

	// Build result: exclude hexes occupied by troops or structures
	// (can pass through friendly troops but not stop on them/structures)
	result := make(map[hex.Coord]int, len(reached))
	for pos, cost := range reached {
		if pos == start {
			continue // don't include start in "reachable destinations"
		}
		if gs.TroopAtHex(pos) != nil {
			continue // can't stop on occupied hex
		}
		if gs.StructureAtHex(pos) != nil {
			continue // can't stop on structure
		}
		result[pos] = cost
	}

	return result
}

// CanReach checks if a troop can move to the target hex.
// Returns the movement cost if reachable, or -1 if not.
//
// Unlike ReachableHexes it only searches until the target is found, so it is
// much cheaper when a single destination is all that is needed (move
// validation and execution paths).
func CanReach(gs *GameState, troop *model.Troop, target hex.Coord) int {
	return movementCostTo(gs, troop, target)
}

// MoveCostTo is an alias for CanReach. It exists so call sites that only need
// the movement cost (rather than a boolean reachability check) read clearly.
func MoveCostTo(gs *GameState, troop *model.Troop, target hex.Coord) int {
	return movementCostTo(gs, troop, target)
}

// movementCostTo runs a targeted Dijkstra search from the troop's position to
// target, stopping as soon as the target is dequeued with its final cost.
func movementCostTo(gs *GameState, troop *model.Troop, target hex.Coord) int {
	start := troop.Hex
	mobility := troop.RemainingMobility
	ownerID := troop.OwnerID

	if target == start {
		return -1
	}

	// The target must be a legal stopping hex.
	if !model.IsPassable(gs.GetTerrainAt(target)) {
		return -1
	}
	if gs.TroopAtHex(target) != nil || gs.StructureAtHex(target) != nil {
		return -1
	}

	reached := make(map[hex.Coord]int, 32)
	reached[start] = 0

	frontier := &costQueue{}
	heap.Init(frontier)
	heap.Push(frontier, costEntry{pos: start, cost: 0})

	for frontier.Len() > 0 {
		current := heap.Pop(frontier).(costEntry)

		if current.cost > reached[current.pos] {
			continue // stale entry
		}
		if current.pos == target {
			return current.cost
		}

		neighbors := gs.Grid.Neighbors(current.pos)
		for _, neighbor := range neighbors {
			terrain := gs.GetTerrainAt(neighbor)
			if !model.IsPassable(terrain) {
				continue
			}

			// Enemies block passage entirely.
			if gs.IsHexOccupiedByEnemy(neighbor, ownerID) {
				continue
			}

			moveCost := model.MovementCost(terrain)
			totalCost := current.cost + moveCost

			// Minimum movement rule: if adjacent to start and mobility is
			// non-zero, allow reaching at least one cell even if the terrain
			// cost exceeds the remaining mobility.
			if current.pos == start && mobility > 0 && totalCost > mobility {
				totalCost = mobility
			}

			if totalCost > mobility {
				continue
			}

			if prevCost, visited := reached[neighbor]; !visited || totalCost < prevCost {
				reached[neighbor] = totalCost
				heap.Push(frontier, costEntry{pos: neighbor, cost: totalCost})
			}
		}
	}

	return -1
}

// HexDistance returns the hex distance between two cube coordinates.
func HexDistance(a, b hex.Coord) int {
	return a.Distance(b)
}

// CanAttackTarget checks if a troop can attack a target at the given hex.
// Only checks range, not whether the target is valid.
func CanAttackTarget(troop *model.Troop, targetHex hex.Coord) bool {
	dist := HexDistance(troop.Hex, targetHex)
	return dist >= 1 && dist <= troop.Range
}

// --- Priority queue for Dijkstra ---

type costEntry struct {
	pos  hex.Coord
	cost int
}

type costQueue []costEntry

func (q costQueue) Len() int            { return len(q) }
func (q costQueue) Less(i, j int) bool  { return q[i].cost < q[j].cost }
func (q costQueue) Swap(i, j int)       { q[i], q[j] = q[j], q[i] }
func (q *costQueue) Push(x interface{}) { *q = append(*q, x.(costEntry)) }
func (q *costQueue) Pop() interface{} {
	old := *q
	n := len(old)
	item := old[n-1]
	*q = old[:n-1]
	return item
}
