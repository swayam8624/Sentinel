#!/bin/bash

# Comprehensive release script for Sentinel
# Usage: ./scripts/release.sh <version>

set -e

if [ $# -ne 1 ]; then
    echo "Usage: $0 <version>"
    exit 1
fi

VERSION=$1

echo "Starting release process for version $VERSION"

# 1. Build and test
echo "1. Building and testing..."
make test

# 2. Build Docker image
echo "2. Building Docker image..."
docker build -t sentinel/gateway:$VERSION -t sentinel/gateway:latest .

# 3. Package Helm chart
echo "3. Packaging Helm chart..."
./scripts/publish-helm.sh

# 4. Create GitHub release
echo "4. Creating GitHub release..."
# Note: This requires GITHUB_TOKEN environment variable
./scripts/create-release.sh "v$VERSION" "Release $VERSION"

# 5. Publish Docker image (requires DOCKER_USERNAME and DOCKER_PASSWORD)
echo "5. Publishing Docker image..."
./scripts/publish-docker.sh "$VERSION"

# 6. Publish Helm chart
echo "6. Publishing Helm chart..."
# This would typically involve pushing to a gh-pages branch
# For now, we've prepared the files in the docs/charts directory

echo "Release process completed!"
echo "Manual steps required:"
echo "1. Push the docs/charts directory to GitHub Pages"
echo "2. Verify the release on GitHub"