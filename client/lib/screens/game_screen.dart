import 'package:flutter/material.dart';
import 'package:flame/game.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import '../../game/hex_game.dart';
import '../../providers/game_state_provider.dart';
import '../../providers/selection_provider.dart';
import '../../providers/session_provider.dart';
import '../../providers/settings_provider.dart';
import '../../providers/core_providers.dart';
import '../../models/enums.dart';
import '../widgets/hud/top_bar.dart';
import '../widgets/hud/bottom_bar.dart';
import '../widgets/hud/troop_popup.dart';
import '../widgets/hud/shop_panel.dart';
import '../widgets/hud/combat_log_overlay.dart';
import '../widgets/hud/game_over_overlay.dart';

class GameScreen extends ConsumerStatefulWidget {
  final String roomId;

  const GameScreen({super.key, required this.roomId});

  @override
  ConsumerState<GameScreen> createState() => _GameScreenState();
}

class _GameScreenState extends ConsumerState<GameScreen> {
  late HexGame game;
  bool _connected = false;

  @override
  void initState() {
    super.initState();
    game = HexGame();
    game.onHexTap = (hex) {
      final session = ref.read(sessionProviderProvider).value;
      if (session != null) {
        ref
            .read(selectionStateProvider.notifier)
            .handleHexTap(hex, session.id);
      }
    };

    // Wire up attack arrow callback so combat events show a projectile
    ref.read(gameStateProvider.notifier).setAttackArrowCallback(
        (from, to, targetId, isStructure, hit, killed, captured, newOwner) {
      game.showAttackArrow(from, to,
          targetId: targetId,
          isStructure: isStructure,
          hit: hit,
          killed: killed,
          captured: captured,
          newOwner: newOwner);
    });

    print('GameScreen: initState called for room ${widget.roomId}');

    WidgetsBinding.instance.addPostFrameCallback((_) {
      print('GameScreen: addPostFrameCallback running');
      _connectToGame();
    });
  }

  Future<void> _connectToGame() async {
    try {
      print('GameScreen: _connectToGame starting');
      final session = await ref.read(sessionProviderProvider.future);

      if (session == null) {
        print('GameScreen: ERROR - session is null!');
        return;
      }

      final wsService = ref.read(wsServiceProvider);
      final storedGameId = ref.read(settingsProvider).gameId;

      print(
          'GameScreen: Connecting to game with roomId: ${widget.roomId}, token: ${session.token.substring(0, 10)}...');
      await wsService.connect(session.token);

      if (storedGameId == widget.roomId) {
        print('GameScreen: Reconnecting to existing game ${widget.roomId}...');
        wsService.sendReconnect(widget.roomId, session.token);
      } else {
        print(
            'GameScreen: Joining room ${widget.roomId} for the first time...');
        wsService.sendJoinGame(widget.roomId);
        // Persist game ID for future reconnection
        await ref
            .read(settingsProvider.notifier)
            .setReconnectData(widget.roomId, session.token, session.id);
      }

      setState(() {
        _connected = true;
      });
    } catch (e, st) {
      print('GameScreen: ERROR in _connectToGame: $e');
      print(st);
    }
  }

  @override
  Widget build(BuildContext context) {
    // Listen to game state updates
    ref.listen(gameStateProvider, (prev, next) {
      if (next != null) {
        game.updateGameState(next);
      }
    });

    // Listen to selection updates
    ref.listen(selectionStateProvider, (prev, next) {
      game.updateSelection(next.highlightedMoves, next.highlightedAttacks);
    });

    final selection = ref.watch(selectionStateProvider);
    final gameState = ref.watch(gameStateProvider);

    return Scaffold(
      body: Stack(
        children: [
          GameWidget(game: game),
          const Positioned(top: 0, left: 0, right: 0, child: TopBar()),
          const Positioned(bottom: 0, left: 0, right: 0, child: BottomBar()),
          const Positioned(top: 60, left: 16, child: CombatLogOverlay()),
          if (selection.state == SelectionFSM.troopSelected)
            const Positioned(
              right: 16,
              top: 80,
              width: 200,
              child: TroopPopup(),
            ),
          if (selection.state == SelectionFSM.structureSelected)
            const Positioned(
              bottom: 80,
              left: 0,
              right: 0,
              child: ShopPanel(),
            ),
          if (selection.state == SelectionFSM.confirmAttack)
            Center(
              child: Card(
                color: Colors.black87,
                child: Padding(
                  padding: const EdgeInsets.all(16.0),
                  child: Column(
                    mainAxisSize: MainAxisSize.min,
                    children: [
                      const Text(
                        'CONFIRM ATTACK?',
                        style: TextStyle(
                            color: Colors.white, fontWeight: FontWeight.bold),
                      ),
                      const SizedBox(height: 16),
                      Row(
                        mainAxisSize: MainAxisSize.min,
                        children: [
                          TextButton(
                            onPressed: () => ref
                                .read(selectionStateProvider.notifier)
                                .clearSelection(),
                            child: const Text('CANCEL',
                                style: TextStyle(color: Colors.red)),
                          ),
                          const SizedBox(width: 16),
                          ElevatedButton(
                            onPressed: () {
                              ref
                                  .read(selectionStateProvider.notifier)
                                  .handleHexTap(
                                      selection.targetHex!,
                                      ref
                                          .read(sessionProviderProvider)
                                          .value!
                                          .id);
                            },
                            child: const Text('ATTACK'),
                          ),
                        ],
                      ),
                    ],
                  ),
                ),
              ),
            ),
          if (gameState?.phase == GamePhase.gameOver)
            GameOverOverlay(state: gameState!),
        ],
      ),
    );
  }
}
