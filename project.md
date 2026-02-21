# Hex & Dice — Complete Project Specification

## 1. Overview

**Name:** Hex & Dice
**Genre:** Turn-based tactical strategy
**Theme:** Sci-fi military
**Players:** 2 (online multiplayer)
**Platforms:** Web, Android, iOS (single codebase)
**Art Style:** 32px low-res pixel art, bright & vibrant palette, Red vs Blue team colors

---

## 2. Tech Stack

| Layer | Technology |
|---|---|
| **Backend** | Go (goroutines, strict typing) |
| **Frontend** | Flutter + Flame engine (Dart) |
| **Protocol** | WebSocket (gameplay) + REST (lobby/matchmaking) |
| **Database** | Redis (game state) + PostgreSQL (future persistence) |
| **Hosting** | Single VPS (backend + Flutter web static files) |
| **Repo** | Monorepo: `/server` (Go), `/client` (Flutter) |
| **Mobile dist.** | Test builds only (APK sideload, TestFlight) |
| **Testing** | Server-side unit tests for core game logic |

---

## 3. Authentication & Matchmaking

- **Auth:** Guest only (pick a nickname, play immediately)
- **Persistence:** None for MVP (no accounts, no stats)
- **Matchmaking:** Both room codes (invite a friend) and queue matchmaking (random opponent)
- **Room creation options:** Map size, turn timer duration, turn mode (alternating vs simultaneous)

---

## 4. Map

- **Hex orientation:** Pointy-top
- **Map boundary:** Hexagonal
- **Map sizes:** Small (~7 hex radius, ~127 hexes), Medium (~10 hex radius, ~271 hexes), Large (~13 hex radius, ~469 hexes) — player chooses when creating room
- **Generation:** Procedural with rotational symmetry (180°)
- **Terrain distribution:** Clustered biomes (natural-looking zones of similar terrain)
- **Fog of war:** None for MVP (full visibility). Planned for later with per-unit vision range and terrain memory

### Terrain Types (5)

| Terrain | Movement Cost | Combat Modifier |
|---|---|---|
| Plains | 1 | None |
| Forest | 2 | +2 DEF |
| Hills | 2 | +1 ATK, +1 DEF |
| Water | Impassable | — |
| Mountains | Impassable | — |

- No elevation or line-of-sight mechanics for MVP

---

## 5. Structures

Three structure categories exist. Structures use dice-based combat (same formula as troops). Neutral structures attack all players in range but with weaker stats than player-owned ones.

### Structure Stats

| Structure | HP | ATK | DEF | Range | Damage | Income | Spawn | Notes |
|---|---|---|---|---|---|---|---|---|
| **Outpost** | 8 | +2 | 12 | 2 | 1D4 | +50/turn | Yes | Weak, numerous |
| **Command Center** | 15 | +4 | 15 | 3 | 1D6+2 | +50/turn | Yes | Strong, rare |
| **HQ** (per player) | 20 | +3 | 16 | 2 | 1D6 | — | Yes | Starting base, win condition target |

- **Neutral structure stats:** Same types but ATK modifier reduced by 2, damage dice one step lower
- **Capture:** Attack structure to 0 HP → instant capture at 1 HP by the attacking player
- **Repair:** Passive regen (+2 HP/turn for owned structures, up to max)
- **Placement:** Fixed count per map size (Small=3, Medium=5, Large=7 neutral structures + 2 HQs). Even spread with contested structures in the center

---

## 6. Troops

Four unit types. Sci-fi themed. One troop per hex. Enemies block movement, friendlies can be passed through.

### Troop Stats

| Unit | Cost | HP | ATK | DEF | Mobility | Range | Damage | Role |
|---|---|---|---|---|---|---|---|---|
| **Marine** | 100 | 10 | +3 | 14 | 3 | 1 | 1D6+1 | Frontline tank |
| **Sniper** | 150 | 6 | +4 | 11 | 2 | 3 | 1D8 | Ranged glass cannon |
| **Hoverbike** | 200 | 8 | +4 | 12 | 5 | 1 | 1D8+1 | Fast flanker |
| **Mech** | 350 | 12 | +5 | 10 | 1 | 3 | 2D6+2 | Slow, anti-structure (2x damage vs structures) |

### Troop Behavior

- **Purchase:** Buy anytime during your turn at HQ or owned structures. Cannot act until next turn
- **Actions per turn:** Move THEN attack (can do one, both, or neither)
- **Facing:** None (no flanking)
- **Stacking:** One troop per hex
- **Healing:** +2 HP/turn if the troop did not attack and was not attacked this turn
- **Death:** Permanent. No revival

---

## 7. Combat System

### Hit Resolution

- **Formula:** Attacker rolls D20 + ATK modifier vs defender's static DEF score
- **Hit:** If roll >= DEF, the attack hits
- **Critical hit:** Natural 20 = automatic hit + double damage dice
- **Fumble:** Natural 1 = automatic miss + defender gets a free counterattack opportunity

### Damage Resolution

- On hit, roll the unit's damage dice (variable per unit type)
- Subtract result from target's HP
- If HP <= 0, unit is destroyed

### Counterattack (Melee Only)

- When a melee unit (range=1) attacks another melee unit, the defender counter-rolls
- **Counter formula:** Defender rolls D20 + ATK vs attacker's static DEF
- **Counter damage:** Half damage dice (e.g., if unit normally deals 1D6+1, counter deals 1D3 or 1D4, rounded down)
- Ranged attacks never trigger counterattacks
- Fumble counterattacks use the same half-damage formula

### Structure Combat

- Structures attack automatically using the same D20 + ATK vs DEF formula
- Structures attack during their controlling player's turn (neutral structures attack at the start of each round)
- Player-owned structures attack only enemy troops in range
- Neutral structures attack all troops in range (with reduced stats)

### Kill Rewards

- None. Killing a unit has no coin/XP reward beyond the tactical advantage

### Dice Rolls

- All dice rolls generated server-side (anti-cheat)
- Animated dice roll shown to both players on the client

---

## 8. Economy

| Parameter | Value |
|---|---|
| Starting coins | 1000 |
| Passive income | 100/turn |
| Structure income | +50/turn per owned Outpost or Command Center |
| Troop upkeep | None |
| Troop pricing | Tiered: Marine=100, Sniper=150, Hoverbike=200, Mech=350 |

---

## 9. Turn Modes

The room creator selects one of two turn modes.

### Mode A: Full Turn Alternation

- Player A performs all actions (buy, move, attack for all units), then ends turn
- Player B performs all actions, then ends turn
- **First turn mitigation:** Player 1 (randomly chosen) cannot attack on their first turn (move only)

### Mode B: Simultaneous Turns

- **Planning phase:** Both players see last-known enemy positions and queue orders for all their units (move targets, attack targets)
- **Resolution — Movement phase:** All moves resolve simultaneously. Conflicts resolved by mobility stat (faster unit arrives first, slower unit's move is cancelled). Animated for both players
- **Resolution — Attack phase:** All attacks resolve simultaneously. Attacks target a hex; if the target moved to an adjacent hex, attack still hits with a penalty (reduced ATK modifier by 2). Animated for both players
- **No first-turn advantage** in simultaneous mode (both players act at once)

### Turn Timer

- Configurable per game by room creator (60s / 90s / 120s)
- If timer expires, turn auto-ends (alternating mode) or orders auto-submit as-is (simultaneous mode)

---

## 10. Win Conditions (Multiple Paths)

A player wins by achieving **any one** of:

1. **Destroy enemy HQ** — Reduce opponent's HQ to 0 HP
2. **Structure dominance** — Control a majority of all structures (including HQs) for 3 consecutive turns

If neither condition is met by the sudden death resolution, the player controlling more structures wins. If tied, the player with more total troop HP wins. If still tied, draw.

---

## 11. Sudden Death

- **Trigger:** Scales with map size (Small=20 turns, Medium=30 turns, Large=40 turns)
- **Mechanic:** Shrinking safe zone from map edges toward center
- **HQ relocation:** Both player HQs are moved inside the safe zone boundary each time it shrinks
- **Damage:** Escalating — troops in the storm take 1 damage on the first sudden-death turn, 2 on the second, 3 on the third, etc.
- **Shrink rate:** Safe zone radius decreases by 1 hex per turn after sudden death activates

---

## 12. Networking

| Concern | Solution |
|---|---|
| Game state authority | Server is authoritative; client validates locally for responsiveness |
| Logic sharing | Dual implementation: Go (server) + Dart (client). Kept in sync manually |
| Disconnect | 60-second reconnect window, then forfeit |
| Surrender | Available after turn 5 |
| Reconnect | Client receives full game state snapshot on reconnect |
| Protocol | WebSocket for gameplay, REST for lobby/rooms/matchmaking |

---

## 13. UI & Screens

### Screen Flow

```
Title Screen → Play → Create Room / Join Room (code) / Quick Match → Game → End Screen → Title
                → How to Play (rules text + diagrams)
```

### In-Game HUD (Minimal)

- **Top bar:** Turn number, current player indicator, timer countdown
- **Bottom bar:** Coin count, End Turn button
- **Troop info:** Tap a troop for a popup showing all stats
- **Move/attack preview:** Tap own troop → reachable hexes highlighted blue, attackable hexes highlighted red
- **Buy troops:** Tap owned spawn structure → shop panel slides up with troop cards (cost, stats). Tap to purchase, troop appears on hex
- **Emotes:** Quick-access emote button with predefined phrases ("Good move!", "GG", "Oops", etc.)

### End-of-Game Screen

- Winner announcement
- Basic stats: troops killed, structures held, turns played
- Buttons: Play Again / Return to Menu

---

## 14. Visual & Audio

### Art

- **Resolution:** 32x32px sprites, pointy-top hex tiles (~36x32)
- **Style:** Low-res pixel art, sci-fi theme
- **Palette:** Bright & vibrant with strong Red vs Blue team colors
- **Source:** AI-generated pixel art
- **Animations:** Basic (2-3 frame idle, walk, attack per unit)
- **Camera:** Top-down, pannable/zoomable
- **Dice animation:** Animated D20 / damage dice rolls on combat

### Audio

- **Music:** Chiptune BGM (menu + gameplay tracks). Custom/AI-generated
- **SFX:** Core set only — attack hit, miss, troop move, death, dice roll, turn notification, structure capture, coin gain. Custom/AI-generated

---

## 15. Scope Phasing

### MVP (Phase 1)

- Guest-only auth (nickname)
- Room codes + queue matchmaking
- Full turn alternation mode only
- Procedural map generation (rotationally symmetric)
- 4 troop types, 3 structure types
- Full combat system (D20, damage dice, crits, fumbles, counterattacks)
- Full terrain system (5 types, movement + combat modifiers)
- Economy (income, structure bonuses, tiered troop pricing)
- Multiple win conditions (HQ destruction + structure dominance)
- Sudden death (shrinking zone)
- Core SFX, chiptune BGM
- AI-generated placeholder sprites (32px)
- Minimal HUD, tap interactions
- Emote communication
- Rules/help screen
- Server unit tests
- Single VPS deployment

### Phase 2 (Post-MVP)

- Simultaneous turn mode
- Fog of war (per-unit vision, terrain memory)
- Guest + optional accounts (OAuth)
- Basic stats persistence (wins/losses)
- Polished AI-generated or commissioned sprites
- Interactive tutorial

### Phase 3 (Future)

- ELO/ranking, leaderboards
- More troop types
- More structure types
- Terrain: roads, swamp
- Elevation + line of sight
- Replay system
- App store distribution
