import pandas as pd
import json
import os

def transform_web():
    print("Loading Web Traffic (NAB) dataset...")
    # Use a specific file known for anomalies: 'art_daily_jumpsup.csv' or realAWSCloudwatch
    # Let's use 'ec2_cpu_utilization_5f5533.csv' (real AWS data)
    path = "test-data/web_traffic/realAWSCloudwatch/realAWSCloudwatch/ec2_cpu_utilization_5f5533.csv"
    if not os.path.exists(path):
        print(f"File not found: {path}")
        return

    df = pd.read_csv(path)
    
    output = []
    print("Transforming records...")
    # Limit to 2000
    df = df.head(2000)
    
    for i, row in df.iterrows():
        # timestamp, value
        # value -> amount_usd (load)
        # processing_ms -> random noise around a baseline, or correlate with value?
        # Let's say higher CPU = slower processing
        
        val = float(row['value'])
        proc_ms = int(20 + val * 2) # scaling factor
        
        record = {
            "timestamp": row['timestamp'],
            "transaction_id": f"cpu_{i}",
            "amount_usd": val,
            "processing_ms": proc_ms,
            "origin_country": "us-east-1",
            "api_endpoint": "i-5f5533",
            "status": "high_load" if val > 60 else "nominal"
        }
        output.append(record)
    
    print(f"Saving {len(output)} records to test-data/web_traffic/driftlock_ready.json...")
    with open("test-data/web_traffic/driftlock_ready.json", "w") as f:
        json.dump(output, f, indent=2)
    
    print("Done.")

if __name__ == "__main__":
    transform_web()
