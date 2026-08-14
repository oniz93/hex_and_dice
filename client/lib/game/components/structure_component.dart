import 'package:flutter/material.dart';
import 'package:flame/components.dart';
import '../../models/structure.dart';
import '../../models/enums.dart';
import '../hex/hex_layout.dart';

class StructureComponent extends PositionComponent {
  Structure structure;
  final HexLayout layout;
  final Sprite? sprite;
  Color teamColor;

  // Flash animation state
  bool _isFlashing = false;
  double _flashTimer = 0;
  VoidCallback? _onFlashComplete;
  static const double _flashDuration = 0.5; // 2 full red flashes
  static const double _flashHalfCycle = 0.125; // each flash on/off cycle

  StructureComponent({
    required this.structure,
    required this.layout,
    required this.teamColor,
    this.sprite,
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

  /// Start a red flash animation (2 flashes). Calls [onComplete] when done.
  void startFlash(VoidCallback? onComplete) {
    _isFlashing = true;
    _flashTimer = 0;
    _onFlashComplete = onComplete;
  }

  bool get isFlashing => _isFlashing;

  @override
  void update(double dt) {
    super.update(dt);
    if (_isFlashing) {
      _flashTimer += dt;
      if (_flashTimer >= _flashDuration) {
        _isFlashing = false;
        _flashTimer = 0;
        _onFlashComplete?.call();
        _onFlashComplete = null;
      }
    }
  }

  /// Returns true if the component should render with a red tint right now.
  bool get _showRedFlash {
    if (!_isFlashing) return false;
    final halfCycleIndex = (_flashTimer / _flashHalfCycle).floor();
    return halfCycleIndex % 2 == 0;
  }

  @override
  void render(Canvas canvas) {
    super.render(canvas);

    final isRed = _showRedFlash;

    if (sprite != null) {
      final tint = isRed ? Colors.red : teamColor;
      sprite!.render(
        canvas,
        size: size,
        overridePaint: Paint()
          ..colorFilter = ColorFilter.mode(tint, BlendMode.modulate),
      );
    } else {
      // Fallback rendering for when a sprite is missing: colored rect and a
      // structure letter.
      canvas.drawRect(
        Rect.fromLTWH(0, 0, size.x, size.y),
        Paint()
          ..color = isRed ? const Color(0xFFFF0000) : const Color(0xFFEEEEEE)
          ..style = PaintingStyle.fill,
      );
      canvas.drawRect(
        Rect.fromLTWH(0, 0, size.x, size.y),
        Paint()
          ..color = const Color(0xFF000000)
          ..style = PaintingStyle.stroke
          ..strokeWidth = 1,
      );

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
            color: isRed ? Colors.white : teamColor,
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
    }

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
