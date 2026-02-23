import 'package:flutter/material.dart';
import 'package:flame/game.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import '../../game/hex_game.dart';
import '../../providers/game_state_provider.dart';
import '../../providers/selection_provider.dart';
import '../../providers/session_provider.dart';
import '../../providers/core_providers.dart';
import '../widgets/hud/top_bar.dart';
import '../widgets/hud/bottom_bar.dart';
import '../widgets/hud/troop_popup.dart';
import '../widgets/hud/shop_panel.dart';
import '../widgets/hud/combat_log_overlay.dart';

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
            .read(selectionStateNotifierProvider.notifier)
            .handleHexTap(hex, session.id);
      }
    };

    print('GameScreen: initState called for room ${widget.roomId}');

    WidgetsBinding.instance.addPostFrameCallback((_) {
      print('GameScreen: addPostFrameCallback running');
      _connectToGame();
    });
  }

  Future<void> _connectToGame() async {
    try {
      print('GameScreen: _connectToGame starting');
      final sessionAsync = ref.read(sessionProviderProvider);
      print('GameScreen: Session state: $sessionAsync');

      final session = sessionAsync.value;
      if (session == null) {
        print('GameScreen: ERROR - session is null!');
        return;
      }

      print(
          'GameScreen: Connecting to game with roomId: ${widget.roomId}, token: ${session.token.substring(0, 10)}...');
      final wsService = ref.read(wsServiceProvider);
      await wsService.connect(session.token);

      print(
          'GameScreen: Connected to WS, sending join_game for room ${widget.roomId}...');
      wsService.sendJoinGame(widget.roomId);

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
    ref.listen(gameStateNotifierProvider, (prev, next) {
      if (next != null) {
        game.updateGameState(next);
      }
    });

    // Listen to selection updates
    ref.listen(selectionStateNotifierProvider, (prev, next) {
      game.updateSelection(next.highlightedMoves, next.highlightedAttacks);
    });

    final selection = ref.watch(selectionStateNotifierProvider);

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
                                .read(selectionStateNotifierProvider.notifier)
                                .clearSelection(),
                            child: const Text('CANCEL',
                                style: TextStyle(color: Colors.red)),
                          ),
                          const SizedBox(width: 16),
                          ElevatedButton(
                            onPressed: () {
                              ref
                                  .read(selectionStateNotifierProvider.notifier)
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
        ],
      ),
    );
  }
}
