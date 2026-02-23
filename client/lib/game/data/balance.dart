import '../../models/enums.dart';

class TroopStatBlock {
  final int cost;
  final int hp;
  final int atk;
  final int def;
  final int mobility;
  final int range;
  final String damage;

  const TroopStatBlock({
    required this.cost,
    required this.hp,
    required this.atk,
    required this.def,
    required this.mobility,
    required this.range,
    required this.damage,
  });
}

const troopStats = {
  TroopType.marine: TroopStatBlock(
    cost: 100,
    hp: 10,
    atk: 3,
    def: 14,
    mobility: 3,
    range: 1,
    damage: '1D6+1',
  ),
  TroopType.sniper: TroopStatBlock(
    cost: 150,
    hp: 6,
    atk: 4,
    def: 11,
    mobility: 2,
    range: 3,
    damage: '1D8',
  ),
  TroopType.hoverbike: TroopStatBlock(
    cost: 200,
    hp: 8,
    atk: 4,
    def: 12,
    mobility: 5,
    range: 1,
    damage: '1D8+1',
  ),
  TroopType.mech: TroopStatBlock(
    cost: 350,
    hp: 12,
    atk: 5,
    def: 10,
    mobility: 1,
    range: 3,
    damage: '2D6+2',
  ),
};
