#!/bin/bash
#
# post_create_command.sh hooks into the devcontainer 'postCreateCommand'
#
# - Runs inside the container once, after both onCreateCommand and updateContentCommand, and only when the container is first assigned to a user.
# - It's the last of the three creation hooks and can safely reference user‐specific secrets, credentials, or environment variables (e.g. installing packages from a private registry, initializing dotfiles). 
#

# Exit immediately if a command exits with a non-zero status
set -e

# Lets ensure that we have the latest code from the remote repository
git pull origin
echo "Code updated from remote repository."

# Decode Github codespace user secrets - see https://github.com/settings/codespaces

# Process NATS_CREDS_* environment variables
for var in $(env | grep '^NATS_CREDS_' | cut -d'=' -f1); do
    echo "Processing environment variable: $var"
    # Get the value of the environment variable
    value="${!var}"
    if [ -n "$value" ]; then
        # Create the credentials file path
        creds_file="$HOME/${var}.creds"
        echo "Decoding $var and saving to $creds_file"
        # Clean the base64 string (remove spaces and newlines) then decode
        printf '%s' "$value" | tr -d ' \n' | base64 -d > "$creds_file"
        # Set appropriate permissions
        chmod 600 "$creds_file"
        echo "Created credentials file: $creds_file"
        # Add the NATS context using the credentials file
        nats context add "$var" --server "tls://connect.ngs.global" --creds "$creds_file"
        echo "Added NATS context for $var"
    else
        echo "Warning: $var is empty"
    fi
done


echo "Waiting for the docker socket to be ready..."
until test -S /var/run/docker.sock; do
    echo "Waiting for Docker socket '/var/run/docker.sock' to exist..."
    sleep 2
done

echo "Waiting for docker daemon to be ready..."
until docker info >/dev/null 2>&1 && docker ps >/dev/null 2>&1; do
    echo "Waiting for Docker daemon to be fully ready..."
    sleep 2
done

sleep 2  # Additional wait to ensure everything is stable

echo "Docker daemon is ready!"

make restart-docker-compose
echo "post_create_command.sh completed successfully."
