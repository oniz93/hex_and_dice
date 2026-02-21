import 'dart:ui';
import 'package:flame/components.dart';
import '../../models/troop.dart';
import '../hex/hex_layout.dart';

class TroopComponent extends PositionComponent {
  final Troop troop;
  final HexLayout layout;

  TroopComponent({required this.troop, required this.layout}) {
    final pos = layout.hexToPixel(troop.hex);
    position = Vector2(pos.dx, pos.dy);
    anchor = Anchor.center;
    size = Vector2(24, 24);
    priority = 5;
  }

  @override
  void render(Canvas canvas) {
    super.render(canvas);

    // Draw unit circle based on owner
    // A real implementation would use sprites and check game state for current player
    // For now we'll just draw a white circle
    canvas.drawCircle(
      Offset(size.x / 2, size.y / 2),
      size.x / 2,
      Paint()..color = const Color(0xFFFFFFFF)..style = PaintingStyle.fill,
    );
    canvas.drawCircle(
      Offset(size.x / 2, size.y / 2),
      size.x / 2,
      Paint()..color = const Color(0xFF000000)..style = PaintingStyle.stroke..strokeWidth = 2,
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
        ..style = PaintingStyle.stroke,
    );
  }
}
