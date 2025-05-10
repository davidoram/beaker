#!/bin/bash
#
# on_create_command.sh hooks into the devcontainer 'onCreateCommand'
#
# - Runs inside the container once, immediately after the container is built and started for the first time.
# - It’s the first of three “finalize‐creation” commands (onCreateCommand → updateContentCommand → postCreateCommand).
# - Intended for container‐level setup that doesn’t require user‐scoped assets or secrets (e.g. installing global tooling, compiling code) and that cloud prebuilds can cache.
#
make setup