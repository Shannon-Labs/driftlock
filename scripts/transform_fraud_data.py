import pandas as pd
import json
import random
import time

def transform():
    print("Loading fraud dataset...")
    df = pd.read_csv("test-data/fraud/fraud_data.csv")
    
    # Limit to 5000 rows for the demo run to be fast but meaningful
    df = df.head(5000)
    
    output = []
    
    print("Transforming records...")
    for _, row in df.iterrows():
        # Simulate processing time based on amount (larger amount = slightly longer check) + noise
        proc_ms = int(20 + (row['amt'] / 100.0) + random.random() * 50)
        
        record = {
            "timestamp": row['trans_date_trans_time'],
            "transaction_id": row['trans_num'],
            "amount_usd": float(row['amt']),
            "processing_ms": proc_ms,
            "origin_country": f"{row['city']}, {row['state']}",
            "api_endpoint": row['merchant'].replace('"', '').strip(), # Clean up quotes
            "status": "approved" if row['is_fraud'] == 0 else "review_flagged"
        }
        output.append(record)
    
    # Save as JSON array
    print(f"Saving {len(output)} records to test-data/fraud/driftlock_ready.json...")
    with open("test-data/fraud/driftlock_ready.json", "w") as f:
        json.dump(output, f, indent=2)
    
    print("Done.")

if __name__ == "__main__":
    transform()
