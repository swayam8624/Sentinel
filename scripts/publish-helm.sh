#!/bin/bash

# Script to package and publish Helm chart
# Usage: ./scripts/publish-helm.sh

set -e

# Package Helm chart
helm package charts/sentinel

# Create index.yaml
helm repo index . --url https://swayam8624.github.io/Sentinel

# Move packaged chart to docs directory for GitHub Pages
mkdir -p docs/charts
mv sentinel-*.tgz docs/charts/

# Update index.yaml
helm repo index docs/charts --url https://swayam8624.github.io/Sentinel/charts

echo "Helm chart packaged and index updated successfully!"
echo "Charts are now available in the docs/charts directory for GitHub Pages."