import 'dart:async';
import 'package:riverpod_annotation/riverpod_annotation.dart';
import 'game_state_provider.dart';
import 'session_provider.dart';
import 'core_providers.dart';

part 'turn_timer_provider.g.dart';

@Riverpod(keepAlive: true)
class TurnTimer extends _$TurnTimer {
  Timer? _timer;
  String? _lastActivePlayerId;
  int? _lastTurnNumber;

  DateTime? _localTurnStartTime;

  @override
  int build() {
    final gameState = ref.watch(gameStateNotifierProvider);
    final sessionAsync = ref.watch(sessionProviderProvider);
    final session = sessionAsync.value;

    ref.onDispose(() {
      _stopTimer();
    });

    if (gameState == null || session == null) {
      _stopTimer();
      _localTurnStartTime = null;
      return 0;
    }

    // If turn changed, reset local logic
    if (_lastActivePlayerId != gameState.activePlayerState.id ||
        _lastTurnNumber != gameState.turnNumber) {
      _lastActivePlayerId = gameState.activePlayerState.id;
      _lastTurnNumber = gameState.turnNumber;
      // Anchor the timer to when the client perceived the turn start
      _localTurnStartTime = DateTime.now();
    }

    _startTimer(session.id);

    return _calculateRemaining(gameState);
  }

  void _startTimer(String myId) {
    _timer?.cancel();
    _timer = Timer.periodic(const Duration(seconds: 1), (timer) {
      if (!ref.exists(turnTimerProvider)) {
        timer.cancel();
        return;
      }

      final currentGS = ref.read(gameStateNotifierProvider);
      if (currentGS == null) return;

      final remaining = _calculateRemaining(currentGS);
      state = remaining;

      if (remaining <= 0) {
        timer.cancel();
        _onTimeout(currentGS, myId);
      }
    });
  }

  void _stopTimer() {
    _timer?.cancel();
    _timer = null;
  }

  int _calculateRemaining(gameState) {
    if (_localTurnStartTime == null) return gameState.turnTimer;

    final elapsed = DateTime.now().difference(_localTurnStartTime!).inSeconds;
    return (gameState.turnTimer - elapsed).clamp(0, gameState.turnTimer);
  }

  void _onTimeout(gameState, String myId) {
    if (gameState.isActivePlayer(myId)) {
      print('TurnTimer: Timeout reached, sending end_turn');
      // Use ref.read to get the service without watching
      ref.read(wsServiceProvider).sendEndTurn();
    }
  }
}
