// GENERATED CODE - DO NOT MODIFY BY HAND

part of 'troop.dart';

// **************************************************************************
// JsonSerializableGenerator
// **************************************************************************

Troop _$TroopFromJson(Map<String, dynamic> json) => Troop(
      id: json['id'] as String,
      type: $enumDecode(_$TroopTypeEnumMap, json['type']),
      ownerId: json['owner_id'] as String,
      hex: CubeCoord.fromJson(json['hex']),
      currentHp: (json['current_hp'] as num).toInt(),
      maxHp: (json['max_hp'] as num).toInt(),
      atk: (json['atk'] as num).toInt(),
      def: (json['def'] as num).toInt(),
      mobility: (json['mobility'] as num).toInt(),
      range: (json['range'] as num).toInt(),
      damage: json['damage'] as String,
      isReady: json['is_ready'] as bool,
      hasMoved: json['has_moved'] as bool,
      hasAttacked: json['has_attacked'] as bool,
      wasInCombat: json['was_in_combat'] as bool,
      remainingMobility: (json['remaining_mobility'] as num).toInt(),
    );

Map<String, dynamic> _$TroopToJson(Troop instance) => <String, dynamic>{
      'id': instance.id,
      'type': _$TroopTypeEnumMap[instance.type]!,
      'owner_id': instance.ownerId,
      'hex': instance.hex,
      'current_hp': instance.currentHp,
      'max_hp': instance.maxHp,
      'atk': instance.atk,
      'def': instance.def,
      'mobility': instance.mobility,
      'range': instance.range,
      'damage': instance.damage,
      'is_ready': instance.isReady,
      'has_moved': instance.hasMoved,
      'has_attacked': instance.hasAttacked,
      'was_in_combat': instance.wasInCombat,
      'remaining_mobility': instance.remainingMobility,
    };

const _$TroopTypeEnumMap = {
  TroopType.marine: 'marine',
  TroopType.sniper: 'sniper',
  TroopType.hoverbike: 'hoverbike',
  TroopType.mech: 'mech',
};
