import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:go_router/go_router.dart';
import '../../models/game_state.dart';
import '../../models/messages.dart';
import '../../models/enums.dart';

class GameOverOverlay extends ConsumerWidget {
  final GameState state;

  const GameOverOverlay({super.key, required this.state});

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final data = state.gameOverData;
    if (data == null) return const SizedBox.shrink();

    String title = "GAME OVER";
    Color winnerColor = Colors.yellow;
    String winnerName = "";

    if (data.winnerId.isEmpty) {
      title = "IT'S A DRAW!";
    } else {
      final winnerIndex =
          state.players.indexWhere((p) => p.id == data.winnerId);
      if (winnerIndex >= 0) {
        winnerName = state.players[winnerIndex].nickname;
        title = "WINNER: $winnerName";
        winnerColor = winnerIndex == 0 ? Colors.red : Colors.blue;
      } else {
        title = "WINNER: ${data.winnerId}";
      }
    }

    return Container(
      color: Colors.black54,
      child: Center(
        child: Card(
          elevation: 8,
          child: Padding(
            padding: const EdgeInsets.all(24.0),
            child: Column(
              mainAxisSize: MainAxisSize.min,
              children: [
                Text(
                  title,
                  textAlign: TextAlign.center,
                  style: Theme.of(context).textTheme.headlineMedium?.copyWith(
                        fontWeight: FontWeight.bold,
                        color: winnerColor,
                      ),
                ),
                const SizedBox(height: 16),
                Text(
                  "Reason: ${data.reason.name.toUpperCase().replaceAll('_', ' ')}",
                  style: Theme.of(context).textTheme.bodyLarge,
                ),
                const SizedBox(height: 24),
                ElevatedButton(
                  onPressed: () {
                    context.go('/');
                  },
                  child: const Text("RETURN TO TITLE"),
                ),
              ],
            ),
          ),
        ),
      ),
    );
  }
}
