import 'dart:async';
import 'package:riverpod_annotation/riverpod_annotation.dart';
import '../models/game_state.dart';
import '../models/enums.dart';
import '../models/messages.dart';
import '../models/troop.dart';
import '../models/structure.dart';
import '../game/hex/cube_coord.dart';
import '../game/data/balance.dart';
import '../services/ws_service.dart';
import '../services/message_parser.dart';
import 'core_providers.dart';
import 'combat_log_provider.dart';

part 'game_state_provider.g.dart';

@Riverpod(keepAlive: true)
class GameStateNotifier extends _$GameStateNotifier {
  StreamSubscription<ServerMessage>? _sub;

  /// Callback to fire an attack arrow on the game canvas.
  /// Set by the GameScreen to call HexGame.showAttackArrow().
  Function(CubeCoord from, CubeCoord to)? _onAttackArrow;

  void setAttackArrowCallback(Function(CubeCoord from, CubeCoord to) cb) {
    _onAttackArrow = cb;
  }

  /// Look up a player's nickname from their UUID.
  String _playerName(String? ownerId) {
    if (state == null || ownerId == null || ownerId.isEmpty) return 'Neutral';
    try {
      return state!.players.firstWhere((p) => p.id == ownerId).nickname;
    } catch (_) {
      return 'Unknown';
    }
  }

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
        _handleCombatResult(msg.data as CombatResultData);
        break;
      case 'structure_attacked':
        _handleStructureAttacked(msg.data as StructureAttackedData);
        break;
      case 'structure_fires':
        _handleStructureFires(msg.data as StructureFiresData);
        break;
      case 'troop_destroyed':
        _handleTroopDestroyed(msg.data as TroopDestroyedData);
        break;
      case 'turn_start':
        _handleTurnStart(msg.data as TurnStartData);
        break;
      case 'game_over':
        _handleGameOver(msg.data as GameOverData);
        break;
    }
  }

  void _handleGameOver(GameOverData data) {
    if (state == null) return;
    state = state!.copyWith(
      phase: GamePhase.gameOver,
      gameOverData: data,
    );
    print('GameStateProvider: Game Over! Winner: ${data.winnerId}');
  }

  void _handleTroopMoved(TroopMovedData data) {
    if (state == null) return;

    final tMap = Map<String, Troop>.from(state!.troops);
    final troop = tMap[data.unitId];
    if (troop != null) {
      tMap[data.unitId] = troop.copyWith(
        hex: CubeCoord(data.toQ, data.toR, data.toS),
        hasMoved: true,
        remainingMobility: data.remainingMobility,
      );
      state = state!.copyWith(troops: tMap);
      print(
          'GameStateProvider: Troop ${data.unitId} moved to (${data.toQ}, ${data.toR})');
    }
  }

  void _handleTroopPurchased(TroopPurchasedData data) {
    if (state == null) return;

    final stats = troopStats[data.unitType]!;
    final newTroop = Troop(
      id: data.unitId,
      type: data.unitType,
      ownerId: data.owner,
      hex: CubeCoord(data.hexQ, data.hexR, data.hexS),
      currentHp: stats.hp,
      maxHp: stats.hp,
      atk: stats.atk,
      def: stats.def,
      mobility: stats.mobility,
      range: stats.range,
      damage: stats.damage,
      isReady: false,
      hasMoved: false,
      hasAttacked: false,
      wasInCombat: false,
      remainingMobility: 0,
    );

    final newTroops = Map<String, Troop>.from(state!.troops);
    newTroops[data.unitId] = newTroop;

    final newPlayers = state!.players.map((p) {
      if (p.id == data.owner) {
        return p.copyWith(coins: data.coinsRemaining);
      }
      return p;
    }).toList();

    state = state!.copyWith(
      troops: newTroops,
      players: newPlayers,
    );

    print(
        'GameStateProvider: Troop purchased: ${data.unitId} at (${data.hexQ}, ${data.hexR}), coins left: ${data.coinsRemaining}');
  }

  void _handleCombatResult(CombatResultData data) {
    if (state == null) return;

    final tMap = Map<String, Troop>.from(state!.troops);
    final attacker = tMap[data.attackerId];
    final defender = tMap[data.defenderId];

    if (attacker != null) {
      tMap[data.attackerId] = attacker.copyWith(
        currentHp: data.attackerHp,
        hasAttacked: true,
        hasMoved: true,
      );
    }
    if (defender != null) {
      tMap[data.defenderId] = defender.copyWith(currentHp: data.defenderHp);
    }

    state = state!.copyWith(troops: tMap);

    // Log results with player nicknames
    final log = ref.read(combatLogProvider.notifier);
    final atkOwner = _playerName(attacker?.ownerId);
    final defOwner = _playerName(defender?.ownerId);
    final atkUnit = attacker?.type.name.toUpperCase() ?? 'Unit';
    final defUnit = defender?.type.name.toUpperCase() ?? 'Unit';

    // Fire attack arrow: attacker -> defender
    if (attacker != null && defender != null) {
      _onAttackArrow?.call(attacker.hex, defender.hex);
    }

    if (data.hit) {
      log.addEntry(
          '$atkOwner\'s $atkUnit hit $defOwner\'s $defUnit (Roll: ${data.hitRoll}, Dmg: ${data.damage})');
      if (data.killed) {
        log.addEntry('💀 $defOwner\'s $defUnit was destroyed!');
      }
    } else {
      log.addEntry(
          '$atkOwner\'s $atkUnit missed $defOwner\'s $defUnit (Roll: ${data.hitRoll})');
    }

    // Counterattack
    if (data.hasCounter) {
      // Fire counter-attack arrow: defender -> attacker
      if (defender != null && attacker != null) {
        _onAttackArrow?.call(defender.hex, attacker.hex);
      }

      if (data.counterHit == true) {
        log.addEntry(
            '⚡ $defOwner\'s $defUnit countered $atkOwner\'s $atkUnit! (Roll: ${data.counterHitRoll}, Dmg: ${data.counterDamage})');
        if (data.attackerKilled) {
          log.addEntry('💀 $atkOwner\'s $atkUnit was destroyed by counter!');
        }
      } else {
        log.addEntry(
            '⚡ $defOwner\'s $defUnit counter-attack missed $atkOwner\'s $atkUnit.');
      }
    }
  }

  void _handleStructureAttacked(StructureAttackedData data) {
    if (state == null) return;

    final sMap = Map<String, Structure>.from(state!.structures);
    final tMap = Map<String, Troop>.from(state!.troops);
    final structure = sMap[data.structureId];
    final attacker = tMap[data.attackerId];

    if (structure != null) {
      sMap[data.structureId] = structure.copyWith(
        currentHp: data.structureHp,
        ownerId: data.captured ? data.newOwner : structure.ownerId,
      );
    }
    if (attacker != null) {
      tMap[data.attackerId] = attacker.copyWith(
        hasAttacked: true,
        hasMoved: true,
      );
    }

    state = state!.copyWith(structures: sMap, troops: tMap);

    final log = ref.read(combatLogProvider.notifier);
    final atkOwner = _playerName(attacker?.ownerId);
    final atkUnit = attacker?.type.name.toUpperCase() ?? 'Unit';
    final structName = structure?.type.name.toUpperCase() ?? 'Structure';

    // Fire attack arrow: attacker troop -> structure
    if (attacker != null && structure != null) {
      _onAttackArrow?.call(attacker.hex, structure.hex);
    }

    log.addEntry(
        '$atkOwner\'s $atkUnit attacked $structName (Roll: ${data.hitRoll}, Dmg: ${data.damage})');
    if (data.captured) {
      final capturerName = _playerName(data.newOwner);
      log.addEntry('🚩 $structName was captured by $capturerName!');
    }
  }

  void _handleStructureFires(StructureFiresData data) {
    if (state == null) return;

    final sMap = Map<String, Structure>.from(state!.structures);
    final tMap = Map<String, Troop>.from(state!.troops);
    final structure = sMap[data.structureId];
    final target = tMap[data.targetId];

    // Apply damage to the target troop
    if (target != null) {
      tMap[data.targetId] = target.copyWith(currentHp: data.targetHp);
    }

    state = state!.copyWith(structures: sMap, troops: tMap);

    final log = ref.read(combatLogProvider.notifier);
    final structOwner = _playerName(structure?.ownerId);
    final structName = structure?.type.name.toUpperCase() ?? 'Structure';
    final targetOwner = _playerName(target?.ownerId);
    final targetUnit = target?.type.name.toUpperCase() ?? 'Unit';

    // Fire attack arrow: structure -> target troop
    if (structure != null && target != null) {
      _onAttackArrow?.call(structure.hex, target.hex);
    }

    if (data.damage > 0) {
      log.addEntry(
          '$structOwner\'s $structName fired at $targetOwner\'s $targetUnit (Roll: ${data.hitRoll}, Dmg: ${data.damage})');
    } else {
      log.addEntry(
          '$structOwner\'s $structName missed $targetOwner\'s $targetUnit (Roll: ${data.hitRoll})');
    }
    if (data.killed) {
      log.addEntry(
          '💀 $targetOwner\'s $targetUnit was destroyed by $structName!');
    }
  }

  void _handleTroopDestroyed(TroopDestroyedData data) {
    if (state == null) return;
    final tMap = Map.of(state!.troops);
    tMap.remove(data.unitId);
    state = state!.copyWith(troops: tMap);
  }

  void _handleTurnStart(TurnStartData data) {
    if (state == null) return;

    final players = state!.players.map((p) {
      if (p.id == data.activePlayerId) {
        return p.copyWith(
          coins: data.totalCoins,
          // income is usually calculated from structures, but we can update it if needed
        );
      }
      return p;
    }).toList();

    final activePlayerIndex =
        state!.players.indexWhere((p) => p.id == data.activePlayerId);

    // Reset troops for the active player
    final tMap = state!.troops.map((id, troop) {
      if (troop.ownerId == data.activePlayerId) {
        return MapEntry(
            id,
            troop.copyWith(
              isReady: true,
              hasMoved: false,
              hasAttacked: false,
              remainingMobility: troop.mobility,
            ));
      }
      return MapEntry(id, troop);
    });

    state = state!.copyWith(
      turnNumber: data.turnNumber,
      activePlayer:
          activePlayerIndex != -1 ? activePlayerIndex : state!.activePlayer,
      turnTimer: data.timerSeconds,
      turnStartedAt: DateTime.now().toUtc(),
      players: players,
      troops: tMap,
    );

    print(
        'GameStateProvider: Turn started for ${data.activePlayerId}, timer: ${data.timerSeconds}s');
  }
}
