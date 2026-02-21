import 'package:json_annotation/json_annotation.dart';
import 'enums.dart';

part 'room.g.dart';

@JsonSerializable()
class RoomSettings {
  @JsonKey(name: 'map_size')
  final MapSize mapSize;
  @JsonKey(name: 'turn_timer')
  final int turnTimer;
  @JsonKey(name: 'turn_mode')
  final TurnMode turnMode;

  const RoomSettings({
    required this.mapSize,
    required this.turnTimer,
    required this.turnMode,
  });

  factory RoomSettings.fromJson(Map<String, dynamic> json) =>
      _$RoomSettingsFromJson(json);
  Map<String, dynamic> toJson() => _$RoomSettingsToJson(this);
}

@JsonSerializable()
class Room {
  final String id;
  final RoomState state;
  final RoomSettings settings;
  final List<String> players; // Player IDs
  @JsonKey(name: 'created_at')
  final DateTime createdAt;

  const Room({
    required this.id,
    required this.state,
    required this.settings,
    required this.players,
    required this.createdAt,
  });

  factory Room.fromJson(Map<String, dynamic> json) => _$RoomFromJson(json);
  Map<String, dynamic> toJson() => _$RoomToJson(this);
}
