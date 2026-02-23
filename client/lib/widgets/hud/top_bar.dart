import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import '../../providers/game_state_provider.dart';
import '../../providers/session_provider.dart';
import '../../providers/turn_timer_provider.dart';

class TopBar extends ConsumerWidget {
  const TopBar({super.key});

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final gameState = ref.watch(gameStateNotifierProvider);
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
          Text(
            'Turn ${gameState.turnNumber}',
            style: const TextStyle(color: Colors.white, fontSize: 16),
          ),
          Text(
            turnText,
            style: TextStyle(
              color: isMyTurn ? Colors.green : Colors.red,
              fontSize: 16,
              fontWeight: FontWeight.bold,
            ),
          ),
          Text(
            '⏱ $timerStr',
            style: TextStyle(
              color: remainingSeconds < 10 ? Colors.red : Colors.white,
              fontSize: 16,
              fontWeight:
                  remainingSeconds < 10 ? FontWeight.bold : FontWeight.normal,
            ),
          ),
        ],
      ),
    );
  }
}
