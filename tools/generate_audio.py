#!/usr/bin/env python3
"""Procedurally generate the Hex & Dice sound effects and music.

Run from the repository root:

    python3 tools/generate_audio.py

Requires numpy. Music loops (.mp3) additionally require ffmpeg on PATH for
the WAV -> MP3 conversion; if ffmpeg is missing, WAVs are kept next to the
output directory.

Output layout (loaded by FlameAudio with the ``assets/audio/`` prefix):

    client/assets/audio/sfx/*.wav
    client/assets/audio/music/*.mp3
"""

from __future__ import annotations

import os
import shutil
import struct
import subprocess
import sys
import wave

import numpy as np

ROOT = os.path.dirname(os.path.dirname(os.path.abspath(__file__)))
AUDIO_DIR = os.path.join(ROOT, "client", "assets", "audio")
SFX_DIR = os.path.join(AUDIO_DIR, "sfx")
MUSIC_DIR = os.path.join(AUDIO_DIR, "music")

SFX_RATE = 22050
MUSIC_RATE = 44100


# ---------------------------------------------------------------------------
# WAV helpers
# ---------------------------------------------------------------------------

def write_wav(path: str, samples: np.ndarray, rate: int) -> None:
    """Write mono 16-bit PCM WAV, soft-limiting peaks to avoid clipping."""
    peak = np.abs(samples).max()
    if peak > 0.95:
        samples = samples / peak * 0.95
    samples = np.clip(samples, -1.0, 1.0)
    pcm = (samples * 32767).astype(np.int16)
    with wave.open(path, "wb") as w:
        w.setnchannels(1)
        w.setsampwidth(2)
        w.setframerate(rate)
        w.writeframes(pcm.tobytes())


def t(dur: float, rate: int) -> np.ndarray:
    return np.arange(int(dur * rate)) / rate


def sine(freq: float, dur: float, rate: int, phase: float = 0.0) -> np.ndarray:
    return np.sin(2 * np.pi * freq * t(dur, rate) + phase)


def square(freq: float, dur: float, rate: int, duty: float = 0.5) -> np.ndarray:
    ph = (freq * t(dur, rate)) % 1.0
    return np.where(ph < duty, 1.0, -1.0)


def saw(freq: float, dur: float, rate: int) -> np.ndarray:
    return 2.0 * ((freq * t(dur, rate)) % 1.0) - 1.0


def noise(dur: float, rate: int, seed: int = 0) -> np.ndarray:
    rng = np.random.default_rng(seed)
    return rng.uniform(-1.0, 1.0, int(dur * rate))


def env_adsr(length: int, a: float, d: float, s: float, r: float,
             sustain_level: float = 0.7) -> np.ndarray:
    """Attack/decay/sustain/release envelope over ``length`` samples."""
    e = np.ones(length) * sustain_level
    na, nd, nr = int(a * length), int(d * length), int(r * length)
    if na > 0:
        e[:na] = np.linspace(0, 1, na)
    if nd > 0:
        e[na:na + nd] = np.linspace(1, sustain_level, nd)
    if nr > 0:
        e[-nr:] = np.linspace(e[-nr], 0, nr)
    return e


def lowpass(x: np.ndarray, window: int) -> np.ndarray:
    if window <= 1:
        return x
    kernel = np.ones(window) / window
    return np.convolve(x, kernel, mode="same")


def sweep_noise(dur: float, rate: int, start_win: int, end_win: int,
                seed: int = 1) -> np.ndarray:
    """White noise with a time-varying moving-average lowpass (whoosh)."""
    x = noise(dur, rate, seed)
    n = len(x)
    out = np.empty(n)
    chunk = 64
    wins = np.linspace(start_win, end_win, n // chunk + 1)
    for i in range(0, n, chunk):
        j = min(n, i + chunk)
        w = max(1, int(wins[i // chunk]))
        k = np.ones(w) / w
        seg = np.convolve(x[max(0, i - w):min(n, j + w)], k, mode="same")
        off = min(i, w)
        out[i:j] = seg[off:off + (j - i)]
    return out


def place(track: np.ndarray, offset: int, sig: np.ndarray,
          gain: float = 1.0) -> None:
    end = min(len(track), offset + len(sig))
    track[offset:end] += sig[:end - offset] * gain


# ---------------------------------------------------------------------------
# Sound effects
# ---------------------------------------------------------------------------

def sfx_attack_hit() -> np.ndarray:
    r = SFX_RATE
    dur = 0.25
    hit = noise(dur, r, 11) * env_adsr(int(dur * r), 0.01, 0.2, 0.6, 0.5, 0.3)
    thump = sine(110, dur, r) * env_adsr(int(dur * r), 0.01, 0.3, 0.5, 0.4, 0.2)
    return lowpass(hit * 0.6, 3) + thump * 0.8


def sfx_attack_miss() -> np.ndarray:
    r = SFX_RATE
    dur = 0.3
    whoosh = sweep_noise(dur, r, 2, 24)
    return whoosh * env_adsr(int(dur * r), 0.2, 0.2, 0.6, 0.4, 0.5) * 0.7


def sfx_troop_move() -> np.ndarray:
    r = SFX_RATE
    dur = 0.12
    out = np.zeros(int(dur * r))
    tick = lowpass(noise(0.03, r, 3), 2) * env_adsr(int(0.03 * r), 0.05, 0.2, 0.4, 0.6, 0.3)
    place(out, 0, tick, 0.8)
    place(out, int(0.05 * r), tick, 0.5)
    return out


def sfx_troop_death() -> np.ndarray:
    r = SFX_RATE
    dur = 0.5
    n = int(dur * r)
    # descending tone
    tt = t(dur, r)
    freq = 420 * np.exp(-tt * 6)
    tone = np.sin(2 * np.pi * np.cumsum(freq) / r)
    burst = noise(dur, r, 21) * env_adsr(n, 0.01, 0.3, 0.5, 0.4, 0.2)
    return lowpass(burst * 0.5, 4) + tone * env_adsr(n, 0.01, 0.2, 0.5, 0.5, 0.3) * 0.6


def sfx_dice_roll() -> np.ndarray:
    r = SFX_RATE
    dur = 0.6
    out = np.zeros(int(dur * r))
    tick = lowpass(noise(0.04, r, 5), 3) * env_adsr(int(0.04 * r), 0.05, 0.3, 0.4, 0.5, 0.2)
    rng = np.random.default_rng(7)
    pos = 0
    while pos < int(0.52 * r):
        place(out, pos, tick, rng.uniform(0.4, 0.9))
        pos += int(rng.uniform(0.04, 0.09) * r)
    return out


def sfx_turn_start() -> np.ndarray:
    r = SFX_RATE
    out = np.zeros(int(0.35 * r))
    n1 = sine(620, 0.08, r) * env_adsr(int(0.08 * r), 0.05, 0.2, 0.5, 0.4, 0.4)
    n2 = sine(930, 0.12, r) * env_adsr(int(0.12 * r), 0.05, 0.2, 0.5, 0.4, 0.4)
    place(out, 0, n1, 0.7)
    place(out, int(0.1 * r), n2, 0.7)
    return out


def sfx_structure_capture() -> np.ndarray:
    r = SFX_RATE
    out = np.zeros(int(0.6 * r))
    notes = [523.25, 659.25, 783.99, 1046.5]  # C5 E5 G5 C6
    for i, f in enumerate(notes):
        tone = square(f, 0.12, r, 0.5) * env_adsr(int(0.12 * r), 0.05, 0.2, 0.5, 0.5, 0.4)
        place(out, int(i * 0.11 * r), tone, 0.4)
    return out


def sfx_coin_gain() -> np.ndarray:
    r = SFX_RATE
    out = np.zeros(int(0.3 * r))
    n1 = square(987.77, 0.06, r) * env_adsr(int(0.06 * r), 0.05, 0.2, 0.5, 0.5, 0.3)
    n2 = square(1318.5, 0.18, r) * env_adsr(int(0.18 * r), 0.05, 0.1, 0.5, 0.5, 0.4)
    place(out, 0, n1, 0.5)
    place(out, int(0.07 * r), n2, 0.5)
    return out


def sfx_purchase() -> np.ndarray:
    r = SFX_RATE
    tone = square(660, 0.16, r) * env_adsr(int(0.16 * r), 0.03, 0.2, 0.5, 0.5, 0.4)
    return tone * 0.5


def sfx_emote_pop() -> np.ndarray:
    r = SFX_RATE
    dur = 0.12
    tt = t(dur, r)
    freq = 800 * np.exp(-tt * 10) + 200
    pop = np.sin(2 * np.pi * np.cumsum(freq) / r)
    return pop * env_adsr(int(dur * r), 0.02, 0.2, 0.4, 0.5, 0.3) * 0.7


# ---------------------------------------------------------------------------
# Music (simple chiptune sequencer)
# ---------------------------------------------------------------------------

def midi_to_freq(m: int) -> float:
    return 440.0 * 2 ** ((m - 69) / 12.0)


def render_song(events, dur: float, rate: int = MUSIC_RATE) -> np.ndarray:
    """events: list of (wave, midi_or_None, start_sec, dur_sec, gain)."""
    track = np.zeros(int(dur * rate))
    for wave_fn, m, start, length, gain in events:
        if m is None:
            continue
        f = midi_to_freq(m)
        n = int(length * rate)
        sig = wave_fn(f, length, rate)
        e = env_adsr(n, 0.02, 0.1, 0.6, 0.2, 0.6)
        place(track, int(start * rate), sig * e, gain)
    # gentle master limiter
    peak = np.abs(track).max()
    if peak > 0.95:
        track = track / peak * 0.95
    return track


def menu_theme() -> np.ndarray:
    r = MUSIC_RATE
    bpm = 96
    beat = 60.0 / bpm
    bars = 8
    total = bars * 4 * beat

    # Bass line: A2 F2 G2 E2 pattern (midi 45, 41, 43, 40)
    bass_prog = [45, 41, 43, 40]
    events = []
    for bar in range(bars):
        root = bass_prog[bar % 4]
        for b in range(4):
            events.append((square, root, bar * 4 * beat + b * beat, beat * 0.9, 0.25))

    # Arpeggio: minor arps above root
    arp_offsets = [12, 15, 19, 24, 19, 15]
    for bar in range(bars):
        root = bass_prog[bar % 4]
        for i in range(16):  # 16th notes
            m = root + 12 + arp_offsets[i % len(arp_offsets)]
            events.append((square, m, bar * 4 * beat + i * beat / 4, beat / 4 * 0.8, 0.12))

    # Slow pad every 2 bars
    for bar in range(0, bars, 2):
        root = bass_prog[bar % 4]
        events.append((saw, root + 12, bar * 4 * beat, beat * 8, 0.06))
    return render_song(events, total, r)


def battle_theme() -> np.ndarray:
    r = MUSIC_RATE
    bpm = 140
    beat = 60.0 / bpm
    bars = 8
    total = bars * 4 * beat

    bass_prog = [40, 40, 43, 45]  # E2 E2 G2 A2
    events = []
    for bar in range(bars):
        root = bass_prog[bar % 4]
        for i in range(8):  # driving 8ths
            events.append((square, root, bar * 4 * beat + i * beat / 2, beat / 2 * 0.9, 0.3))

    # Lead stabs on offbeats
    lead_prog = [64, 64, 67, 69]
    for bar in range(bars):
        root = lead_prog[bar % 4]
        for i in range(4):
            events.append((square, root, bar * 4 * beat + i * beat + beat / 2, beat * 0.4, 0.14))

    # Render tones first
    track = render_song(events, total, r)
    # Then hi-hat-ish noise ticks every 8th

    for bar in range(bars):
        for i in range(8):
            start = bar * 4 * beat + i * beat / 2
            tick = noise(0.03, r, bar * 8 + i) * env_adsr(int(0.03 * r), 0.02, 0.2, 0.4, 0.5, 0.2)
            place(track, int(start * r), tick, 0.08)
    return track


def convert_to_mp3(wav_path: str, mp3_path: str) -> bool:
    ffmpeg = shutil.which("ffmpeg")
    if not ffmpeg:
        return False
    subprocess.run(
        [ffmpeg, "-y", "-loglevel", "error", "-i", wav_path,
         "-c:a", "libmp3lame", "-q:a", "4", mp3_path],
        check=True,
    )
    return True


def main() -> None:
    os.makedirs(SFX_DIR, exist_ok=True)
    os.makedirs(MUSIC_DIR, exist_ok=True)

    sfx = {
        "attack_hit.wav": sfx_attack_hit,
        "attack_miss.wav": sfx_attack_miss,
        "troop_move.wav": sfx_troop_move,
        "troop_death.wav": sfx_troop_death,
        "dice_roll.wav": sfx_dice_roll,
        "turn_start.wav": sfx_turn_start,
        "structure_capture.wav": sfx_structure_capture,
        "coin_gain.wav": sfx_coin_gain,
        "purchase.wav": sfx_purchase,
        "emote_pop.wav": sfx_emote_pop,
    }
    for filename, factory in sfx.items():
        write_wav(os.path.join(SFX_DIR, filename), factory(), SFX_RATE)
        print(f"generated sfx/{filename}")

    music = {
        "menu_theme": menu_theme,
        "battle_theme": battle_theme,
    }
    for name, factory in music.items():
        wav_path = os.path.join(MUSIC_DIR, f"{name}.wav")
        write_wav(wav_path, factory(), MUSIC_RATE)
        mp3_path = os.path.join(MUSIC_DIR, f"{name}.mp3")
        if convert_to_mp3(wav_path, mp3_path):
            os.remove(wav_path)
            print(f"generated music/{name}.mp3")
        else:
            print(f"ffmpeg not found; kept music/{name}.wav", file=sys.stderr)


if __name__ == "__main__":
    main()
