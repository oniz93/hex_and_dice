// GENERATED CODE - DO NOT MODIFY BY HAND

part of 'structure.dart';

// **************************************************************************
// JsonSerializableGenerator
// **************************************************************************

Structure _$StructureFromJson(Map<String, dynamic> json) => Structure(
      id: json['id'] as String,
      type: $enumDecode(_$StructureTypeEnumMap, json['type']),
      ownerId: json['owner_id'] as String,
      hex: CubeCoord.fromJson(json['hex']),
      currentHp: (json['current_hp'] as num).toInt(),
      maxHp: (json['max_hp'] as num).toInt(),
      atk: (json['atk'] as num).toInt(),
      def: (json['def'] as num).toInt(),
      range: (json['range'] as num).toInt(),
      damage: json['damage'] as String,
      income: (json['income'] as num).toInt(),
      canSpawn: json['can_spawn'] as bool,
    );

Map<String, dynamic> _$StructureToJson(Structure instance) => <String, dynamic>{
      'id': instance.id,
      'type': _$StructureTypeEnumMap[instance.type]!,
      'owner_id': instance.ownerId,
      'hex': instance.hex,
      'current_hp': instance.currentHp,
      'max_hp': instance.maxHp,
      'atk': instance.atk,
      'def': instance.def,
      'range': instance.range,
      'damage': instance.damage,
      'income': instance.income,
      'can_spawn': instance.canSpawn,
    };

const _$StructureTypeEnumMap = {
  StructureType.outpost: 'outpost',
  StructureType.commandCenter: 'command_center',
  StructureType.hq: 'hq',
};
