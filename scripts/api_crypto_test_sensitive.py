#!/usr/bin/env python3
"""
Driftlock API Crypto Test - SENSITIVE MODE
Same as api_crypto_test.py but with more sensitive detection settings
to increase the likelihood of detecting anomalies for demos/recordings.

This version uses:
- Smaller window_size (20 instead of 50)
- Smaller baseline_lines (40 instead of 100)
- Lower NCD threshold (0.25 instead of 0.3)
- Higher p-value threshold (0.1 instead of 0.05)

Use this when you want to capture anomalies on Loom!
"""

import asyncio
import time
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

# Data sources
DEFAULT_SOURCE = os.environ.get("CRYPTO_SOURCE", "coingecko").lower()

# Binance WebSocket URL
BINANCE_WS_URL = os.environ.get("BINANCE_WS_URL", "wss://stream.binance.us:9443/ws")

# Default crypto pairs to monitor (Binance)
DEFAULT_STREAMS = [
    "btcusdt@trade",
    "ethusdt@trade",
    "solusdt@trade",
    "linkusdt@trade",
    "avaxusdt@trade",
    "dogeusdt@trade",
    "ltcusdt@trade",
]

def parse_coin_ids(raw: Optional[str]) -> list[str]:
    if not raw:
        return []
    return [c.strip().lower() for c in raw.split(",") if c.strip()]

# CoinGecko polling (no key required)
_DEFAULT_COINGECKO_IDS = [
    "bitcoin",
    "ethereum",
    "solana",
    "chainlink",
    "avalanche-2",
    "dogecoin",
    "litecoin",
]
DEFAULT_COINGECKO_IDS = parse_coin_ids(os.environ.get("COINGECKO_IDS")) or _DEFAULT_COINGECKO_IDS
COINGECKO_URL = "https://api.coingecko.com/api/v3/simple/price"
COINGECKO_VS = os.environ.get("COINGECKO_VS", "usd")
COINGECKO_INTERVAL = float(os.environ.get("COINGECKO_INTERVAL", "5"))
COINGECKO_API_KEY = os.environ.get("COINGECKO_API_KEY")

# API Configuration
API_KEY = os.environ.get("DRIFTLOCK_API_KEY")
API_URL = os.environ.get("DRIFTLOCK_API_URL", "https://driftlock-api-o6kjgrsowq-uc.a.run.app/v1")
DETECT_ENDPOINT = f"{API_URL}/detect"

# Batching settings
BATCH_SIZE = 10  # Send events in batches
BATCH_INTERVAL = 5.0  # seconds between batches

class SensitiveCryptoAPITester:
    """Crypto tester with sensitive anomaly detection settings."""
    
    def __init__(
        self,
        api_key: str,
        api_url: str,
        streams: list,
        source: str,
        coingecko_ids: list[str],
        coingecko_interval: float,
        coingecko_api_key: Optional[str],
    ):
        self.api_key = api_key
        self.api_url = api_url
        self.detect_endpoint = f"{api_url}/detect"
        self.streams = streams
        self.source = source
        self.coingecko_ids = coingecko_ids
        self.coingecko_interval = max(coingecko_interval, 2.5)
        self.coingecko_api_key = coingecko_api_key
        self.event_buffer = []
        self.last_send_time = time.monotonic()
        self.total_events = 0
        self.total_anomalies = 0
        self.session = requests.Session()

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

    def build_coingecko_events(self) -> list[dict]:
        """Fetch prices from CoinGecko and normalize for Driftlock."""
        params = {
            "ids": ",".join(self.coingecko_ids),
            "vs_currencies": COINGECKO_VS,
            "include_24hr_vol": "true",
            "include_24hr_change": "true",
            "precision": "8",
        }

        headers = {}
        if self.coingecko_api_key:
            headers["x-cg-pro-api-key"] = self.coingecko_api_key

        response = self.session.get(
            COINGECKO_URL,
            params=params,
            headers=headers,
            timeout=15,
        )
        response.raise_for_status()
        payload = response.json()

        now = datetime.now(timezone.utc)
        events: list[dict] = []

        for coin_id in self.coingecko_ids:
            data = payload.get(coin_id)
            if not data:
                continue

            price = float(data.get(COINGECKO_VS) or 0.0)
            volume = float(data.get(f"{COINGECKO_VS}_24h_vol") or 0.0)
            change = float(data.get(f"{COINGECKO_VS}_24h_change") or 0.0)

            symbol = coin_id.replace("-", "").upper()
            events.append({
                "timestamp": now.isoformat(),
                "id": f"{coin_id}-{int(now.timestamp())}",
                "type": "crypto_price",
                "symbol": symbol,
                "price": price,
                "volume_usd": volume,
                "change_24h": change,
                "message": f"{symbol} price ${price:.4f} (24h {change:+.2f}%)",
            })

        return events

    async def send_batch(self, events: list):
        """Send a batch with SENSITIVE detection settings."""
        if not events:
            return

        # SENSITIVE SETTINGS for higher anomaly detection rate
        payload = {
            "events": events,
            "window_size": 20,      # Smaller window = more sensitive
            "baseline_lines": 40,    # Smaller baseline = faster to detect changes
            "ncd_threshold": 0.25,   # Lower threshold = more anomalies
            "p_value_threshold": 0.1, # Higher p-value = less strict
        }

        try:
            logger.info(f"📤 Sending batch: {len(events)} events (SENSITIVE MODE)")
            response = self.session.post(
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
                    logger.info(f"🚨🚨🚨 ANOMALY DETECTED! {len(events)} events → {len(anomalies)} anomalies")
                    for i, anomaly in enumerate(anomalies[:5], 1):  # Log first 5
                        explanation = anomaly.get('explanation', anomaly.get('why', 'N/A'))
                        ncd = anomaly.get('metrics', {}).get('ncd', 'N/A')
                        confidence = anomaly.get('metrics', {}).get('confidence', 'N/A')
                        logger.info(f"   🎯 Anomaly #{i}: NCD={ncd}, Confidence={confidence}")
                        logger.info(f"      {explanation[:100]}")
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
            if self.total_events % 50 == 0:
                logger.info(f"ℹ️  Received {self.total_events} events (buffer={len(self.event_buffer)})")

            current_time = time.monotonic()
            should_send = (
                len(self.event_buffer) >= BATCH_SIZE or
                (self.event_buffer and current_time - self.last_send_time >= BATCH_INTERVAL)
            )

            if should_send:
                batch = self.event_buffer.copy()
                self.event_buffer.clear()
                self.last_send_time = current_time
                await self.send_batch(batch)

    async def stream_binance(self):
        """Connect to Binance WebSocket and stream data."""
        backoff = 1
        ssl_context = ssl.create_default_context(cafile=certifi.where())

        logger.info("🚀 Starting Driftlock Crypto API Test (SENSITIVE MODE)")
        logger.info(f"   Source: Binance WebSocket")
        logger.info(f"   API: {self.api_url}")
        logger.info(f"   Settings: window_size=20, baseline_lines=40, ncd_threshold=0.25")
        logger.info("")

        while True:
            try:
                logger.info("Connecting to Binance WebSocket...")
                async with websockets.connect(BINANCE_WS_URL, ssl=ssl_context) as websocket:
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

                    backoff = 1

                    async for message in websocket:
                        try:
                            data = json.loads(message)
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

    async def poll_coingecko(self):
        """Poll CoinGecko every few seconds and send price batches."""
        backoff = self.coingecko_interval

        logger.info("🚀 Starting Driftlock Crypto API Test (SENSITIVE MODE)")
        logger.info("   Source: CoinGecko (REST polling, no API key required)")
        logger.info(f"   API: {self.api_url}")
        logger.info(f"   Coins: {', '.join(self.coingecko_ids)}")
        logger.info(f"   Poll interval: {self.coingecko_interval}s")
        logger.info(f"   Settings: window_size=20, baseline_lines=40, ncd_threshold=0.25")
        logger.info("")

        while True:
            try:
                events = self.build_coingecko_events()
                if events:
                    self.total_events += len(events)
                    self.last_send_time = time.monotonic()
                    await self.send_batch(events)
                else:
                    logger.warning("CoinGecko returned no data; check coin IDs")
                backoff = self.coingecko_interval
                await asyncio.sleep(self.coingecko_interval)
            except requests.exceptions.HTTPError as err:
                status = err.response.status_code if err.response else None
                backoff = min(backoff * 2, 90)
                if status == 429:
                    logger.warning(f"CoinGecko rate limited (429). Backing off {backoff:.0f}s")
                else:
                    logger.warning(f"CoinGecko HTTP error {status}: {err}")
                await asyncio.sleep(backoff)
            except requests.exceptions.RequestException as err:
                backoff = min(backoff * 2, 60)
                logger.warning(f"CoinGecko request failed: {err}. Retrying in {backoff:.0f}s")
                await asyncio.sleep(backoff)
            except Exception as err:
                backoff = min(backoff * 2, 60)
                logger.error(f"Unexpected CoinGecko error: {err}")
                await asyncio.sleep(backoff)

    async def run(self):
        if self.source == "binance":
            await self.stream_binance()
        else:
            await self.poll_coingecko()

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
        description="Stream live crypto data to Driftlock API (SENSITIVE MODE for demos)"
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
        "--source",
        default=DEFAULT_SOURCE,
        choices=["coingecko", "binance"],
        help="Data source to use (coingecko|binance). Default: coingecko"
    )
    parser.add_argument(
        "--streams",
        default=",".join(DEFAULT_STREAMS),
        help="Comma-separated list of Binance streams"
    )
    parser.add_argument(
        "--coins",
        default=",".join(DEFAULT_COINGECKO_IDS),
        help="Comma-separated list of CoinGecko coin IDs"
    )
    parser.add_argument(
        "--poll-interval",
        type=float,
        default=COINGECKO_INTERVAL,
        help="Seconds between CoinGecko polls"
    )
    parser.add_argument(
        "--cg-api-key",
        default=COINGECKO_API_KEY,
        help="CoinGecko API key (optional)"
    )

    args = parser.parse_args()

    if not args.api_key:
        logger.error("❌ API key required. Set DRIFTLOCK_API_KEY env var or use --api-key")
        sys.exit(1)

    streams = [s.strip().lower() for s in args.streams.split(",") if s.strip()]
    coingecko_ids = [c.strip().lower() for c in args.coins.split(",") if c.strip()]

    tester = SensitiveCryptoAPITester(
        args.api_key,
        args.api_url,
        streams,
        args.source.lower(),
        coingecko_ids,
        args.poll_interval,
        args.cg_api_key,
    )

    try:
        asyncio.run(tester.run())
    except KeyboardInterrupt:
        logger.info("")
        logger.info("Stopping test...")
        asyncio.run(tester.shutdown())

if __name__ == "__main__":
    main()

