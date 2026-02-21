// GENERATED CODE - DO NOT MODIFY BY HAND

part of 'player_state.dart';

// **************************************************************************
// JsonSerializableGenerator
// **************************************************************************

PlayerState _$PlayerStateFromJson(Map<String, dynamic> json) => PlayerState(
      id: json['id'] as String,
      nickname: json['nickname'] as String,
      coins: (json['coins'] as num).toInt(),
      dominanceTurnCounter: (json['dominance_turn_counter'] as num).toInt(),
      isDisconnected: json['is_disconnected'] as bool,
    );

Map<String, dynamic> _$PlayerStateToJson(PlayerState instance) =>
    <String, dynamic>{
      'id': instance.id,
      'nickname': instance.nickname,
      'coins': instance.coins,
      'dominance_turn_counter': instance.dominanceTurnCounter,
      'is_disconnected': instance.isDisconnected,
    };

GameOverStats _$GameOverStatsFromJson(Map<String, dynamic> json) =>
    GameOverStats(
      turnsPlayed: (json['turns_played'] as num).toInt(),
      troopsKilled: (json['troops_killed'] as num).toInt(),
      troopsLost: (json['troops_lost'] as num).toInt(),
      structuresHeld: (json['structures_held'] as num).toInt(),
      totalDamageDealt: (json['total_damage_dealt'] as num).toInt(),
    );

Map<String, dynamic> _$GameOverStatsToJson(GameOverStats instance) =>
    <String, dynamic>{
      'turns_played': instance.turnsPlayed,
      'troops_killed': instance.troopsKilled,
      'troops_lost': instance.troopsLost,
      'structures_held': instance.structuresHeld,
      'total_damage_dealt': instance.totalDamageDealt,
    };
