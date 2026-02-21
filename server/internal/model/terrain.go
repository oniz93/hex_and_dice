package model

// TerrainInfo holds the gameplay properties of a terrain type.
type TerrainInfo struct {
	Type         TerrainType `json:"type"`
	MovementCost int         `json:"movement_cost"` // 0 means impassable
	ATKModifier  int         `json:"atk_modifier"`
	DEFModifier  int         `json:"def_modifier"`
	Passable     bool        `json:"passable"`
}

// TerrainTable maps terrain types to their properties.
var TerrainTable = map[TerrainType]TerrainInfo{
	TerrainPlains: {
		Type:         TerrainPlains,
		MovementCost: 1,
		ATKModifier:  0,
		DEFModifier:  0,
		Passable:     true,
	},
	TerrainForest: {
		Type:         TerrainForest,
		MovementCost: 2,
		ATKModifier:  0,
		DEFModifier:  2,
		Passable:     true,
	},
	TerrainHills: {
		Type:         TerrainHills,
		MovementCost: 2,
		ATKModifier:  1,
		DEFModifier:  1,
		Passable:     true,
	},
	TerrainWater: {
		Type:         TerrainWater,
		MovementCost: 0,
		ATKModifier:  0,
		DEFModifier:  0,
		Passable:     false,
	},
	TerrainMountains: {
		Type:         TerrainMountains,
		MovementCost: 0,
		ATKModifier:  0,
		DEFModifier:  0,
		Passable:     false,
	},
}

// GetTerrainInfo returns the terrain info for the given type.
// Returns plains info as default for unknown types.
func GetTerrainInfo(t TerrainType) TerrainInfo {
	info, ok := TerrainTable[t]
	if !ok {
		return TerrainTable[TerrainPlains]
	}
	return info
}

// IsPassable returns whether the terrain type can be traversed.
func IsPassable(t TerrainType) bool {
	return GetTerrainInfo(t).Passable
}

// MovementCost returns the movement cost for the terrain type.
// Returns 0 for impassable terrain.
func MovementCost(t TerrainType) int {
	return GetTerrainInfo(t).MovementCost
}
