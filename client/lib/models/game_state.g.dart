// GENERATED CODE - DO NOT MODIFY BY HAND

part of 'game_state.dart';

// **************************************************************************
// JsonSerializableGenerator
// **************************************************************************

GameState _$GameStateFromJson(Map<String, dynamic> json) => GameState(
      id: json['id'] as String,
      phase: $enumDecode(_$GamePhaseEnumMap, json['phase']),
      mapSize: $enumDecode(_$MapSizeEnumMap, json['map_size']),
      turnMode: $enumDecode(_$TurnModeEnumMap, json['turn_mode']),
      turnTimer: (json['turn_timer'] as num).toInt(),
      turnNumber: (json['turn_number'] as num).toInt(),
      activePlayer: (json['active_player'] as num).toInt(),
      players: (json['players'] as List<dynamic>)
          .map((e) => PlayerState.fromJson(e as Map<String, dynamic>))
          .toList(),
      troops: (json['troops'] as Map<String, dynamic>?)?.map(
            (k, e) => MapEntry(k, Troop.fromJson(e as Map<String, dynamic>)),
          ) ??
          {},
      structures: (json['structures'] as Map<String, dynamic>?)?.map(
            (k, e) =>
                MapEntry(k, Structure.fromJson(e as Map<String, dynamic>)),
          ) ??
          {},
      terrain: json['terrain'] == null
          ? {}
          : GameState._terrainFromJson(json['terrain'] as Map<String, dynamic>),
      seed: (json['seed'] as num).toInt(),
      createdAt: DateTime.parse(json['created_at'] as String),
      turnStartedAt: DateTime.parse(json['turn_started_at'] as String),
      suddenDeathActive: json['sudden_death_active'] as bool,
      suddenDeathTurn: (json['sudden_death_turn'] as num).toInt(),
      safeZoneRadius: (json['safe_zone_radius'] as num).toInt(),
      stats: (json['stats'] as List<dynamic>)
          .map((e) => GameOverStats.fromJson(e as Map<String, dynamic>))
          .toList(),
      firstTurnRestriction: json['first_turn_restriction'] as bool,
    );

Map<String, dynamic> _$GameStateToJson(GameState instance) => <String, dynamic>{
      'id': instance.id,
      'phase': _$GamePhaseEnumMap[instance.phase]!,
      'map_size': _$MapSizeEnumMap[instance.mapSize]!,
      'turn_mode': _$TurnModeEnumMap[instance.turnMode]!,
      'turn_timer': instance.turnTimer,
      'turn_number': instance.turnNumber,
      'active_player': instance.activePlayer,
      'players': instance.players,
      'troops': instance.troops,
      'structures': instance.structures,
      'terrain': GameState._terrainToJson(instance.terrain),
      'seed': instance.seed,
      'created_at': instance.createdAt.toIso8601String(),
      'turn_started_at': instance.turnStartedAt.toIso8601String(),
      'sudden_death_active': instance.suddenDeathActive,
      'sudden_death_turn': instance.suddenDeathTurn,
      'safe_zone_radius': instance.safeZoneRadius,
      'stats': instance.stats,
      'first_turn_restriction': instance.firstTurnRestriction,
    };

const _$GamePhaseEnumMap = {
  GamePhase.waitingForPlayers: 'waiting_for_players',
  GamePhase.generatingMap: 'generating_map',
  GamePhase.gameStarted: 'game_started',
  GamePhase.turnStart: 'turn_start',
  GamePhase.structureCombat: 'structure_combat',
  GamePhase.playerAction: 'player_action',
  GamePhase.turnTransition: 'turn_transition',
  GamePhase.gameOver: 'game_over',
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
