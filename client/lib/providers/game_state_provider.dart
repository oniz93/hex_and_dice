import 'dart:async';
import 'package:riverpod_annotation/riverpod_annotation.dart';
import '../models/game_state.dart';
import '../models/messages.dart';
import '../services/ws_service.dart';
import '../services/message_parser.dart';
import 'core_providers.dart';

part 'game_state_provider.g.dart';

@Riverpod(keepAlive: true)
class GameStateNotifier extends _$GameStateNotifier {
  StreamSubscription<ServerMessage>? _sub;

  @override
  GameState? build() {
    final wsService = ref.watch(wsServiceProvider);

    _sub?.cancel();
    _sub = wsService.messages.listen((msg) {
      final parsed = parseMessage(msg);
      if (parsed != null) {
        _handleMessage(parsed);
      }
    });

    ref.onDispose(() {
      _sub?.cancel();
    });

    return null;
  }

  void _handleMessage(ParsedMessage msg) {
    switch (msg.type) {
      case 'game_state':
        state = msg.data as GameState;
        break;
      // Handle deltas
      case 'troop_moved':
        _handleTroopMoved(msg.data as TroopMovedData);
        break;
      case 'troop_purchased':
        _handleTroopPurchased(msg.data as TroopPurchasedData);
        break;
      case 'combat_result':
        // Wait for animation queue logic, but update state directly for now
        break;
      case 'troop_destroyed':
        _handleTroopDestroyed(msg.data as TroopDestroyedData);
        break;
      case 'turn_start':
        _handleTurnStart(msg.data as TurnStartData);
        break;
      case 'game_over':
        // Update game state
        break;
    }
  }

  void _handleTroopMoved(TroopMovedData data) {
    if (state == null) return;

    final tMap = Map.of(state!.troops);
    final troop = tMap[data.unitId];
    if (troop != null) {
      // Create new troop position
      // In Dart, we copy it
      // tMap[data.unitId] = troop.copyWith(hex: CubeCoord(data.toQ, data.toR, data.toS), hasMoved: true, remainingMobility: data.remainingMobility);
    }
    // state = state!.copyWith(troops: tMap);
  }

  void _handleTroopPurchased(TroopPurchasedData data) {
    //
  }

  void _handleTroopDestroyed(TroopDestroyedData data) {
    if (state == null) return;
    final tMap = Map.of(state!.troops);
    tMap.remove(data.unitId);
    state = state!.copyWith(troops: tMap);
  }

  void _handleTurnStart(TurnStartData data) {
    if (state == null) return;
    // state = state!.copyWith(turnNumber: data.turnNumber, ...);
  }
}
