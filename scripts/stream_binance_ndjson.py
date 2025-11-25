#!/usr/bin/env python3
"""
Stream Binance trades as NDJSON for driftlock scan --stdin
Connects to Binance WebSocket and outputs normalized trade data.
Optional synthetic spike injection for demo reliability.
"""

import sys
import json
import time
import random
from datetime import datetime
import websocket
import argparse
import threading

# Statistics
trade_count = 0
start_time = time.time()
synthetic_injected = 0


def print_stats():
    """Print statistics every 60 seconds"""
    global trade_count, start_time, synthetic_injected
    while True:
        time.sleep(60)
        elapsed = time.time() - start_time
        rate = trade_count / elapsed if elapsed > 0 else 0
        print(f"📊 Stats: {trade_count} trades | {synthetic_injected} synthetic | {rate:.1f}/sec", 
              file=sys.stderr)


def on_message(ws, message):
    """Handle incoming WebSocket messages"""
    global trade_count, synthetic_injected
    
    try:
        data = json.loads(message)
        
        # Skip subscription confirmations
        if 'result' in data or 'id' in data:
            return
            
        # Parse trade data
        if 'e' in data and data['e'] == 'aggTrade':
            trade = {
                "ts": data["E"],
                "price": float(data["p"]),
                "qty": float(data["q"]),
                "id": data["a"],
                "is_buyer_maker": data["m"]
            }
            
            # Output as NDJSON
            sys.stdout.write(json.dumps(trade) + "\n")
            sys.stdout.flush()
            trade_count += 1
            
            # Optional: Inject synthetic anomaly
            synthetic_interval = getattr(on_message, 'synthetic_interval', 0)
            if synthetic_interval > 0 and trade_count > 0:
                if trade_count % synthetic_interval == 0:
                    # Create a synthetic spike (10x volume)
                    synthetic_trade = trade.copy()
                    synthetic_trade["qty"] = synthetic_trade["qty"] * 10
                    synthetic_trade["synthetic"] = True
                    synthetic_trade["id"] = f"SYNTH_{trade_count}"
                    synthetic_trade["ts"] = synthetic_trade["ts"] + 1
                    
                    sys.stdout.write(json.dumps(synthetic_trade) + "\n")
                    sys.stdout.flush()
                    synthetic_injected += 1
                    print(f"🎯 Injected synthetic spike (10x volume)", 
                          file=sys.stderr)

    except Exception as e:
        print(f"Error processing message: {e}", file=sys.stderr)


def on_error(ws, error):
    print(f"WebSocket error: {error}", file=sys.stderr)


def on_close(ws, close_status_code, close_msg):
    print(f"WebSocket closed: {close_status_code} - {close_msg}", file=sys.stderr)


def on_open(ws):
    print("✅ Connected to Binance WebSocket", file=sys.stderr)
    
    # Subscribe to symbol aggregated trades
    subscribe_msg = {
        "method": "SUBSCRIBE",
        "params": [f"{on_open.symbol}@aggTrade"],
        "id": 1
    }
    ws.send(json.dumps(subscribe_msg))


def main():
    parser = argparse.ArgumentParser(description='Stream Binance trades as NDJSON')
    parser.add_argument('--symbol', default='ethusdt', help='Binance symbol (e.g., ethusdt, btcusdt)')
    parser.add_argument('--synthetic-every', type=int, default=0,
                      help='Inject synthetic spike every N trades (0=disabled)')
    args = parser.parse_args()
    
    # Store synthetic interval in the message handler
    on_message.synthetic_interval = args.synthetic_every
    on_open.symbol = args.symbol
    
    websocket_url = f"wss://stream.binance.com:9443/ws/{args.symbol}@aggTrade"
    
    print(f"🚀 Connecting to {websocket_url}", file=sys.stderr)
    print(f"   Symbol: {args.symbol}", file=sys.stderr)
    if args.synthetic_every > 0:
        print(f"🎯 Synthetic spikes: every {args.synthetic_every} trades", 
              file=sys.stderr)
    print(f"💡 Output format: NDJSON (ts, price, qty, id)", file=sys.stderr)
    print("", file=sys.stderr)
    
    # Start stats thread
    stats_thread = threading.Thread(target=print_stats, daemon=True)
    stats_thread.start()
    
    # Connect to WebSocket
    ws = websocket.WebSocketApp(websocket_url,
                               on_open=on_open,
                               on_message=on_message,
                               on_error=on_error,
                               on_close=on_close)
    
    # Run forever
    ws.run_forever()


if __name__ == "__main__":
    main()
