package game

import (
	"fmt"

	"github.com/teomiscia/hexbattle/internal/config"
	"github.com/teomiscia/hexbattle/internal/hex"
	"github.com/teomiscia/hexbattle/internal/model"
)

// TestBuilder is a helper for constructing game states programmatically.
type TestBuilder struct {
	state *GameState
}

// NewTestGame creates a new test game builder with reasonable defaults.
func NewTestGame() *TestBuilder {
	// Provide a default dummy balance for tests
	if Balance == nil {
		LoadBalance(&config.BalanceData{
			Economy: config.EconomyConfig{
				StartingCoins:   1000,
				PassiveIncome:   100,
				StructureIncome: 50,
			},
			Troops: map[string]config.TroopConfig{
				"marine":    {Cost: 100, HP: 10, ATK: 3, DEF: 14, Mobility: 3, Range: 1, Damage: "1D6+1"},
				"sniper":    {Cost: 150, HP: 6, ATK: 4, DEF: 11, Mobility: 2, Range: 3, Damage: "1D8"},
				"hoverbike": {Cost: 200, HP: 8, ATK: 4, DEF: 12, Mobility: 5, Range: 1, Damage: "1D8+1"},
				"mech":      {Cost: 350, HP: 12, ATK: 5, DEF: 10, Mobility: 1, Range: 3, Damage: "2D6+2", AntiStructureMultiplier: 2},
			},
			Structures: map[string]config.StructureConfig{
				"outpost":        {HP: 8, ATK: 2, DEF: 12, Range: 2, Damage: "1D4", Income: 50, Spawn: true},
				"command_center": {HP: 15, ATK: 4, DEF: 15, Range: 3, Damage: "1D6+2", Income: 50, Spawn: true},
				"hq":             {HP: 20, ATK: 3, DEF: 16, Range: 2, Damage: "1D6", Income: 0, Spawn: true},
			},
			NeutralMod: config.NeutralModConfig{
				ATKReduction:   2,
				DamageStepDown: 1,
			},
			Terrain: map[string]config.TerrainConfig{
				"plains":    {MovementCost: 1, ATKModifier: 0, DEFModifier: 0},
				"forest":    {MovementCost: 2, ATKModifier: 0, DEFModifier: 2},
				"hills":     {MovementCost: 2, ATKModifier: 1, DEFModifier: 1},
				"water":     {Passable: ptr(false)},
				"mountains": {Passable: ptr(false)},
			},
			Healing: config.HealingConfig{
				PassiveRate: 2,
			},
			SuddenDeath: config.SuddenDeathConfig{
				TurnThresholds: map[string]int{"small": 20, "medium": 30, "large": 40},
				ShrinkRate:     1,
			},
			WinCond: config.WinCondConfig{
				DominanceTurnsRequired: 3,
			},
		})
	}

	settings := model.RoomSettings{
		MapSize:   model.MapSizeSmall,
		TurnTimer: 90,
		TurnMode:  model.TurnModeAlternating,
	}

	p1 := model.PlayerState{ID: "p1", Nickname: "Player1", Coins: 1000}
	p2 := model.PlayerState{ID: "p2", Nickname: "Player2", Coins: 1000}

	state := NewGameState("test_game", settings, p1, p2, 42)
	state.Phase = model.PhasePlayerAction
	state.TurnNumber = 1
	state.FirstTurnRestriction = false // Disable for easier testing unless explicitly testing it

	return &TestBuilder{state: state}
}

func (b *TestBuilder) WithMapSize(size model.MapSize) *TestBuilder {
	b.state.MapSize = size
	b.state.Grid = hex.NewGrid(size.Radius())
	b.state.SafeZoneRadius = size.Radius()
	return b
}

func (b *TestBuilder) WithTroop(ownerID string, troopType model.TroopType, pos hex.Coord, ready bool) *TestBuilder {
	id := fmt.Sprintf("unit_%d_%d_%d", pos.Q, pos.R, pos.S)
	t, _ := NewTroopFromBalance(id, troopType, ownerID, pos)
	t.IsReady = ready
	b.state.AddTroop(t)
	return b
}

func (b *TestBuilder) WithStructure(sType model.StructureType, ownerID string, pos hex.Coord) *TestBuilder {
	id := fmt.Sprintf("struct_%d_%d_%d", pos.Q, pos.R, pos.S)
	s, _ := NewStructureFromBalance(id, sType, ownerID, pos)
	b.state.AddStructure(s)
	return b
}

func (b *TestBuilder) WithTerrain(pos hex.Coord, tType model.TerrainType) *TestBuilder {
	b.state.Terrain[pos] = tType
	return b
}

func (b *TestBuilder) WithCoins(playerID string, coins int) *TestBuilder {
	for i, p := range b.state.Players {
		if p.ID == playerID {
			b.state.Players[i].Coins = coins
			break
		}
	}
	return b
}

func (b *TestBuilder) WithTurn(turn int) *TestBuilder {
	b.state.TurnNumber = turn
	return b
}

func (b *TestBuilder) WithActivePlayer(playerID string) *TestBuilder {
	if b.state.Players[1].ID == playerID {
		b.state.ActivePlayer = 1
	} else {
		b.state.ActivePlayer = 0
	}
	return b
}

func (b *TestBuilder) Build() *GameState {
	// Ensure grid bounds are respected if terrain isn't explicitly set
	for _, c := range b.state.Grid.AllHexes() {
		if _, ok := b.state.Terrain[c]; !ok {
			b.state.Terrain[c] = model.TerrainPlains
		}
	}
	return b.state
}

func ptr(b bool) *bool {
	return &b
}
