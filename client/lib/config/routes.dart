import 'package:go_router/go_router.dart';
import '../screens/title_screen.dart';
import '../screens/play_screen.dart';
import '../screens/game_screen.dart';
import '../screens/matchmaking_screen.dart';
import '../screens/create_room_screen.dart';
import '../screens/join_room_screen.dart';
import '../screens/lobby_screen.dart';

import 'package:flutter_riverpod/flutter_riverpod.dart';
import '../providers/settings_provider.dart';

final goRouter = GoRouter(
  initialLocation: '/',
  redirect: (context, state) {
    // Only redirect to active game if the user is trying to access the root title screen.
    // This prevents stuck states if they manually navigate or we try to leave.
    if (state.matchedLocation == '/') {
      final container = ProviderScope.containerOf(context, listen: false);
      final settings = container.read(settingsProvider);
      if (settings.gameId != null && settings.gameId!.isNotEmpty) {
        print('Router: Active game found (${settings.gameId}), redirecting...');
        return '/game/${settings.gameId}';
      }
    }
    return null;
  },
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
    GoRoute(
      path: '/create-room',
      builder: (context, state) => const CreateRoomScreen(),
    ),
    GoRoute(
      path: '/join-room',
      builder: (context, state) => const JoinRoomScreen(),
    ),
    GoRoute(
      path: '/lobby/:code',
      builder: (context, state) =>
          LobbyScreen(roomCode: state.pathParameters['code']!),
    ),
  ],
);
