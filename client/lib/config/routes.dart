import 'package:flutter/material.dart';
import 'package:go_router/go_router.dart';
import '../screens/title_screen.dart';
import '../screens/play_screen.dart';
import '../screens/game_screen.dart';
import '../screens/matchmaking_screen.dart';

final goRouter = GoRouter(
  initialLocation: '/',
  routes: [
    GoRoute(path: '/', builder: (context, state) => const TitleScreen()),
    GoRoute(path: '/play', builder: (context, state) => const PlayScreen()),
    GoRoute(
      path: '/game/:id',
      builder: (context, state) =>
          GameScreen(roomId: state.pathParameters['id']!),
    ),
    GoRoute(
      path: '/matchmaking',
      builder: (context, state) => const MatchmakingScreen(),
    ),
  ],
);
