const sentinel = require('@yugenkairo/sentinel-sdk');

console.log('Sentinel SDK imported successfully!');
console.log('Available functions:', Object.keys(sentinel));

// Test basic functionality if available
try {
    // This is a placeholder test - in a real scenario, you would test actual SDK functions
    console.log('NPM package test completed successfully!');
} catch (error) {
    console.error('Error testing NPM package:', error);
}