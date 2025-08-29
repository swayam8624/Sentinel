#!/bin/bash

# Script to create a GitHub release
# Usage: ./scripts/create-release.sh <tag_name> <release_name>

set -e

if [ $# -ne 2 ]; then
    echo "Usage: $0 <tag_name> <release_name>"
    exit 1
fi

TAG_NAME=$1
RELEASE_NAME=$2

# Read release notes
RELEASE_NOTES=$(cat RELEASE_NOTES.md)

# Create release using GitHub API
curl -X POST \
  -H "Accept: application/vnd.github.v3+json" \
  -H "Authorization: token $GITHUB_TOKEN" \
  https://api.github.com/repos/swayam8624/Sentinel/releases \
  -d "{
    \"tag_name\": \"$TAG_NAME\",
    \"name\": \"$RELEASE_NAME\",
    \"body\": \"$RELEASE_NOTES\",
    \"draft\": false,
    \"prerelease\": false
  }"

echo "Release created successfully!"