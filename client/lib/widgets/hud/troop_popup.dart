import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import '../../providers/game_state_provider.dart';
import '../../providers/selection_provider.dart';
import '../../models/troop.dart';
import '../../models/enums.dart';

class TroopPopup extends ConsumerWidget {
  const TroopPopup({super.key});

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final selection = ref.watch(selectionStateNotifierProvider);
    final gameState = ref.watch(gameStateNotifierProvider);

    if (gameState == null || selection.selectedUnitId == null) {
      return const SizedBox.shrink();
    }

    final troop = gameState.troops[selection.selectedUnitId];
    if (troop == null) return const SizedBox.shrink();

    return Card(
      color: Colors.black87,
      child: Padding(
        padding: const EdgeInsets.all(12.0),
        child: Column(
          mainAxisSize: MainAxisSize.min,
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            Row(
              mainAxisAlignment: MainAxisAlignment.spaceBetween,
              children: [
                Text(
                  troop.type.name.toUpperCase(),
                  style: const TextStyle(
                    color: Colors.white,
                    fontWeight: FontWeight.bold,
                    fontSize: 16,
                  ),
                ),
                Text(
                  'HP ${troop.currentHp}/${troop.maxHp}',
                  style: TextStyle(
                    color: troop.currentHp < troop.maxHp / 2
                        ? Colors.red
                        : Colors.green,
                    fontWeight: FontWeight.bold,
                  ),
                ),
              ],
            ),
            const Divider(color: Colors.white24),
            _StatRow(label: 'ATK', value: '+${troop.atk}'),
            _StatRow(label: 'DEF', value: '${troop.def}'),
            _StatRow(
                label: 'MOB',
                value: '${troop.remainingMobility}/${troop.mobility}'),
            _StatRow(label: 'RNG', value: '${troop.range}'),
            _StatRow(label: 'DMG', value: troop.damage),
            const SizedBox(height: 8),
            Text(
              'Status: ${troop.isReady ? "Ready" : "Waiting"}',
              style: TextStyle(
                color: troop.isReady ? Colors.blue : Colors.orange,
                fontSize: 12,
              ),
            ),
          ],
        ),
      ),
    );
  }
}

class _StatRow extends StatelessWidget {
  final String label;
  final String value;

  const _StatRow({required this.label, required this.value});

  @override
  Widget build(BuildContext context) {
    return Padding(
      padding: const EdgeInsets.symmetric(vertical: 2.0),
      child: Row(
        mainAxisAlignment: MainAxisAlignment.spaceBetween,
        children: [
          Text(label,
              style: const TextStyle(color: Colors.white70, fontSize: 12)),
          Text(value,
              style: const TextStyle(color: Colors.white, fontSize: 12)),
        ],
      ),
    );
  }
}
