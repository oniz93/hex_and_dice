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

    return Scaffold(
      body: Stack(
        children: [
          GameWidget(game: game),
          const Positioned(top: 0, left: 0, right: 0, child: TopBar()),
          const Positioned(bottom: 0, left: 0, right: 0, child: BottomBar()),
          // TroopPopup, ShopPanel, ConfirmAttack can be added here
        ],
      ),
    );
  }
}
