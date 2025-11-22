import pandas as pd
import json
import random

def transform_nasa():
    print("Loading NASA Turbofan dataset (FD001)...")
    # train_FD001.txt: C-MAPSS data
    # Columns: unit, time_cycles, setting1, setting2, setting3, s1...s21
    cols = ["unit", "time", "os1", "os2", "os3"] + [f"s{i}" for i in range(1, 22)]
    df = pd.read_csv("test-data/nasa_turbofan/CMaps/train_FD001.txt", sep=r"\s+", header=None, names=cols)
    
    # We'll focus on Unit 1's run-to-failure for a clear story
    df_unit1 = df[df['unit'] == 1]
    
    output = []
    
    print("Transforming records...")
    for _, row in df_unit1.iterrows():
        # Map sensor readings to the demo schema
        # amount_usd -> Sensor 4 (Fan Inlet Temp? High variance usually)
        # processing_ms -> Sensor 9 (Core speed?)
        
        record = {
            "timestamp": f"Cycle {int(row['time'])}",
            "transaction_id": f"u1_c{int(row['time'])}",
            "amount_usd": float(row['s4']), # Proxy for a sensor value
            "processing_ms": int(row['s9']), # Proxy for another sensor
            "origin_country": f"Engine_{int(row['unit'])}",
            "api_endpoint": "Sensor_Telemetry",
            "status": "nominal"
        }
        output.append(record)
    
    print(f"Saving {len(output)} records to test-data/nasa_turbofan/driftlock_ready.json...")
    with open("test-data/nasa_turbofan/driftlock_ready.json", "w") as f:
        json.dump(output, f, indent=2)
    
    print("Done.")

if __name__ == "__main__":
    transform_nasa()
