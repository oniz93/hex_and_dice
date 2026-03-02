import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import '../../providers/combat_log_provider.dart';

class CombatLogOverlay extends ConsumerStatefulWidget {
  const CombatLogOverlay({super.key});

  @override
  ConsumerState<CombatLogOverlay> createState() => _CombatLogOverlayState();
}

class _CombatLogOverlayState extends ConsumerState<CombatLogOverlay> {
  bool _isExpanded = true;

  @override
  Widget build(BuildContext context) {
    final logs = ref.watch(combatLogProvider);

    if (logs.isEmpty) return const SizedBox.shrink();

    return Container(
      width: 250,
      padding: const EdgeInsets.all(8.0),
      decoration: BoxDecoration(
        color: Colors.black.withOpacity(0.6),
        borderRadius: BorderRadius.circular(8),
      ),
      child: Column(
        mainAxisSize: MainAxisSize.min,
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          GestureDetector(
            onTap: () => setState(() => _isExpanded = !_isExpanded),
            child: Row(
              mainAxisAlignment: MainAxisAlignment.spaceBetween,
              children: [
                const Text(
                  'COMBAT LOG',
                  style: TextStyle(
                    color: Colors.white,
                    fontWeight: FontWeight.bold,
                    fontSize: 10,
                  ),
                ),
                Icon(
                  _isExpanded ? Icons.expand_less : Icons.expand_more,
                  color: Colors.white,
                  size: 16,
                ),
              ],
            ),
          ),
          if (_isExpanded) ...[
            const Divider(color: Colors.white24, height: 8),
            Flexible(
              child: ListView.builder(
                shrinkWrap: true,
                itemCount: logs.length,
                itemBuilder: (context, index) {
                  return Padding(
                    padding: const EdgeInsets.symmetric(vertical: 2.0),
                    child: Text(
                      logs[index].message,
                      style: const TextStyle(color: Colors.white, fontSize: 11),
                    ),
                  );
                },
              ),
            ),
          ],
        ],
      ),
    );
  }
}
