import 'package:flutter/material.dart';
import 'package:flame/components.dart';
import '../../models/troop.dart';
import '../../models/enums.dart';
import '../hex/hex_layout.dart';

class TroopComponent extends PositionComponent {
  Troop troop;
  final HexLayout layout;
  Color teamColor;

  TroopComponent({
    required this.troop,
    required this.layout,
    required this.teamColor,
  }) {
    _updatePosition();
    anchor = Anchor.center;
    size = Vector2(24, 24);
    priority = 5;
  }

  void _updatePosition() {
    final pos = layout.hexToPixel(troop.hex);
    position = Vector2(pos.dx, pos.dy);
  }

  void updateTroop(Troop newTroop, Color newColor) {
    troop = newTroop;
    teamColor = newColor;
    _updatePosition();
  }

  @override
  void render(Canvas canvas) {
    super.render(canvas);

    // Draw unit circle
    canvas.drawCircle(
      Offset(size.x / 2, size.y / 2),
      size.x / 2,
      Paint()
        ..color = const Color(0xFFFFFFFF)
        ..style = PaintingStyle.fill,
    );
    canvas.drawCircle(
      Offset(size.x / 2, size.y / 2),
      size.x / 2,
      Paint()
        ..color = const Color(0xFF000000)
        ..style = PaintingStyle.stroke
        ..strokeWidth = 1,
    );

    // Render troop letter
    String letter = '';
    switch (troop.type) {
      case TroopType.marine:
        letter = 'M';
        break;
      case TroopType.sniper:
        letter = 'S';
        break;
      case TroopType.hoverbike:
        letter = 'H';
        break;
      case TroopType.mech:
        letter = 'R';
        break;
    }

    final textPainter = TextPainter(
      text: TextSpan(
        text: letter,
        style: TextStyle(
          color: teamColor,
          fontSize: 18,
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
    final hpPct = troop.currentHp / troop.maxHp;
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
