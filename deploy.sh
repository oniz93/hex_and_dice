#!/bin/bash

# Exit on error
set -e

echo "--- Building images ---"
docker compose build --build-arg BUILD_VERSION=$(date +%s)

echo "--- Deploying stack: hexdice ---"
docker stack deploy --resolve-image always -c docker-compose.yml hexdice

echo "--- Force updating services to ensure latest local images are used ---"
docker service update --force hexdice_client
docker service update --force hexdice_server

echo "--- Deployment status ---"
docker stack services hexdice
