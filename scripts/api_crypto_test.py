#!/usr/bin/env python3
"""
Driftlock API Crypto Test - Live Binance Data Streaming
This script streams live Binance cryptocurrency trade data and sends it to the Driftlock API
for real-time anomaly detection. Perfect for demonstrating the platform's capabilities.

Usage:
    export DRIFTLOCK_API_KEY="dlk_..."
    export DRIFTLOCK_API_URL="https://driftlock.web.app/api/v1"
    python3 scripts/api_crypto_test.py

Or specify API key and URL:
    python3 scripts/api_crypto_test.py --api-key "dlk_..." --api-url "https://driftlock.web.app/api/v1"
"""

import asyncio
import json
import sys
import logging
import ssl
import os
import websockets
import certifi
import requests
import argparse
from datetime import datetime, timezone
from typing import Optional

# Configure logging
logging.basicConfig(
    level=logging.INFO,
    format='%(asctime)s - %(levelname)s - %(message)s'
)
logger = logging.getLogger(__name__)

# Binance WebSocket URL
BINANCE_WS_URL = os.environ.get("BINANCE_WS_URL", "wss://stream.binance.us:9443/ws")

# Default crypto pairs to monitor
DEFAULT_STREAMS = [
    "btcusdt@trade",
    "ethusdt@trade",
    "solusdt@trade",
    "linkusdt@trade",
    "avaxusdt@trade",
    "dogeusdt@trade",
    "ltcusdt@trade",
]

# API Configuration
API_KEY = os.environ.get("DRIFTLOCK_API_KEY")
API_URL = os.environ.get("DRIFTLOCK_API_URL", "https://driftlock.web.app/api/v1")
DETECT_ENDPOINT = f"{API_URL}/detect"

# Batching settings
BATCH_SIZE = 10  # Send events in batches
BATCH_INTERVAL = 5.0  # seconds between batches

class CryptoAPITester:
    def __init__(self, api_key: str, api_url: str, streams: list):
        self.api_key = api_key
        self.api_url = api_url
        self.detect_endpoint = f"{api_url}/detect"
        self.streams = streams
        self.event_buffer = []
        self.last_send_time = asyncio.get_event_loop().time()
        self.total_events = 0
        self.total_anomalies = 0

    def normalize_trade(self, data: dict) -> Optional[dict]:
        """Normalize Binance trade event for Driftlock API."""
        try:
            if "e" not in data or data["e"] != "trade":
                return None

            price = float(data["p"])
            quantity = float(data["q"])
            symbol = data["s"]
            action = "SELL" if data["m"] else "BUY"

            normalized = {
                "timestamp": datetime.fromtimestamp(data["T"] / 1000.0, timezone.utc).isoformat(),
                "id": str(data["t"]),
                "type": "crypto_trade",
                "symbol": symbol,
                "price": price,
                "quantity": quantity,
                "volume_usd": price * quantity,
                "side": action,
                "message": f"{action} {quantity:.4f} {symbol} @ {price:.8f}"
            }
            return normalized
        except Exception as e:
            logger.debug(f"Error normalizing trade: {e}")
            return None

    async def send_batch(self, events: list):
        """Send a batch of events to the Driftlock API."""
        if not events:
            return

        payload = {
            "events": events,
            "window_size": 50,  # Standard window size
            "baseline_lines": 100  # Baseline for comparison
        }

        try:
            response = requests.post(
                self.detect_endpoint,
                headers={
                    "X-Api-Key": self.api_key,
                    "Content-Type": "application/json"
                },
                json=payload,
                timeout=10
            )

            if response.status_code == 200:
                result = response.json()
                anomalies = result.get("anomalies", [])
                if anomalies:
                    self.total_anomalies += len(anomalies)
                    logger.info(f"✅ Batch sent: {len(events)} events → {len(anomalies)} anomalies detected")
                    for anomaly in anomalies[:3]:  # Log first 3
                        logger.info(f"   Anomaly: {anomaly.get('explanation', 'N/A')[:80]}")
                else:
                    logger.info(f"✅ Batch sent: {len(events)} events → No anomalies")
            elif response.status_code == 401:
                logger.error("❌ Authentication failed. Check your API key.")
                sys.exit(1)
            else:
                logger.warning(f"⚠️  API returned {response.status_code}: {response.text[:200]}")
        except requests.exceptions.RequestException as e:
            logger.error(f"❌ Failed to send batch: {e}")

    async def process_event(self, event: dict):
        """Process a single event and add to buffer."""
        normalized = self.normalize_trade(event)
        if normalized:
            self.event_buffer.append(normalized)
            self.total_events += 1

            # Check if we should send a batch
            current_time = asyncio.get_event_loop().time()
            should_send = (
                len(self.event_buffer) >= BATCH_SIZE or
                (self.event_buffer and current_time - self.last_send_time >= BATCH_INTERVAL)
            )

            if should_send:
                batch = self.event_buffer.copy()
                self.event_buffer.clear()
                self.last_send_time = current_time
                await self.send_batch(batch)

    async def connect_and_stream(self):
        """Connect to Binance WebSocket and stream data."""
        backoff = 1
        ssl_context = ssl.create_default_context(cafile=certifi.where())

        logger.info(f"🚀 Starting Driftlock Crypto API Test")
        logger.info(f"   API: {self.api_url}")
        logger.info(f"   Streams: {self.streams}")
        logger.info(f"   Batch size: {BATCH_SIZE} events")
        logger.info(f"   Batch interval: {BATCH_INTERVAL}s")
        logger.info("")

        while True:
            try:
                logger.info(f"Connecting to Binance WebSocket...")
                async with websockets.connect(BINANCE_WS_URL, ssl=ssl_context) as websocket:
                    # Subscribe to streams
                    payload = {
                        "method": "SUBSCRIBE",
                        "params": self.streams,
                        "id": 1
                    }
                    await websocket.send(json.dumps(payload))
                    logger.info(f"✅ Subscribed to {len(self.streams)} streams")
                    logger.info("📊 Streaming live crypto data to Driftlock API...")
                    logger.info("   (Press Ctrl+C to stop)")
                    logger.info("")

                    backoff = 1  # Reset backoff on success

                    async for message in websocket:
                        try:
                            data = json.loads(message)
                            # Handle subscription response
                            if "result" in data and data["id"] == 1:
                                continue

                            await self.process_event(data)

                        except json.JSONDecodeError:
                            continue

            except (websockets.exceptions.ConnectionClosed, OSError) as e:
                logger.error(f"Connection lost: {e}")
                logger.info(f"Retrying in {backoff} seconds...")
                await asyncio.sleep(backoff)
                backoff = min(backoff * 2, 60)
            except Exception as e:
                logger.error(f"Unexpected error: {e}")
                await asyncio.sleep(5)

    async def shutdown(self):
        """Send any remaining events before shutdown."""
        if self.event_buffer:
            logger.info("Sending final batch...")
            await self.send_batch(self.event_buffer)
            self.event_buffer.clear()

        logger.info("")
        logger.info("📊 Test Summary:")
        logger.info(f"   Total events processed: {self.total_events}")
        logger.info(f"   Total anomalies detected: {self.total_anomalies}")
        if self.total_events > 0:
            anomaly_rate = (self.total_anomalies / self.total_events) * 100
            logger.info(f"   Anomaly rate: {anomaly_rate:.2f}%")

def main():
    parser = argparse.ArgumentParser(
        description="Stream live Binance crypto data to Driftlock API for anomaly detection"
    )
    parser.add_argument(
        "--api-key",
        default=API_KEY,
        help="Driftlock API key (or set DRIFTLOCK_API_KEY env var)"
    )
    parser.add_argument(
        "--api-url",
        default=API_URL,
        help="Driftlock API base URL (or set DRIFTLOCK_API_URL env var)"
    )
    parser.add_argument(
        "--streams",
        default=",".join(DEFAULT_STREAMS),
        help="Comma-separated list of Binance streams (e.g., btcusdt@trade,ethusdt@trade)"
    )

    args = parser.parse_args()

    if not args.api_key:
        logger.error("❌ API key required. Set DRIFTLOCK_API_KEY env var or use --api-key")
        sys.exit(1)

    streams = [s.strip().lower() for s in args.streams.split(",") if s.strip()]

    tester = CryptoAPITester(args.api_key, args.api_url, streams)

    try:
        asyncio.run(tester.connect_and_stream())
    except KeyboardInterrupt:
        logger.info("")
        logger.info("Stopping test...")
        asyncio.run(tester.shutdown())

if __name__ == "__main__":
    main()

