import 'dart:ui';
import 'package:flame/events.dart';
import '../hex/hex_layout.dart';
import '../hex/cube_coord.dart';

class InputSystem {
  final HexLayout layout;
  final Function(CubeCoord) onHexTap;

  InputSystem({required this.layout, required this.onHexTap});

  void handleTap(TapDownEvent event) {
    final pos = event.canvasPosition;
    final hex = layout.pixelToHex(Offset(pos.x, pos.y));
    onHexTap(hex);
  }
}
