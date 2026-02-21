import 'package:json_annotation/json_annotation.dart';
import '../game/hex/cube_coord.dart';
import 'enums.dart';

part 'troop.g.dart';

@JsonSerializable()
class Troop {
  final String id;
  final TroopType type;
  @JsonKey(name: 'owner_id')
  final String ownerId;
  final CubeCoord hex;
  @JsonKey(name: 'current_hp')
  final int currentHp;
  @JsonKey(name: 'max_hp')
  final int maxHp;
  final int atk;
  final int def;
  final int mobility;
  final int range;
  final String damage;
  @JsonKey(name: 'is_ready')
  final bool isReady;
  @JsonKey(name: 'has_moved')
  final bool hasMoved;
  @JsonKey(name: 'has_attacked')
  final bool hasAttacked;
  @JsonKey(name: 'was_in_combat')
  final bool wasInCombat;
  @JsonKey(name: 'remaining_mobility')
  final int remainingMobility;

  const Troop({
    required this.id,
    required this.type,
    required this.ownerId,
    required this.hex,
    required this.currentHp,
    required this.maxHp,
    required this.atk,
    required this.def,
    required this.mobility,
    required this.range,
    required this.damage,
    required this.isReady,
    required this.hasMoved,
    required this.hasAttacked,
    required this.wasInCombat,
    required this.remainingMobility,
  });

  bool get isAlive => currentHp > 0;
  bool get canAct => isAlive && isReady;
  bool get canMove => canAct && !hasMoved && remainingMobility > 0;
  bool get canAttack => canAct && !hasAttacked;
  bool get isMelee => range == 1;

  factory Troop.fromJson(Map<String, dynamic> json) => _$TroopFromJson(json);
  Map<String, dynamic> toJson() => _$TroopToJson(this);

  Troop copyWith({
    String? id,
    TroopType? type,
    String? ownerId,
    CubeCoord? hex,
    int? currentHp,
    int? maxHp,
    int? atk,
    int? def,
    int? mobility,
    int? range,
    String? damage,
    bool? isReady,
    bool? hasMoved,
    bool? hasAttacked,
    bool? wasInCombat,
    int? remainingMobility,
  }) {
    return Troop(
      id: id ?? this.id,
      type: type ?? this.type,
      ownerId: ownerId ?? this.ownerId,
      hex: hex ?? this.hex,
      currentHp: currentHp ?? this.currentHp,
      maxHp: maxHp ?? this.maxHp,
      atk: atk ?? this.atk,
      def: def ?? this.def,
      mobility: mobility ?? this.mobility,
      range: range ?? this.range,
      damage: damage ?? this.damage,
      isReady: isReady ?? this.isReady,
      hasMoved: hasMoved ?? this.hasMoved,
      hasAttacked: hasAttacked ?? this.hasAttacked,
      wasInCombat: wasInCombat ?? this.wasInCombat,
      remainingMobility: remainingMobility ?? this.remainingMobility,
    );
  }
}
