import 'package:json_annotation/json_annotation.dart';
import 'enums.dart';

part 'messages.g.dart';

@JsonSerializable()
class ErrorData {
  final ErrorCode code;
  final String message;

  const ErrorData({required this.code, required this.message});

  factory ErrorData.fromJson(Map<String, dynamic> json) =>
      _$ErrorDataFromJson(json);
  Map<String, dynamic> toJson() => _$ErrorDataToJson(this);
}

@JsonSerializable()
class AckData {
  final int seq;
  @JsonKey(name: 'action_type')
  final String actionType;

  const AckData({required this.seq, required this.actionType});

  factory AckData.fromJson(Map<String, dynamic> json) =>
      _$AckDataFromJson(json);
  Map<String, dynamic> toJson() => _$AckDataToJson(this);
}

@JsonSerializable()
class NackData {
  final int seq;
  @JsonKey(name: 'action_type')
  final String actionType;
  final ErrorData error;

  const NackData({
    required this.seq,
    required this.actionType,
    required this.error,
  });

  factory NackData.fromJson(Map<String, dynamic> json) =>
      _$NackDataFromJson(json);
  Map<String, dynamic> toJson() => _$NackDataToJson(this);
}

@JsonSerializable()
class TroopMovedData {
  @JsonKey(name: 'unit_id')
  final String unitId;
  @JsonKey(name: 'from_q')
  final int fromQ;
  @JsonKey(name: 'from_r')
  final int fromR;
  @JsonKey(name: 'from_s')
  final int fromS;
  @JsonKey(name: 'to_q')
  final int toQ;
  @JsonKey(name: 'to_r')
  final int toR;
  @JsonKey(name: 'to_s')
  final int toS;
  @JsonKey(name: 'remaining_mobility')
  final int remainingMobility;

  const TroopMovedData({
    required this.unitId,
    required this.fromQ,
    required this.fromR,
    required this.fromS,
    required this.toQ,
    required this.toR,
    required this.toS,
    required this.remainingMobility,
  });

  factory TroopMovedData.fromJson(Map<String, dynamic> json) =>
      _$TroopMovedDataFromJson(json);
  Map<String, dynamic> toJson() => _$TroopMovedDataToJson(this);
}

@JsonSerializable()
class CombatResultData {
  @JsonKey(name: 'attacker_id')
  final String attackerId;
  @JsonKey(name: 'defender_id')
  final String defenderId;
  @JsonKey(name: 'hit_roll')
  final int hitRoll;
  @JsonKey(name: 'natural_roll')
  final int naturalRoll;
  final bool hit;
  @JsonKey(name: 'damage_roll')
  final int damageRoll;
  final int damage;
  @JsonKey(name: 'defender_hp')
  final int defenderHp;
  final bool killed;
  final bool crit;
  final bool fumble;

  @JsonKey(name: 'has_counter')
  final bool hasCounter;
  @JsonKey(name: 'counter_hit_roll')
  final int? counterHitRoll;
  @JsonKey(name: 'counter_natural_roll')
  final int? counterNaturalRoll;
  @JsonKey(name: 'counter_hit')
  final bool? counterHit;
  @JsonKey(name: 'counter_damage')
  final int? counterDamage;

  @JsonKey(name: 'attacker_hp')
  final int attackerHp;
  @JsonKey(name: 'attacker_killed')
  final bool attackerKilled;

  const CombatResultData({
    required this.attackerId,
    required this.defenderId,
    required this.hitRoll,
    required this.naturalRoll,
    required this.hit,
    required this.damageRoll,
    required this.damage,
    required this.defenderHp,
    required this.killed,
    required this.crit,
    required this.fumble,
    required this.hasCounter,
    this.counterHitRoll,
    this.counterNaturalRoll,
    this.counterHit,
    this.counterDamage,
    required this.attackerHp,
    required this.attackerKilled,
  });

  factory CombatResultData.fromJson(Map<String, dynamic> json) =>
      _$CombatResultDataFromJson(json);
  Map<String, dynamic> toJson() => _$CombatResultDataToJson(this);
}

@JsonSerializable()
class TroopPurchasedData {
  @JsonKey(name: 'unit_id')
  final String unitId;
  @JsonKey(name: 'unit_type')
  final TroopType unitType;
  @JsonKey(name: 'hex_q')
  final int hexQ;
  @JsonKey(name: 'hex_r')
  final int hexR;
  @JsonKey(name: 'hex_s')
  final int hexS;
  final String owner;
  @JsonKey(name: 'coins_remaining')
  final int coinsRemaining;

  const TroopPurchasedData({
    required this.unitId,
    required this.unitType,
    required this.hexQ,
    required this.hexR,
    required this.hexS,
    required this.owner,
    required this.coinsRemaining,
  });

  factory TroopPurchasedData.fromJson(Map<String, dynamic> json) =>
      _$TroopPurchasedDataFromJson(json);
  Map<String, dynamic> toJson() => _$TroopPurchasedDataToJson(this);
}

@JsonSerializable()
class TroopDestroyedData {
  @JsonKey(name: 'unit_id')
  final String unitId;
  @JsonKey(name: 'hex_q')
  final int hexQ;
  @JsonKey(name: 'hex_r')
  final int hexR;
  @JsonKey(name: 'hex_s')
  final int hexS;
  final String cause;

  const TroopDestroyedData({
    required this.unitId,
    required this.hexQ,
    required this.hexR,
    required this.hexS,
    required this.cause,
  });

  factory TroopDestroyedData.fromJson(Map<String, dynamic> json) =>
      _$TroopDestroyedDataFromJson(json);
  Map<String, dynamic> toJson() => _$TroopDestroyedDataToJson(this);
}

@JsonSerializable()
class StructureAttackedData {
  @JsonKey(name: 'structure_id')
  final String structureId;
  @JsonKey(name: 'attacker_id')
  final String attackerId;
  @JsonKey(name: 'hit_roll')
  final int hitRoll;
  final int damage;
  @JsonKey(name: 'structure_hp')
  final int structureHp;
  final bool captured;
  @JsonKey(name: 'new_owner')
  final String? newOwner;

  const StructureAttackedData({
    required this.structureId,
    required this.attackerId,
    required this.hitRoll,
    required this.damage,
    required this.structureHp,
    required this.captured,
    this.newOwner,
  });

  factory StructureAttackedData.fromJson(Map<String, dynamic> json) =>
      _$StructureAttackedDataFromJson(json);
  Map<String, dynamic> toJson() => _$StructureAttackedDataToJson(this);
}

@JsonSerializable()
class StructureFiresData {
  @JsonKey(name: 'structure_id')
  final String structureId;
  @JsonKey(name: 'target_id')
  final String targetId;
  @JsonKey(name: 'hit_roll')
  final int hitRoll;
  final int damage;
  @JsonKey(name: 'target_hp')
  final int targetHp;
  final bool killed;

  const StructureFiresData({
    required this.structureId,
    required this.targetId,
    required this.hitRoll,
    required this.damage,
    required this.targetHp,
    required this.killed,
  });

  factory StructureFiresData.fromJson(Map<String, dynamic> json) =>
      _$StructureFiresDataFromJson(json);
  Map<String, dynamic> toJson() => _$StructureFiresDataToJson(this);
}

@JsonSerializable()
class HealedUnit {
  @JsonKey(name: 'unit_id')
  final String unitId;
  @JsonKey(name: 'hp_before')
  final int hpBefore;
  @JsonKey(name: 'hp_after')
  final int hpAfter;

  const HealedUnit({
    required this.unitId,
    required this.hpBefore,
    required this.hpAfter,
  });

  factory HealedUnit.fromJson(Map<String, dynamic> json) =>
      _$HealedUnitFromJson(json);
  Map<String, dynamic> toJson() => _$HealedUnitToJson(this);
}

@JsonSerializable()
class StructureRegen {
  @JsonKey(name: 'structure_id')
  final String structureId;
  @JsonKey(name: 'hp_before')
  final int hpBefore;
  @JsonKey(name: 'hp_after')
  final int hpAfter;

  const StructureRegen({
    required this.structureId,
    required this.hpBefore,
    required this.hpAfter,
  });

  factory StructureRegen.fromJson(Map<String, dynamic> json) =>
      _$StructureRegenFromJson(json);
  Map<String, dynamic> toJson() => _$StructureRegenToJson(this);
}

@JsonSerializable()
class SuddenDeathDamage {
  @JsonKey(name: 'unit_id')
  final String unitId;
  final int damage;
  @JsonKey(name: 'hp_after')
  final int hpAfter;
  final bool killed;

  const SuddenDeathDamage({
    required this.unitId,
    required this.damage,
    required this.hpAfter,
    required this.killed,
  });

  factory SuddenDeathDamage.fromJson(Map<String, dynamic> json) =>
      _$SuddenDeathDamageFromJson(json);
  Map<String, dynamic> toJson() => _$SuddenDeathDamageToJson(this);
}

@JsonSerializable()
class TurnStartData {
  @JsonKey(name: 'turn_number')
  final int turnNumber;
  @JsonKey(name: 'active_player_id')
  final String activePlayerId;
  @JsonKey(name: 'timer_seconds')
  final int timerSeconds;
  @JsonKey(name: 'income_gained')
  final int incomeGained;
  @JsonKey(name: 'structure_income')
  final int structureIncome;
  @JsonKey(name: 'total_coins')
  final int totalCoins;
  @JsonKey(name: 'healed_units', defaultValue: [])
  final List<HealedUnit> healedUnits;
  @JsonKey(name: 'structure_regens', defaultValue: [])
  final List<StructureRegen> structureRegens;
  @JsonKey(name: 'sudden_death_damage', defaultValue: [])
  final List<SuddenDeathDamage> suddenDeathDamages;

  const TurnStartData({
    required this.turnNumber,
    required this.activePlayerId,
    required this.timerSeconds,
    required this.incomeGained,
    required this.structureIncome,
    required this.totalCoins,
    required this.healedUnits,
    required this.structureRegens,
    required this.suddenDeathDamages,
  });

  factory TurnStartData.fromJson(Map<String, dynamic> json) =>
      _$TurnStartDataFromJson(json);
  Map<String, dynamic> toJson() => _$TurnStartDataToJson(this);
}

@JsonSerializable()
class GameOverData {
  @JsonKey(name: 'winner_id')
  final String winnerId;
  final WinReason reason;
  // Note: in Go it's map[string]model.GameOverStats. In Dart we'll need to use dynamic if we didn't import PlayerState's GameOverStats.
  // Let's import player_state.dart
  // Wait, I will use dynamic here to avoid cyclical dependency, or just import it.

  const GameOverData({required this.winnerId, required this.reason});

  factory GameOverData.fromJson(Map<String, dynamic> json) =>
      _$GameOverDataFromJson(json);
  Map<String, dynamic> toJson() => _$GameOverDataToJson(this);
}

@JsonSerializable()
class MatchFoundData {
  @JsonKey(name: 'room_id')
  final String roomId;

  const MatchFoundData({required this.roomId});

  factory MatchFoundData.fromJson(Map<String, dynamic> json) =>
      _$MatchFoundDataFromJson(json);
  Map<String, dynamic> toJson() => _$MatchFoundDataToJson(this);
}
