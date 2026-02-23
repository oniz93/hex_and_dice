import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:go_router/go_router.dart';
import 'dart:async';
import '../providers/core_providers.dart';

class MatchmakingScreen extends ConsumerStatefulWidget {
  const MatchmakingScreen({super.key});

  @override
  ConsumerState<MatchmakingScreen> createState() => _MatchmakingScreenState();
}

class _MatchmakingScreenState extends ConsumerState<MatchmakingScreen> {
  Timer? _timer;
  int _elapsed = 0;

  @override
  void initState() {
    super.initState();
    _timer = Timer.periodic(const Duration(seconds: 1), (timer) {
      setState(() {
        _elapsed++;
      });
      _checkStatus();
    });
  }

  Future<void> _checkStatus() async {
    try {
      final status = await ref.read(apiServiceProvider).getMatchmakingStatus();
      print('MatchmakingScreen: status = $status');
      print(
          'MatchmakingScreen: queued = ${status['queued']}, type = ${status['queued'].runtimeType}');
      print(
          'MatchmakingScreen: has room_id = ${status.containsKey('room_id')}, value = ${status['room_id']}');

      if (status['queued'] == false && status.containsKey('room_id')) {
        print('MatchmakingScreen: Navigating to game...');
        _timer?.cancel();
        if (mounted) context.go('/game/${status['room_id']}');
      }
    } catch (e) {
      print('MatchmakingScreen: Error: $e');
    }
  }

  @override
  void dispose() {
    _timer?.cancel();
    super.dispose();
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      body: Center(
        child: Column(
          mainAxisAlignment: MainAxisAlignment.center,
          children: [
            const CircularProgressIndicator(),
            const SizedBox(height: 20),
            const Text('Searching for an opponent...'),
            const SizedBox(height: 10),
            Text('Elapsed: ${_elapsed}s'),
            const SizedBox(height: 30),
            ElevatedButton(
              onPressed: () async {
                await ref.read(apiServiceProvider).leaveMatchmaking();
                if (mounted) context.go('/play');
              },
              child: const Text('Cancel'),
            ),
          ],
        ),
      ),
    );
  }
}
