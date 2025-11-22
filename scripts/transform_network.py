import pandas as pd
import json
import random

def transform_network():
    print("Loading Network Intrusion dataset...")
    # Requires pyarrow
    try:
        df = pd.read_parquet("test-data/network/UNSW_NB15_testing-set.parquet")
    except Exception as e:
        print(f"Parquet read failed ({e}), skipping network test for now.")
        return

    output = []
    
    print("Transforming records...")
    # Limit to 2000 for demo
    df = df.head(2000)
    
    for i, row in df.iterrows():
        # Map fields
        # dur -> processing_ms
        # sbytes + dbytes -> amount_usd (proxy for data volume)
        
        record = {
            "timestamp": f"Tick {i}",
            "transaction_id": f"net_{i}",
            "amount_usd": float(row['sbytes'] + row['dbytes']),
            "processing_ms": int(row['dur'] * 1000) if row['dur'] > 0 else 1,
            "origin_country": row['proto'],
            "api_endpoint": row['service'],
            "status": "attack" if row['label'] == 1 else "normal"
        }
        output.append(record)
    
    print(f"Saving {len(output)} records to test-data/network/driftlock_ready.json...")
    with open("test-data/network/driftlock_ready.json", "w") as f:
        json.dump(output, f, indent=2)
    
    print("Done.")

if __name__ == "__main__":
    transform_network()
