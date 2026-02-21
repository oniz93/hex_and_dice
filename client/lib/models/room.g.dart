// GENERATED CODE - DO NOT MODIFY BY HAND

part of 'room.dart';

// **************************************************************************
// JsonSerializableGenerator
// **************************************************************************

RoomSettings _$RoomSettingsFromJson(Map<String, dynamic> json) => RoomSettings(
      mapSize: $enumDecode(_$MapSizeEnumMap, json['map_size']),
      turnTimer: (json['turn_timer'] as num).toInt(),
      turnMode: $enumDecode(_$TurnModeEnumMap, json['turn_mode']),
    );

Map<String, dynamic> _$RoomSettingsToJson(RoomSettings instance) =>
    <String, dynamic>{
      'map_size': _$MapSizeEnumMap[instance.mapSize]!,
      'turn_timer': instance.turnTimer,
      'turn_mode': _$TurnModeEnumMap[instance.turnMode]!,
    };

const _$MapSizeEnumMap = {
  MapSize.small: 'small',
  MapSize.medium: 'medium',
  MapSize.large: 'large',
};

const _$TurnModeEnumMap = {
  TurnMode.alternating: 'alternating',
  TurnMode.simultaneous: 'simultaneous',
};

Room _$RoomFromJson(Map<String, dynamic> json) => Room(
      id: json['id'] as String,
      state: $enumDecode(_$RoomStateEnumMap, json['state']),
      settings: RoomSettings.fromJson(json['settings'] as Map<String, dynamic>),
      players:
          (json['players'] as List<dynamic>).map((e) => e as String).toList(),
      createdAt: DateTime.parse(json['created_at'] as String),
    );

Map<String, dynamic> _$RoomToJson(Room instance) => <String, dynamic>{
      'id': instance.id,
      'state': _$RoomStateEnumMap[instance.state]!,
      'settings': instance.settings,
      'players': instance.players,
      'created_at': instance.createdAt.toIso8601String(),
    };

const _$RoomStateEnumMap = {
  RoomState.waitingForOpponent: 'waiting_for_opponent',
  RoomState.ready: 'ready',
  RoomState.gameInProgress: 'game_in_progress',
  RoomState.gameOver: 'game_over',
};
