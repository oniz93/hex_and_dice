#!/usr/bin/env python3
"""Regenerate the Hex & Dice game sprites.

Run from the repository root:

    python3 tools/generate_sprites.py

Requires Pillow (``pip install Pillow``).

The troop and structure sprites are drawn in grayscale and tinted with the
player color at runtime by the Flame components, so white pixels become the
team color, gray pixels become darker shades, and near-black stays dark.

The hills tile is derived from the existing plains tile with a hue shift so it
keeps the same texture while reading as a different terrain type.
"""

from __future__ import annotations

import colorsys
import os
from PIL import Image, ImageDraw

ROOT = os.path.dirname(os.path.dirname(os.path.abspath(__file__)))
SPRITE_DIR = os.path.join(ROOT, "client", "assets", "images", "sprites")

# Grayscale palette (tinted at runtime).
OUTLINE = (30, 32, 40, 255)
BRIGHT = (255, 255, 255, 255)
LIGHT = (200, 200, 200, 255)
MID = (140, 140, 140, 255)
DARK = (80, 80, 80, 255)
BLACK = (10, 10, 10, 255)
TRANSPARENT = (0, 0, 0, 0)


def new_canvas(size: int = 32) -> tuple[Image.Image, ImageDraw.ImageDraw]:
    img = Image.new("RGBA", (size, size), TRANSPARENT)
    return img, ImageDraw.Draw(img)


def compose(parts: list[tuple[int, int, int, int, tuple[int, int, int, int]]],
            disks: list[tuple[int, int, int, tuple[int, int, int, int]]] | None = None,
            ) -> Image.Image:
    """Draw rectangles (and optional disks) with a 1px dark outline."""
    img, draw = new_canvas()
    disks = disks or []
    for x0, y0, x1, y1, _color in parts:
        draw.rectangle([x0 - 1, y0 - 1, x1 + 1, y1 + 1], fill=OUTLINE)
    for cx, cy, r, _color in disks:
        draw.ellipse([cx - r - 1, cy - r - 1, cx + r + 1, cy + r + 1], fill=OUTLINE)
    for x0, y0, x1, y1, color in parts:
        draw.rectangle([x0, y0, x1, y1], fill=color)
    for cx, cy, r, color in disks:
        draw.ellipse([cx - r, cy - r, cx + r, cy + r], fill=color)
    return img


def _rect(img: Image.Image, x0: int, y0: int, x1: int, y1: int,
          color: tuple[int, int, int, int]) -> None:
    ImageDraw.Draw(img).rectangle([x0, y0, x1, y1], fill=color)


def marine() -> Image.Image:
    img = compose(
        [
            (12, 4, 20, 9, BRIGHT),    # helmet
            (13, 10, 19, 11, MID),     # face
            (12, 12, 20, 19, BRIGHT),  # torso
            (9, 14, 12, 17, LIGHT),    # left arm
            (20, 14, 23, 17, LIGHT),   # right arm
            (22, 13, 30, 14, DARK),    # rifle barrel
            (23, 15, 25, 18, DARK),    # rifle magazine
            (12, 20, 15, 26, LIGHT),   # left leg
            (17, 20, 20, 26, LIGHT),   # right leg
            (11, 27, 16, 28, DARK),    # left boot
            (16, 27, 21, 28, DARK),    # right boot
        ],
        disks=[(16, 6, 4, BRIGHT)],
    )
    _rect(img, 14, 5, 18, 7, LIGHT)
    _rect(img, 14, 13, 18, 15, MID)
    return img


def sniper() -> Image.Image:
    img = compose(
        [
            (13, 5, 19, 9, BRIGHT),    # helmet
            (14, 10, 18, 11, MID),     # face
            (13, 12, 19, 17, BRIGHT),  # torso
            (12, 18, 15, 24, LIGHT),   # left leg
            (17, 18, 20, 24, LIGHT),   # right leg
            (11, 25, 16, 26, DARK),    # left boot
            (16, 25, 21, 26, DARK),    # right boot
            (10, 13, 13, 16, LIGHT),   # left arm
            (19, 13, 22, 16, LIGHT),   # right arm
            (21, 12, 30, 13, DARK),    # long rifle barrel
            (22, 14, 24, 16, DARK),    # scope
        ],
    )
    _rect(img, 15, 6, 17, 8, LIGHT)
    return img


def hoverbike() -> Image.Image:
    img = compose(
        [
            (10, 6, 16, 9, BRIGHT),    # rider helmet
            (11, 10, 15, 11, MID),     # rider face
            (10, 12, 16, 16, BRIGHT),  # rider torso
            (6, 17, 25, 19, DARK),     # bike body
            (7, 20, 13, 21, DARK),     # front fork
            (19, 20, 25, 21, DARK),    # rear fork
            (5, 22, 9, 24, DARK),      # front wheel
            (23, 22, 27, 24, DARK),    # rear wheel
        ],
        disks=[(5, 23, 1, MID), (9, 23, 1, MID),
               (23, 23, 1, MID), (27, 23, 1, MID)],
    )
    _rect(img, 12, 7, 14, 9, LIGHT)
    return img


def mech() -> Image.Image:
    img = compose(
        [
            (11, 3, 21, 8, BRIGHT),    # head
            (10, 9, 22, 19, BRIGHT),   # torso
            (8, 11, 11, 16, LIGHT),    # left arm
            (21, 11, 24, 16, LIGHT),   # right arm
            (9, 12, 12, 13, DARK),     # left cannon
            (20, 12, 23, 13, DARK),    # right cannon
            (12, 20, 16, 27, LIGHT),   # left leg
            (17, 20, 21, 27, LIGHT),   # right leg
            (11, 27, 17, 28, DARK),    # left foot
            (16, 27, 22, 28, DARK),    # right foot
        ],
    )
    _rect(img, 13, 5, 19, 7, MID)
    _rect(img, 12, 14, 20, 16, MID)
    _rect(img, 14, 11, 18, 13, DARK)
    return img


def outpost() -> Image.Image:
    img = compose(
        [
            (6, 16, 26, 27, BRIGHT),   # bunker body
            (4, 12, 28, 15, DARK),     # roof slab
            (8, 20, 24, 23, MID),      # door
            (5, 10, 8, 11, DARK),      # left antenna
            (24, 9, 26, 11, DARK),     # right antenna
        ],
    )
    _rect(img, 7, 17, 25, 19, LIGHT)
    return img


def command_center() -> Image.Image:
    img = compose(
        [
            (4, 14, 28, 27, BRIGHT),   # main building
            (3, 10, 29, 13, DARK),     # roof
            (7, 17, 25, 23, MID),      # door
            (10, 8, 22, 9, DARK),      # antenna bar
        ],
        disks=[(16, 8, 3, DARK)],
    )
    _rect(img, 5, 15, 27, 17, LIGHT)
    _rect(img, 11, 18, 21, 22, DARK)
    return img


def hq() -> Image.Image:
    img = compose(
        [
            (5, 12, 27, 28, BRIGHT),   # main base
            (2, 8, 10, 11, DARK),      # left tower
            (22, 8, 30, 11, DARK),     # right tower
            (4, 6, 8, 7, DARK),        # left antenna
            (24, 6, 28, 7, DARK),      # right antenna
            (9, 20, 23, 26, MID),      # gate
        ],
        disks=[(6, 5, 1, DARK), (26, 5, 1, DARK)],
    )
    _rect(img, 6, 13, 26, 15, LIGHT)
    _rect(img, 3, 9, 9, 10, MID)
    _rect(img, 23, 9, 29, 10, MID)
    _rect(img, 11, 21, 21, 25, DARK)
    return img


def hills_from_plains() -> Image.Image:
    """Create a warm hills tile from the plains texture."""
    plains_path = os.path.join(SPRITE_DIR, "plains.png")
    plains = Image.open(plains_path).convert("RGBA")
    pixels = plains.load()
    result = plains.copy()
    rp = result.load()
    for y in range(plains.height):
        for x in range(plains.width):
            r, g, b, a = pixels[x, y]
            if a == 0:
                continue
            h, s, v = colorsys.rgb_to_hsv(r / 255, g / 255, b / 255)
            # Warm orange hue, keeping the original luminance variation.
            h = 0.08 + (h - 0.5) * 0.05
            s = min(1.0, s * 1.15)
            nr, ng, nb = colorsys.hsv_to_rgb(h, s, v)
            rp[x, y] = (
                round(nr * 255),
                round(ng * 255),
                round(nb * 255),
                a,
            )
    return result


def main() -> None:
    os.makedirs(SPRITE_DIR, exist_ok=True)
    hills_from_plains().save(os.path.join(SPRITE_DIR, "hills.png"))
    units = {
        "troop_marine.png": marine,
        "troop_sniper.png": sniper,
        "troop_hoverbike.png": hoverbike,
        "troop_mech.png": mech,
        "structure_outpost.png": outpost,
        "structure_command_center.png": command_center,
        "structure_hq.png": hq,
    }
    for filename, factory in units.items():
        factory().save(os.path.join(SPRITE_DIR, filename))
        print(f"generated {filename}")


if __name__ == "__main__":
    main()
