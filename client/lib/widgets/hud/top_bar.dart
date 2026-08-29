import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:go_router/go_router.dart';
import '../../providers/core_providers.dart';
import '../../providers/game_state_provider.dart';
import '../../providers/session_provider.dart';
import '../../providers/settings_provider.dart';
import '../../providers/turn_timer_provider.dart';
import '../../services/app_updater.dart';

class TopBar extends ConsumerWidget {
  const TopBar({super.key});

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final gameState = ref.watch(gameStateProvider);
    final sessionAsync = ref.watch(sessionProviderProvider);
    final session = sessionAsync.value;
    final remainingSeconds = ref.watch(turnTimerProvider);

    if (gameState == null || session == null) {
      return const SizedBox.shrink();
    }

    final isMyTurn = gameState.isActivePlayer(session.id);
    final turnText = isMyTurn ? 'Your Turn' : "Opponent's Turn";

    // Format MM:SS
    final minutes = (remainingSeconds / 60).floor();
    final seconds = remainingSeconds % 60;
    final timerStr =
        '${minutes.toString().padLeft(2, '0')}:${seconds.toString().padLeft(2, '0')}';

    return Container(
      color: Colors.black54,
      padding: const EdgeInsets.symmetric(horizontal: 16.0, vertical: 8.0),
      child: Row(
        mainAxisAlignment: MainAxisAlignment.spaceBetween,
        children: [
          IconButton(
            icon: const Icon(Icons.exit_to_app, color: Colors.white, size: 20),
            tooltip: 'Exit Game',
            onPressed: () {
              showDialog(
                context: context,
                builder: (dialogContext) => AlertDialog(
                  title: const Text('Exit Game?'),
                  content: const Text(
                      'Are you sure you want to leave the game? You can still reconnect later unless the game expires.'),
                  actions: [
                    TextButton(
                      onPressed: () => Navigator.pop(dialogContext),
                      child: const Text('CANCEL'),
                    ),
                    TextButton(
                      onPressed: () {
                        // 1. Close dialog
                        Navigator.pop(dialogContext);
                        // 2. Disconnect WebSocket
                        ref.read(wsServiceProvider).disconnect();
                        // 3. Clear reconnect data
                        ref.read(settingsProvider.notifier).clearReconnectData();
                        // 4. Go to title screen
                        context.go('/');
                      },
                      child: const Text('EXIT',
                          style: TextStyle(color: Colors.red)),
                    ),
                  ],
                ),
              );
            },
          ),
          Text(
            'Turn ${gameState.turnNumber}',
            style: const TextStyle(color: Colors.white, fontSize: 14),
          ),
          Text(
            turnText,
            style: TextStyle(
              color: isMyTurn ? Colors.green : Colors.red,
              fontSize: 14,
              fontWeight: FontWeight.bold,
            ),
          ),
          Text(
            '⏱ $timerStr',
            style: TextStyle(
              color: remainingSeconds < 10 ? Colors.red : Colors.white,
              fontSize: 14,
              fontWeight:
                  remainingSeconds < 10 ? FontWeight.bold : FontWeight.normal,
            ),
          ),
          TextButton(
            onPressed: () => forceUpdateApp(),
            style: TextButton.styleFrom(
              padding: const EdgeInsets.symmetric(horizontal: 8, vertical: 4),
              minimumSize: Size.zero,
              tapTargetSize: MaterialTapTargetSize.shrinkWrap,
            ),
            child: const Text(
              'UPDATE APP',
              style: TextStyle(color: Colors.orange, fontSize: 10),
            ),
          ),
        ],
      ),
    );
  }
}
