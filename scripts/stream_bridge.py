#!/usr/bin/env python3
import json
import sys
import time
import requests
import logging

# Configure logging to stderr so stdout is kept clean for data piping
logging.basicConfig(stream=sys.stderr, level=logging.INFO, format='%(asctime)s - %(levelname)s - %(message)s')
logger = logging.getLogger(__name__)

URL = "https://stream.wikimedia.org/v2/stream/recentchange"

def stream_events():
    """
    Connects to Wikimedia EventStreams and yields event data dictionaries.
    Handles reconnection with backoff.
    """
    backoff = 1
    headers = {
        "User-Agent": "DriftlockSoakTest/1.0 (https://github.com/factory/driftlock; contact@driftlock.net)"
    }
    while True:
        try:
            logger.info(f"Connecting to {URL}...")
            with requests.get(URL, headers=headers, stream=True, timeout=30) as response:
                response.raise_for_status()
                client = response.iter_lines(chunk_size=None)
                
                for line in client:
                    if line:
                        decoded_line = line.decode('utf-8')
                        
                        # SSE format: "data: {json payload}"
                        if decoded_line.startswith("data: "):
                            try:
                                data_str = decoded_line[6:]
                                data = json.loads(data_str)
                                yield data
                                # Reset backoff on successful data
                                backoff = 1
                            except json.JSONDecodeError:
                                logger.warning("Failed to decode JSON line")
                                continue
        except Exception as e:
            logger.error(f"Connection error: {e}")
            logger.info(f"Retrying in {backoff} seconds...")
            time.sleep(backoff)
            backoff = min(backoff * 2, 60) # Cap backoff at 60 seconds

def normalize_event(event):
    """
    Extracts relevant fields and flattens the structure for Driftlock.
    """
    # We want a clean flat JSON object
    try:
        meta = event.get('meta', {})
        normalized = {
            "timestamp": meta.get('dt') or event.get('timestamp'),
            "id": meta.get('id') or event.get('id'),
            "type": event.get('type'),
            "user": event.get('user'),
            "bot": event.get('bot'),
            "title": event.get('title'),
            "comment": event.get('comment'),
            "server_url": event.get('server_url'),
            "wiki": event.get('wiki'),
            "length_new": event.get('length', {}).get('new'),
            "length_old": event.get('length', {}).get('old'),
            "message": f"Edit by {event.get('user')} on {event.get('title')}: {event.get('comment')}" 
        }
        
        # Ensure we have a message field for Driftlock to analyze if it defaults to message
        if not normalized.get('comment'):
             normalized['message'] = f"Edit by {event.get('user')} on {event.get('title')}"
             
        return normalized
    except Exception as e:
        logger.warning(f"Error normalizing event: {e}")
        return None

def main():
    try:
        for event in stream_events():
            normalized = normalize_event(event)
            if normalized:
                # Write NDJSON to stdout
                print(json.dumps(normalized), flush=True)
    except KeyboardInterrupt:
        logger.info("Bridge stopped by user")

if __name__ == "__main__":
    main()
