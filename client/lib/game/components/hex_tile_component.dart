import 'dart:ui';
import 'package:flame/components.dart';
import '../../models/enums.dart';
import '../hex/cube_coord.dart';
import '../hex/hex_layout.dart';
import 'dart:math' as math;

enum HighlightType { none, move, attack }

class HexTileComponent extends PositionComponent {
  final CubeCoord coord;
  final TerrainType terrain;
  final HexLayout layout;

  HighlightType highlight = HighlightType.none;

  HexTileComponent({
    required this.coord,
    required this.terrain,
    required this.layout,
  }) {
    final pos = layout.hexToPixel(coord);
    position = Vector2(pos.dx, pos.dy);
    anchor = Anchor.center;
    size = Vector2(layout.hexSize * 2, layout.hexSize * 2);
  }

  @override
  void render(Canvas canvas) {
    super.render(canvas);

    // Render hex shape
    final path = Path();
    for (int i = 0; i < 6; i++) {
      final angle = 2 * math.pi / 6 * (i + 0.5);
      final x = layout.hexSize + layout.hexSize * math.cos(angle);
      final y = layout.hexSize + layout.hexSize * math.sin(angle);
      if (i == 0) {
        path.moveTo(x, y);
      } else {
        path.lineTo(x, y);
      }
    }
    path.close();

    Color terrainColor;
    switch (terrain) {
      case TerrainType.plains:
        terrainColor = const Color(0xFFF0E68C); // Yellow (Khaki)
        break;
      case TerrainType.forest:
        terrainColor = const Color(0xFF228B22); // Green (ForestGreen)
        break;
      case TerrainType.hills:
        terrainColor = const Color(0xFFFFA500); // Orange
        break;
      case TerrainType.water:
        terrainColor = const Color(0xFF4169E1); // Blue (RoyalBlue)
        break;
      case TerrainType.mountains:
        terrainColor = const Color(0xFF8B4513); // Brown (SaddleBrown)
        break;
    }

    // Fill
    canvas.drawPath(
        path,
        Paint()
          ..color = terrainColor
          ..style = PaintingStyle.fill);

    // Stroke
    canvas.drawPath(
        path,
        Paint()
          ..color = const Color(0xFF000000)
          ..style = PaintingStyle.stroke
          ..strokeWidth = 1);

    // Render highlight overlay
    if (highlight != HighlightType.none) {
      final highlightColor = highlight == HighlightType.move
          ? const Color(0x660000FF)
          : const Color(0x66FF0000);
      canvas.drawPath(
          path,
          Paint()
            ..color = highlightColor
            ..style = PaintingStyle.fill);
    }
  }
}
