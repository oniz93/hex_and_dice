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
      return 0;
    }

    // If turn changed, reset local logic if needed
    if (_lastActivePlayerId != gameState.activePlayerState.id ||
        _lastTurnNumber != gameState.turnNumber) {
      _lastActivePlayerId = gameState.activePlayerState.id;
      _lastTurnNumber = gameState.turnNumber;
    }

    _startTimer(gameState, session.id);

    return _calculateRemaining(gameState);
  }

  void _startTimer(gameState, String myId) {
    _timer?.cancel();
    _timer = Timer.periodic(const Duration(seconds: 1), (timer) {
      if (!ref.exists(turnTimerProvider)) {
        timer.cancel();
        return;
      }
      final remaining = _calculateRemaining(gameState);
      state = remaining;

      if (remaining <= 0) {
        timer.cancel();
        _onTimeout(gameState, myId);
      }
    });
  }

  void _stopTimer() {
    _timer?.cancel();
    _timer = null;
  }

  int _calculateRemaining(gameState) {
    final now = DateTime.now().toUtc();
    // Assuming turnStartedAt is in UTC.
    // If it's not, we might need to adjust.
    final startedAt = gameState.turnStartedAt.isUtc
        ? gameState.turnStartedAt
        : gameState.turnStartedAt.toUtc();

    final elapsed = now.difference(startedAt).inSeconds;
    return (gameState.turnTimer - elapsed).clamp(0, gameState.turnTimer);
  }

  void _onTimeout(gameState, String myId) {
    if (gameState.isActivePlayer(myId)) {
      print('TurnTimer: Timeout reached, sending end_turn');
      ref.read(wsServiceProvider).sendEndTurn();
    }
  }
}
