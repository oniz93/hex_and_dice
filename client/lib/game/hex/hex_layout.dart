import 'dart:math';
import 'package:flutter/material.dart';
import 'cube_coord.dart';

class HexLayout {
  final double hexSize; // Distance from center to vertex

  const HexLayout(this.hexSize);

  Offset hexToPixel(CubeCoord hex) {
    final x = hexSize * (sqrt(3) * hex.q + sqrt(3) / 2 * hex.r);
    final y = hexSize * (3.0 / 2 * hex.r);
    return Offset(x, y);
  }

  CubeCoord pixelToHex(Offset point) {
    final q = (sqrt(3) / 3 * point.dx - 1.0 / 3 * point.dy) / hexSize;
    final r = (2.0 / 3 * point.dy) / hexSize;
    return _cubeRound(q, r, -q - r);
  }

  CubeCoord _cubeRound(double fracQ, double fracR, double fracS) {
    int q = fracQ.round();
    int r = fracR.round();
    int s = fracS.round();

    final qDiff = (q - fracQ).abs();
    final rDiff = (r - fracR).abs();
    final sDiff = (s - fracS).abs();

    if (qDiff > rDiff && qDiff > sDiff) {
      q = -r - s;
    } else if (rDiff > sDiff) {
      r = -q - s;
    } else {
      s = -q - r;
    }

    return CubeCoord(q, r, s);
  }
}
