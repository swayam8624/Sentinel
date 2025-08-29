#!/bin/bash

# Script to package and publish Helm charts to GitHub Pages
# Usage: ./scripts/publish-helm-charts.sh

set -e

echo "Packaging Helm charts..."
cd /Users/swayamsingal/Desktop/Programming/Sentinel

# Create charts directory if it doesn't exist
mkdir -p docs/charts

# Package the Helm chart
helm package charts/sentinel -d docs/charts

# Update the index
helm repo index docs/charts --url https://swayam8624.github.io/Sentinel/charts

echo "Helm charts packaged successfully!"
echo "To publish to GitHub Pages:"
echo "1. Commit and push the changes to the docs/charts directory"
echo "2. GitHub Pages will automatically serve the charts from https://swayam8624.github.io/Sentinel/charts"