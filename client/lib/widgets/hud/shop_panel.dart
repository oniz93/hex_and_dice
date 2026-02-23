import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import '../../providers/game_state_provider.dart';
import '../../providers/selection_provider.dart';
import '../../providers/core_providers.dart';
import '../../models/enums.dart';
import '../../game/data/balance.dart';

class ShopPanel extends ConsumerWidget {
  const ShopPanel({super.key});

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final selection = ref.watch(selectionStateNotifierProvider);
    final gameState = ref.watch(gameStateNotifierProvider);

    if (gameState == null ||
        selection.state != SelectionFSM.structureSelected ||
        selection.targetHex == null) {
      return const SizedBox.shrink();
    }

    final structure = gameState.structureAt(selection.targetHex!);
    if (structure == null) return const SizedBox.shrink();

    final player =
        gameState.players.firstWhere((p) => p.id == structure.ownerId);

    return Container(
      height: 240,
      color: Colors.black87,
      child: Column(
        children: [
          Padding(
            padding: const EdgeInsets.all(8.0),
            child: Row(
              mainAxisAlignment: MainAxisAlignment.spaceBetween,
              children: [
                Text(
                  'Buy Troops at ${structure.type.name.toUpperCase()}',
                  style: const TextStyle(
                      color: Colors.white, fontWeight: FontWeight.bold),
                ),
                IconButton(
                  icon: const Icon(Icons.close, color: Colors.white),
                  onPressed: () => ref
                      .read(selectionStateNotifierProvider.notifier)
                      .clearSelection(),
                ),
              ],
            ),
          ),
          Expanded(
            child: ListView(
              scrollDirection: Axis.horizontal,
              padding: const EdgeInsets.symmetric(horizontal: 8.0),
              children: TroopType.values.map((type) {
                final stats = troopStats[type]!;
                final canAfford = player.coins >= stats.cost;
                final isOccupied =
                    gameState.troopAt(selection.targetHex!) != null;

                return Container(
                  width: 120,
                  margin: const EdgeInsets.all(4.0),
                  child: Card(
                    color: Colors.grey[900],
                    child: Padding(
                      padding: const EdgeInsets.all(8.0),
                      child: Column(
                        mainAxisAlignment: MainAxisAlignment.spaceBetween,
                        children: [
                          Text(
                            type.name.toUpperCase(),
                            style: const TextStyle(
                                color: Colors.white,
                                fontSize: 12,
                                fontWeight: FontWeight.bold),
                          ),
                          Text('${stats.cost}🪙',
                              style: const TextStyle(
                                  color: Colors.yellow, fontSize: 14)),
                          Column(
                            children: [
                              _SmallStatRow(label: 'HP', value: '${stats.hp}'),
                              _SmallStatRow(
                                  label: 'ATK', value: '+${stats.atk}'),
                            ],
                          ),
                          ElevatedButton(
                            style: ElevatedButton.styleFrom(
                              backgroundColor: (canAfford && !isOccupied)
                                  ? Colors.blue
                                  : Colors.grey,
                              padding:
                                  const EdgeInsets.symmetric(horizontal: 8),
                            ),
                            onPressed: (canAfford && !isOccupied)
                                ? () {
                                    ref
                                        .read(wsServiceProvider)
                                        .sendBuy(type, structure.id);
                                    ref
                                        .read(selectionStateNotifierProvider
                                            .notifier)
                                        .clearSelection();
                                  }
                                : null,
                            child: const Text('BUY',
                                style: TextStyle(
                                    fontSize: 10, color: Colors.white)),
                          ),
                        ],
                      ),
                    ),
                  ),
                );
              }).toList(),
            ),
          ),
        ],
      ),
    );
  }
}

class _SmallStatRow extends StatelessWidget {
  final String label;
  final String value;

  const _SmallStatRow({required this.label, required this.value});

  @override
  Widget build(BuildContext context) {
    return Row(
      mainAxisAlignment: MainAxisAlignment.spaceBetween,
      children: [
        Text(label,
            style: const TextStyle(color: Colors.white70, fontSize: 10)),
        Text(value, style: const TextStyle(color: Colors.white, fontSize: 10)),
      ],
    );
  }
}
