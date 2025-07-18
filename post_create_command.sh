#!/bin/bash
#
# post_create_command.sh hooks into the devcontainer 'postCreateCommand'
#
# - Runs inside the container once, after both onCreateCommand and updateContentCommand, and only when the container is first assigned to a user.
# - It’s the last of the three creation hooks and can safely reference user‐specific secrets, credentials, or environment variables (e.g. installing packages from a private registry, initializing dotfiles). 
#

# Lets ensure that we have the latest code from the remote repository
git pull origin


# Wait for the docker server to be ready
while ! docker info > /dev/null 2>&1; do
  echo "Waiting for Docker to be ready..."
  sleep 5
done
make restart-docker-compose
