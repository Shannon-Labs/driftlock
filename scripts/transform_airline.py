import pandas as pd
import json
import random

def transform_airline():
    print("Loading Airline dataset...")
    try:
        df = pd.read_csv("test-data/airline/Airline_Delay_Cause.csv")
    except Exception as e:
        print(f"Read failed ({e}), skipping airline test.")
        return

    output = []
    print("Transforming records...")
    # Limit to 2000 for demo speed
    df = df.head(2000)
    
    for i, row in df.iterrows():
        # year, month, carrier, carrier_name, airport, ... arr_delay
        # Map:
        # arr_delay -> processing_ms (latency)
        # arr_flights -> amount_usd (volume)
        
        # Filter out negative delays (early arrivals) for simpler demo math
        delay = row['arr_delay']
        if pd.isna(delay): delay = 0
        if delay < 0: delay = 0
        
        record = {
            "timestamp": f"{int(row['year'])}-{int(row['month']):02d}",
            "transaction_id": f"flt_{row['carrier']}_{row['airport']}_{i}",
            "amount_usd": float(row['arr_flights']) if not pd.isna(row['arr_flights']) else 0.0,
            "processing_ms": int(delay),
            "origin_country": row['airport_name'] if isinstance(row['airport_name'], str) else "Unknown",
            "api_endpoint": row['carrier_name'] if isinstance(row['carrier_name'], str) else "Unknown",
            "status": "delayed" if delay > 15 else "on_time"
        }
        output.append(record)
    
    print(f"Saving {len(output)} records to test-data/airline/driftlock_ready.json...")
    with open("test-data/airline/driftlock_ready.json", "w") as f:
        json.dump(output, f, indent=2)
    
    print("Done.")

if __name__ == "__main__":
    transform_airline()
