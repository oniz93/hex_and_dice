import 'package:riverpod_annotation/riverpod_annotation.dart';
import '../models/enums.dart';
import '../models/game_state.dart';
import '../game/hex/cube_coord.dart';
import '../game/hex/pathfinding.dart';
import 'game_state_provider.dart';
import 'core_providers.dart';

part 'selection_provider.g.dart';

enum SelectionFSM { idle, troopSelected, confirmAttack, structureSelected }

class SelectionState {
  final SelectionFSM state;
  final String? selectedUnitId;
  final CubeCoord? targetHex;
  final Set<CubeCoord> highlightedMoves;
  final Set<CubeCoord> highlightedAttacks;

  const SelectionState({
    this.state = SelectionFSM.idle,
    this.selectedUnitId,
    this.targetHex,
    this.highlightedMoves = const {},
    this.highlightedAttacks = const {},
  });

  SelectionState copyWith({
    SelectionFSM? state,
    String? selectedUnitId,
    CubeCoord? targetHex,
    Set<CubeCoord>? highlightedMoves,
    Set<CubeCoord>? highlightedAttacks,
  }) {
    return SelectionState(
      state: state ?? this.state,
      selectedUnitId: selectedUnitId ?? this.selectedUnitId,
      targetHex: targetHex ?? this.targetHex,
      highlightedMoves: highlightedMoves ?? this.highlightedMoves,
      highlightedAttacks: highlightedAttacks ?? this.highlightedAttacks,
    );
  }
}

@Riverpod(keepAlive: true)
class SelectionStateNotifier extends _$SelectionStateNotifier {
  @override
  SelectionState build() {
    return const SelectionState();
  }

  void handleHexTap(CubeCoord hex, String playerId) {
    final gameState = ref.read(gameStateNotifierProvider);
    if (gameState == null) return;

    if (!gameState.isActivePlayer(playerId)) {
      clearSelection();
      return; // Cannot interact when not active player
    }

    final troop = gameState.troopAt(hex);
    final structure = gameState.structureAt(hex);

    switch (state.state) {
      case SelectionFSM.idle:
        if (troop != null && troop.ownerId == playerId) {
          _selectTroop(troop.id, gameState, playerId);
        } else if (structure != null && structure.ownerId == playerId) {
          state = state.copyWith(
            state: SelectionFSM.structureSelected,
            targetHex: hex,
          );
        } else {
          clearSelection();
        }
        break;

      case SelectionFSM.structureSelected:
        if (structure != null && structure.ownerId == playerId) {
          state = state.copyWith(targetHex: hex);
        } else if (troop != null && troop.ownerId == playerId) {
          _selectTroop(troop.id, gameState, playerId);
        } else {
          clearSelection();
        }
        break;

      case SelectionFSM.troopSelected:
        if (state.highlightedMoves.contains(hex)) {
          _sendMove(state.selectedUnitId!, hex);
          clearSelection();
        } else if (state.highlightedAttacks.contains(hex)) {
          state = state.copyWith(
            state: SelectionFSM.confirmAttack,
            targetHex: hex,
          );
        } else if (troop != null && troop.ownerId == playerId) {
          _selectTroop(troop.id, gameState, playerId);
        } else {
          clearSelection();
        }
        break;

      case SelectionFSM.confirmAttack:
        if (hex == state.targetHex) {
          _sendAttack(state.selectedUnitId!, hex);
        }
        clearSelection();
        break;
    }
  }

  void _sendMove(String unitId, CubeCoord target) {
    ref.read(wsServiceProvider).sendMove(unitId, target);
  }

  void _sendAttack(String unitId, CubeCoord target) {
    ref.read(wsServiceProvider).sendAttack(unitId, target);
  }

  void _selectTroop(String unitId, GameState gameState, String playerId) {
    final troop = gameState.troops[unitId];
    if (troop == null) return;

    Set<CubeCoord> moves = {};
    Set<CubeCoord> attacks = {};

    if (troop.canMove) {
      moves = Pathfinding.reachableHexes(
        troop.hex,
        troop.remainingMobility,
        gameState,
        playerId,
      );
    }

    if (troop.canAttack) {
      // Highlight all hexes in range that have an enemy
      final inRange = troop.hex.spiral(troop.range);
      for (final h in inRange) {
        if (h == troop.hex) continue; // Cannot attack self

        final targetTroop = gameState.troopAt(h);
        final targetStructure = gameState.structureAt(h);

        final isEnemyTroop =
            targetTroop != null && targetTroop.ownerId != playerId;
        final isEnemyStructure =
            targetStructure != null && targetStructure.ownerId != playerId;

        if (isEnemyTroop || isEnemyStructure) {
          attacks.add(h);
        }
      }
    }

    state = SelectionState(
      state: SelectionFSM.troopSelected,
      selectedUnitId: unitId,
      highlightedMoves: moves,
      highlightedAttacks: attacks,
    );
  }

  void clearSelection() {
    state = const SelectionState();
  }
}
