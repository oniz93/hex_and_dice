import 'dart:async';

import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:go_router/go_router.dart';

import '../providers/core_providers.dart';

class LobbyScreen extends ConsumerStatefulWidget {
  final String roomCode;

  const LobbyScreen({super.key, required this.roomCode});

  @override
  ConsumerState<LobbyScreen> createState() => _LobbyScreenState();
}

class _LobbyScreenState extends ConsumerState<LobbyScreen> {
  Timer? _pollTimer;
  String? _guestNickname;
  String? _error;

  @override
  void initState() {
    super.initState();
    _poll();
    _pollTimer = Timer.periodic(const Duration(milliseconds: 1500), (_) {
      _poll();
    });
  }

  @override
  void dispose() {
    _pollTimer?.cancel();
    super.dispose();
  }

  Future<void> _poll() async {
    try {
      final status = await ref
          .read(apiServiceProvider)
          .getRoomStatus(widget.roomCode);

      if (!mounted) return;

      if (status.gameId != null && status.gameId!.isNotEmpty) {
        _pollTimer?.cancel();
        context.go('/game/${status.gameId}');
        return;
      }

      setState(() {
        _guestNickname = status.guestNickname;
        _error = null;
      });
    } catch (_) {
      if (!mounted) return;
      setState(() {
        _error = 'Could not reach the lobby. Retrying…';
      });
    }
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(title: const Text('Waiting for Opponent')),
      body: Center(
        child: Column(
          mainAxisAlignment: MainAxisAlignment.center,
          children: [
            const Text(
              'Share this room code:',
              style: TextStyle(fontSize: 16),
            ),
            const SizedBox(height: 12),
            SelectableText(
              widget.roomCode,
              style: const TextStyle(
                fontSize: 48,
                fontWeight: FontWeight.bold,
                letterSpacing: 8,
              ),
            ),
            const SizedBox(height: 24),
            if (_guestNickname != null)
              Text('Opponent joined: $_guestNickname',
                  style: const TextStyle(color: Colors.green))
            else
              const Text('Waiting for opponent to join…'),
            const SizedBox(height: 12),
            const SizedBox(
              width: 24,
              height: 24,
              child: CircularProgressIndicator(strokeWidth: 2),
            ),
            if (_error != null) ...[
              const SizedBox(height: 16),
              Text(_error!, style: const TextStyle(color: Colors.orange)),
            ],
          ],
        ),
      ),
    );
  }
}
