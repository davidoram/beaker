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
        # Base64 decode and save to file
        echo "$value" | base64 -d > "$creds_file"
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


# Try to actually use Docker with a simple operation
echo "Docker daemon responding, testing full functionality..."
until docker info >/dev/null 2>&1 && docker ps >/dev/null 2>&1; do
    echo "Waiting for Docker daemon to be fully ready..."
    sleep 2
done

# Final test - try a simple image operation to ensure daemon can handle image pulls
echo "Testing Docker image operations..."
max_retries=30
retry_count=0
while [ $retry_count -lt $max_retries ]; do
    if docker pull hello-world:latest >/dev/null 2>&1; then
        echo "Docker daemon is fully ready!"
        docker rmi hello-world:latest >/dev/null 2>&1  # Clean up test image
        break
    else
        echo "Docker daemon not ready for image operations, retrying... ($((retry_count + 1))/$max_retries)"
        sleep 2
        retry_count=$((retry_count + 1))
    fi
done

if [ $retry_count -eq $max_retries ]; then
    echo "Error: Docker daemon failed to become ready after $max_retries attempts"
    exit 1
fi

echo "Docker daemon is ready!"

make restart-docker-compose
echo "post_create_command.sh completed successfully."
