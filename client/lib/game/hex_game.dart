import 'dart:ui';
import 'package:flame/game.dart';
import 'package:flame/components.dart';
import 'package:flame/events.dart';
import '../models/game_state.dart';
import 'components/hex_map_component.dart';
import 'components/troop_component.dart';
import 'components/structure_component.dart';
import 'hex/cube_coord.dart';
import 'hex/hex_layout.dart';

class HexGame extends FlameGame with TapCallbacks, PanDetector, ScrollDetector {
  late HexMapComponent hexMap;
  final HexLayout layout = const HexLayout(24.0); // 48px hexes

  final Map<String, TroopComponent> _troopComponents = {};
  final Map<String, StructureComponent> _structureComponents = {};

  // Set from Flutter
  GameState? gameState;
  Function(CubeCoord)? onHexTap;

  HexGame() : super() {
    // Center the camera at 0,0
    camera.viewfinder.position = Vector2.zero();
    camera.viewfinder.zoom = 1.0;
  }

  @override
  Future<void> onLoad() async {
    hexMap = HexMapComponent(layout);
    world.add(hexMap);
  }

  void updateGameState(GameState state) {
    gameState = state;
    hexMap.updateTerrain(state.terrain);

    // Update structures
    for (final s in state.structures.values) {
      if (!_structureComponents.containsKey(s.id)) {
        final sc = StructureComponent(structure: s, layout: layout);
        _structureComponents[s.id] = sc;
        world.add(sc);
      }
    }

    // Update troops
    for (final t in state.troops.values) {
      if (!_troopComponents.containsKey(t.id)) {
        final tc = TroopComponent(troop: t, layout: layout);
        _troopComponents[t.id] = tc;
        world.add(tc);
      }
    }

    // Remove dead troops
    final toRemove = _troopComponents.keys
        .where((id) => !state.troops.containsKey(id))
        .toList();
    for (final id in toRemove) {
      _troopComponents[id]?.removeFromParent();
      _troopComponents.remove(id);
    }
  }

  void updateSelection(Set<CubeCoord> moves, Set<CubeCoord> attacks) {
    hexMap.updateHighlights(moves, attacks);
  }

  @override
  void onTapDown(TapDownEvent event) {
    if (onHexTap != null) {
      // Need to convert global screen tap to world coordinates
      final pos = camera.globalToLocal(event.canvasPosition);
      final hex = layout.pixelToHex(Offset(pos.x, pos.y));
      onHexTap!(hex);
    }
    super.onTapDown(event);
  }

  @override
  void onPanUpdate(DragUpdateInfo info) {
    camera.viewfinder.position -= info.delta.global;
  }

  @override
  void onScroll(PointerScrollInfo info) {
    var zoom = camera.viewfinder.zoom;
    zoom += info.scrollDelta.global.y > 0 ? -0.1 : 0.1;
    zoom = zoom.clamp(0.2, 3.0);
    camera.viewfinder.zoom = zoom;
  }
}
