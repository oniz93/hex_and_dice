# Hex & Dice: Production Deployment Path Summary

This document tracks the surgical changes and architectural decisions made to transition the "Hex & Dice" project from local development to a production-ready Docker Swarm environment on an `aarch64` VPS.

## 1. Infrastructure & Orchestration
- **Docker Swarm Migration:** Upgraded `docker-compose.yml` to version `3.8`. Implemented `deploy` configurations (replicas, restart policies) and an `overlay` network for secure internal communication.
- **Port Management:**
    - Published the **Client (Nginx)** on host port `8082`.
    - Removed host port publishing for the **Server (Go)** to resolve conflicts with `github-trending` (port 8080), routing all API traffic internally via the overlay network.
- **System-Level Reverse Proxy:** Corrected the system-level Nginx (`/etc/nginx/sites-enabled/api.hexdice.teomiscia.com`) to proxy traffic to port `8082` instead of `8080`.

## 2. Multi-Architecture Support
- **ARM64 Compatibility:** Updated the `client/Dockerfile` to use `ghcr.io/cirruslabs/flutter:stable` as the builder image, resolving `exec format error` issues on the `aarch64` VPS.

## 3. Network & Security
- **Domain Routing:** Configured separate `server_name` blocks in `nginx.conf` for `hexdice.teomiscia.com` (Frontend) and `api.hexdice.teomiscia.com` (API/WebSocket).
- **HTTPS Enforcement:** Hardcoded `https://` and `wss://` for production URLs in `client/lib/providers/core_providers.dart` to eliminate "Mixed Content" browser errors.
- **CORS Resolution:** 
    - Simplified `CORS_ORIGINS` to `*` in the backend environment.
    - Refactored Go `CORSMiddleware` for robust `*` and `Vary: Origin` handling.
    - Removed redundant `add_header` calls in Docker Nginx to fix the "multiple values '*, *'" CORS policy error.

## 4. Deployment Automation & Reliability
- **`deploy.sh` Enhancements:**
    - Added `--resolve-image always` and `docker service update --force` to guarantee that local image builds trigger rolling updates in the Swarm.
    - Implemented a `BUILD_VERSION` argument in the `Dockerfile` to bust the Docker build cache and ensure fresh Flutter compilation.
- **Cache Invalidation:** Configured `nginx.conf` with `Cache-Control: no-store` for critical entry points (`index.html`, `main.dart.js`, service workers) to prevent stale code issues.

## 5. Verification & Versioning
- **Incremental Versioning:** Iterated from `v1.1` to `v1.5` in the application title and console logs to empirically verify that new builds were correctly successfully deployed and served by the browser.
