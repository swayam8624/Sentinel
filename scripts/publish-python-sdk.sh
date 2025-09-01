#!/bin/bash

# Script to publish the Python SDK to PyPI
# Usage: ./scripts/publish-python-sdk.sh

set -e

echo "Publishing Sentinel Python SDK to PyPI..."

# Navigate to the Python SDK directory
cd "$(dirname "$0")/../sdk/python"

# Check if we're on the main branch
BRANCH=$(git rev-parse --abbrev-ref HEAD)
if [ "$BRANCH" != "main" ]; then
    echo "Warning: You are not on the main branch. Current branch: $BRANCH"
    read -p "Continue publishing from this branch? (y/N): " -n 1 -r
    echo
    if [[ ! $REPLY =~ ^[Yy]$ ]]; then
        echo "Publishing cancelled."
        exit 1
    fi
fi

# Check if there are uncommitted changes
if [[ -n $(git status --porcelain) ]]; then
    echo "Warning: You have uncommitted changes."
    read -p "Continue publishing with uncommitted changes? (y/N): " -n 1 -r
    echo
    if [[ ! $REPLY =~ ^[Yy]$ ]]; then
        echo "Publishing cancelled."
        exit 1
    fi
fi

# Install build tools if not already installed
echo "Installing build tools..."
pip install build twine --quiet

# Clean previous builds
echo "Cleaning previous builds..."
rm -rf dist/ build/ *.egg-info/

# Build the package
echo "Building the package..."
python -m build

# Check if we have the built files
if [ ! -f "dist/sentinel_sdk-0.1.1-py3-none-any.whl" ] || [ ! -f "dist/sentinel_sdk-0.1.1.tar.gz" ]; then
    echo "Error: Package build failed. Required files not found."
    exit 1
fi

# Show what we're about to publish
echo "Files to be published:"
ls -la dist/

# Check if we have the PyPI API token
if [ -z "$PYPI_API_TOKEN" ]; then
    echo "Error: PYPI_API_TOKEN environment variable is not set."
    echo "Please set it with: export PYPI_API_TOKEN=your_token_here"
    exit 1
fi

# Publish to PyPI
echo "Publishing to PyPI..."
twine upload --username __token__ --password "$PYPI_API_TOKEN" dist/*

echo "Python SDK published successfully!"

# Clean up
rm -rf dist/ build/ *.egg-info/

echo "Publishing complete!"