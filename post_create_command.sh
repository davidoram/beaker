#!/bin/bash
#
# post_create_command.sh hooks into the devcontainer 'postCreateCommand'
#
# - Runs inside the container once, after both onCreateCommand and updateContentCommand, and only when the container is first assigned to a user.
# - It’s the last of the three creation hooks and can safely reference user‐specific secrets, credentials, or environment variables (e.g. installing packages from a private registry, initializing dotfiles). 
#

# Lets ensure that we have the latest code from the remote repository
git pull origin
echo "Code updated from remote repository."


# Wait for Docker socket to appear
while [ ! -S /var/run/docker.sock ]; do
    echo "Waiting for Docker socket..."
    sleep 1
done

# Check if Docker responds
until docker version >/dev/null 2>&1; do
    echo "Waiting for Docker daemon to be ready..."
    sleep 1
done
echo "Docker daemon is ready!"

make restart-docker-compose
echo "post_create_command.sh completed successfully."
