#!/usr/bin/env python3
"""Generate the Hex & Dice web icon set (favicon + PWA icons).

Run from the repository root:

    python3 tools/generate_icons.py

Requires Pillow. The logo is a pointy-top hex tile containing a die with red
pips, matching the game's name and terrain palette.
"""

from __future__ import annotations

import os
from PIL import Image, ImageDraw

ROOT = os.path.dirname(os.path.dirname(os.path.abspath(__file__)))
WEB_DIR = os.path.join(ROOT, "client", "web")
ICON_DIR = os.path.join(WEB_DIR, "icons")

TEAL = (70, 95, 118, 255)          # matches the plains tile palette
BORDER = (44, 23, 48, 255)         # deep purple, matches the mountain tile
WHITE = (255, 255, 255, 255)
RED = (229, 57, 53, 255)
BG = (26, 26, 46, 255)             # app background (maskable icons)


def draw_logo(size: int, background: tuple[int, int, int, int] | None = None,
              pad_ratio: float = 0.0) -> Image.Image:
    img = Image.new("RGBA", (size, size), background or (0, 0, 0, 0))
    d = ImageDraw.Draw(img)
    s = size / 64.0

    def px(v: float) -> float:
        return v * s

    # Optional padding shrinks the artwork for maskable safe zones.
    pad = pad_ratio * 64.0

    def P(v: float) -> float:
        return (v - pad) / (64.0 - 2 * pad) * 64.0

    hex_pts = [(32, 4), (58, 19), (58, 45), (32, 60), (6, 45), (6, 19)]
    inner_pts = [(32, 7), (55, 20), (55, 44), (32, 57), (9, 44), (9, 20)]

    def scale_pts(pts):
        return [(px(P(x)), px(P(y))) for x, y in pts]

    d.polygon(scale_pts(hex_pts), fill=BORDER)
    d.polygon(scale_pts(inner_pts), fill=TEAL)

    # Die
    x0, y0, x1, y1 = 22, 22, 42, 42
    d.rounded_rectangle(
        [px(P(x0 - 2)), px(P(y0 - 2)), px(P(x1 + 2)), px(P(y1 + 2))],
        radius=px(4), fill=BORDER)
    d.rounded_rectangle(
        [px(P(x0)), px(P(y0)), px(P(x1)), px(P(y1))],
        radius=px(3), fill=WHITE)

    # Five pips
    cx, cy, off, r = 32, 32, 6, 2.2
    for dx, dy in [(-1, -1), (1, -1), (0, 0), (-1, 1), (1, 1)]:
        cxp = px(P(cx + dx * off))
        cyp = px(P(cy + dy * off))
        rr = px(r)
        d.ellipse([cxp - rr, cyp - rr, cxp + rr, cyp + rr], fill=RED)
    return img


def main() -> None:
    os.makedirs(ICON_DIR, exist_ok=True)

    draw_logo(64).save(os.path.join(WEB_DIR, "favicon.png"))
    print("generated web/favicon.png")

    for size in (192, 512):
        draw_logo(size).save(os.path.join(ICON_DIR, f"Icon-{size}.png"))
        print(f"generated web/icons/Icon-{size}.png")

    for size in (192, 512):
        draw_logo(size, background=BG, pad_ratio=0.18).save(
            os.path.join(ICON_DIR, f"Icon-maskable-{size}.png"))
        print(f"generated web/icons/Icon-maskable-{size}.png")


if __name__ == "__main__":
    main()
