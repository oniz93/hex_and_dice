import 'package:flutter/material.dart';

import '../models/enums.dart';

class BotGameConfig {
  final String difficulty;
  final MapSize mapSize;
  final int turnTimer;

  const BotGameConfig({
    required this.difficulty,
    required this.mapSize,
    required this.turnTimer,
  });
}

/// Shows a bottom sheet for configuring a bot match.
/// Returns null if the user dismisses it.
Future<BotGameConfig?> showBotSetupSheet(BuildContext context) {
  return showModalBottomSheet<BotGameConfig>(
    context: context,
    builder: (context) => const _BotSetupSheet(),
  );
}

class _BotSetupSheet extends StatefulWidget {
  const _BotSetupSheet();

  @override
  State<_BotSetupSheet> createState() => _BotSetupSheetState();
}

class _BotSetupSheetState extends State<_BotSetupSheet> {
  String _difficulty = 'easy';
  MapSize _mapSize = MapSize.small;
  int _turnTimer = 90;

  @override
  Widget build(BuildContext context) {
    return SafeArea(
      child: Padding(
        padding: const EdgeInsets.all(24),
        child: Column(
          mainAxisSize: MainAxisSize.min,
          crossAxisAlignment: CrossAxisAlignment.stretch,
          children: [
            const Text(
              'Play vs Bot',
              textAlign: TextAlign.center,
              style: TextStyle(fontSize: 20, fontWeight: FontWeight.bold),
            ),
            const SizedBox(height: 20),
            DropdownButtonFormField<String>(
              initialValue: _difficulty,
              decoration: const InputDecoration(
                labelText: 'Difficulty',
                border: OutlineInputBorder(),
              ),
              items: const [
                DropdownMenuItem(value: 'easy', child: Text('Easy')),
                DropdownMenuItem(value: 'medium', child: Text('Medium')),
                DropdownMenuItem(value: 'hard', child: Text('Hard')),
              ],
              onChanged: (value) {
                if (value != null) setState(() => _difficulty = value);
              },
            ),
            const SizedBox(height: 12),
            DropdownButtonFormField<MapSize>(
              initialValue: _mapSize,
              decoration: const InputDecoration(
                labelText: 'Map Size',
                border: OutlineInputBorder(),
              ),
              items: const [
                DropdownMenuItem(value: MapSize.small, child: Text('Small')),
                DropdownMenuItem(value: MapSize.medium, child: Text('Medium')),
                DropdownMenuItem(value: MapSize.large, child: Text('Large')),
              ],
              onChanged: (value) {
                if (value != null) setState(() => _mapSize = value);
              },
            ),
            const SizedBox(height: 12),
            DropdownButtonFormField<int>(
              initialValue: _turnTimer,
              decoration: const InputDecoration(
                labelText: 'Turn Timer',
                border: OutlineInputBorder(),
              ),
              items: const [
                DropdownMenuItem(value: 60, child: Text('60 seconds')),
                DropdownMenuItem(value: 90, child: Text('90 seconds')),
                DropdownMenuItem(value: 120, child: Text('120 seconds')),
              ],
              onChanged: (value) {
                if (value != null) setState(() => _turnTimer = value);
              },
            ),
            const SizedBox(height: 24),
            ElevatedButton(
              onPressed: () {
                Navigator.of(context).pop(
                  BotGameConfig(
                    difficulty: _difficulty,
                    mapSize: _mapSize,
                    turnTimer: _turnTimer,
                  ),
                );
              },
              child: const Text('START MATCH'),
            ),
          ],
        ),
      ),
    );
  }
}
