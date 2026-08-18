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

Rect = tuple[int, int, int, int, tuple[int, int, int, int]]
Disk = tuple[int, int, int, tuple[int, int, int, int]]


def compose(
    parts: list[Rect] = (),
    disks: list[Disk] = (),
    details: list[tuple] = (),
) -> Image.Image:
    """Draw rectangles (and optional disks) with a 1px dark outline, then
    apply detail primitives (rect/line/px/disk) on top."""
    img = Image.new("RGBA", (32, 32), TRANSPARENT)
    draw = ImageDraw.Draw(img)

    for x0, y0, x1, y1, _ in parts:
        draw.rectangle([x0 - 1, y0 - 1, x1 + 1, y1 + 1], fill=OUTLINE)
    for cx, cy, r, _ in disks:
        draw.ellipse([cx - r - 1, cy - r - 1, cx + r + 1, cy + r + 1], fill=OUTLINE)
    for x0, y0, x1, y1, color in parts:
        draw.rectangle([x0, y0, x1, y1], fill=color)
    for cx, cy, r, color in disks:
        draw.ellipse([cx - r, cy - r, cx + r, cy + r], fill=color)

    for detail in details:
        kind, *args = detail
        if kind == "rect":
            draw.rectangle(list(args[:-1]), fill=args[-1])
        elif kind == "line":
            draw.line(list(args[:-1]), fill=args[-1], width=1)
        elif kind == "px":
            draw.point(args[:2], fill=args[-1])
        elif kind == "disk":
            cx, cy, r, color = args
            draw.ellipse([cx - r, cy - r, cx + r, cy + r], fill=color)
        else:
            raise ValueError(f"unknown detail kind: {kind}")
    return img


def marine() -> Image.Image:
    return compose(
        parts=[
            (12, 6, 20, 9, BRIGHT),    # helmet base
            (12, 10, 20, 10, MID),     # face/neck line
            (11, 12, 21, 19, BRIGHT),  # torso
            (9, 13, 11, 17, LIGHT),    # left arm
            (21, 13, 23, 17, LIGHT),   # right arm
            (22, 12, 29, 13, DARK),    # rifle barrel
            (24, 14, 26, 18, DARK),    # rifle magazine
            (12, 20, 15, 26, LIGHT),   # left leg
            (17, 20, 20, 26, LIGHT),   # right leg
            (11, 27, 16, 28, DARK),    # left boot
            (16, 27, 21, 28, DARK),    # right boot
        ],
        disks=[(16, 5, 3, BRIGHT)],
        details=[
            ("rect", 13, 7, 19, 7, BLACK),    # visor slit
            ("rect", 13, 14, 19, 16, MID),    # chest plate
            ("rect", 14, 14, 18, 14, LIGHT),  # chest highlight
        ],
    )


def sniper() -> Image.Image:
    return compose(
        parts=[
            (14, 4, 19, 8, BRIGHT),    # hood
            (14, 9, 18, 10, MID),      # face line
            (13, 11, 19, 17, BRIGHT),  # slim torso
            (12, 18, 15, 24, LIGHT),   # left leg
            (17, 18, 20, 24, LIGHT),   # right leg
            (11, 25, 16, 26, DARK),    # left boot
            (16, 25, 21, 26, DARK),    # right boot
            (11, 12, 13, 15, LIGHT),   # left arm
            (19, 12, 21, 15, LIGHT),   # right arm (forward grip)
            (21, 11, 29, 12, DARK),    # long barrel
            (20, 9, 22, 10, DARK),     # scope
        ],
        disks=[(16, 5, 3, BRIGHT)],
        details=[
            ("rect", 14, 7, 18, 7, BLACK),    # hood shadow/eyes
            ("rect", 14, 13, 18, 14, MID),    # chest strap
            ("px", 29, 11, BLACK),            # muzzle tip
        ],
    )


def hoverbike() -> Image.Image:
    return compose(
        parts=[
            (14, 3, 19, 7, BRIGHT),    # rider helmet
            (13, 9, 20, 13, BRIGHT),   # rider torso
            (8, 15, 25, 17, DARK),     # bike chassis
            (20, 13, 24, 14, LIGHT),   # handlebar/windshield
        ],
        disks=[(16, 5, 3, BRIGHT), (16, 9, 1, BLACK)],
        details=[
            ("rect", 14, 6, 18, 6, BLACK),    # visor
            ("rect", 14, 10, 17, 11, MID),    # jacket detail
            ("rect", 10, 18, 23, 18, LIGHT),  # hover glow
            ("line", 7, 18, 7, 20, DARK),     # left fork
            ("line", 26, 18, 26, 20, DARK),   # right fork
            ("disk", 7, 23, 4, DARK),         # left wheel ring
            ("disk", 7, 23, 3, LIGHT),
            ("disk", 7, 23, 1, DARK),
            ("disk", 26, 23, 4, DARK),        # right wheel ring
            ("disk", 26, 23, 3, LIGHT),
            ("disk", 26, 23, 1, DARK),
        ],
    )


def mech() -> Image.Image:
    return compose(
        parts=[
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
        details=[
            ("rect", 13, 5, 19, 5, BLACK),    # eye slit
            ("rect", 15, 2, 17, 2, DARK),     # antenna base
            ("line", 16, 0, 16, 1, DARK),     # antenna
            ("px", 16, 0, BLACK),
            ("rect", 12, 14, 20, 16, MID),    # chest core
            ("rect", 14, 12, 18, 13, DARK),   # core top
            ("rect", 15, 15, 17, 15, BLACK),  # core dot
        ],
    )


def outpost() -> Image.Image:
    return compose(
        parts=[
            (6, 17, 26, 27, BRIGHT),   # bunker body
            (4, 14, 28, 16, DARK),     # roof slab
            (8, 21, 24, 24, MID),      # door
            (24, 10, 26, 13, DARK),    # antenna mast
        ],
        details=[
            ("rect", 7, 18, 25, 20, LIGHT),   # body highlight
            ("rect", 10, 22, 14, 23, DARK),   # door left
            ("rect", 18, 22, 22, 23, DARK),   # door right
            ("line", 25, 6, 25, 9, DARK),     # antenna up
            ("px", 25, 5, BLACK),             # beacon tip
        ],
    )


def command_center() -> Image.Image:
    return compose(
        parts=[
            (4, 14, 28, 27, BRIGHT),   # main building
            (3, 11, 29, 13, DARK),     # roof
            (7, 17, 25, 23, MID),      # door
            (12, 8, 20, 9, DARK),      # antenna bar
        ],
        disks=[(16, 7, 3, DARK)],
        details=[
            ("rect", 5, 15, 27, 16, LIGHT),   # top highlight
            ("rect", 11, 18, 21, 22, DARK),   # door inner
            ("rect", 13, 19, 19, 21, MID),    # door panel
            ("line", 16, 4, 16, 6, DARK),     # mast
            ("px", 16, 3, BLACK),             # beacon
            ("disk", 16, 7, 2, LIGHT),        # dish face
            ("px", 16, 7, BLACK),             # dish center
        ],
    )


def hq() -> Image.Image:
    return compose(
        parts=[
            (5, 12, 27, 28, BRIGHT),   # main base
            (2, 8, 10, 11, DARK),      # left tower
            (22, 8, 30, 11, DARK),     # right tower
            (9, 20, 23, 26, MID),      # gate
        ],
        disks=[(6, 6, 2, DARK), (26, 6, 2, DARK)],
        details=[
            ("rect", 6, 13, 26, 15, LIGHT),   # top highlight
            ("rect", 11, 21, 21, 25, DARK),   # gate inner
            ("rect", 13, 22, 19, 24, MID),    # gate panel
            ("rect", 4, 9, 8, 10, MID),       # left tower window
            ("rect", 24, 9, 28, 10, MID),     # right tower window
            ("line", 6, 3, 6, 5, DARK),       # left antenna
            ("line", 26, 3, 26, 5, DARK),     # right antenna
            ("px", 6, 2, BLACK),
            ("px", 26, 2, BLACK),
            ("rect", 14, 17, 18, 19, MID),    # crest
            ("px", 16, 18, BLACK),
        ],
    )


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
