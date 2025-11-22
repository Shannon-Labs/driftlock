#!/usr/bin/env python3
import asyncio
import json
import sys
import logging
import ssl
import os
import websockets
import certifi
from datetime import datetime, timezone

# Configure logging to stderr
logging.basicConfig(stream=sys.stderr, level=logging.INFO, format='%(asctime)s - %(levelname)s - %(message)s')
logger = logging.getLogger(__name__)

# Binance WebSocket URL (default: binance.us to avoid geo blocks; override with BINANCE_WS_URL for global)
URI = os.environ.get("BINANCE_WS_URL", "wss://stream.binance.us:9443/ws")

# Pairs to monitor (lowercase required for streams). Can override with BINANCE_STREAMS="btcusdt,ethusdt,..."
# Keep to pairs listed on binance.us by default to ensure subscriptions succeed.
DEFAULT_STREAMS = [
    "btcusdt@trade",
    "ethusdt@trade",
    "solusdt@trade",
    "linkusdt@trade",
    "avaxusdt@trade",
    "dogeusdt@trade",
    "ltcusdt@trade",
]
STREAMS = [s.strip().lower() for s in os.environ.get("BINANCE_STREAMS", ",".join(DEFAULT_STREAMS)).split(",") if s.strip()]

async def subscribe(websocket):
    payload = {
        "method": "SUBSCRIBE",
        "params": STREAMS,
        "id": 1
    }
    await websocket.send(json.dumps(payload))
    logger.info(f"Subscribed to streams: {STREAMS}")

def normalize_trade(data):
    """
    Normalize Binance trade event to flat JSON for Driftlock.
    Binance Trade Event:
    {
      "e": "trade",     // Event type
      "E": 123456789,   // Event time
      "s": "BNBBTC",    // Symbol
      "t": 12345,       // Trade ID
      "p": "0.001",     // Price
      "q": "100",       // Quantity
      "b": 88,          // Buyer order ID
      "a": 50,          // Seller order ID
      "T": 123456785,   // Trade time
      "m": true,        // Is the buyer the market maker?
      "M": true         // Ignore
    }
    """
    try:
        if "e" not in data or data["e"] != "trade":
            return None

        price = float(data["p"])
        quantity = float(data["q"])
        symbol = data["s"]
        
        # Construct a human-readable message for Driftlock's explanation
        action = "SELL" if data["m"] else "BUY" # m=true means buyer is maker (limit order), so taker is seller -> SELL
        
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
        # logger.warning(f"Error normalizing: {e}")
        return None

async def connect_and_stream():
    backoff = 1
    ssl_context = ssl.create_default_context(cafile=certifi.where())
    while True:
        try:
            logger.info(f"Connecting to {URI}...")
            async with websockets.connect(URI, ssl=ssl_context) as websocket:
                await subscribe(websocket)
                backoff = 1 # Reset backoff on success
                
                async for message in websocket:
                    try:
                        data = json.loads(message)
                        # Handle subscription response
                        if "result" in data and data["id"] == 1:
                            continue
                            
                        normalized = normalize_trade(data)
                        if normalized:
                            # Write to stdout for Driftlock
                            print(json.dumps(normalized), flush=True)
                            
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

if __name__ == "__main__":
    try:
        asyncio.run(connect_and_stream())
    except KeyboardInterrupt:
        logger.info("Bridge stopped by user")
