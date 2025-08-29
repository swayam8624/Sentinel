#!/bin/bash

# Script to build and publish Docker image
# Usage: ./scripts/publish-docker.sh <version>

set -e

if [ $# -ne 1 ]; then
    echo "Usage: $0 <version>"
    exit 1
fi

VERSION=$1

# Build Docker image
docker build -t sentinel/gateway:$VERSION -t sentinel/gateway:latest .

# Login to Docker Hub (requires DOCKER_USERNAME and DOCKER_PASSWORD environment variables)
echo "$DOCKER_PASSWORD" | docker login -u "$DOCKER_USERNAME" --password-stdin

# Push Docker image
docker push sentinel/gateway:$VERSION
docker push sentinel/gateway:latest

echo "Docker image published successfully!"