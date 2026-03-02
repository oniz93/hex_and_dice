import 'package:flutter/material.dart' hide Animation;
import 'package:flame/game.dart';
import 'package:flame/components.dart';
import 'package:flame/events.dart';
import '../models/game_state.dart';
import '../models/enums.dart';
import 'components/hex_map_component.dart';
import 'components/troop_component.dart';
import 'components/structure_component.dart';
import 'components/attack_arrow_component.dart';
import 'hex/cube_coord.dart';
import 'hex/hex_layout.dart';

class HexGame extends FlameGame with TapCallbacks, ScaleDetector {
  late HexMapComponent hexMap;
  final HexLayout layout =
      const HexLayout(32.0); // 64px hex width (2 * hexSize)

  final Map<TerrainType, Sprite> tileSprites = {};
  final Map<String, TroopComponent> _troopComponents = {};
  final Map<String, StructureComponent> _structureComponents = {};

  bool _loaded = false;
  GameState? _pendingState;
  final Set<String> _animatingEntities = {};

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
    // Load sprites
    tileSprites[TerrainType.plains] = await loadSprite('sprites/plains.png');
    tileSprites[TerrainType.forest] = await loadSprite('sprites/forest.png');
    tileSprites[TerrainType.water] = await loadSprite('sprites/water.png');
    tileSprites[TerrainType.mountains] =
        await loadSprite('sprites/mountain.png');
    // Hills fallback to plains for now as requested
    tileSprites[TerrainType.hills] = tileSprites[TerrainType.plains]!;

    hexMap = HexMapComponent(layout, tileSprites);
    world.add(hexMap);

    _loaded = true;

    // If a game state arrived before onLoad finished, apply it now
    if (_pendingState != null) {
      updateGameState(_pendingState!);
      _pendingState = null;
    }
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
    // Buffer state if the game engine hasn't finished loading yet
    if (!_loaded) {
      _pendingState = state;
      return;
    }

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
        // Don't update structures that are mid-animation
        if (!_animatingEntities.contains(s.id)) {
          existing.updateStructure(s, color);
        }
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
        // Don't update troops that are mid-animation
        if (!_animatingEntities.contains(t.id)) {
          existing.updateTroop(t, color);
        }
      }
    }

    // Remove dead troops -- but skip those currently animating (flash in progress)
    final toRemove = _troopComponents.keys
        .where((id) =>
            !state.troops.containsKey(id) && !_animatingEntities.contains(id))
        .toList();
    for (final id in toRemove) {
      _troopComponents[id]?.removeFromParent();
      _troopComponents.remove(id);
    }
  }

  void updateSelection(Set<CubeCoord> moves, Set<CubeCoord> attacks) {
    hexMap.updateHighlights(moves, attacks);
  }

  /// Shows an animated projectile arrow from [fromHex] to [toHex].
  /// When the projectile impacts, if [hit] is true, the target entity flashes
  /// red twice. If [killed] is true the troop is removed after the flash.
  /// If [captured] is true the structure owner changes after the flash.
  void showAttackArrow(
    CubeCoord fromHex,
    CubeCoord toHex, {
    required String targetId,
    required bool isStructure,
    required bool hit,
    required bool killed,
    required bool captured,
    String? newOwner,
  }) {
    final fromPixel = layout.hexToPixel(fromHex);
    final toPixel = layout.hexToPixel(toHex);

    // If the target will be visually affected, mark it as animating so
    // updateGameState won't remove/update it mid-flight.
    if (hit && (killed || captured)) {
      _animatingEntities.add(targetId);
    }

    final arrow = AttackArrowComponent(
      from: Vector2(fromPixel.dx, fromPixel.dy),
      to: Vector2(toPixel.dx, toPixel.dy),
      onImpact: () {
        if (!hit) {
          // Miss -- nothing to flash, just clean up
          _animatingEntities.remove(targetId);
          return;
        }

        if (isStructure) {
          _handleStructureImpact(targetId, killed, captured, newOwner);
        } else {
          _handleTroopImpact(targetId, killed);
        }
      },
    );
    world.add(arrow);
  }

  void _handleTroopImpact(String troopId, bool killed) {
    final tc = _troopComponents[troopId];
    if (tc == null) {
      _animatingEntities.remove(troopId);
      return;
    }

    tc.startFlash(() {
      _animatingEntities.remove(troopId);
      if (killed) {
        tc.removeFromParent();
        _troopComponents.remove(troopId);
      } else {
        // Re-apply latest game state to this troop if it still exists
        if (gameState != null && gameState!.troops.containsKey(troopId)) {
          final t = gameState!.troops[troopId]!;
          tc.updateTroop(t, getPlayerColor(t.ownerId));
        }
      }
    });
  }

  void _handleStructureImpact(
      String structureId, bool killed, bool captured, String? newOwner) {
    final sc = _structureComponents[structureId];
    if (sc == null) {
      _animatingEntities.remove(structureId);
      return;
    }

    sc.startFlash(() {
      _animatingEntities.remove(structureId);
      // Re-apply latest game state (updated owner, HP, etc.)
      if (gameState != null && gameState!.structures.containsKey(structureId)) {
        final s = gameState!.structures[structureId]!;
        sc.updateStructure(s, getPlayerColor(s.ownerId));
      }
    });
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

  bool _isScaling = false;

  @override
  void onScaleStart(ScaleStartInfo info) {
    _isScaling = false;
  }

  @override
  void onScaleUpdate(ScaleUpdateInfo info) {
    // Detect pinch zoom (two fingers)
    final scaleDelta = info.scale.global.x;
    if ((scaleDelta - 1.0).abs() > 0.01) {
      _isScaling = true;
      var zoom = camera.viewfinder.zoom;
      // Dampen the scale factor to reduce sensitivity
      final dampened = 1.0 + (scaleDelta - 1.0) * 0.3;
      zoom *= dampened;
      zoom = zoom.clamp(0.1, 3.0);
      camera.viewfinder.zoom = zoom;
    }

    // Handle pan (single finger drag)
    if (!_isScaling) {
      final delta = info.delta.global;
      camera.viewfinder.position -= delta / camera.viewfinder.zoom;
    }
  }

  @override
  void onScaleEnd(ScaleEndInfo info) {
    _isScaling = false;
  }
}
