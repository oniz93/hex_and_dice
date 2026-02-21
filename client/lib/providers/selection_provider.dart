import 'package:riverpod_annotation/riverpod_annotation.dart';
import '../models/enums.dart';
import '../models/game_state.dart';
import '../game/hex/cube_coord.dart';
import '../game/hex/pathfinding.dart';
import 'game_state_provider.dart';

part 'selection_provider.g.dart';

enum SelectionFSM { idle, troopSelected, confirmAttack }

class SelectionState {
  final SelectionFSM state;
  final String? selectedUnitId;
  final Set<CubeCoord> highlightedMoves;
  final Set<CubeCoord> highlightedAttacks;

  const SelectionState({
    this.state = SelectionFSM.idle,
    this.selectedUnitId,
    this.highlightedMoves = const {},
    this.highlightedAttacks = const {},
  });

  SelectionState copyWith({
    SelectionFSM? state,
    String? selectedUnitId,
    Set<CubeCoord>? highlightedMoves,
    Set<CubeCoord>? highlightedAttacks,
  }) {
    return SelectionState(
      state: state ?? this.state,
      selectedUnitId: selectedUnitId ?? this.selectedUnitId,
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

    switch (state.state) {
      case SelectionFSM.idle:
        if (troop != null && troop.ownerId == playerId) {
          _selectTroop(troop.id, gameState, playerId);
        } else {
          clearSelection();
        }
        break;

      case SelectionFSM.troopSelected:
        if (state.highlightedMoves.contains(hex)) {
          clearSelection();
        } else if (state.highlightedAttacks.contains(hex)) {
          state = state.copyWith(state: SelectionFSM.confirmAttack);
        } else if (troop != null && troop.ownerId == playerId) {
          _selectTroop(troop.id, gameState, playerId);
        } else {
          clearSelection();
        }
        break;

      case SelectionFSM.confirmAttack:
        clearSelection();
        break;
    }
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
      final ring = troop.hex.spiral(troop.range);
      for (final h in ring) {
        if (gameState.structureAt(h)?.ownerId != playerId ||
            gameState.troopAt(h)?.ownerId != playerId) {
          attacks.add(h);
        }
      }
      attacks.remove(troop.hex);
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
