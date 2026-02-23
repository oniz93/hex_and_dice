import 'package:riverpod_annotation/riverpod_annotation.dart';

part 'combat_log_provider.g.dart';

class CombatLogEntry {
  final String message;
  final DateTime timestamp;

  CombatLogEntry(this.message) : timestamp = DateTime.now();
}

@Riverpod(keepAlive: true)
class CombatLog extends _$CombatLog {
  @override
  List<CombatLogEntry> build() => [];

  void addEntry(String msg) {
    state = [CombatLogEntry(msg), ...state].take(10).toList();
  }

  void clear() {
    state = [];
  }
}
