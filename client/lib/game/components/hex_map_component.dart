import 'package:flame/components.dart';
import '../../models/enums.dart';
import '../hex/cube_coord.dart';
import '../hex/hex_layout.dart';
import 'hex_tile_component.dart';

class HexMapComponent extends Component {
  final HexLayout layout;
  final Map<CubeCoord, HexTileComponent> tiles = {};

  HexMapComponent(this.layout);

  void updateTerrain(Map<CubeCoord, TerrainType> terrain) {
    for (final entry in terrain.entries) {
      if (!tiles.containsKey(entry.key)) {
        final tile = HexTileComponent(
          coord: entry.key,
          terrain: entry.value,
          layout: layout,
        );
        tiles[entry.key] = tile;
        add(tile);
      }
    }
  }

  void updateHighlights(Set<CubeCoord> moves, Set<CubeCoord> attacks) {
    for (final entry in tiles.entries) {
      final tile = entry.value;
      if (attacks.contains(entry.key)) {
        tile.highlight = HighlightType.attack;
      } else if (moves.contains(entry.key)) {
        tile.highlight = HighlightType.move;
      } else {
        tile.highlight = HighlightType.none;
      }
    }
  }

  void clearHighlights() {
    for (final tile in tiles.values) {
      tile.highlight = HighlightType.none;
    }
  }
}
