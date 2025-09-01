#!/bin/bash

# Script to verify the Python SDK installation and functionality
# Usage: ./scripts/verify-python-sdk.sh

set -e

echo "Verifying Sentinel Python SDK..."

# Create a temporary directory for testing
TEMP_DIR=$(mktemp -d)
cd "$TEMP_DIR"

echo "Creating test environment..."
python -m venv test-env
source test-env/bin/activate

echo "Installing sentinel-sdk from local source..."
pip install -e "/Users/swayamsingal/Desktop/Programming/Sentinel/sdk/python"

echo "Testing basic import..."
python -c "import sentinel; print('✓ SDK imported successfully')"

echo "Testing client creation..."
python -c "
from sentinel import SentinelClient
client = SentinelClient(base_url='http://localhost:8080', api_key='test-key')
print('✓ Client created successfully')
print(f'  Base URL: {client.base_url}')
print(f'  API Key: {\"*\" * len(client.api_key) if client.api_key else \"None\"}')
"

echo "Testing available methods..."
python -c "
from sentinel import SentinelClient
client = SentinelClient()
methods = [method for method in dir(client) if not method.startswith('_')]
print('✓ Available client methods:')
for method in methods:
    print(f'  - {method}')
"

echo "Testing chat completions availability..."
python -c "
from sentinel import SentinelClient
client = SentinelClient()
print('✓ Chat completions available:', hasattr(client, 'chat_completions'))
"

echo "Running example script..."
python "/Users/swayamsingal/Desktop/Programming/Sentinel/sdk/python/examples/offline_example.py"

echo "Cleaning up..."
deactivate
cd -
rm -rf "$TEMP_DIR"

echo "✅ Python SDK verification completed successfully!"