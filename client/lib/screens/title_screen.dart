import 'dart:html' as html;
import 'package:flutter/foundation.dart';
import 'package:flutter/material.dart';
import 'package:go_router/go_router.dart';

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

class TitleScreen extends StatelessWidget {
  const TitleScreen({super.key});

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      body: Stack(
        children: [
          Center(
            child: Column(
              mainAxisAlignment: MainAxisAlignment.center,
              children: [
                const Text('Hex & Dice', style: TextStyle(fontSize: 48)),
                const SizedBox(height: 50),
                ElevatedButton(
                  onPressed: () {
                    context.go('/play');
                  },
                  child: const Text('Play'),
                ),
                ElevatedButton(
                  onPressed: () {
                    // Navigate to how to play
                  },
                  child: const Text('How to Play'),
                ),
              ],
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
