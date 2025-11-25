#!/usr/bin/env python3
"""
Simple NDJSON streamer for crypto data (REST-based to avoid WebSocket restrictions)
Uses CoinGecko for data and outputs NDJSON format
"""

import sys
import json
import time
import requests
import argparse
from datetime import datetime

# Default coins - high volume ones
COINS = ["bitcoin", "ethereum", "solana", "dogecoin", "litecoin"]


def stream_coingecko(interval=5, synthetic_every=0):
    """Stream crypto data from CoinGecko"""
    trade_id = 0
    synthetic_count = 0
    
    print(f"🚀 Streaming crypto data from CoinGecko", file=sys.stderr)
    print(f"⏱️  Interval: {interval}s", file=sys.stderr)
    if synthetic_every > 0:
        print(f"🎯 Synthetic spikes: every {synthetic_every} batches", file=sys.stderr)
    print(f"💡 Coins: {', '.join(COINS)}", file=sys.stderr)
    print("", file=sys.stderr)
    
    while True:
        try:
            # Fetch current prices
            url = f"https://api.coingecko.com/api/v3/simple/price"
            params = {
                "ids": ",".join(COINS),
                "vs_currencies": "usd",
                "include_24hr_vol": "true",
                "include_24hr_change": "true",
                "precision": "8"
            }
            
            response = requests.get(url, params=params, timeout=10)
            
            if response.status_code == 200:
                data = response.json()
                
                # Output each coin as separate NDJSON record
                batch_size = 0
                for coin in COINS:
                    if coin in data:
                        coin_data = data[coin]
                        trade_id += 1
                        record = {
                            "ts": int(time.time() * 1000),
                            "coin": coin,
                            "price": coin_data["usd"],
                            "volume": coin_data.get("usd_24h_vol", 0),
                            "change_24h": coin_data.get("usd_24h_change", 0),
                            "id": trade_id
                        }
                        
                        sys.stdout.write(json.dumps(record) + "\n")
                        batch_size += 1
                
                sys.stdout.flush()
                
                # Inject synthetic anomaly if configured
                if synthetic_every > 0 and batch_size > 0:
                    if (trade_id // batch_size) % synthetic_every == 0:
                        synthetic_count += 1
                        # Create a synthetic spike (10x price change)
                        first_coin = list(data.keys())[0] if data else "bitcoin"
                        if first_coin in data:
                            synthetic_record = {
                                "ts": int(time.time() * 1000) + 1,
                                "coin": first_coin,
                                "price": data[first_coin]["usd"] * 0.1,  # 90% drop
                                "volume": data[first_coin].get("usd_24h_vol", 0),
                                "change_24h": -50.0,
                                "id": f"SYNTH_{trade_id}_{synthetic_count}",
                                "synthetic": True
                            }
                            
                            sys.stdout.write(json.dumps(synthetic_record) + "\n")
                            sys.stdout.flush()
                            print(f"🎯 Injected synthetic spike #{synthetic_count}:",
                                  f"{first_coin} price drop to {synthetic_record['price']:.2f}",
                                  file=sys.stderr)
                
                print(f"📤 Batch: {batch_size} records | Total: {trade_id}", 
                      file=sys.stderr)
                
            elif response.status_code == 429:
                print(f"⚠️  Rate limited, cooling down...", file=sys.stderr)
                time.sleep(10)
            else:
                print(f"⚠️  Error {response.status_code}: {response.text[:100]}",
                      file=sys.stderr)
                time.sleep(5)
                
        except Exception as e:
            print(f"⚠️  Error: {e}", file=sys.stderr)
            time.sleep(5)
        
        time.sleep(interval)


def main():
    parser = argparse.ArgumentParser(description='Stream crypto data as NDJSON')
    parser.add_argument('--interval', type=int, default=5,
                      help='Polling interval in seconds (default: 5)')
    parser.add_argument('--synthetic-every', type=int, default=0,
                      help='Inject synthetic anomaly every N batches')
    args = parser.parse_args()
    
    try:
        stream_coingecko(args.interval, args.synthetic_every)
    except KeyboardInterrupt:
        print(f"\n✅ Total trades: {trade_id}", file=sys.stderr)
        print(f"✅ Synthetic spikes: {synthetic_count}", file=sys.stderr)


if __name__ == "__main__":
    trade_id = 0
    synthetic_count = 0
    main()