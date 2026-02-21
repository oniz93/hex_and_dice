// GENERATED CODE - DO NOT MODIFY BY HAND

part of 'cube_coord.dart';

// **************************************************************************
// JsonSerializableGenerator
// **************************************************************************

CubeCoord _$CubeCoordFromJson(Map<String, dynamic> json) => CubeCoord(
      (json['q'] as num).toInt(),
      (json['r'] as num).toInt(),
      (json['s'] as num).toInt(),
    );

Map<String, dynamic> _$CubeCoordToJson(CubeCoord instance) => <String, dynamic>{
      'q': instance.q,
      'r': instance.r,
      's': instance.s,
    };
