import 'package:json_annotation/json_annotation.dart';
import '../game/hex/cube_coord.dart';
import 'enums.dart';

part 'structure.g.dart';

@JsonSerializable()
class Structure {
  final String id;
  final StructureType type;
  @JsonKey(name: 'owner_id')
  final String ownerId;
  @JsonKey(fromJson: CubeCoord.fromJson)
  final CubeCoord hex;
  @JsonKey(name: 'current_hp')
  final int currentHp;
  @JsonKey(name: 'max_hp')
  final int maxHp;
  final int atk;
  final int def;
  final int range;
  final String damage;
  final int income;
  @JsonKey(name: 'can_spawn')
  final bool canSpawn;

  const Structure({
    required this.id,
    required this.type,
    required this.ownerId,
    required this.hex,
    required this.currentHp,
    required this.maxHp,
    required this.atk,
    required this.def,
    required this.range,
    required this.damage,
    required this.income,
    required this.canSpawn,
  });

  bool get isNeutral => ownerId.isEmpty;
  bool isOwnedBy(String playerId) => ownerId == playerId;
  bool get isAlive => currentHp > 0;

  factory Structure.fromJson(Map<String, dynamic> json) =>
      _$StructureFromJson(json);
  Map<String, dynamic> toJson() => _$StructureToJson(this);

  Structure copyWith({
    String? id,
    StructureType? type,
    String? ownerId,
    CubeCoord? hex,
    int? currentHp,
    int? maxHp,
    int? atk,
    int? def,
    int? range,
    String? damage,
    int? income,
    bool? canSpawn,
  }) {
    return Structure(
      id: id ?? this.id,
      type: type ?? this.type,
      ownerId: ownerId ?? this.ownerId,
      hex: hex ?? this.hex,
      currentHp: currentHp ?? this.currentHp,
      maxHp: maxHp ?? this.maxHp,
      atk: atk ?? this.atk,
      def: def ?? this.def,
      range: range ?? this.range,
      damage: damage ?? this.damage,
      income: income ?? this.income,
      canSpawn: canSpawn ?? this.canSpawn,
    );
  }
}
