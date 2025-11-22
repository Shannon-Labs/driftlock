import pandas as pd
import json
import os

def transform_terra():
    print("Loading Terra/Luna dataset...")
    # We focus on terra-luna.csv as the primary signal
    df = pd.read_csv("test-data/terra_luna/terra-luna.csv")
    
    # Sort by timestamp to simulate a stream
    df['timestamp'] = pd.to_datetime(df['timestamp'])
    df = df.sort_values('timestamp')
    
    output = []
    
    print("Transforming records...")
    for _, row in df.iterrows():
        # Create a transaction-like structure for the demo
        # We map "price" to "amount_usd"
        # Since volume is missing in this CSV, we'll use a random processing_ms + price derivative as proxy for "load"
        # or just simple random noise to satisfy the schema.
        
        import random
        record = {
            "timestamp": row['timestamp'].strftime('%Y-%m-%d %H:%M:%S'),
            "transaction_id": f"tx_{int(row['timestamp'].timestamp())}",
            "amount_usd": float(row['price']), # Actual price
            "processing_ms": int(20 + random.random() * 50), # Simulated load
            "origin_country": "Blockchain",
            "api_endpoint": "LUNA-USD",
            "status": "processed"
        }
        output.append(record)
    
    # Limit to 2000 points around the crash for the visual demo
    # The crash was May 2022.
    # Let's take a slice if it's too big, or just dump it all if reasonable.
    if len(output) > 5000:
       output = output[-5000:] # Take the last 5000 which usually includes the crash in these datasets

    print(f"Saving {len(output)} records to test-data/terra_luna/driftlock_ready.json...")
    with open("test-data/terra_luna/driftlock_ready.json", "w") as f:
        json.dump(output, f, indent=2)
    
    print("Done.")

if __name__ == "__main__":
    transform_terra()
