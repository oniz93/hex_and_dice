import 'package:json_annotation/json_annotation.dart';
import 'dart:math' as math;

part 'cube_coord.g.dart';

@JsonSerializable()
class CubeCoord {
  final int q;
  final int r;
  final int s;

  const CubeCoord(this.q, this.r, this.s)
    : assert(q + r + s == 0, 'q+r+s must be 0');

  factory CubeCoord.qr(int q, int r) => CubeCoord(q, r, -q - r);

  factory CubeCoord.fromJson(Map<String, dynamic> json) =>
      _$CubeCoordFromJson(json);
  Map<String, dynamic> toJson() => _$CubeCoordToJson(this);

  static const CubeCoord origin = CubeCoord(0, 0, 0);

  CubeCoord operator +(CubeCoord other) {
    return CubeCoord(q + other.q, r + other.r, s + other.s);
  }

  CubeCoord operator -(CubeCoord other) {
    return CubeCoord(q - other.q, r - other.r, s - other.s);
  }

  CubeCoord operator *(int k) {
    return CubeCoord(q * k, r * k, s * k);
  }

  CubeCoord rotate180() {
    return CubeCoord(-q, -r, -s);
  }

  int distance(CubeCoord other) {
    return ((q - other.q).abs() + (r - other.r).abs() + (s - other.s).abs()) ~/
        2;
  }

  int get length => distance(origin);

  CubeCoord neighbor(int dir) {
    const dirs = [
      CubeCoord(1, 0, -1),
      CubeCoord(1, -1, 0),
      CubeCoord(0, -1, 1),
      CubeCoord(-1, 0, 1),
      CubeCoord(-1, 1, 0),
      CubeCoord(0, 1, -1),
    ];
    return this + dirs[dir];
  }

  List<CubeCoord> ring(int radius) {
    if (radius == 0) return [this];
    var current =
        this + CubeCoord(-1, 1, 0) * radius; // southwest direction equivalent
    final results = <CubeCoord>[];
    for (int i = 0; i < 6; i++) {
      for (int j = 0; j < radius; j++) {
        results.add(current);
        current = current.neighbor(i);
      }
    }
    return results;
  }

  List<CubeCoord> spiral(int radius) {
    final results = <CubeCoord>[this];
    for (int r = 1; r <= radius; r++) {
      results.addAll(ring(r));
    }
    return results;
  }

  @override
  bool operator ==(Object other) =>
      identical(this, other) ||
      other is CubeCoord &&
          runtimeType == other.runtimeType &&
          q == other.q &&
          r == other.r &&
          s == other.s;

  @override
  int get hashCode => q.hashCode ^ r.hashCode ^ s.hashCode;

  @override
  String toString() => '($q, $r, $s)';
}
