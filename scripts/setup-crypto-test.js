#!/usr/bin/env node
/**
 * Driftlock Crypto Test Setup Script (Node.js)
 * 
 * Date: 2025-11-22
 * Pre-req: gcloud logged in with access to driftlock project; repo at /Volumes/VIXinSSD/driftlock
 * 
 * This script:
 * 1. Ensures .env has an API key (Cloud Run job path; safe with Cloud SQL socket URL)
 * 2. Loads environment variables
 * 3. Performs a quick API smoke test
 * 4. Starts the 4-hour Binance stream when ready
 * 5. Provides monitoring commands
 */

import fs from 'fs';
import path from 'path';
import { spawn } from 'child_process';
import https from 'https';
import http from 'http';
import { fileURLToPath } from 'url';
import { dirname } from 'path';

const __filename = fileURLToPath(import.meta.url);
const __dirname = dirname(__filename);
const REPO_ROOT = path.resolve(__dirname, '..');
const ENV_FILE = path.join(REPO_ROOT, '.env');
const LOGS_DIR = path.join(REPO_ROOT, 'logs');

// ANSI color codes for terminal output
const colors = {
  reset: '\x1b[0m',
  green: '\x1b[32m',
  yellow: '\x1b[33m',
  red: '\x1b[31m',
  blue: '\x1b[34m',
  cyan: '\x1b[36m',
};

function log(message, color = 'reset') {
  console.log(`${colors[color]}${message}${colors.reset}`);
}

function logStep(step, message) {
  log(`\n${step}. ${message}`, 'cyan');
}

/**
 * Check if .env file exists and contains DRIFTLOCK_API_KEY
 */
function checkEnvFile() {
  if (!fs.existsSync(ENV_FILE)) {
    return false;
  }
  
  const envContent = fs.readFileSync(ENV_FILE, 'utf8');
  return envContent.includes('DRIFTLOCK_API_KEY=');
}

/**
 * Load environment variables from .env file
 */
function loadEnvFile() {
  if (!fs.existsSync(ENV_FILE)) {
    return {};
  }
  
  const envContent = fs.readFileSync(ENV_FILE, 'utf8');
  const env = {};
  
  envContent.split('\n').forEach(line => {
    const trimmed = line.trim();
    if (trimmed && !trimmed.startsWith('#')) {
      const [key, ...valueParts] = trimmed.split('=');
      if (key && valueParts.length > 0) {
        env[key] = valueParts.join('=');
      }
    }
  });
  
  return env;
}

/**
 * Create API key using Cloud Run job script
 */
async function createApiKey() {
  return new Promise((resolve, reject) => {
    log('   Attempting to create via Cloud Run Job...', 'yellow');
    
    const scriptPath = path.join(__dirname, 'create-test-api-key-cloudrun.sh');
    const child = spawn('bash', [scriptPath], {
      cwd: REPO_ROOT,
      stdio: 'inherit',
      env: { ...process.env }
    });
    
    child.on('close', (code) => {
      if (code === 0) {
        log('✅ API key created!', 'green');
        resolve(true);
      } else {
        log('❌ Failed to create API key', 'red');
        reject(new Error(`Script exited with code ${code}`));
      }
    });
    
    child.on('error', (err) => {
      log(`❌ Error running script: ${err.message}`, 'red');
      reject(err);
    });
  });
}

/**
 * Perform API smoke test
 */
async function smokeTest(apiUrl, apiKey) {
  return new Promise((resolve, reject) => {
    const testPayload = {
      events: [{
        id: "1",
        type: "crypto_trade",
        timestamp: new Date().toISOString(),
        symbol: "BTCUSDT",
        price: 43000,
        quantity: 0.1,
        volume_usd: 4300,
        side: "BUY",
        message: "BUY 0.1 BTCUSDT @ 43000"
      }],
      window_size: 50,
      baseline_lines: 100
    };
    
    const url = new URL(`${apiUrl}/detect`);
    const isHttps = url.protocol === 'https:';
    const client = isHttps ? https : http;
    
    const options = {
      hostname: url.hostname,
      port: url.port || (isHttps ? 443 : 80),
      path: url.pathname + url.search,
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
        'X-Api-Key': apiKey,
      },
    };
    
    log('   Testing API endpoint...', 'yellow');
    
    const req = client.request(options, (res) => {
      let data = '';
      
      res.on('data', (chunk) => {
        data += chunk;
        // Limit response to first 500 chars
        if (data.length > 500) {
          req.destroy();
        }
      });
      
      res.on('end', () => {
        // All response codes are acceptable for smoke test (even 401/403 means API is reachable)
        log(`   ✅ API responded with status ${res.statusCode}`, 'green');
        log(`   Response preview: ${data.substring(0, 200)}...`, 'cyan');
        resolve(true);
      });
    });
    
    req.on('error', (err) => {
      log(`   ❌ Request failed: ${err.message}`, 'red');
      reject(err);
    });
    
    req.setTimeout(10000, () => {
      req.destroy();
      reject(new Error('Request timeout'));
    });
    
    req.write(JSON.stringify(testPayload));
    req.end();
  });
}

/**
 * Start the crypto test script
 */
function startCryptoTest() {
  return new Promise((resolve, reject) => {
    const scriptPath = path.join(__dirname, 'start_crypto_test.sh');
    
    log('\n📊 Starting 4-hour Binance stream...', 'cyan');
    log('   Logs will be saved to: logs/crypto-api-test-*.log', 'cyan');
    
    const child = spawn('bash', [scriptPath], {
      cwd: REPO_ROOT,
      stdio: 'inherit',
      env: { ...process.env }
    });
    
    child.on('close', (code) => {
      if (code === 0) {
        resolve(true);
      } else {
        reject(new Error(`Script exited with code ${code}`));
      }
    });
    
    child.on('error', (err) => {
      reject(err);
    });
    
    // Store child process for potential cleanup
    process.on('SIGINT', () => {
      log('\n⚠️  Interrupted. Stopping crypto test...', 'yellow');
      child.kill('SIGINT');
      process.exit(0);
    });
  });
}

/**
 * Main execution
 */
async function main() {
  log('🚀 Driftlock Crypto Test Setup', 'blue');
  log('   Date: 2025-11-22', 'cyan');
  log('   Repo: ' + REPO_ROOT, 'cyan');
  
  // Step 1: Ensure .env has an API key
  logStep('1', 'Checking for API key in .env file');
  
  let env = loadEnvFile();
  let apiKey = env.DRIFTLOCK_API_KEY || process.env.DRIFTLOCK_API_KEY;
  
  if (!apiKey || !checkEnvFile()) {
    log('⚠️  No API key found in .env file', 'yellow');
    
    try {
      await createApiKey();
      // Reload env after creating key
      env = loadEnvFile();
      apiKey = env.DRIFTLOCK_API_KEY;
    } catch (err) {
      log('\n❌ Failed to create API key automatically.', 'red');
      log('\nPlease either:', 'yellow');
      log('  1. Run manually: ./scripts/create-test-api-key-cloudrun.sh', 'yellow');
      log('  2. Or sign up at https://driftlock.web.app and set:', 'yellow');
      log('     export DRIFTLOCK_API_KEY=\'dlk_...\'', 'yellow');
      process.exit(1);
    }
  }
  
  if (!apiKey) {
    log('❌ API key still not found after creation attempt', 'red');
    process.exit(1);
  }
  
  // Step 2: Load environment variables
  logStep('2', 'Loading environment variables');
  
  // Merge loaded env with process.env
  Object.assign(process.env, env);
  
  const apiUrl = env.DRIFTLOCK_API_URL || process.env.DRIFTLOCK_API_URL || 'https://driftlock.web.app/api/v1';
  
  log(`   ✅ API Key: ${apiKey.substring(0, 20)}...`, 'green');
  log(`   ✅ API URL: ${apiUrl}`, 'green');
  
  // Step 3: Quick API smoke test
  logStep('3', 'Performing API smoke test');
  
  try {
    await smokeTest(apiUrl, apiKey);
  } catch (err) {
    log(`\n⚠️  Smoke test failed: ${err.message}`, 'yellow');
    log('   Continuing anyway...', 'yellow');
  }
  
  // Step 4: Start the 4-hour Binance stream
  logStep('4', 'Starting 4-hour Binance stream');
  log('   (Press Ctrl+C to stop)', 'yellow');
  
  try {
    await startCryptoTest();
  } catch (err) {
    log(`\n❌ Failed to start crypto test: ${err.message}`, 'red');
    process.exit(1);
  }
  
  // Step 5: Monitoring instructions
  logStep('5', 'Monitoring');
  log('\nTo monitor progress:', 'cyan');
  log('   tail -f logs/crypto-api-test-*.log', 'yellow');
  log('\nTo stop:', 'cyan');
  log('   kill $(cat logs/crypto-api-test-*.pid)', 'yellow');
}

// Run main function
main().catch((err) => {
  log(`\n❌ Fatal error: ${err.message}`, 'red');
  console.error(err);
  process.exit(1);
});
