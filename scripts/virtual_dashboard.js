#!/usr/bin/env node
// Render a lightweight dashboard onto the virtual display using Playwright + Chromium

import fs from 'fs';
import path from 'path';
import { fileURLToPath } from 'url';

let chromium;
try {
    ({ chromium } = await import('playwright'));
} catch (err) {
    console.error('Playwright is not installed. Install with: npm install --save-dev playwright');
    process.exit(1);
}

const __dirname = path.dirname(fileURLToPath(import.meta.url));
const repoRoot = path.resolve(__dirname, '..');

const alertLog = path.resolve(process.env.ALERT_LOG ?? path.join(repoRoot, 'logs/anomaly_alerts.log'));
const refreshMs = Number(process.env.DASHBOARD_REFRESH_MS ?? 4000);
const resolution = (process.env.DASHBOARD_RESOLUTION ?? '1600x900').split('x');
const viewport = { width: Number(resolution[0]) || 1600, height: Number(resolution[1]) || 900 };

const escapeHtml = (value) =>
    value
        .replace(/&/g, '&amp;')
        .replace(/</g, '&lt;')
        .replace(/>/g, '&gt;')
        .replace(/"/g, '&quot;')
        .replace(/'/g, '&#39;');

const tailLog = (limit = 24) => {
    if (!fs.existsSync(alertLog)) return 'Waiting for anomalies...';
    const lines = fs.readFileSync(alertLog, 'utf8').trim().split('\n');
    return lines.slice(-limit).join('\n');
};

const htmlForLog = (logText) => {
    const lastAnomaly = logText
        .split('\n')
        .reverse()
        .find((line) => line.toLowerCase().includes('anomaly id')) || 'No anomalies yet';
    const timestamp = new Date().toISOString();
    return `<!doctype html>
<html lang="en">
<head>
    <meta charset="utf-8" />
    <title>Driftlock Virtual Dashboard</title>
    <style>
        body { margin: 0; font-family: "Inter", "Helvetica Neue", Arial, sans-serif; background: #0b1220; color: #d9e1ff; }
        .wrap { display: flex; height: 100vh; padding: 24px; box-sizing: border-box; gap: 24px; }
        .panel { flex: 1; background: radial-gradient(circle at 20% 20%, #122440 0%, #0b1220 55%); border: 1px solid #1f2f4f; border-radius: 16px; padding: 20px; box-shadow: 0 18px 45px rgba(0,0,0,0.35); }
        h1 { margin: 0 0 8px 0; font-size: 28px; letter-spacing: 0.4px; }
        .meta { color: #8fa3c9; margin-bottom: 16px; font-size: 14px; }
        .badge { display: inline-block; background: #1d3557; color: #b9d0ff; padding: 6px 10px; border-radius: 10px; font-size: 12px; margin-right: 8px; }
        pre { background: #0b1528; color: #c8e0ff; border-radius: 12px; padding: 16px; font-size: 13px; line-height: 1.45; overflow: hidden; border: 1px solid #1f2f4f; white-space: pre-wrap; }
    </style>
</head>
<body>
    <div class="wrap">
        <div class="panel">
            <h1>Driftlock Virtual Dashboard</h1>
            <div class="meta">
                <span class="badge">Display ${escapeHtml(process.env.DISPLAY || '')}</span>
                <span class="badge">Last anomaly: ${escapeHtml(lastAnomaly)}</span>
                <span class="badge">Updated: ${escapeHtml(timestamp)}</span>
            </div>
            <pre>${escapeHtml(logText)}</pre>
        </div>
    </div>
</body>
</html>`;
};

const renderLoop = async (page) => {
    const refresh = async () => {
        const logText = tailLog();
        await page.setContent(htmlForLog(logText));
    };

    await refresh();
    setInterval(refresh, refreshMs);
};

const main = async () => {
    const browser = await chromium.launch({ headless: false });
    const page = await browser.newPage({ viewport });

    await renderLoop(page);

    const close = async () => {
        await browser.close();
        process.exit(0);
    };

    process.on('SIGINT', close);
    process.on('SIGTERM', close);

    // Keep process alive while intervals run
    await new Promise(() => {});
};

main().catch((err) => {
    console.error('Failed to start virtual dashboard', err);
    process.exit(1);
});

