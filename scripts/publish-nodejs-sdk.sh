#!/bin/bash

# Script to publish Node.js SDK to npm
# Usage: ./scripts/publish-nodejs-sdk.sh

set -e

echo "Publishing Node.js SDK to npm..."

cd /Users/swayamsingal/Desktop/Programming/Sentinel/sdk/nodejs

# Check if package.json exists
if [ ! -f package.json ]; then
    echo "Error: package.json not found in sdk/nodejs directory"
    exit 1
fi

# Check if we're logged in to npm
if ! npm whoami > /dev/null 2>&1; then
    echo "Error: Not logged in to npm. Please run 'npm login' first."
    exit 1
fi

# Publish to npm
echo "Publishing to npm..."
npm publish

echo "Node.js SDK published successfully!"
echo "Package is now available at: https://www.npmjs.com/package/@sentinel-platform/sentinel-sdk"