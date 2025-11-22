import pandas as pd
import json
import random

def transform_safety():
    print("Loading AI Safety dataset...")
    # malignant.csv has 'prompt' and 'label' (0=safe, 1=malignant)
    try:
        df = pd.read_csv("test-data/ai_safety/malignant.csv")
    except Exception as e:
        print(f"Read failed: {e}")
        return

    output = []
    print("Transforming records...")
    # Take a mix of safe and malignant to test detection
    # The dataset might be huge, take 1000 of each if possible, or just head(2000)
    df = df.head(2000)
    
    for i, row in df.iterrows():
        # Map:
        # text -> encoded into 'origin_country' or 'api_endpoint' to test string compression?
        # Actually, Driftlock demo uses JSON stringification of the whole record.
        # So putting the prompt in any field works.
        # Let's put it in 'api_endpoint' as if it's the requested resource path/query.
        
        prompt = str(row['text'])
        # Truncate insane lengths for the demo json
        if len(prompt) > 500: prompt = prompt[:500] + "..."
        
        # Determine status based on category/base_class if label is missing
        # malignant.csv typically has 'category' or 'base_class' for non-safe
        status = "flagged"
        # If we have a way to know safe vs unsafe, we use it. 
        # Looking at the file, it seems all are malignant?
        # The dataset is "Prompt Injection Malignant". So all are bad.
        
        record = {
            "timestamp": f"Req_{i}",
            "transaction_id": f"llm_{i}",
            "amount_usd": 0.0, # Free tier?
            "processing_ms": int(len(prompt) / 2), # Longer prompts take longer
            "origin_country": "User_Input",
            "api_endpoint": prompt, # The payload to check!
            "status": status
        }
        output.append(record)
    
    print(f"Saving {len(output)} records to test-data/ai_safety/driftlock_ready.json...")
    with open("test-data/ai_safety/driftlock_ready.json", "w") as f:
        json.dump(output, f, indent=2)
    
    print("Done.")

if __name__ == "__main__":
    transform_safety()
