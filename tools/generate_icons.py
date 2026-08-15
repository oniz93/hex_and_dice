#!/usr/bin/env python3
"""Generate the Hex & Dice icon set for web, Android, and iOS.

Run from the repository root:

    python3 tools/generate_icons.py

Requires Pillow. The logo is a pointy-top hex tile (matching the in-game
terrain palette) containing a white die with red pips.
"""

from __future__ import annotations

import os
from PIL import Image, ImageDraw

ROOT = os.path.dirname(os.path.dirname(os.path.abspath(__file__)))
CLIENT = os.path.join(ROOT, "client")
WEB_DIR = os.path.join(CLIENT, "web")
WEB_ICONS = os.path.join(WEB_DIR, "icons")
ANDROID_RES = os.path.join(CLIENT, "android", "app", "src", "main", "res")
IOS_ICONS = os.path.join(
    CLIENT, "ios", "Runner", "Assets.xcassets", "AppIcon.appiconset")

# Palette (hex fill matches the in-game plains tile).
HEX_FILL = (124, 194, 74, 255)
HEX_BORDER = (46, 94, 42, 255)
WHITE = (255, 255, 255, 255)
RED = (229, 57, 53, 255)
DARK = (30, 32, 40, 255)
BG = (26, 26, 46, 255)  # app background, also manifest theme_color


def draw_logo(size: int, background: tuple[int, int, int, int] | None = None,
              pad_ratio: float = 0.0) -> Image.Image:
    img = Image.new("RGBA", (size, size), background or (0, 0, 0, 0))
    d = ImageDraw.Draw(img)
    s = size / 64.0

    def px(v: float) -> float:
        return v * s

    pad = pad_ratio * 64.0

    def P(v: float) -> float:
        return (v - pad) / (64.0 - 2 * pad) * 64.0

    def scale_pts(pts):
        return [(px(P(x)), px(P(y))) for x, y in pts]

    d.polygon(scale_pts([(32, 4), (58, 19), (58, 45), (32, 60), (6, 45), (6, 19)]),
              fill=HEX_BORDER)
    d.polygon(scale_pts([(32, 7), (55, 20), (55, 44), (32, 57), (9, 44), (9, 20)]),
              fill=HEX_FILL)

    # Die
    x0, y0, x1, y1 = 22, 22, 42, 42
    d.rounded_rectangle(
        [px(P(x0 - 2)), px(P(y0 - 2)), px(P(x1 + 2)), px(P(y1 + 2))],
        radius=px(4), fill=DARK)
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


def save(img: Image.Image, path: str, opaque: bool = False) -> None:
    if opaque:
        img = img.convert("RGB")
    img.save(path)
    print(f"generated {os.path.relpath(path, ROOT)}")


def main() -> None:
    # --- Web ---
    os.makedirs(WEB_ICONS, exist_ok=True)
    save(draw_logo(64), os.path.join(WEB_DIR, "favicon.png"))
    for size in (192, 512):
        save(draw_logo(size), os.path.join(WEB_ICONS, f"Icon-{size}.png"))
        save(draw_logo(size, background=BG, pad_ratio=0.18),
             os.path.join(WEB_ICONS, f"Icon-maskable-{size}.png"))

    # --- Android launcher icons (legacy square icons) ---
    android_sizes = {
        "mipmap-mdpi": 48,
        "mipmap-hdpi": 72,
        "mipmap-xhdpi": 96,
        "mipmap-xxhdpi": 144,
        "mipmap-xxxhdpi": 192,
    }
    for folder, size in android_sizes.items():
        path = os.path.join(ANDROID_RES, folder, "ic_launcher.png")
        if os.path.exists(path):
            save(draw_logo(size, background=BG, pad_ratio=0.10), path)

    # --- iOS AppIcon set (must be opaque) ---
    ios_icons = {
        "Icon-App-20x20@1x.png": 20, "Icon-App-20x20@2x.png": 40,
        "Icon-App-20x20@3x.png": 60, "Icon-App-29x29@1x.png": 29,
        "Icon-App-29x29@2x.png": 58, "Icon-App-29x29@3x.png": 87,
        "Icon-App-40x40@1x.png": 40, "Icon-App-40x40@2x.png": 80,
        "Icon-App-40x40@3x.png": 120, "Icon-App-60x60@2x.png": 120,
        "Icon-App-60x60@3x.png": 180, "Icon-App-76x76@1x.png": 76,
        "Icon-App-76x76@2x.png": 152, "Icon-App-83.5x83.5@2x.png": 167,
        "Icon-App-1024x1024@1x.png": 1024,
    }
    for name, size in ios_icons.items():
        path = os.path.join(IOS_ICONS, name)
        if os.path.exists(path):
            save(draw_logo(size, background=BG, pad_ratio=0.10), path,
                 opaque=True)


if __name__ == "__main__":
    main()
