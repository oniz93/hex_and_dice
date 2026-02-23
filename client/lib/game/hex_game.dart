import 'package:flutter/material.dart' hide Animation;
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

  Color getPlayerColor(String? playerId) {
    if (playerId == null || playerId.isEmpty) {
      return Colors.black;
    }
    if (gameState == null) return Colors.black;

    final index = gameState!.players.indexWhere((p) => p.id == playerId);
    if (index == 0) return Colors.red;
    if (index == 1) return Colors.blue;
    return Colors.black;
  }

  void updateGameState(GameState state) {
    gameState = state;
    hexMap.updateTerrain(state.terrain);

    // Update structures
    for (final s in state.structures.values) {
      final existing = _structureComponents[s.id];
      final color = getPlayerColor(s.ownerId);
      if (existing == null) {
        final sc = StructureComponent(
          structure: s,
          layout: layout,
          teamColor: color,
        );
        _structureComponents[s.id] = sc;
        world.add(sc);
      } else {
        existing.updateStructure(s, color);
      }
    }

    // Update troops
    for (final t in state.troops.values) {
      final existing = _troopComponents[t.id];
      final color = getPlayerColor(t.ownerId);
      if (existing == null) {
        final tc = TroopComponent(
          troop: t,
          layout: layout,
          teamColor: color,
        );
        _troopComponents[t.id] = tc;
        world.add(tc);
      } else {
        existing.updateTroop(t, color);
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
