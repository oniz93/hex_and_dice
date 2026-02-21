import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import '../../providers/game_state_provider.dart';
import '../../providers/session_provider.dart';
import '../../providers/core_providers.dart';

class BottomBar extends ConsumerWidget {
  const BottomBar({super.key});

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final gameState = ref.watch(gameStateNotifierProvider);
    final sessionAsync = ref.watch(sessionProviderProvider);
    final session = sessionAsync.value;

    if (gameState == null || session == null) {
      return const SizedBox.shrink();
    }

    final isMyTurn = gameState.isActivePlayer(session.id);
    final coins = gameState.players.firstWhere((p) => p.id == session.id).coins;

    return Container(
      color: Colors.black54,
      padding: const EdgeInsets.all(16.0),
      child: Row(
        mainAxisAlignment: MainAxisAlignment.spaceBetween,
        children: [
          Row(
            children: [
              const Icon(Icons.monetization_on, color: Colors.yellow),
              const SizedBox(width: 8),
              Text(
                '$coins',
                style: const TextStyle(
                  color: Colors.white,
                  fontSize: 18,
                  fontWeight: FontWeight.bold,
                ),
              ),
            ],
          ),
          IconButton(
            icon: const Icon(Icons.emoji_emotions, color: Colors.white),
            onPressed: () {
              // Open emote bar
            },
          ),
          ElevatedButton(
            style: ElevatedButton.styleFrom(
              backgroundColor: isMyTurn ? Colors.blue : Colors.grey,
            ),
            onPressed: isMyTurn
                ? () {
                    ref.read(wsServiceProvider).sendEndTurn();
                  }
                : null,
            child: const Text(
              'END TURN',
              style: TextStyle(color: Colors.white),
            ),
          ),
        ],
      ),
    );
  }
}
