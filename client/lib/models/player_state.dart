import 'package:json_annotation/json_annotation.dart';

part 'player_state.g.dart';

@JsonSerializable()
class PlayerState {
  final String id;
  final String nickname;
  final int coins;
  @JsonKey(name: 'dominance_turn_counter')
  final int dominanceTurnCounter;
  @JsonKey(name: 'is_disconnected')
  final bool isDisconnected;

  const PlayerState({
    required this.id,
    required this.nickname,
    required this.coins,
    required this.dominanceTurnCounter,
    required this.isDisconnected,
  });

  factory PlayerState.fromJson(Map<String, dynamic> json) =>
      _$PlayerStateFromJson(json);
  Map<String, dynamic> toJson() => _$PlayerStateToJson(this);

  PlayerState copyWith({
    String? id,
    String? nickname,
    int? coins,
    int? dominanceTurnCounter,
    bool? isDisconnected,
  }) {
    return PlayerState(
      id: id ?? this.id,
      nickname: nickname ?? this.nickname,
      coins: coins ?? this.coins,
      dominanceTurnCounter: dominanceTurnCounter ?? this.dominanceTurnCounter,
      isDisconnected: isDisconnected ?? this.isDisconnected,
    );
  }
}

@JsonSerializable()
class GameOverStats {
  @JsonKey(name: 'turns_played')
  final int turnsPlayed;
  @JsonKey(name: 'troops_killed')
  final int troopsKilled;
  @JsonKey(name: 'troops_lost')
  final int troopsLost;
  @JsonKey(name: 'structures_held')
  final int structuresHeld;
  @JsonKey(name: 'total_damage_dealt')
  final int totalDamageDealt;

  const GameOverStats({
    required this.turnsPlayed,
    required this.troopsKilled,
    required this.troopsLost,
    required this.structuresHeld,
    required this.totalDamageDealt,
  });

  factory GameOverStats.fromJson(Map<String, dynamic> json) =>
      _$GameOverStatsFromJson(json);
  Map<String, dynamic> toJson() => _$GameOverStatsToJson(this);
}
