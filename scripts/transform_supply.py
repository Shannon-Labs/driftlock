import pandas as pd
import json
import random

def transform_supply():
    print("Loading Supply Chain dataset...")
    df = pd.read_csv("test-data/supply_chain/dynamic_supply_chain_logistics_dataset_with_country.csv")
    
    output = []
    
    print("Transforming records...")
    # Limit to 2000
    df = df.head(2000)
    
    for i, row in df.iterrows():
        # Map fields
        # delivery_time_deviation -> processing_ms (proxy for delay)
        # shipping_costs -> amount_usd
        
        record = {
            "timestamp": f"Shipment_{i}",
            "transaction_id": f"ship_{row['product_id']}_{i}",
            "amount_usd": float(row['shipping_costs']),
            "processing_ms": int(row['delivery_time_deviation'] * 100) if 'delivery_time_deviation' in row else 100,
            "origin_country": row['supplier_country'],
            "api_endpoint": f"Route_Risk_{row['risk_classification']}",
            "status": row['risk_classification']
        }
        # Clean up potential NaNs
        if pd.isna(record["amount_usd"]): record["amount_usd"] = 0.0
        if pd.isna(record["processing_ms"]): record["processing_ms"] = 100
        
        output.append(record)
    
    print(f"Saving {len(output)} records to test-data/supply_chain/driftlock_ready.json...")
    with open("test-data/supply_chain/driftlock_ready.json", "w") as f:
        json.dump(output, f, indent=2)
    
    print("Done.")

if __name__ == "__main__":
    transform_supply()
