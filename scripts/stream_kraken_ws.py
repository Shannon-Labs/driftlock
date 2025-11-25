#!/usr/bin/env python3
"""
High-frequency trade streamer from Kraken WebSocket -> NDJSON.

Usage:
  python3 -u scripts/stream_kraken_ws.py | driftlock scan --stdin --follow --format ndjson --baseline-lines 120 --algo entropy --show-all
"""

import json
import os
import sys
import time
import traceback

import websocket


PAIR = os.environ.get("KRAKEN_PAIR", "XBT/USD")
WS_URL = "wss://ws.kraken.com"


def stream_trades() -> None:
    """Connect to Kraken trades feed and emit NDJSON rows."""
    trade_counter = 0
    
    while True:
        ws = None
        try:
            print(f"🔌 Connecting to Kraken WebSocket ({PAIR})...", file=sys.stderr, flush=True)
            ws = websocket.create_connection(WS_URL, timeout=10)
            ws.settimeout(10)
            subscribe = {
                "event": "subscribe",
                "pair": [PAIR],
                "subscription": {"name": "trade"},
            }
            ws.send(json.dumps(subscribe))
            print(f"✅ Connected to Kraken - streaming {PAIR} trades", file=sys.stderr, flush=True)

            while True:
                raw = ws.recv()
                if not raw or not raw.startswith("["):
                    continue  # ignore heartbeats/status

                msg = json.loads(raw)
                if len(msg) < 4 or msg[2] != "trade":
                    continue

                trades = msg[1]
                for idx, trade in enumerate(trades):
                    # trade layout: [price, volume, time, side, orderType, misc]
                    try:
                        price = float(trade[0])
                        qty = float(trade[1])
                        ts = float(trade[2])
                        side = "buy" if trade[3] == "b" else "sell"
                        out = {
                            "ts": ts,
                            "price": price,
                            "qty": qty,
                            "side": side,
                            "pair": PAIR,
                            "id": f"kraken-{int(ts * 1000)}-{idx}",
                            "source": "kraken",
                        }
                        json.dump(out, sys.stdout)
                        sys.stdout.write("\n")
                        sys.stdout.flush()
                        
                        trade_counter += 1
                        if trade_counter % 100 == 0:
                            print(f"📊 Health: {trade_counter} trades processed", file=sys.stderr, flush=True)
                    except Exception:
                        traceback.print_exc(file=sys.stderr)
        except KeyboardInterrupt:
            print("⏹️  Stopped by user", file=sys.stderr, flush=True)
            break
        except Exception as e:
            print(f"❌ Connection error: {e}", file=sys.stderr, flush=True)
            traceback.print_exc(file=sys.stderr)
            print("🔄 Reconnecting in 3 seconds...", file=sys.stderr, flush=True)
            time.sleep(3)
        finally:
            try:
                if ws is not None:
                    ws.close()
                    print("🔌 WebSocket closed", file=sys.stderr, flush=True)
            except Exception:
                pass


if __name__ == "__main__":
    stream_trades()
