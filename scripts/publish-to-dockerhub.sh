#!/bin/bash

# Script to publish Docker images to Docker Hub
# Usage: ./scripts/publish-to-dockerhub.sh <dockerhub_username>

set -e

if [ $# -ne 1 ]; then
    echo "Usage: $0 <dockerhub_username>"
    echo "Example: $0 yourusername"
    exit 1
fi

DOCKERHUB_USERNAME=$1
IMAGE_NAME="sentinel-gateway"
VERSION="0.1.0"

echo "Publishing Docker images to Docker Hub..."
echo "Username: $DOCKERHUB_USERNAME"
echo "Image: $IMAGE_NAME"
echo "Version: $VERSION"

# Tag images for Docker Hub
echo "Tagging images for Docker Hub..."
docker tag sentinel/gateway:$VERSION $DOCKERHUB_USERNAME/$IMAGE_NAME:$VERSION
docker tag sentinel/gateway:latest $DOCKERHUB_USERNAME/$IMAGE_NAME:latest

# Push images to Docker Hub
echo "Pushing images to Docker Hub..."
docker push $DOCKERHUB_USERNAME/$IMAGE_NAME:$VERSION
docker push $DOCKERHUB_USERNAME/$IMAGE_NAME:latest

echo "Docker images published successfully!"
echo "Images are now available at:"
echo "  - https://hub.docker.com/r/$DOCKERHUB_USERNAME/$IMAGE_NAME/tags"