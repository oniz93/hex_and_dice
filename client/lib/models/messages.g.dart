// GENERATED CODE - DO NOT MODIFY BY HAND

part of 'messages.dart';

// **************************************************************************
// JsonSerializableGenerator
// **************************************************************************

ErrorData _$ErrorDataFromJson(Map<String, dynamic> json) => ErrorData(
      code: $enumDecode(_$ErrorCodeEnumMap, json['code']),
      message: json['message'] as String,
    );

Map<String, dynamic> _$ErrorDataToJson(ErrorData instance) => <String, dynamic>{
      'code': _$ErrorCodeEnumMap[instance.code]!,
      'message': instance.message,
    };

const _$ErrorCodeEnumMap = {
  ErrorCode.notYourTurn: 'NOT_YOUR_TURN',
  ErrorCode.invalidMove: 'INVALID_MOVE',
  ErrorCode.invalidAttack: 'INVALID_ATTACK',
  ErrorCode.insufficientFunds: 'INSUFFICIENT_FUNDS',
  ErrorCode.spawnOccupied: 'SPAWN_OCCUPIED',
  ErrorCode.spawnNotOwned: 'SPAWN_NOT_OWNED',
  ErrorCode.unitAlreadyActed: 'UNIT_ALREADY_ACTED',
  ErrorCode.unitNotReady: 'UNIT_NOT_READY',
  ErrorCode.unitNotFound: 'UNIT_NOT_FOUND',
  ErrorCode.gameNotFound: 'GAME_NOT_FOUND',
  ErrorCode.roomNotFound: 'ROOM_NOT_FOUND',
  ErrorCode.roomFull: 'ROOM_FULL',
  ErrorCode.roomExpired: 'ROOM_EXPIRED',
  ErrorCode.invalidMessage: 'INVALID_MESSAGE',
  ErrorCode.rateLimited: 'RATE_LIMITED',
};

AckData _$AckDataFromJson(Map<String, dynamic> json) => AckData(
      seq: (json['seq'] as num).toInt(),
      actionType: json['action_type'] as String,
    );

Map<String, dynamic> _$AckDataToJson(AckData instance) => <String, dynamic>{
      'seq': instance.seq,
      'action_type': instance.actionType,
    };

NackData _$NackDataFromJson(Map<String, dynamic> json) => NackData(
      seq: (json['seq'] as num).toInt(),
      actionType: json['action_type'] as String,
      error: ErrorData.fromJson(json['error'] as Map<String, dynamic>),
    );

Map<String, dynamic> _$NackDataToJson(NackData instance) => <String, dynamic>{
      'seq': instance.seq,
      'action_type': instance.actionType,
      'error': instance.error,
    };

TroopMovedData _$TroopMovedDataFromJson(Map<String, dynamic> json) =>
    TroopMovedData(
      unitId: json['unit_id'] as String,
      fromQ: (json['from_q'] as num).toInt(),
      fromR: (json['from_r'] as num).toInt(),
      fromS: (json['from_s'] as num).toInt(),
      toQ: (json['to_q'] as num).toInt(),
      toR: (json['to_r'] as num).toInt(),
      toS: (json['to_s'] as num).toInt(),
      remainingMobility: (json['remaining_mobility'] as num).toInt(),
    );

Map<String, dynamic> _$TroopMovedDataToJson(TroopMovedData instance) =>
    <String, dynamic>{
      'unit_id': instance.unitId,
      'from_q': instance.fromQ,
      'from_r': instance.fromR,
      'from_s': instance.fromS,
      'to_q': instance.toQ,
      'to_r': instance.toR,
      'to_s': instance.toS,
      'remaining_mobility': instance.remainingMobility,
    };

CombatResultData _$CombatResultDataFromJson(Map<String, dynamic> json) =>
    CombatResultData(
      attackerId: json['attacker_id'] as String,
      defenderId: json['defender_id'] as String,
      hitRoll: (json['hit_roll'] as num).toInt(),
      naturalRoll: (json['natural_roll'] as num).toInt(),
      hit: json['hit'] as bool,
      damageRoll: (json['damage_roll'] as num).toInt(),
      damage: (json['damage'] as num).toInt(),
      defenderHp: (json['defender_hp'] as num).toInt(),
      killed: json['killed'] as bool,
      crit: json['crit'] as bool,
      fumble: json['fumble'] as bool,
      hasCounter: json['has_counter'] as bool,
      counterHitRoll: (json['counter_hit_roll'] as num?)?.toInt(),
      counterNaturalRoll: (json['counter_natural_roll'] as num?)?.toInt(),
      counterHit: json['counter_hit'] as bool?,
      counterDamage: (json['counter_damage'] as num?)?.toInt(),
      attackerHp: (json['attacker_hp'] as num).toInt(),
      attackerKilled: json['attacker_killed'] as bool,
    );

Map<String, dynamic> _$CombatResultDataToJson(CombatResultData instance) =>
    <String, dynamic>{
      'attacker_id': instance.attackerId,
      'defender_id': instance.defenderId,
      'hit_roll': instance.hitRoll,
      'natural_roll': instance.naturalRoll,
      'hit': instance.hit,
      'damage_roll': instance.damageRoll,
      'damage': instance.damage,
      'defender_hp': instance.defenderHp,
      'killed': instance.killed,
      'crit': instance.crit,
      'fumble': instance.fumble,
      'has_counter': instance.hasCounter,
      'counter_hit_roll': instance.counterHitRoll,
      'counter_natural_roll': instance.counterNaturalRoll,
      'counter_hit': instance.counterHit,
      'counter_damage': instance.counterDamage,
      'attacker_hp': instance.attackerHp,
      'attacker_killed': instance.attackerKilled,
    };

TroopPurchasedData _$TroopPurchasedDataFromJson(Map<String, dynamic> json) =>
    TroopPurchasedData(
      unitId: json['unit_id'] as String,
      unitType: $enumDecode(_$TroopTypeEnumMap, json['unit_type']),
      hexQ: (json['hex_q'] as num).toInt(),
      hexR: (json['hex_r'] as num).toInt(),
      hexS: (json['hex_s'] as num).toInt(),
      owner: json['owner'] as String,
      coinsRemaining: (json['coins_remaining'] as num).toInt(),
    );

Map<String, dynamic> _$TroopPurchasedDataToJson(TroopPurchasedData instance) =>
    <String, dynamic>{
      'unit_id': instance.unitId,
      'unit_type': _$TroopTypeEnumMap[instance.unitType]!,
      'hex_q': instance.hexQ,
      'hex_r': instance.hexR,
      'hex_s': instance.hexS,
      'owner': instance.owner,
      'coins_remaining': instance.coinsRemaining,
    };

const _$TroopTypeEnumMap = {
  TroopType.marine: 'marine',
  TroopType.sniper: 'sniper',
  TroopType.hoverbike: 'hoverbike',
  TroopType.mech: 'mech',
};

TroopDestroyedData _$TroopDestroyedDataFromJson(Map<String, dynamic> json) =>
    TroopDestroyedData(
      unitId: json['unit_id'] as String,
      hexQ: (json['hex_q'] as num).toInt(),
      hexR: (json['hex_r'] as num).toInt(),
      hexS: (json['hex_s'] as num).toInt(),
      cause: json['cause'] as String,
    );

Map<String, dynamic> _$TroopDestroyedDataToJson(TroopDestroyedData instance) =>
    <String, dynamic>{
      'unit_id': instance.unitId,
      'hex_q': instance.hexQ,
      'hex_r': instance.hexR,
      'hex_s': instance.hexS,
      'cause': instance.cause,
    };

StructureAttackedData _$StructureAttackedDataFromJson(
        Map<String, dynamic> json) =>
    StructureAttackedData(
      structureId: json['structure_id'] as String,
      attackerId: json['attacker_id'] as String,
      hitRoll: (json['hit_roll'] as num).toInt(),
      damage: (json['damage'] as num).toInt(),
      structureHp: (json['structure_hp'] as num).toInt(),
      captured: json['captured'] as bool,
      newOwner: json['new_owner'] as String?,
    );

Map<String, dynamic> _$StructureAttackedDataToJson(
        StructureAttackedData instance) =>
    <String, dynamic>{
      'structure_id': instance.structureId,
      'attacker_id': instance.attackerId,
      'hit_roll': instance.hitRoll,
      'damage': instance.damage,
      'structure_hp': instance.structureHp,
      'captured': instance.captured,
      'new_owner': instance.newOwner,
    };

StructureFiresData _$StructureFiresDataFromJson(Map<String, dynamic> json) =>
    StructureFiresData(
      structureId: json['structure_id'] as String,
      targetId: json['target_id'] as String,
      hitRoll: (json['hit_roll'] as num).toInt(),
      damage: (json['damage'] as num).toInt(),
      targetHp: (json['target_hp'] as num).toInt(),
      killed: json['killed'] as bool,
    );

Map<String, dynamic> _$StructureFiresDataToJson(StructureFiresData instance) =>
    <String, dynamic>{
      'structure_id': instance.structureId,
      'target_id': instance.targetId,
      'hit_roll': instance.hitRoll,
      'damage': instance.damage,
      'target_hp': instance.targetHp,
      'killed': instance.killed,
    };

HealedUnit _$HealedUnitFromJson(Map<String, dynamic> json) => HealedUnit(
      unitId: json['unit_id'] as String,
      hpBefore: (json['hp_before'] as num).toInt(),
      hpAfter: (json['hp_after'] as num).toInt(),
    );

Map<String, dynamic> _$HealedUnitToJson(HealedUnit instance) =>
    <String, dynamic>{
      'unit_id': instance.unitId,
      'hp_before': instance.hpBefore,
      'hp_after': instance.hpAfter,
    };

StructureRegen _$StructureRegenFromJson(Map<String, dynamic> json) =>
    StructureRegen(
      structureId: json['structure_id'] as String,
      hpBefore: (json['hp_before'] as num).toInt(),
      hpAfter: (json['hp_after'] as num).toInt(),
    );

Map<String, dynamic> _$StructureRegenToJson(StructureRegen instance) =>
    <String, dynamic>{
      'structure_id': instance.structureId,
      'hp_before': instance.hpBefore,
      'hp_after': instance.hpAfter,
    };

SuddenDeathDamage _$SuddenDeathDamageFromJson(Map<String, dynamic> json) =>
    SuddenDeathDamage(
      unitId: json['unit_id'] as String,
      damage: (json['damage'] as num).toInt(),
      hpAfter: (json['hp_after'] as num).toInt(),
      killed: json['killed'] as bool,
    );

Map<String, dynamic> _$SuddenDeathDamageToJson(SuddenDeathDamage instance) =>
    <String, dynamic>{
      'unit_id': instance.unitId,
      'damage': instance.damage,
      'hp_after': instance.hpAfter,
      'killed': instance.killed,
    };

TurnStartData _$TurnStartDataFromJson(Map<String, dynamic> json) =>
    TurnStartData(
      turnNumber: (json['turn_number'] as num).toInt(),
      activePlayerId: json['active_player_id'] as String,
      timerSeconds: (json['timer_seconds'] as num).toInt(),
      incomeGained: (json['income_gained'] as num).toInt(),
      structureIncome: (json['structure_income'] as num).toInt(),
      totalCoins: (json['total_coins'] as num).toInt(),
      healedUnits: (json['healed_units'] as List<dynamic>?)
              ?.map((e) => HealedUnit.fromJson(e as Map<String, dynamic>))
              .toList() ??
          [],
      structureRegens: (json['structure_regens'] as List<dynamic>?)
              ?.map((e) => StructureRegen.fromJson(e as Map<String, dynamic>))
              .toList() ??
          [],
      suddenDeathDamages: (json['sudden_death_damage'] as List<dynamic>?)
              ?.map(
                  (e) => SuddenDeathDamage.fromJson(e as Map<String, dynamic>))
              .toList() ??
          [],
    );

Map<String, dynamic> _$TurnStartDataToJson(TurnStartData instance) =>
    <String, dynamic>{
      'turn_number': instance.turnNumber,
      'active_player_id': instance.activePlayerId,
      'timer_seconds': instance.timerSeconds,
      'income_gained': instance.incomeGained,
      'structure_income': instance.structureIncome,
      'total_coins': instance.totalCoins,
      'healed_units': instance.healedUnits,
      'structure_regens': instance.structureRegens,
      'sudden_death_damage': instance.suddenDeathDamages,
    };

GameOverData _$GameOverDataFromJson(Map<String, dynamic> json) => GameOverData(
      winnerId: json['winner_id'] as String,
      reason: $enumDecode(_$WinReasonEnumMap, json['reason']),
    );

Map<String, dynamic> _$GameOverDataToJson(GameOverData instance) =>
    <String, dynamic>{
      'winner_id': instance.winnerId,
      'reason': _$WinReasonEnumMap[instance.reason]!,
    };

const _$WinReasonEnumMap = {
  WinReason.hqDestroyed: 'HQ_DESTROYED',
  WinReason.structureDominance: 'STRUCTURE_DOMINANCE',
  WinReason.suddenDeath: 'SUDDEN_DEATH',
  WinReason.forfeit: 'FORFEIT',
  WinReason.disconnect: 'DISCONNECT',
  WinReason.draw: 'DRAW',
};

MatchFoundData _$MatchFoundDataFromJson(Map<String, dynamic> json) =>
    MatchFoundData(
      roomId: json['room_id'] as String,
    );

Map<String, dynamic> _$MatchFoundDataToJson(MatchFoundData instance) =>
    <String, dynamic>{
      'room_id': instance.roomId,
    };
