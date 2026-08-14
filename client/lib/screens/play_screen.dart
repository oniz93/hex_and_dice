import 'dart:html' as html;
import 'package:flutter/foundation.dart';
import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:go_router/go_router.dart';
import '../providers/session_provider.dart';
import '../providers/settings_provider.dart';
import '../providers/core_providers.dart';

void _forceUpdateApp() {
  if (kIsWeb) {
    try {
      final sw = html.window.navigator.serviceWorker;
      if (sw != null) {
        sw.getRegistrations().then((registrations) {
          for (final reg in registrations) {
            reg.unregister();
          }
          html.window.location.reload();
        });
      } else {
        html.window.location.reload();
      }
    } catch (_) {
      html.window.location.reload();
    }
  }
}

class PlayScreen extends ConsumerWidget {
  const PlayScreen({super.key});

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final sessionAsync = ref.watch(sessionProviderProvider);

    return Scaffold(
      appBar: AppBar(title: const Text('Play')),
      body: Stack(
        children: [
          Center(
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
                final storedGameId = ref.watch(settingsProvider).gameId;

                return Column(
                  mainAxisAlignment: MainAxisAlignment.center,
                  children: [
                    Text('Welcome, ${session.nickname}!'),
                    const SizedBox(height: 20),
                    if (storedGameId != null)
                      Padding(
                        padding: const EdgeInsets.only(bottom: 20.0),
                        child: ElevatedButton(
                          onPressed: () => context.go('/game/$storedGameId'),
                          style: ElevatedButton.styleFrom(
                            backgroundColor: Colors.orange,
                            foregroundColor: Colors.white,
                          ),
                          child: const Text('REJOIN IN-PROGRESS GAME'),
                        ),
                      ),
                    ElevatedButton(
                      onPressed: () async {
                        try {
                          final res = await ref
                              .read(apiServiceProvider)
                              .joinMatchmaking();
                          if (res.status == 'matched') {
                            context.go('/game/${res.roomId}');
                          } else {
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
                        try {
                          final res = await ref
                              .read(apiServiceProvider)
                              .createBotGame(difficulty: 'easy');
                          if (res.roomId.isNotEmpty) {
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
                        try {
                          final res = await ref
                              .read(apiServiceProvider)
                              .createBotGame(difficulty: 'hard');
                          if (res.roomId.isNotEmpty) {
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
                      onPressed: () => context.go('/create-room'),
                      child: const Text('Create Room'),
                    ),
                    const SizedBox(height: 10),
                    ElevatedButton(
                      onPressed: () => context.go('/join-room'),
                      child: const Text('Join Room'),
                    ),
                  ],
                );
              },
              loading: () => const CircularProgressIndicator(),
              error: (e, st) => Text('Error: $e'),
            ),
          ),
          Positioned(
            top: 16,
            right: 16,
            child: TextButton(
              onPressed: _forceUpdateApp,
              child: const Text(
                'UPDATE APP',
                style: TextStyle(color: Colors.orange, fontSize: 12),
              ),
            ),
          ),
        ],
      ),
    );
  }
}
