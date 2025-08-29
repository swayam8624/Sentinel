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
./scripts/publish-helm-charts.sh

# 4. Create GitHub release
echo "4. Creating GitHub release..."
# Note: This requires GITHUB_TOKEN environment variable
./scripts/create-release.sh "v$VERSION" "Release $VERSION"

echo "Release process completed!"
echo "Manual steps required:"
echo "1. Publish Docker images: ./scripts/publish-to-dockerhub.sh <dockerhub_username>"
echo "2. Publish Node.js SDK: ./scripts/publish-nodejs-sdk.sh"
echo "3. Publish Python SDK: ./scripts/publish-python-sdk.sh"
echo "4. Push Helm chart updates to GitHub Pages"