package bot

import (
	"log/slog"
	"math/rand"
	"sort"

	"github.com/teomiscia/hexbattle/internal/game"
	"github.com/teomiscia/hexbattle/internal/hex"
	"github.com/teomiscia/hexbattle/internal/model"
)

// Difficulty controls the bot's decision quality.
type Difficulty string

const (
	DifficultyEasy   Difficulty = "easy"
	DifficultyMedium Difficulty = "medium"
	DifficultyHard   Difficulty = "hard"
)

// Bot implements game.BotPlayer with a simple greedy AI.
type Bot struct {
	id         string
	difficulty Difficulty
	rng        *rand.Rand
	logger     *slog.Logger

	// Per-turn state to track which actions have been yielded.
	turnNumber int
	phase      botPhase
	acted      map[string]bool // unit IDs that have already been processed this turn
}

type botPhase int

const (
	phaseBuy botPhase = iota
	phaseAttack
	phaseMove
	phaseDone
)

// New creates a new Bot with the given player ID and difficulty.
func New(playerID string, difficulty Difficulty, seed int64) *Bot {
	return &Bot{
		id:         playerID,
		difficulty: difficulty,
		rng:        rand.New(rand.NewSource(seed)),
		logger:     slog.Default().With("component", "bot", "player_id", playerID, "difficulty", difficulty),
		acted:      make(map[string]bool),
	}
}

// PlayerID returns the bot's player ID.
func (b *Bot) PlayerID() string {
	return b.id
}

// NextAction returns the next action the bot wants to take.
// Returns nil when the bot is done (engine will end the turn).
func (b *Bot) NextAction(gs *game.GameState) *game.BotAction {
	// Reset per-turn state if the turn number changed.
	if gs.TurnNumber != b.turnNumber {
		b.turnNumber = gs.TurnNumber
		b.phase = phaseBuy
		b.acted = make(map[string]bool)
	}

	for b.phase != phaseDone {
		switch b.phase {
		case phaseBuy:
			action := b.planBuy(gs)
			if action != nil {
				return action
			}
			b.phase = phaseAttack

		case phaseAttack:
			action := b.planAttack(gs)
			if action != nil {
				return action
			}
			b.phase = phaseMove

		case phaseMove:
			action := b.planMove(gs)
			if action != nil {
				return action
			}
			b.phase = phaseDone
		}
	}

	return nil // done — engine will end the turn
}

// ---------------------------------------------------------------------------
// Buy phase
// ---------------------------------------------------------------------------

func (b *Bot) planBuy(gs *game.GameState) *game.BotAction {
	idx := gs.PlayerIndex(b.id)
	if idx < 0 {
		return nil
	}
	coins := gs.Players[idx].Coins

	// Find spawnable structures owned by the bot.
	var spawners []*model.Structure
	for _, s := range gs.Structures {
		if s.IsOwnedBy(b.id) && s.CanSpawn {
			// Only if the hex is free (ValidatePurchase checks this).
			if gs.TroopAtHex(s.Hex) == nil {
				spawners = append(spawners, s)
			}
		}
	}
	if len(spawners) == 0 {
		return nil
	}

	// Pick which troop to buy.
	troopType, cost := b.chooseTroopType(gs, coins)
	if troopType == "" || cost > coins {
		return nil
	}

	// Pick the spawner closest to the enemy HQ.
	enemyHQ := b.enemyHQ(gs)
	sort.Slice(spawners, func(i, j int) bool {
		return spawners[i].Hex.Distance(enemyHQ) < spawners[j].Hex.Distance(enemyHQ)
	})
	spawner := spawners[0]

	b.logger.Debug("bot buy", "troop_type", troopType, "structure_id", spawner.ID, "coins", coins)
	return &game.BotAction{
		Type:        game.BotActionBuy,
		TroopType:   troopType,
		StructureID: spawner.ID,
	}
}

func (b *Bot) chooseTroopType(gs *game.GameState, coins int) (model.TroopType, int) {
	// Simple strategy: buy the most expensive troop we can afford.
	// Easy bot: only marines. Medium: marines + snipers. Hard: all types.
	type troopOption struct {
		t    model.TroopType
		cost int
	}

	var options []troopOption
	switch b.difficulty {
	case DifficultyEasy:
		options = []troopOption{
			{model.TroopMarine, game.TroopCost(model.TroopMarine)},
		}
	case DifficultyMedium:
		options = []troopOption{
			{model.TroopMarine, game.TroopCost(model.TroopMarine)},
			{model.TroopSniper, game.TroopCost(model.TroopSniper)},
			{model.TroopHoverbike, game.TroopCost(model.TroopHoverbike)},
		}
	default: // hard
		options = []troopOption{
			{model.TroopMarine, game.TroopCost(model.TroopMarine)},
			{model.TroopSniper, game.TroopCost(model.TroopSniper)},
			{model.TroopHoverbike, game.TroopCost(model.TroopHoverbike)},
			{model.TroopMech, game.TroopCost(model.TroopMech)},
		}
	}

	// Filter affordable options.
	var affordable []troopOption
	for _, opt := range options {
		if opt.cost > 0 && opt.cost <= coins {
			affordable = append(affordable, opt)
		}
	}
	if len(affordable) == 0 {
		return "", 0
	}

	// For easy: always marine. For medium/hard: weighted random.
	if b.difficulty == DifficultyEasy {
		return affordable[0].t, affordable[0].cost
	}

	// Pick a random affordable troop, weighted toward cheaper ones.
	pick := b.rng.Intn(len(affordable))
	return affordable[pick].t, affordable[pick].cost
}

// ---------------------------------------------------------------------------
// Attack phase
// ---------------------------------------------------------------------------

func (b *Bot) planAttack(gs *game.GameState) *game.BotAction {
	troops := gs.PlayerTroops(b.id)
	for _, troop := range troops {
		if b.acted[troop.ID+"_atk"] {
			continue
		}
		if !troop.CanAttack() {
			b.acted[troop.ID+"_atk"] = true
			continue
		}

		// First-turn restriction check.
		if gs.FirstTurnRestriction && gs.TurnNumber == 1 && gs.ActivePlayer == 0 {
			b.acted[troop.ID+"_atk"] = true
			continue
		}

		target := b.findAttackTarget(gs, troop)
		if target == (hex.Coord{}) && !b.hasTarget(gs, troop) {
			b.acted[troop.ID+"_atk"] = true
			continue
		}
		if target != (hex.Coord{}) {
			b.acted[troop.ID+"_atk"] = true
			b.logger.Debug("bot attack", "unit_id", troop.ID, "target", target)
			return &game.BotAction{
				Type:   game.BotActionAttack,
				UnitID: troop.ID,
				Target: target,
			}
		}
		b.acted[troop.ID+"_atk"] = true
	}
	return nil
}

func (b *Bot) hasTarget(gs *game.GameState, troop *model.Troop) bool {
	_, found := b.findBestAttackTarget(gs, troop)
	return found
}

func (b *Bot) findAttackTarget(gs *game.GameState, troop *model.Troop) hex.Coord {
	target, found := b.findBestAttackTarget(gs, troop)
	if found {
		return target
	}
	return hex.Coord{}
}

func (b *Bot) findBestAttackTarget(gs *game.GameState, troop *model.Troop) (hex.Coord, bool) {
	type candidate struct {
		pos      hex.Coord
		priority int // lower is better
	}
	var candidates []candidate

	// Check all hexes in attack range.
	for _, h := range troop.Hex.Spiral(troop.Range) {
		if h == troop.Hex {
			continue
		}
		dist := troop.Hex.Distance(h)
		if dist < 1 || dist > troop.Range {
			continue
		}

		// Enemy troop?
		enemy := gs.TroopAtHex(h)
		if enemy != nil && enemy.OwnerID != b.id {
			prio := enemy.CurrentHP // prefer low-HP targets
			candidates = append(candidates, candidate{pos: h, priority: prio})
			continue
		}

		// Enemy/neutral structure?
		structure := gs.StructureAtHex(h)
		if structure != nil && !structure.IsOwnedBy(b.id) {
			prio := 100 + structure.CurrentHP // troops first, then structures
			candidates = append(candidates, candidate{pos: h, priority: prio})
		}
	}

	if len(candidates) == 0 {
		return hex.Coord{}, false
	}

	sort.Slice(candidates, func(i, j int) bool {
		return candidates[i].priority < candidates[j].priority
	})
	return candidates[0].pos, true
}

// ---------------------------------------------------------------------------
// Move phase
// ---------------------------------------------------------------------------

func (b *Bot) planMove(gs *game.GameState) *game.BotAction {
	troops := gs.PlayerTroops(b.id)
	for _, troop := range troops {
		if b.acted[troop.ID+"_mv"] {
			continue
		}
		if !troop.CanMove() {
			b.acted[troop.ID+"_mv"] = true
			continue
		}

		target := b.findMoveTarget(gs, troop)
		if target == (hex.Coord{}) {
			b.acted[troop.ID+"_mv"] = true
			continue
		}

		b.acted[troop.ID+"_mv"] = true
		b.logger.Debug("bot move", "unit_id", troop.ID, "target", target)
		return &game.BotAction{
			Type:   game.BotActionMove,
			UnitID: troop.ID,
			Target: target,
		}
	}
	return nil
}

func (b *Bot) findMoveTarget(gs *game.GameState, troop *model.Troop) hex.Coord {
	reachable := game.ReachableHexes(gs, troop)
	if len(reachable) == 0 {
		return hex.Coord{}
	}

	// Determine objective: move toward the nearest high-value target.
	objective := b.chooseObjective(gs, troop)

	// Find the reachable hex that minimizes distance to the objective.
	bestHex := hex.Coord{}
	bestDist := 999999
	for h := range reachable {
		dist := h.Distance(objective)
		if dist < bestDist {
			bestDist = dist
			bestHex = h
		}
	}

	// Don't move if we're already adjacent or closer to objective than our best move.
	currentDist := troop.Hex.Distance(objective)
	if bestDist >= currentDist {
		return hex.Coord{} // no improvement
	}

	return bestHex
}

func (b *Bot) chooseObjective(gs *game.GameState, troop *model.Troop) hex.Coord {
	// Priority list:
	// 1. Nearby enemy troops (if within reasonable range)
	// 2. Uncaptured/enemy structures
	// 3. Enemy HQ

	// Look for nearby enemy troops.
	nearestEnemyDist := 999
	nearestEnemyPos := hex.Coord{}
	for _, enemy := range gs.Troops {
		if enemy.OwnerID == b.id || !enemy.IsAlive() {
			continue
		}
		dist := troop.Hex.Distance(enemy.Hex)
		if dist < nearestEnemyDist {
			nearestEnemyDist = dist
			nearestEnemyPos = enemy.Hex
		}
	}

	// If an enemy is within 5 hexes, go for it.
	if nearestEnemyDist <= 5 {
		return nearestEnemyPos
	}

	// Look for uncaptured structures.
	nearestStructDist := 999
	nearestStructPos := hex.Coord{}
	for _, s := range gs.Structures {
		if s.IsOwnedBy(b.id) {
			continue
		}
		dist := troop.Hex.Distance(s.Hex)
		if dist < nearestStructDist {
			nearestStructDist = dist
			nearestStructPos = s.Hex
		}
	}
	if nearestStructDist < 999 {
		return nearestStructPos
	}

	// Fallback: enemy HQ.
	return b.enemyHQ(gs)
}

// ---------------------------------------------------------------------------
// Helpers
// ---------------------------------------------------------------------------

func (b *Bot) enemyHQ(gs *game.GameState) hex.Coord {
	for _, s := range gs.Structures {
		if s.Type == model.StructureHQ && !s.IsOwnedBy(b.id) {
			return s.Hex
		}
	}
	return hex.Coord{}
}
