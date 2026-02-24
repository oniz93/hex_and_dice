import 'dart:collection';
import '../../models/game_state.dart';
import '../../models/enums.dart';
import 'cube_coord.dart';

class Pathfinding {
  static Set<CubeCoord> reachableHexes(
    CubeCoord start,
    int mobility,
    GameState state,
    String playerId,
  ) {
    final visited = <CubeCoord>{};
    final queue = Queue<_PathNode>();

    queue.add(_PathNode(start, mobility));
    visited.add(start);

    while (queue.isNotEmpty) {
      final current = queue.removeFirst();

      for (int i = 0; i < 6; i++) {
        final neighbor = _neighbor(current.coord, i);

        // check bounds / passable
        if (!state.terrain.containsKey(neighbor)) continue;
        final terrain = state.terrain[neighbor]!;
        var cost = _movementCost(terrain);
        if (cost == 0) continue; // Impassable

        // Check enemy troop (cannot pass through)
        final troopAtNeighbor = state.troopAt(neighbor);
        if (troopAtNeighbor != null && troopAtNeighbor.ownerId != playerId) {
          continue;
        }

        // Minimum movement rule: if adjacent to start and have mobility left,
        // always allow moving to at least one cell (unless impassable).
        if (current.coord == start &&
            mobility > 0 &&
            cost > current.remainingMobility) {
          cost = current.remainingMobility;
        }

        final remaining = current.remainingMobility - cost;
        if (remaining >= 0 && !visited.contains(neighbor)) {
          visited.add(neighbor);
          queue.add(_PathNode(neighbor, remaining));
        }
      }
    }

    // Exclude start position and occupied cells
    return visited.where((h) {
      if (h == start) return false;
      // Cannot end move on another troop or a structure
      if (state.troopAt(h) != null) return false;
      if (state.structureAt(h) != null) return false;
      return true;
    }).toSet();
  }

  static CubeCoord _neighbor(CubeCoord coord, int dir) {
    final dirs = [
      CubeCoord(1, 0, -1),
      CubeCoord(1, -1, 0),
      CubeCoord(0, -1, 1),
      CubeCoord(-1, 0, 1),
      CubeCoord(-1, 1, 0),
      CubeCoord(0, 1, -1),
    ];
    return coord + dirs[dir];
  }

  static int _movementCost(TerrainType t) {
    switch (t) {
      case TerrainType.plains:
        return 1;
      case TerrainType.forest:
        return 2;
      case TerrainType.hills:
        return 2;
      case TerrainType.water:
        return 0;
      case TerrainType.mountains:
        return 0;
    }
  }
}

class _PathNode {
  final CubeCoord coord;
  final int remainingMobility;
  _PathNode(this.coord, this.remainingMobility);
}
