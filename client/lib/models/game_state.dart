import 'package:json_annotation/json_annotation.dart';
import '../game/hex/cube_coord.dart';
import 'enums.dart';
import 'player_state.dart';
import 'structure.dart';
import 'troop.dart';
import 'messages.dart';

part 'game_state.g.dart';

@JsonSerializable()
class GameState {
  final String id;
  final GamePhase phase;
  @JsonKey(name: 'map_size')
  final MapSize mapSize;
  @JsonKey(name: 'turn_mode')
  final TurnMode turnMode;
  @JsonKey(name: 'turn_timer')
  final int turnTimer;
  @JsonKey(name: 'turn_number')
  final int turnNumber;
  @JsonKey(name: 'active_player')
  final int activePlayer;
  final List<PlayerState> players;
  @JsonKey(defaultValue: {})
  final Map<String, Troop> troops;
  @JsonKey(defaultValue: {})
  final Map<String, Structure> structures;

  // We need to parse terrain which is a map with stringified CubeCoord keys.
  // We'll write custom converters if necessary or use string map keys.
  @JsonKey(fromJson: _terrainFromJson, toJson: _terrainToJson, defaultValue: {})
  final Map<CubeCoord, TerrainType> terrain;

  final int seed;
  @JsonKey(name: 'created_at')
  final DateTime createdAt;
  @JsonKey(name: 'turn_started_at')
  final DateTime turnStartedAt;

  @JsonKey(name: 'sudden_death_active')
  final bool suddenDeathActive;
  @JsonKey(name: 'sudden_death_turn')
  final int suddenDeathTurn;
  @JsonKey(name: 'safe_zone_radius')
  final int safeZoneRadius;

  final List<GameOverStats> stats;
  @JsonKey(name: 'first_turn_restriction')
  final bool firstTurnRestriction;

  @JsonKey(includeFromJson: false, includeToJson: false)
  final GameOverData? gameOverData;

  const GameState({
    required this.id,
    required this.phase,
    required this.mapSize,
    required this.turnMode,
    required this.turnTimer,
    required this.turnNumber,
    required this.activePlayer,
    required this.players,
    required this.troops,
    required this.structures,
    required this.terrain,
    required this.seed,
    required this.createdAt,
    required this.turnStartedAt,
    required this.suddenDeathActive,
    required this.suddenDeathTurn,
    required this.safeZoneRadius,
    required this.stats,
    required this.firstTurnRestriction,
    this.gameOverData,
  });

  factory GameState.fromJson(Map<String, dynamic> json) =>
      _$GameStateFromJson(json);
  Map<String, dynamic> toJson() => _$GameStateToJson(this);

  static Map<CubeCoord, TerrainType> _terrainFromJson(
    Map<String, dynamic> json,
  ) {
    return json.map((key, value) {
      final typeStr = value.toString();
      final type = TerrainType.values.firstWhere(
        (e) => e.name == typeStr,
        orElse: () => TerrainType.plains,
      );
      return MapEntry(CubeCoord.fromJson(key), type);
    });
  }

  static Map<String, dynamic> _terrainToJson(
    Map<CubeCoord, TerrainType> terrain,
  ) {
    return terrain.map((key, value) {
      String strKey = '${key.q},${key.r},${key.s}';
      return MapEntry(strKey, value.name);
    });
  }

  PlayerState get activePlayerState => players[activePlayer];
  PlayerState get inactivePlayerState => players[1 - activePlayer];

  bool isActivePlayer(String playerId) => activePlayerState.id == playerId;

  Troop? troopAt(CubeCoord hex) {
    try {
      return troops.values.firstWhere((t) => t.hex == hex && t.isAlive);
    } catch (_) {
      return null;
    }
  }

  Structure? structureAt(CubeCoord hex) {
    try {
      return structures.values.firstWhere((s) => s.hex == hex);
    } catch (_) {
      return null;
    }
  }

  GameState copyWith({
    String? id,
    GamePhase? phase,
    MapSize? mapSize,
    TurnMode? turnMode,
    int? turnTimer,
    int? turnNumber,
    int? activePlayer,
    List<PlayerState>? players,
    Map<String, Troop>? troops,
    Map<String, Structure>? structures,
    Map<CubeCoord, TerrainType>? terrain,
    int? seed,
    DateTime? createdAt,
    DateTime? turnStartedAt,
    bool? suddenDeathActive,
    int? suddenDeathTurn,
    int? safeZoneRadius,
    List<GameOverStats>? stats,
    bool? firstTurnRestriction,
    GameOverData? gameOverData,
  }) {
    return GameState(
      id: id ?? this.id,
      phase: phase ?? this.phase,
      mapSize: mapSize ?? this.mapSize,
      turnMode: turnMode ?? this.turnMode,
      turnTimer: turnTimer ?? this.turnTimer,
      turnNumber: turnNumber ?? this.turnNumber,
      activePlayer: activePlayer ?? this.activePlayer,
      players: players ?? this.players,
      troops: troops ?? this.troops,
      structures: structures ?? this.structures,
      terrain: terrain ?? this.terrain,
      seed: seed ?? this.seed,
      createdAt: createdAt ?? this.createdAt,
      turnStartedAt: turnStartedAt ?? this.turnStartedAt,
      suddenDeathActive: suddenDeathActive ?? this.suddenDeathActive,
      suddenDeathTurn: suddenDeathTurn ?? this.suddenDeathTurn,
      safeZoneRadius: safeZoneRadius ?? this.safeZoneRadius,
      stats: stats ?? this.stats,
      firstTurnRestriction: firstTurnRestriction ?? this.firstTurnRestriction,
      gameOverData: gameOverData ?? this.gameOverData,
    );
  }
}
