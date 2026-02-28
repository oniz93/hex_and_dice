import 'dart:math';
import 'package:flutter/material.dart';
import 'package:flame/components.dart';

/// Animated projectile arrow that travels from attacker to defender position.
/// Renders as a bright yellow "bullet" with a fading trail, simulating
/// an ammunition visual effect.
class AttackArrowComponent extends PositionComponent {
  final Vector2 from;
  final Vector2 to;
  final double duration; // seconds for the projectile to travel
  final VoidCallback? onComplete;

  double _elapsed = 0;
  bool _done = false;

  // Trail history: list of recent positions for the tail effect
  final List<Vector2> _trail = [];
  static const int _maxTrailLength = 8;

  AttackArrowComponent({
    required this.from,
    required this.to,
    this.duration = 0.35,
    this.onComplete,
  }) {
    priority = 100; // Draw above everything
  }

  @override
  void update(double dt) {
    super.update(dt);
    if (_done) return;

    _elapsed += dt;
    final t = (_elapsed / duration).clamp(0.0, 1.0);

    // Current projectile position (lerp from -> to)
    final current = Vector2(
      from.x + (to.x - from.x) * t,
      from.y + (to.y - from.y) * t,
    );

    // Record trail
    _trail.add(current.clone());
    if (_trail.length > _maxTrailLength) {
      _trail.removeAt(0);
    }

    if (t >= 1.0) {
      _done = true;
      // Keep visible briefly for the impact flash, then remove
      Future.delayed(const Duration(milliseconds: 150), () {
        removeFromParent();
        onComplete?.call();
      });
    }
  }

  @override
  void render(Canvas canvas) {
    super.render(canvas);

    final t = (_elapsed / duration).clamp(0.0, 1.0);
    final current = Vector2(
      from.x + (to.x - from.x) * t,
      from.y + (to.y - from.y) * t,
    );

    // Draw trail (fading segments)
    if (_trail.length >= 2) {
      for (int i = 1; i < _trail.length; i++) {
        final opacity = (i / _trail.length) * 0.7;
        final width = 1.5 + (i / _trail.length) * 2.0;
        final trailPaint = Paint()
          ..color = Color.fromRGBO(255, 215, 0, opacity)
          ..strokeWidth = width
          ..strokeCap = StrokeCap.round
          ..style = PaintingStyle.stroke;
        canvas.drawLine(
          Offset(_trail[i - 1].x, _trail[i - 1].y),
          Offset(_trail[i].x, _trail[i].y),
          trailPaint,
        );
      }
    }

    // Draw projectile "bullet" (bright yellow glow + white core)
    if (!_done) {
      // Outer glow
      canvas.drawCircle(
        Offset(current.x, current.y),
        6,
        Paint()
          ..color = const Color.fromRGBO(255, 215, 0, 0.5)
          ..maskFilter = const MaskFilter.blur(BlurStyle.normal, 4),
      );
      // Middle ring
      canvas.drawCircle(
        Offset(current.x, current.y),
        4,
        Paint()..color = const Color(0xFFFFFF00),
      );
      // White core
      canvas.drawCircle(
        Offset(current.x, current.y),
        2,
        Paint()..color = const Color(0xFFFFFFFF),
      );
    } else {
      // Impact flash at destination
      final flashOpacity =
          (1.0 - ((_elapsed - duration) / 0.15)).clamp(0.0, 1.0);
      canvas.drawCircle(
        Offset(to.x, to.y),
        10 * flashOpacity,
        Paint()
          ..color = Color.fromRGBO(255, 255, 0, flashOpacity * 0.6)
          ..maskFilter = const MaskFilter.blur(BlurStyle.normal, 6),
      );
    }

    // Draw small arrowhead at the projectile tip (pointing in travel direction)
    if (!_done) {
      final dx = to.x - from.x;
      final dy = to.y - from.y;
      final angle = atan2(dy, dx);
      const arrowLen = 8.0;
      const arrowSpread = 0.5; // radians

      final tipX = current.x;
      final tipY = current.y;
      final leftX = tipX - arrowLen * cos(angle - arrowSpread);
      final leftY = tipY - arrowLen * sin(angle - arrowSpread);
      final rightX = tipX - arrowLen * cos(angle + arrowSpread);
      final rightY = tipY - arrowLen * sin(angle + arrowSpread);

      final arrowPath = Path()
        ..moveTo(tipX, tipY)
        ..lineTo(leftX, leftY)
        ..lineTo(rightX, rightY)
        ..close();

      canvas.drawPath(
        arrowPath,
        Paint()
          ..color = const Color(0xFFFFFF00)
          ..style = PaintingStyle.fill,
      );
    }
  }
}
