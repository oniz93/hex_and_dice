# Hex & Dice

**Hex & Dice** is a turn-based tactical strategy game with a sci-fi military theme. It features a hexagonal grid, procedural map generation, and a dice-based combat system. Built with a Go backend and a Flutter frontend, it supports cross-platform play across Web, Android, and iOS.

## 🚀 Deployment

The project is designed to run on a Docker Swarm cluster, with a system-level Nginx acting as the primary reverse proxy.

### Prerequisites
- Docker & Docker Compose
- Docker Swarm initialized (`docker swarm init`)
- System-level Nginx (for SSL and external routing)

### Domain Configuration
- **Frontend:** `hexdice.teomiscia.com` (proxies to port `8555`)
- **Backend/API:** `api.hexdice.teomiscia.com` (proxies to port `8550`)

### Quick Start (Build & Deploy)
We provide a deployment script that handles the full build and update cycle:

```bash
# Make the script executable (first time only)
chmod +x deploy.sh

# Build and deploy/update the stack
./deploy.sh
```

### Management Commands

| Action | Command |
|---|---|
| **View Services** | `docker stack services hexdice` |
| **Check Health** | `docker stack ps hexdice` |
| **Server Logs** | `docker service logs -f hexdice_server` |
| **Client Logs** | `docker service logs -f hexdice_client` |
| **Force Restart** | `docker service update --force hexdice_server` |
| **Remove Stack** | `docker stack rm hexdice` |

---

## 🏗 Architecture

### Backend (`/server`)
- **Language:** Go
- **Authoritative Source of Truth:** Manages all game state and logic.
- **Communication:** WebSockets for real-time gameplay, REST for lobby/matchmaking.
- **Persistence:** Redis for game state snapshots and recovery.
- **Concurrency:** One goroutine per active game for high-performance state isolation.

### Frontend (`/client`)
- **Framework:** Flutter + Flame Engine
- **State Management:** Riverpod
- **Rendering:** 32px low-res pixel art on a pointy-top hex grid.
- **Responsiveness:** Single codebase for Web (CanvasKit), Android, and iOS.

---

## 🎮 Game Mechanics

- **Combat:** D20-based hit resolution + variable damage dice. Natural 20s are crits; Natural 1s are fumbles triggering counterattacks.
- **Economy:** Capture structures (Outposts, Command Centers) to increase your per-turn credit income.
- **Map:** Procedurally generated with 180° rotational symmetry to ensure fairness.
- **Win Conditions:**
  1. Destroy the enemy HQ.
  2. Maintain structure dominance for 3 consecutive turns.
- **Sudden Death:** After a set number of turns, the map begins shrinking towards the center, forcing a final confrontation.

---

## 📁 Project Structure

```
.
├── client/              # Flutter frontend application
├── server/              # Go backend application
├── docker-compose.yml   # Swarm stack definition
├── nginx.conf           # Internal Nginx routing (Docker)
├── deploy.sh            # Build & Deployment automation script
├── HLD_BE.md            # Backend High-Level Design
├── HLD_FE.md            # Frontend High-Level Design
└── project.md           # Full Project Specification
```

---

## 🛠 Development

### Running Locally
1. **Backend:**
   ```bash
   cd server
   go run cmd/server/main.go
   ```
2. **Frontend:**
   ```bash
   cd client
   flutter run -d chrome
   ```

### Client Base URL Logic
The client automatically detects if it is running on `localhost` and adjusts its API endpoints accordingly. For production, it points to `api.hexdice.teomiscia.com`.
