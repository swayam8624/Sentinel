#!/bin/bash

# Script to publish Python SDK to PyPI
# Usage: ./scripts/publish-python-sdk.sh

set -e

echo "Publishing Python SDK to PyPI..."

cd /Users/swayamsingal/Desktop/Programming/Sentinel/sdk/python

# Check if setup.py exists
if [ ! -f setup.py ]; then
    echo "Error: setup.py not found in sdk/python directory"
    exit 1
fi

# Check if we're logged in to PyPI
if ! python -m twine check dist/* > /dev/null 2>&1; then
    echo "Warning: twine check failed. Make sure you have twine installed and configured."
fi

# Build the package
echo "Building Python package..."
python setup.py sdist bdist_wheel

# Publish to PyPI
echo "Publishing to PyPI..."
python -m twine upload dist/*

echo "Python SDK published successfully!"
echo "Package is now available at: https://pypi.org/project/sentinel-sdk/"