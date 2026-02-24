#!/bin/bash

# Exit on error
set -e

echo "--- Building images ---"
docker compose build

echo "--- Deploying stack: hexdice ---"
# --resolve-image always ensures the local image digest is updated in the service definition
docker stack deploy --resolve-image always -c docker-compose.yml hexdice

echo "--- Deployment status ---"
docker stack services hexdice
