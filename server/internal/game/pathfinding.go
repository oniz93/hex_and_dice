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

	reached := make(map[hex.Coord]int)
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

			// Skip hexes occupied by enemy troops
			if gs.IsHexOccupiedByEnemy(neighbor, ownerID) {
				continue
			}

			moveCost := model.MovementCost(terrain)
			totalCost := current.cost + moveCost

			if totalCost > mobility {
				continue
			}

			if prevCost, visited := reached[neighbor]; !visited || totalCost < prevCost {
				reached[neighbor] = totalCost
				heap.Push(frontier, costEntry{pos: neighbor, cost: totalCost})
			}
		}
	}

	// Build result: exclude hexes occupied by friendly troops (can pass through but not stop on)
	// but keep the start position
	result := make(map[hex.Coord]int)
	for pos, cost := range reached {
		if pos == start {
			continue // don't include start in "reachable destinations"
		}
		occupant := gs.TroopAtHex(pos)
		if occupant != nil {
			continue // can't stop on occupied hex
		}
		result[pos] = cost
	}

	return result
}

// CanReach checks if a troop can move to the target hex.
// Returns the movement cost if reachable, or -1 if not.
func CanReach(gs *GameState, troop *model.Troop, target hex.Coord) int {
	reachable := ReachableHexes(gs, troop)
	cost, ok := reachable[target]
	if !ok {
		return -1
	}
	return cost
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
