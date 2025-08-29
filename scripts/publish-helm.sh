#!/bin/bash

# Script to package and publish Helm chart
# Usage: ./scripts/publish-helm.sh

set -e

# Package Helm chart
helm package charts/sentinel

# Create index.yaml
helm repo index . --url https://swayam8624.github.io/Sentinel

# The chart files would need to be pushed to a GitHub Pages branch
# This is a simplified version - in practice, you'd push to a gh-pages branch
echo "Helm chart packaged successfully!"
echo "To publish, push the .tgz file and index.yaml to your GitHub Pages branch."