#!/usr/bin/env python3
import subprocess
import json
import time
import sys
import os
import requests
import threading
import queue
from datetime import datetime

# Configuration
BRIDGE_SCRIPT = os.path.join(os.path.dirname(__file__), "crypto_bridge.py")
DRIFTLOCK_BIN = os.path.join(os.path.dirname(os.path.dirname(__file__)), "bin", "driftlock")
LOG_DIR = os.path.join(os.path.dirname(os.path.dirname(__file__)), "logs")
STREAM_LOG = os.path.join(LOG_DIR, "live-crypto.ndjson")
GEMINI_LOG = os.path.join(LOG_DIR, "live-gemini.ndjson")
FIREBASE_URL = "https://us-central1-driftlock.cloudfunctions.net/analyzeAnomalies"

# Batching settings
BATCH_SIZE = 5
BATCH_TIMEOUT = 60  # seconds

# Global queue for anomalies to be sent
anomaly_queue = queue.Queue()

def post_anomalies_worker():
    """
    Worker thread that consumes anomalies from the queue and posts them to Firebase.
    """
    batch = []
    last_send_time = time.time()
    
    print(f"Worker started. Targeting {FIREBASE_URL}")

    while True:
        try:
            # Wait for items, with timeout to ensure we send partial batches
            try:
                item = anomaly_queue.get(timeout=1)
                batch.append(item)
            except queue.Empty:
                pass

            current_time = time.time()
            time_diff = current_time - last_send_time

            # Send if batch is full or timeout reached (and batch is not empty)
            if (len(batch) >= BATCH_SIZE) or (batch and time_diff >= BATCH_TIMEOUT):
                send_batch(batch)
                batch = []
                last_send_time = current_time
                
        except Exception as e:
            print(f"Error in worker: {e}", file=sys.stderr)
            time.sleep(5)

def send_batch(batch):
    """
    Sends a batch of anomalies to the Firebase function.
    """
    payload = {
        "query": "Live stream soak test – Crypto Volatility (Binance)",
        "anomalies": batch
    }
    
    print(f"Sending batch of {len(batch)} anomalies...")
    
    try:
        response = requests.post(FIREBASE_URL, json=payload, timeout=30)
        response.raise_for_status()
        
        result = response.json()
        
        # Log the interaction
        log_entry = {
            "timestamp": datetime.utcnow().isoformat(),
            "batch_size": len(batch),
            "response": result
        }
        
        with open(GEMINI_LOG, "a") as f:
            f.write(json.dumps(log_entry) + "\n")
            
        print("Batch sent successfully.")
        
    except Exception as e:
        print(f"Failed to send batch: {e}", file=sys.stderr)
        # Optionally log failure
        with open(GEMINI_LOG, "a") as f:
            f.write(json.dumps({"timestamp": datetime.utcnow().isoformat(), "error": str(e)}) + "\n")

def main():
    # Check binaries
    if not os.path.exists(DRIFTLOCK_BIN):
        print(f"Error: driftlock binary not found at {DRIFTLOCK_BIN}")
        sys.exit(1)

    print(f"Starting Soak Test Runner...")
    print(f"Bridge: {BRIDGE_SCRIPT}")
    print(f"Scanner: {DRIFTLOCK_BIN}")
    print(f"Logs: {LOG_DIR}")

    # Start the worker thread
    worker = threading.Thread(target=post_anomalies_worker, daemon=True)
    worker.start()

    # Start the bridge process
    bridge_proc = subprocess.Popen(
        [sys.executable, BRIDGE_SCRIPT],
        stdout=subprocess.PIPE,
        stderr=sys.stderr, # Passthrough logs
        text=True
    )

    # Start the scan process
    # We use --show-all to capture everything, --stdin to read from bridge
    scan_cmd = [
        DRIFTLOCK_BIN, "scan",
        "--stdin",
        "--format", "ndjson",
        "--output", "ndjson",
        "--show-all",
        "--threshold", "0.35" # Can adjust threshold
    ]
    
    scan_proc = subprocess.Popen(
        scan_cmd,
        stdin=bridge_proc.stdout,
        stdout=subprocess.PIPE,
        stderr=sys.stderr,
        text=True
    )

    # Close our handle to bridge stdout so it can close if scan_proc closes? 
    # Actually we just let them pipe.
    bridge_proc.stdout.close()

    try:
        with open(STREAM_LOG, "a") as log_file:
            for line in scan_proc.stdout:
                # Write raw output to file
                log_file.write(line)
                log_file.flush()
                
                try:
                    record = json.loads(line)
                    # Check for anomaly
                    if record.get("anomaly"):
                        print(f"Anomaly detected! Score: {record.get('anomaly_score')}")
                        anomaly_queue.put(record)
                except json.JSONDecodeError:
                    pass
                    
    except KeyboardInterrupt:
        print("Stopping soak test...")
    finally:
        scan_proc.terminate()
        bridge_proc.terminate()
        scan_proc.wait()
        bridge_proc.wait()
        print("Processes terminated.")

if __name__ == "__main__":
    main()
