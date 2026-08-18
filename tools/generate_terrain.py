#!/usr/bin/env python3
"""Generate the Hex & Dice terrain tiles.

Run from the repository root:

    python3 tools/generate_terrain.py

Requires Pillow.

All five terrain tiles share a single canonical pointy-top hexagon mask sized
to the game's grid (hex size 32 -> 64x64 tile), so tiles tessellate cleanly.
The artwork is drawn at 32x32 and upscaled 2x with nearest-neighbour for a
crisp pixel-art look with a vibrant palette.
"""

from __future__ import annotations

import math
import os
import random
from PIL import Image, ImageDraw

ROOT = os.path.dirname(os.path.dirname(os.path.abspath(__file__)))
SPRITE_DIR = os.path.join(ROOT, "client", "assets", "images", "sprites")

SIZE = 32          # drawing resolution (upscaled 2x on export)
R = SIZE / 2       # hex radius in drawing resolution
EDGE_DARKEN = 0.80 # shading factor on the outermost hex ring

# Vibrant palette
PLAINS   = (124, 194, 74)
PLAINS_T = (158, 222, 106)
PLAINS_D = (94, 162, 51)
FOREST   = (78, 155, 58)
FOREST_D = (58, 130, 47)
TREE     = (37, 107, 50)
TREE_HI  = (62, 143, 66)
WATER    = (58, 150, 222)
WAVE_L   = (176, 224, 250)
WAVE_D   = (36, 110, 178)
HILL     = (222, 178, 94)
HILL_D   = (192, 143, 62)
HILL_L   = (240, 205, 132)
MOUNT    = (142, 142, 156)
MOUNT_SH = (111, 111, 126)
SNOW     = (245, 245, 250)


def hex_mask() -> list[list[bool]]:
    img = Image.new("L", (SIZE, SIZE), 0)
    d = ImageDraw.Draw(img)
    cx = cy = R
    pts = [
        (cx + R * math.cos(math.radians(-90 + 60 * i)),
         cy + R * math.sin(math.radians(-90 + 60 * i)))
        for i in range(6)
    ]
    d.polygon(pts, fill=255)
    px = img.load()
    return [[px[x, y] > 0 for x in range(SIZE)] for y in range(SIZE)]


MASK = hex_mask()


def edge_rings() -> list[list[int]]:
    """0 = outside, 1 = outermost ring, 2 = second ring, 3 = interior."""
    rings = [[0] * SIZE for _ in range(SIZE)]
    for y in range(SIZE):
        for x in range(SIZE):
            if not MASK[y][x]:
                continue
            ring = 3
            for rad in (1, 2):
                hit = False
                for dy in range(-rad, rad + 1):
                    for dx in range(-rad, rad + 1):
                        ny, nx = y + dy, x + dx
                        if 0 <= ny < SIZE and 0 <= nx < SIZE and not MASK[ny][nx]:
                            hit = True
                if hit:
                    ring = rad
                    break
            rings[y][x] = ring
    return rings


RINGS = edge_rings()


def base_tile(color: tuple[int, int, int]) -> Image.Image:
    img = Image.new("RGBA", (SIZE, SIZE), (0, 0, 0, 0))
    for y in range(SIZE):
        for x in range(SIZE):
            if MASK[y][x]:
                img.putpixel((x, y), (*color, 255))
    return img


def scatter(img: Image.Image, kind: str, color: tuple[int, int, int],
            count: int, smin: int, smax: int, rng: random.Random,
            min_gap: int) -> None:
    d = ImageDraw.Draw(img)
    placed: list[tuple[int, int]] = []
    attempts = 0
    while len(placed) < count and attempts < 800:
        attempts += 1
        x = rng.randint(4, SIZE - 5)
        y = rng.randint(4, SIZE - 5)
        if not MASK[y][x] or RINGS[y][x] < 2:
            continue
        s = rng.randint(smin, smax)
        if any(abs(x - px) < min_gap + s and abs(y - py) < min_gap + s
               for px, py in placed):
            continue
        if kind == "px":
            img.putpixel((x, y), (*color, 255))
        elif kind == "blob":
            d.ellipse([x - s, y - s, x + s, y + s], fill=(*color, 255))
        elif kind == "dash":
            d.line([x - s, y, x + s, y], fill=(*color, 255), width=1)
        elif kind == "arc":
            d.arc([x - s, y - s, x + s, y + s], 200, 340, fill=(*color, 255), width=1)
        placed.append((x, y))


def apply_edge_shading(img: Image.Image) -> Image.Image:
    out = img.copy()
    for y in range(SIZE):
        for x in range(SIZE):
            if RINGS[y][x] == 1:
                r, g, b, a = out.getpixel((x, y))
                out.putpixel(
                    (x, y),
                    (round(r * EDGE_DARKEN), round(g * EDGE_DARKEN),
                     round(b * EDGE_DARKEN), a),
                )
    return out


def upscale(img: Image.Image) -> Image.Image:
    return img.resize((SIZE * 2, SIZE * 2), Image.NEAREST)


def plains() -> Image.Image:
    img = base_tile(PLAINS)
    rng = random.Random(11)
    scatter(img, "dash", PLAINS_T, 10, 1, 1, rng, 2)
    scatter(img, "dash", PLAINS_D, 7, 1, 1, rng, 2)
    return apply_edge_shading(img)


def forest() -> Image.Image:
    img = base_tile(FOREST)
    rng = random.Random(22)
    scatter(img, "blob", TREE, 5, 1, 2, rng, 4)
    scatter(img, "px", TREE_HI, 5, 0, 0, rng, 4)
    scatter(img, "dash", FOREST_D, 4, 1, 1, rng, 3)
    return apply_edge_shading(img)


def water() -> Image.Image:
    img = base_tile(WATER)
    rng = random.Random(33)
    scatter(img, "dash", WAVE_L, 6, 1, 2, rng, 3)
    scatter(img, "dash", WAVE_D, 5, 1, 2, rng, 3)
    return apply_edge_shading(img)


def hills() -> Image.Image:
    img = base_tile(HILL)
    rng = random.Random(44)
    scatter(img, "arc", HILL_D, 3, 2, 3, rng, 4)
    scatter(img, "dash", HILL_L, 4, 1, 2, rng, 4)
    scatter(img, "px", HILL_D, 4, 0, 0, rng, 3)
    return apply_edge_shading(img)


def mountain() -> Image.Image:
    img = base_tile(MOUNT)
    d = ImageDraw.Draw(img)
    # Two peaks: lit left face, shaded right face, snow cap.
    for cx, w, h, top in ((11, 7, 9, 9), (21, 8, 10, 10)):
        base_y = top + h
        d.polygon([(cx, top), (cx - w, base_y), (cx, base_y)], fill=(*SNOW, 255))
        d.polygon([(cx, top), (cx, base_y), (cx + w, base_y)], fill=(*MOUNT_SH, 255))
        d.polygon([(cx, top), (cx - 2, top + 3), (cx + 2, top + 3)], fill=(*SNOW, 255))
    # Clip anything that fell outside the mask.
    px = img.load()
    for y in range(SIZE):
        for x in range(SIZE):
            if not MASK[y][x]:
                px[x, y] = (0, 0, 0, 0)
    return apply_edge_shading(img)


def main() -> None:
    os.makedirs(SPRITE_DIR, exist_ok=True)
    tiles = {
        "plains.png": plains,
        "forest.png": forest,
        "water.png": water,
        "mountain.png": mountain,
        "hills.png": hills,
    }
    for filename, factory in tiles.items():
        upscale(factory()).save(os.path.join(SPRITE_DIR, filename))
        print(f"generated {filename}")


if __name__ == "__main__":
    main()
