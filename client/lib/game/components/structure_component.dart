import 'dart:ui';
import 'package:flame/components.dart';
import '../../models/structure.dart';
import '../hex/hex_layout.dart';

class StructureComponent extends PositionComponent {
  final Structure structure;
  final HexLayout layout;

  StructureComponent({required this.structure, required this.layout}) {
    final pos = layout.hexToPixel(structure.hex);
    position = Vector2(pos.dx, pos.dy);
    anchor = Anchor.center;
    size = Vector2(28, 28);
    priority = 4;
  }

  @override
  void render(Canvas canvas) {
    super.render(canvas);

    // Draw structure rect based on owner
    // For now we'll just draw a grey square
    canvas.drawRect(
      Rect.fromLTWH(0, 0, size.x, size.y),
      Paint()..color = const Color(0xFF888888)..style = PaintingStyle.fill,
    );
    canvas.drawRect(
      Rect.fromLTWH(0, 0, size.x, size.y),
      Paint()..color = const Color(0xFF000000)..style = PaintingStyle.stroke..strokeWidth = 2,
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
        ..style = PaintingStyle.stroke,
    );
  }
}
