import 'package:flutter/material.dart';
import 'package:go_router/go_router.dart';

class TitleScreen extends StatelessWidget {
  const TitleScreen({super.key});

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      body: Center(
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
    );
  }
}
