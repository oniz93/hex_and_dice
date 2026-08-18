import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:go_router/go_router.dart';

import '../models/enums.dart';
import '../models/room.dart';
import '../providers/core_providers.dart';

class CreateRoomScreen extends ConsumerStatefulWidget {
  const CreateRoomScreen({super.key});

  @override
  ConsumerState<CreateRoomScreen> createState() => _CreateRoomScreenState();
}

class _CreateRoomScreenState extends ConsumerState<CreateRoomScreen> {
  MapSize _mapSize = MapSize.small;
  int _turnTimer = 90;
  TurnMode _turnMode = TurnMode.alternating;
  bool _creating = false;

  Future<void> _create() async {
    setState(() => _creating = true);
    try {
      final res = await ref.read(apiServiceProvider).createRoom(
            RoomSettings(
              mapSize: _mapSize,
              turnTimer: _turnTimer,
              turnMode: _turnMode,
            ),
          );
      if (!mounted) return;
      context.go('/lobby/${res.roomCode}');
    } catch (e) {
      if (!mounted) return;
      ScaffoldMessenger.of(
        context,
      ).showSnackBar(SnackBar(content: Text('Could not create room: $e')));
    } finally {
      if (mounted) {
        setState(() => _creating = false);
      }
    }
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(title: const Text('Create Room')),
      body: ListView(
        padding: const EdgeInsets.all(24),
        children: [
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
          const SizedBox(height: 16),
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
          const SizedBox(height: 16),
          DropdownButtonFormField<TurnMode>(
            initialValue: _turnMode,
            decoration: const InputDecoration(
              labelText: 'Turn Mode',
              border: OutlineInputBorder(),
            ),
            items: const [
              DropdownMenuItem(
                value: TurnMode.alternating,
                child: Text('Alternating turns'),
              ),
            ],
            onChanged: (value) {
              if (value != null) setState(() => _turnMode = value);
            },
          ),
          const SizedBox(height: 24),
          ElevatedButton(
            onPressed: _creating ? null : _create,
            child: _creating
                ? const SizedBox(
                    height: 20,
                    width: 20,
                    child: CircularProgressIndicator(strokeWidth: 2),
                  )
                : const Text('CREATE ROOM'),
          ),
        ],
      ),
    );
  }
}
