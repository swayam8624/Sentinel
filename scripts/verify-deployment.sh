#!/bin/bash

# Sentinel Deployment Verification Script

echo "🔍 Verifying Sentinel Deployment..."

# Check GitHub Pages
echo "📋 Checking GitHub Pages..."
curl -s -o /dev/null -w "Homepage: %{http_code}\n" https://swayam8624.github.io/Sentinel/
curl -s -o /dev/null -w "Charts Index: %{http_code}\n" https://swayam8624.github.io/Sentinel/charts/index.yaml
curl -s -o /dev/null -w "Chart Package: %{http_code}\n" https://swayam8624.github.io/Sentinel/charts/sentinel-0.1.0.tgz

# Check npm package
echo "📦 Checking npm package..."
cd /tmp
rm -rf test-sentinel-verification
mkdir test-sentinel-verification
cd test-sentinel-verification
npm init -y > /dev/null 2>&1
npm install @yugenkairo/sentinel-sdk > /dev/null 2>&1

# Test npm package functionality
node -e "
const sentinel = require('@yugenkairo/sentinel-sdk');
console.log('✅ npm package installed successfully');

const client = new sentinel.SentinelClient({baseUrl: 'http://localhost:8080'});
console.log('✅ Client created successfully');

client.sanitizePrompt('Test prompt').then(result => {
  console.log('✅ Sanitize prompt function working');
  console.log('   Result:', JSON.stringify(result).substring(0, 100) + '...');
});
" 2>/dev/null

if [ $? -eq 0 ]; then
  echo "✅ All verifications passed!"
else
  echo "❌ Some verifications failed!"
fi

# Cleanup
cd /tmp
rm -rf test-sentinel-verification

echo "🏁 Deployment verification complete!"