import 'package:flutter/material.dart';
import 'package:flame/components.dart';
import '../../models/structure.dart';
import '../../models/enums.dart';
import '../hex/hex_layout.dart';

class StructureComponent extends PositionComponent {
  Structure structure;
  final HexLayout layout;
  Color teamColor;

  StructureComponent({
    required this.structure,
    required this.layout,
    required this.teamColor,
  }) {
    _updatePosition();
    anchor = Anchor.center;
    size = Vector2(28, 28);
    priority = 4;
  }

  void _updatePosition() {
    final pos = layout.hexToPixel(structure.hex);
    position = Vector2(pos.dx, pos.dy);
  }

  void updateStructure(Structure newStructure, Color newColor) {
    structure = newStructure;
    teamColor = newColor;
    _updatePosition();
  }

  @override
  void render(Canvas canvas) {
    super.render(canvas);

    // Draw structure rect
    canvas.drawRect(
      Rect.fromLTWH(0, 0, size.x, size.y),
      Paint()
        ..color = const Color(0xFFEEEEEE)
        ..style = PaintingStyle.fill,
    );
    canvas.drawRect(
      Rect.fromLTWH(0, 0, size.x, size.y),
      Paint()
        ..color = const Color(0xFF000000)
        ..style = PaintingStyle.stroke
        ..strokeWidth = 1,
    );

    // Render structure letter
    String letter = '';
    switch (structure.type) {
      case StructureType.hq:
        letter = 'Q';
        break;
      case StructureType.outpost:
        letter = 'O';
        break;
      case StructureType.commandCenter:
        letter = 'C';
        break;
    }

    final textPainter = TextPainter(
      text: TextSpan(
        text: letter,
        style: TextStyle(
          color: teamColor,
          fontSize: 20,
          fontWeight: FontWeight.bold,
        ),
      ),
      textDirection: TextDirection.ltr,
    );
    textPainter.layout();
    textPainter.paint(
      canvas,
      Offset(
        (size.x - textPainter.width) / 2,
        (size.y - textPainter.height) / 2,
      ),
    );

    // Render HP bar
    final hpPct = structure.currentHp / structure.maxHp;
    final paint = Paint()..color = const Color(0xFF00FF00);
    if (hpPct < 0.5) paint.color = const Color(0xFFFFFF00);
    if (hpPct < 0.25) paint.color = const Color(0xFFFF0000);

    canvas.drawRect(Rect.fromLTWH(0, size.y + 2, size.x * hpPct, 4), paint);
    canvas.drawRect(
      Rect.fromLTWH(0, size.y + 2, size.x, 4),
      Paint()
        ..color = const Color(0xFF000000)
        ..style = PaintingStyle.stroke
        ..strokeWidth = 1,
    );
  }
}
