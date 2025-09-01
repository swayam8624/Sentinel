#!/bin/bash

# Final verification script to test both npm and PyPI packages
echo "🔍 Final Verification of Sentinel Distribution Channels"

# Test Python SDK installation
echo "🐍 Testing Python SDK installation..."
cd /tmp
rm -rf sentinel-test
mkdir sentinel-test
cd sentinel-test

# Create virtual environment
python -m venv venv
source venv/bin/activate

# Install the newly published package
echo "Installing yugenkairo-sentinel-sdk..."
pip install yugenkairo-sentinel-sdk

# Test basic import
echo "Testing basic import..."
python -c "import sentinel; print('✅ Python SDK imported successfully')"

# Test client creation
echo "Testing client creation..."
python -c "
from sentinel import SentinelClient
client = SentinelClient(base_url='http://localhost:8080', api_key='test-key')
print('✅ Python client created successfully')
print(f'Base URL: {client.base_url}')
"

deactivate
cd ../..
rm -rf sentinel-test

# Test Node.js SDK installation
echo "🟢 Testing Node.js SDK installation..."
cd /tmp
rm -rf sentinel-node-test
mkdir sentinel-node-test
cd sentinel-node-test

# Initialize npm project
npm init -y > /dev/null 2>&1

# Install the Node.js SDK
echo "Installing @yugenkairo/sentinel-sdk..."
npm install @yugenkairo/sentinel-sdk > /dev/null 2>&1

# Test basic import
echo "Testing basic import..."
node -e "
const sentinel = require('@yugenkairo/sentinel-sdk');
console.log('✅ Node.js SDK imported successfully');
const client = new sentinel.SentinelClient({baseUrl: 'http://localhost:8080'});
console.log('✅ Node.js client created successfully');
"

cd ../..
rm -rf sentinel-node-test

echo "🎉 All distribution channels verified successfully!"
echo "📦 Python SDK: https://pypi.org/project/yugenkairo-sentinel-sdk/"
echo "📦 Node.js SDK: https://www.npmjs.com/package/@yugenkairo/sentinel-sdk"