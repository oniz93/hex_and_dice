import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:go_router/go_router.dart';
import '../providers/session_provider.dart';
import '../providers/core_providers.dart';

class PlayScreen extends ConsumerWidget {
  const PlayScreen({super.key});

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final sessionAsync = ref.watch(sessionProviderProvider);

    return Scaffold(
      appBar: AppBar(title: const Text('Play')),
      body: Center(
        child: sessionAsync.when(
          data: (session) {
            if (session == null) {
              return Column(
                mainAxisAlignment: MainAxisAlignment.center,
                children: [
                  const Text('Enter Nickname (3-16 chars)'),
                  const SizedBox(height: 10),
                  SizedBox(
                    width: 200,
                    child: TextField(
                      decoration: const InputDecoration(
                        border: OutlineInputBorder(),
                        hintText: 'Nickname',
                      ),
                      onSubmitted: (val) {
                        if (val.trim().length >= 3) {
                          ref
                              .read(sessionProviderProvider.notifier)
                              .registerGuest(val.trim());
                        } else {
                          ScaffoldMessenger.of(context).showSnackBar(
                            const SnackBar(
                              content: Text(
                                'Nickname must be at least 3 characters',
                              ),
                            ),
                          );
                        }
                      },
                    ),
                  ),
                ],
              );
            }
            return Column(
              mainAxisAlignment: MainAxisAlignment.center,
              children: [
                Text('Welcome, ${session.nickname}!'),
                const SizedBox(height: 20),
                ElevatedButton(
                  onPressed: () async {
                    // Call matchmaking
                    try {
                      final res =
                          await ref.read(apiServiceProvider).joinMatchmaking();
                      if (res.status == 'matched') {
                        // Connect to game
                        // ...
                        context.go('/game/${res.roomId}');
                      } else {
                        // Go to matchmaking waiting screen
                        context.go('/matchmaking');
                      }
                    } catch (e) {
                      ScaffoldMessenger.of(
                        context,
                      ).showSnackBar(SnackBar(content: Text('Error: $e')));
                    }
                  },
                  child: const Text('Quick Match'),
                ),
                const SizedBox(height: 10),
                ElevatedButton(
                  onPressed: () async {
                    // Play vs Bot
                    try {
                      final res = await ref
                          .read(apiServiceProvider)
                          .createBotGame(difficulty: 'easy');
                      if (res.roomId.isNotEmpty) {
                        // Connect to game
                        // ignore: use_build_context_synchronously
                        context.go('/game/${res.roomId}');
                      }
                    } catch (e) {
                      ScaffoldMessenger.of(
                        context,
                      ).showSnackBar(SnackBar(content: Text('Error: $e')));
                    }
                  },
                  style: ElevatedButton.styleFrom(
                    backgroundColor: Colors.green,
                    foregroundColor: Colors.white,
                  ),
                  child: const Text('Play vs Bot (Easy)'),
                ),
                const SizedBox(height: 10),
                ElevatedButton(
                  onPressed: () async {
                    // Play vs Bot - Hard
                    try {
                      final res = await ref
                          .read(apiServiceProvider)
                          .createBotGame(difficulty: 'hard');
                      if (res.roomId.isNotEmpty) {
                        // Connect to game
                        // ignore: use_build_context_synchronously
                        context.go('/game/${res.roomId}');
                      }
                    } catch (e) {
                      ScaffoldMessenger.of(
                        context,
                      ).showSnackBar(SnackBar(content: Text('Error: $e')));
                    }
                  },
                  style: ElevatedButton.styleFrom(
                    backgroundColor: Colors.red,
                    foregroundColor: Colors.white,
                  ),
                  child: const Text('Play vs Bot (Hard)'),
                ),
                const SizedBox(height: 20),
                ElevatedButton(
                  onPressed: () {
                    // Create room
                  },
                  child: const Text('Create Room'),
                ),
                const SizedBox(height: 10),
                ElevatedButton(
                  onPressed: () {
                    // Join room
                  },
                  child: const Text('Join Room'),
                ),
              ],
            );
          },
          loading: () => const CircularProgressIndicator(),
          error: (e, st) => Text('Error: $e'),
        ),
      ),
    );
  }
}
