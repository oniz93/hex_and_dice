import 'package:json_annotation/json_annotation.dart';

enum GamePhase {
  @JsonValue('waiting_for_players')
  waitingForPlayers,
  @JsonValue('generating_map')
  generatingMap,
  @JsonValue('game_started')
  gameStarted,
  @JsonValue('turn_start')
  turnStart,
  @JsonValue('structure_combat')
  structureCombat,
  @JsonValue('player_action')
  playerAction,
  @JsonValue('turn_transition')
  turnTransition,
  @JsonValue('game_over')
  gameOver,
}

enum TurnMode {
  @JsonValue('alternating')
  alternating,
  @JsonValue('simultaneous')
  simultaneous,
}

enum MapSize {
  @JsonValue('small')
  small,
  @JsonValue('medium')
  medium,
  @JsonValue('large')
  large,
}

enum TroopType {
  @JsonValue('marine')
  marine,
  @JsonValue('sniper')
  sniper,
  @JsonValue('hoverbike')
  hoverbike,
  @JsonValue('mech')
  mech,
}

enum StructureType {
  @JsonValue('outpost')
  outpost,
  @JsonValue('command_center')
  commandCenter,
  @JsonValue('hq')
  hq,
}

enum TerrainType {
  @JsonValue('plains')
  plains,
  @JsonValue('forest')
  forest,
  @JsonValue('hills')
  hills,
  @JsonValue('water')
  water,
  @JsonValue('mountains')
  mountains,
}

enum RoomState {
  @JsonValue('waiting_for_opponent')
  waitingForOpponent,
  @JsonValue('ready')
  ready,
  @JsonValue('game_in_progress')
  gameInProgress,
  @JsonValue('game_over')
  gameOver,
}

enum WinReason {
  @JsonValue('HQ_DESTROYED')
  hqDestroyed,
  @JsonValue('STRUCTURE_DOMINANCE')
  structureDominance,
  @JsonValue('SUDDEN_DEATH')
  suddenDeath,
  @JsonValue('FORFEIT')
  forfeit,
  @JsonValue('DISCONNECT')
  disconnect,
  @JsonValue('DRAW')
  draw,
}

enum ErrorCode {
  @JsonValue('NOT_YOUR_TURN')
  notYourTurn,
  @JsonValue('INVALID_MOVE')
  invalidMove,
  @JsonValue('INVALID_ATTACK')
  invalidAttack,
  @JsonValue('INSUFFICIENT_FUNDS')
  insufficientFunds,
  @JsonValue('SPAWN_OCCUPIED')
  spawnOccupied,
  @JsonValue('SPAWN_NOT_OWNED')
  spawnNotOwned,
  @JsonValue('UNIT_ALREADY_ACTED')
  unitAlreadyActed,
  @JsonValue('UNIT_NOT_READY')
  unitNotReady,
  @JsonValue('UNIT_NOT_FOUND')
  unitNotFound,
  @JsonValue('GAME_NOT_FOUND')
  gameNotFound,
  @JsonValue('ROOM_NOT_FOUND')
  roomNotFound,
  @JsonValue('ROOM_FULL')
  roomFull,
  @JsonValue('ROOM_EXPIRED')
  roomExpired,
  @JsonValue('INVALID_MESSAGE')
  invalidMessage,
  @JsonValue('RATE_LIMITED')
  rateLimited,
}
