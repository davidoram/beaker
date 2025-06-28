#!/bin/bash
# Script to display the otel collector logs running in the Docker container

# Find the app container ID/name
OTEL_COLLECTOR=$(docker ps | grep devcontainer-otel | awk '{print $1}')

if [ -z "$OTEL_COLLECTOR" ]; then
  echo "Error: Otel collector container not found. Make sure your Docker Compose is running."
  echo "Try running: make docker-compose-up"
  exit 1
fi

# Follow the logs of the otel collector container
docker logs --follow $OTEL_COLLECTOR
