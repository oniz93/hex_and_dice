package game

import (
	"encoding/json"
	"time"

	"github.com/teomiscia/hexbattle/internal/hex"
	"github.com/teomiscia/hexbattle/internal/model"
)

// GameState holds the complete state of an active game.
type GameState struct {
	ID            string                          `json:"id"`
	Phase         model.GamePhase                 `json:"phase"`
	MapSize       model.MapSize                   `json:"map_size"`
	TurnMode      model.TurnMode                  `json:"turn_mode"`
	TurnTimer     int                             `json:"turn_timer"` // seconds per turn
	TurnNumber    int                             `json:"turn_number"`
	ActivePlayer  int                             `json:"active_player"` // 0 or 1 (index into Players)
	Players       [2]model.PlayerState            `json:"players"`
	Troops        map[string]*model.Troop         `json:"troops"`     // unit_id -> troop
	Structures    map[string]*model.Structure     `json:"structures"` // structure_id -> structure
	Terrain       map[hex.Coord]model.TerrainType `json:"terrain"`    // hex -> terrain type
	Grid          *hex.Grid                       `json:"-"`          // not serialized, rebuilt from MapSize
	Seed          int64                           `json:"seed"`
	CreatedAt     time.Time                       `json:"created_at"`
	TurnStartedAt time.Time                       `json:"turn_started_at"`

	// Sudden death state
	SuddenDeathActive bool `json:"sudden_death_active"`
	SuddenDeathTurn   int  `json:"sudden_death_turn"`
	SafeZoneRadius    int  `json:"safe_zone_radius"`

	// Per-player stats tracked during the game
	Stats [2]model.GameOverStats `json:"stats"`

	// First turn restriction: player 0 cannot attack on turn 1
	FirstTurnRestriction bool `json:"first_turn_restriction"`
}

// NewGameState creates an empty game state ready for map generation.
func NewGameState(id string, settings model.RoomSettings, p1, p2 model.PlayerState, seed int64) *GameState {
	p1.Coins = StartingCoins()
	p2.Coins = StartingCoins()

	return &GameState{
		ID:                   id,
		Phase:                model.PhaseWaitingForPlayers,
		MapSize:              settings.MapSize,
		TurnMode:             settings.TurnMode,
		TurnTimer:            settings.TurnTimer,
		TurnNumber:           0,
		ActivePlayer:         0,
		Players:              [2]model.PlayerState{p1, p2},
		Troops:               make(map[string]*model.Troop),
		Structures:           make(map[string]*model.Structure),
		Terrain:              make(map[hex.Coord]model.TerrainType),
		Grid:                 hex.NewGrid(settings.MapSize.Radius()),
		Seed:                 seed,
		CreatedAt:            time.Now(),
		SuddenDeathActive:    false,
		SafeZoneRadius:       settings.MapSize.Radius(),
		FirstTurnRestriction: true,
	}
}

// ActivePlayerState returns the state of the player whose turn it is.
func (gs *GameState) ActivePlayerState() *model.PlayerState {
	return &gs.Players[gs.ActivePlayer]
}

// InactivePlayerState returns the state of the player who is waiting.
func (gs *GameState) InactivePlayerState() *model.PlayerState {
	return &gs.Players[1-gs.ActivePlayer]
}

// ActivePlayerID returns the ID of the active player.
func (gs *GameState) ActivePlayerID() string {
	return gs.Players[gs.ActivePlayer].ID
}

// InactivePlayerID returns the ID of the inactive player.
func (gs *GameState) InactivePlayerID() string {
	return gs.Players[1-gs.ActivePlayer].ID
}

// PlayerIndex returns 0 or 1 for the given player ID, or -1 if not found.
func (gs *GameState) PlayerIndex(playerID string) int {
	for i, p := range gs.Players {
		if p.ID == playerID {
			return i
		}
	}
	return -1
}

// IsActivePlayer returns true if the given player ID is the active player.
func (gs *GameState) IsActivePlayer(playerID string) bool {
	return gs.ActivePlayerID() == playerID
}

// GetTroop returns a troop by ID, or nil if not found.
func (gs *GameState) GetTroop(unitID string) *model.Troop {
	return gs.Troops[unitID]
}

// GetStructure returns a structure by ID, or nil if not found.
func (gs *GameState) GetStructure(structureID string) *model.Structure {
	return gs.Structures[structureID]
}

// TroopAtHex returns the troop at the given hex, or nil if empty.
func (gs *GameState) TroopAtHex(pos hex.Coord) *model.Troop {
	for _, t := range gs.Troops {
		if t.Hex == pos && t.IsAlive() {
			return t
		}
	}
	return nil
}

// StructureAtHex returns the structure at the given hex, or nil if none.
func (gs *GameState) StructureAtHex(pos hex.Coord) *model.Structure {
	for _, s := range gs.Structures {
		if s.Hex == pos {
			return s
		}
	}
	return nil
}

// PlayerTroops returns all living troops belonging to a player.
func (gs *GameState) PlayerTroops(playerID string) []*model.Troop {
	var troops []*model.Troop
	for _, t := range gs.Troops {
		if t.OwnerID == playerID && t.IsAlive() {
			troops = append(troops, t)
		}
	}
	return troops
}

// PlayerStructures returns all structures owned by a player.
func (gs *GameState) PlayerStructures(playerID string) []*model.Structure {
	var structs []*model.Structure
	for _, s := range gs.Structures {
		if s.OwnerID == playerID {
			structs = append(structs, s)
		}
	}
	return structs
}

// AllStructures returns all structures on the map.
func (gs *GameState) AllStructures() []*model.Structure {
	structs := make([]*model.Structure, 0, len(gs.Structures))
	for _, s := range gs.Structures {
		structs = append(structs, s)
	}
	return structs
}

// PlayerHQ returns the HQ structure for the given player, or nil.
func (gs *GameState) PlayerHQ(playerID string) *model.Structure {
	for _, s := range gs.Structures {
		if s.Type == model.StructureHQ && s.OwnerID == playerID {
			return s
		}
	}
	return nil
}

// RemoveTroop removes a dead troop from the game state.
func (gs *GameState) RemoveTroop(unitID string) {
	delete(gs.Troops, unitID)
}

// AddTroop adds a troop to the game state.
func (gs *GameState) AddTroop(troop *model.Troop) {
	gs.Troops[troop.ID] = troop
}

// AddStructure adds a structure to the game state.
func (gs *GameState) AddStructure(s *model.Structure) {
	gs.Structures[s.ID] = s
}

// GetTerrainAt returns the terrain type at the given hex.
func (gs *GameState) GetTerrainAt(pos hex.Coord) model.TerrainType {
	t, ok := gs.Terrain[pos]
	if !ok {
		return model.TerrainPlains // default
	}
	return t
}

// IsHexPassable returns true if the hex has passable terrain and is in bounds.
func (gs *GameState) IsHexPassable(pos hex.Coord) bool {
	if !gs.Grid.Contains(pos) {
		return false
	}
	return model.IsPassable(gs.GetTerrainAt(pos))
}

// IsHexOccupiedByEnemy returns true if an enemy troop is on the hex.
func (gs *GameState) IsHexOccupiedByEnemy(pos hex.Coord, playerID string) bool {
	t := gs.TroopAtHex(pos)
	return t != nil && t.OwnerID != playerID
}

// SwitchActivePlayer toggles the active player index.
func (gs *GameState) SwitchActivePlayer() {
	gs.ActivePlayer = 1 - gs.ActivePlayer
}

// Serialize converts the game state to JSON bytes for persistence.
func (gs *GameState) Serialize() ([]byte, error) {
	return json.Marshal(gs)
}

// DeserializeGameState restores a game state from JSON bytes.
func DeserializeGameState(data []byte) (*GameState, error) {
	var gs GameState
	if err := json.Unmarshal(data, &gs); err != nil {
		return nil, err
	}
	// Rebuild non-serialized fields
	gs.Grid = hex.NewGrid(gs.MapSize.Radius())
	return &gs, nil
}

// StructureCountOwnedBy returns how many structures the player owns.
func (gs *GameState) StructureCountOwnedBy(playerID string) int {
	count := 0
	for _, s := range gs.Structures {
		if s.OwnerID == playerID {
			count++
		}
	}
	return count
}

// TotalStructureCount returns the total number of structures on the map.
func (gs *GameState) TotalStructureCount() int {
	return len(gs.Structures)
}
