#!/usr/bin/env python3
"""
Driftlock Full Stress Test with yfinance
Tests anomaly detection accuracy on real + synthetic financial data

Usage:
    pip install yfinance requests
    python3 scripts/test_yfinance_stress.py
"""
import yfinance as yf
import requests
import json
import random
import copy
from datetime import datetime, timedelta
from typing import List, Dict, Tuple
import time
import sys

API_URL = "https://driftlock-api-o6kjgrsowq-uc.a.run.app"
DEMO_ENDPOINT = f"{API_URL}/v1/demo/detect"

# Test assets - mix of stocks and crypto
TICKERS = ["AAPL", "TSLA", "NVDA", "SPY", "BTC-USD"]

# ANSI colors for terminal output
class Colors:
    GREEN = '\033[92m'
    YELLOW = '\033[93m'
    RED = '\033[91m'
    CYAN = '\033[96m'
    BOLD = '\033[1m'
    END = '\033[0m'


class StressTest:
    def __init__(self, verbose=True):
        self.verbose = verbose
        self.results = {
            "total_events": 0,
            "total_anomalies": 0,
            "injected_anomalies": 0,
            "detected_injected": 0,
            "processing_times": [],
            "ncd_scores": [],
            "errors": [],
            "scenarios": []
        }

    def log(self, msg: str, color: str = ""):
        """Print with optional color"""
        if self.verbose:
            if color:
                print(f"{color}{msg}{Colors.END}")
            else:
                print(msg)

    def fetch_data(self, ticker: str, period="5d", interval="1m") -> List[Dict]:
        """Fetch real market data from yfinance"""
        try:
            t = yf.Ticker(ticker)
            hist = t.history(period=period, interval=interval)
            events = []
            for idx, row in hist.iterrows():
                # Handle NaN values
                if any(pd_isna(row[col]) for col in ['Open', 'High', 'Low', 'Close', 'Volume']):
                    continue
                events.append({
                    "timestamp": idx.isoformat(),
                    "symbol": ticker,
                    "open": float(row['Open']),
                    "high": float(row['High']),
                    "low": float(row['Low']),
                    "close": float(row['Close']),
                    "volume": int(row['Volume'])
                })
            return events
        except Exception as e:
            self.log(f"Error fetching {ticker}: {e}", Colors.RED)
            self.results["errors"].append(f"Fetch {ticker}: {e}")
            return []

    def inject_anomalies(self, events: List[Dict], count=3) -> Tuple[List[Dict], List[int]]:
        """
        Inject synthetic anomalies into a copy of events.
        Returns (modified_events, list of indices where anomalies were injected)
        """
        if len(events) < 30:
            return events, []

        # Work on a copy
        events = copy.deepcopy(events)

        # Pick random indices (avoiding edges)
        max_count = min(count, len(events) // 20)
        if max_count < 1:
            return events, []

        indices = random.sample(range(10, len(events) - 10), max_count)

        for idx in indices:
            anomaly_type = random.choice(["flash_crash", "volume_spike", "gap"])

            if anomaly_type == "flash_crash":
                # Sudden 20% price drop
                events[idx]["close"] *= 0.8
                events[idx]["low"] *= 0.75
                events[idx]["_anomaly_type"] = "flash_crash"

            elif anomaly_type == "volume_spike":
                # 10x volume explosion
                events[idx]["volume"] *= 10
                events[idx]["_anomaly_type"] = "volume_spike"

            else:  # gap
                # 10% gap up
                events[idx]["open"] *= 1.1
                events[idx]["high"] *= 1.15
                events[idx]["_anomaly_type"] = "gap"

            events[idx]["_injected"] = True

        return events, indices

    def send_batch(self, events: List[Dict], batch_start: int = 0) -> Dict:
        """Send a batch of up to 50 events to the demo API"""
        # Remove internal tracking fields before sending
        clean_events = []
        for e in events[:50]:
            clean_e = {k: v for k, v in e.items() if not k.startswith('_')}
            clean_events.append(clean_e)

        start_time = time.time()
        try:
            resp = requests.post(
                DEMO_ENDPOINT,
                json={"events": clean_events},
                headers={"Content-Type": "application/json"},
                timeout=30
            )
            elapsed = time.time() - start_time
            self.results["processing_times"].append(elapsed)

            if resp.status_code == 200:
                return resp.json()
            elif resp.status_code == 429:
                self.log("Rate limited, waiting 60s...", Colors.YELLOW)
                time.sleep(60)
                return self.send_batch(events, batch_start)  # Retry
            else:
                self.log(f"API error: {resp.status_code} - {resp.text[:200]}", Colors.RED)
                return {"success": False, "error": resp.text}

        except Exception as e:
            self.log(f"Request error: {e}", Colors.RED)
            self.results["errors"].append(f"Request: {e}")
            return {"success": False, "error": str(e)}

    def run_scenario(self, name: str, events: List[Dict], inject: bool = False) -> Dict:
        """Run a single test scenario"""
        self.log(f"\n{'='*60}", Colors.CYAN)
        self.log(f"Scenario: {name}", Colors.BOLD)
        self.log(f"Events: {len(events)}")

        scenario_result = {
            "name": name,
            "total_events": 0,
            "anomalies_detected": 0,
            "injected_count": 0,
            "detected_injected": 0
        }

        injected_indices = []
        if inject:
            events, injected_indices = self.inject_anomalies(events, count=5)
            scenario_result["injected_count"] = len(injected_indices)
            self.results["injected_anomalies"] += len(injected_indices)
            if injected_indices:
                self.log(f"Injected {len(injected_indices)} anomalies at indices: {injected_indices[:5]}...", Colors.YELLOW)

        # Process in batches of 50 (demo API limit)
        total_anomalies = 0
        batch_num = 0

        for i in range(0, len(events), 50):
            batch = events[i:i+50]
            batch_num += 1

            result = self.send_batch(batch, i)

            if result.get("success"):
                batch_anomalies = result.get("anomaly_count", 0)
                total_anomalies += batch_anomalies
                scenario_result["total_events"] += len(batch)
                self.results["total_events"] += len(batch)

                # Check each detected anomaly
                for a in result.get("anomalies", []):
                    # Track NCD scores
                    ncd = a.get("metrics", {}).get("ncd", 0)
                    if ncd > 0:
                        self.results["ncd_scores"].append(ncd)

                    # Check if this was an injected anomaly
                    anomaly_idx = a.get("index", -1) + i  # Adjust for batch offset
                    if inject and anomaly_idx in injected_indices:
                        scenario_result["detected_injected"] += 1
                        self.results["detected_injected"] += 1
                        self.log(f"  ✓ Detected injected anomaly at index {anomaly_idx}", Colors.GREEN)

                if batch_anomalies > 0:
                    self.log(f"  Batch {batch_num}: {batch_anomalies} anomalies", Colors.YELLOW)
            else:
                self.log(f"  Batch {batch_num}: Failed", Colors.RED)

            # Small delay between batches to avoid rate limiting
            if i + 50 < len(events):
                time.sleep(0.5)

        scenario_result["anomalies_detected"] = total_anomalies
        self.results["total_anomalies"] += total_anomalies
        self.results["scenarios"].append(scenario_result)

        # Summary for this scenario
        self.log(f"\nScenario Results:", Colors.BOLD)
        self.log(f"  Events processed: {scenario_result['total_events']}")
        self.log(f"  Anomalies detected: {total_anomalies}")
        if inject and scenario_result["injected_count"] > 0:
            rate = scenario_result["detected_injected"] / scenario_result["injected_count"] * 100
            color = Colors.GREEN if rate >= 50 else Colors.YELLOW if rate >= 25 else Colors.RED
            self.log(f"  Detection rate: {rate:.1f}% ({scenario_result['detected_injected']}/{scenario_result['injected_count']})", color)

        return scenario_result

    def run_all(self):
        """Run full stress test suite"""
        self.log("\n" + "="*60, Colors.CYAN)
        self.log("  DRIFTLOCK STRESS TEST", Colors.BOLD)
        self.log("  Testing CBAD on real financial data + synthetic anomalies", Colors.CYAN)
        self.log("="*60 + "\n", Colors.CYAN)

        # Test each ticker
        for ticker in TICKERS:
            self.log(f"\n{'─'*60}")
            self.log(f"Fetching {ticker} data from yfinance...", Colors.CYAN)

            events = self.fetch_data(ticker)

            if not events:
                self.log(f"No data for {ticker}, skipping", Colors.YELLOW)
                continue

            self.log(f"Got {len(events)} data points for {ticker}")

            # Test 1: Normal data (no injection)
            self.run_scenario(f"{ticker} - Normal Market Data", events.copy(), inject=False)

            # Wait to avoid rate limiting
            time.sleep(2)

            # Test 2: Data with injected anomalies
            self.run_scenario(f"{ticker} - With Injected Anomalies", events.copy(), inject=True)

            # Wait between tickers
            time.sleep(3)

        self.print_summary()

    def print_summary(self):
        """Print comprehensive test results summary"""
        self.log("\n" + "="*60, Colors.CYAN)
        self.log("  STRESS TEST RESULTS SUMMARY", Colors.BOLD)
        self.log("="*60, Colors.CYAN)

        # Overall stats
        self.log(f"\n{Colors.BOLD}Overall Statistics:{Colors.END}")
        self.log(f"  Total events processed: {self.results['total_events']:,}")
        self.log(f"  Total anomalies detected: {self.results['total_anomalies']:,}")
        self.log(f"  Total injected anomalies: {self.results['injected_anomalies']}")
        self.log(f"  Correctly detected injections: {self.results['detected_injected']}")

        # Detection accuracy
        if self.results['injected_anomalies'] > 0:
            rate = self.results['detected_injected'] / self.results['injected_anomalies'] * 100
            color = Colors.GREEN if rate >= 50 else Colors.YELLOW if rate >= 25 else Colors.RED
            self.log(f"\n{Colors.BOLD}Detection Accuracy:{Colors.END}")
            self.log(f"  Injected anomaly detection rate: {rate:.1f}%", color)

        # Performance stats
        if self.results['processing_times']:
            avg_time = sum(self.results['processing_times']) / len(self.results['processing_times'])
            max_time = max(self.results['processing_times'])
            min_time = min(self.results['processing_times'])
            self.log(f"\n{Colors.BOLD}Performance:{Colors.END}")
            self.log(f"  Average processing time: {avg_time*1000:.0f}ms")
            self.log(f"  Min/Max: {min_time*1000:.0f}ms / {max_time*1000:.0f}ms")

        # NCD score distribution
        if self.results['ncd_scores']:
            avg_ncd = sum(self.results['ncd_scores']) / len(self.results['ncd_scores'])
            max_ncd = max(self.results['ncd_scores'])
            self.log(f"\n{Colors.BOLD}NCD Scores (Anomaly Metric):{Colors.END}")
            self.log(f"  Average NCD: {avg_ncd:.4f}")
            self.log(f"  Max NCD: {max_ncd:.4f}")

        # Errors
        if self.results['errors']:
            self.log(f"\n{Colors.BOLD}Errors:{Colors.END}", Colors.RED)
            for err in self.results['errors'][:5]:
                self.log(f"  - {err}", Colors.RED)

        # Per-scenario breakdown
        self.log(f"\n{Colors.BOLD}Per-Scenario Breakdown:{Colors.END}")
        for scenario in self.results['scenarios']:
            inj_info = ""
            if scenario['injected_count'] > 0:
                rate = scenario['detected_injected'] / scenario['injected_count'] * 100
                inj_info = f" | Detection: {rate:.0f}%"
            self.log(f"  {scenario['name']}: {scenario['anomalies_detected']} anomalies{inj_info}")

        # Final verdict
        self.log("\n" + "="*60, Colors.CYAN)
        if self.results['total_events'] > 0 and len(self.results['errors']) == 0:
            if self.results['injected_anomalies'] > 0:
                rate = self.results['detected_injected'] / self.results['injected_anomalies'] * 100
                if rate >= 30:
                    self.log("  ✓ STRESS TEST PASSED", Colors.GREEN)
                    self.log("  System is detecting anomalies in financial data", Colors.GREEN)
                else:
                    self.log("  ⚠ DETECTION RATE LOW", Colors.YELLOW)
                    self.log("  Consider tuning CBAD parameters", Colors.YELLOW)
            else:
                self.log("  ✓ API FUNCTIONAL", Colors.GREEN)
                self.log("  No injection test due to data constraints", Colors.YELLOW)
        else:
            self.log("  ✗ STRESS TEST HAD ISSUES", Colors.RED)
            self.log(f"  Errors: {len(self.results['errors'])}", Colors.RED)
        self.log("="*60 + "\n", Colors.CYAN)


def pd_isna(value):
    """Check if value is NaN (works without importing pandas)"""
    try:
        import math
        return math.isnan(value)
    except (TypeError, ValueError):
        return False


def main():
    """Main entry point"""
    print("\n" + "="*60)
    print("  Driftlock Stress Test - yfinance Integration")
    print("="*60)
    print(f"  API URL: {API_URL}")
    print(f"  Tickers: {', '.join(TICKERS)}")
    print("="*60 + "\n")

    # Check dependencies
    try:
        import yfinance
        import requests
    except ImportError as e:
        print(f"Missing dependency: {e}")
        print("Run: pip install yfinance requests")
        sys.exit(1)

    # Run the test
    test = StressTest(verbose=True)
    test.run_all()

    # Return exit code based on results
    if test.results['errors']:
        sys.exit(1)
    sys.exit(0)


if __name__ == "__main__":
    main()
