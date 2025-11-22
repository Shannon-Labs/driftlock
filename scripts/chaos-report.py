#!/usr/bin/env python3
"""Summarize Driftlock benchmark logs for THE_CHAOS_REPORT."""

from __future__ import annotations

import argparse
import json
import re
from pathlib import Path
from typing import Any


ANOMALY_RE = re.compile(r"Found (\d+) anomalies out of (\d+)")
TIME_RE = re.compile(r"Processing time:\s*([0-9.]+)s")
BASELINE_RE = re.compile(
    r"Baseline ready! \(processing_ms median=(\d+)ms, p95=(\d+)ms; amount median=\$([0-9.]+)\)"
)


def parse_log(path: Path) -> dict[str, Any]:
    text = path.read_text(encoding="utf-8", errors="ignore")
    anomaly_match = ANOMALY_RE.search(text)
    time_match = TIME_RE.search(text)
    baseline_match = BASELINE_RE.search(text)

    return {
        "dataset": path.stem.replace("_results", ""),
        "file": str(path),
        "anomalies": int(anomaly_match.group(1)) if anomaly_match else None,
        "total": int(anomaly_match.group(2)) if anomaly_match else None,
        "processing_seconds": float(time_match.group(1)) if time_match else None,
        "baseline_median_ms": int(baseline_match.group(1)) if baseline_match else None,
        "baseline_p95_ms": int(baseline_match.group(2)) if baseline_match else None,
        "amount_median": float(baseline_match.group(3)) if baseline_match else None,
    }


def format_markdown(rows: list[dict[str, Any]]) -> str:
    header = "| Dataset | Anomalies | Total | % Drift | Baseline median (ms) | p95 (ms) | Runtime (s) |"
    line = "| --- | ---: | ---: | ---: | ---: | ---: | ---: |"
    body = []
    for row in rows:
        total = row.get("total") or 0
        anomalies = row.get("anomalies") or 0
        drift_pct = f"{(anomalies / total * 100):.2f}%" if total else "—"
        runtime_val = row.get("processing_seconds")
        runtime_str = f"{runtime_val:.3f}" if runtime_val is not None else "—"
        median = row.get("baseline_median_ms")
        median_str = str(median) if median is not None else "—"
        p95 = row.get("baseline_p95_ms")
        p95_str = str(p95) if p95 is not None else "—"
        body.append(
            "| {dataset} | {anomalies} | {total} | {drift} | {median} | {p95} | {runtime} |".format(
                dataset=row.get("dataset", "?"),
                anomalies=anomalies,
                total=total or "?",
                drift=drift_pct,
                median=median_str,
                p95=p95_str,
                runtime=runtime_str,
            )
        )
    return "\n".join([header, line, *body])


def main() -> None:
    parser = argparse.ArgumentParser(description=__doc__)
    parser.add_argument(
        "logs",
        nargs="*",
        default=[
            "airline_results.log",
            "network_results.log",
            "safety_results.log",
            "web_results.log",
            "supply_results.log",
            "terra_results.log",
            "nasa_results.log",
        ],
        help="Paths to *_results.log files",
    )
    parser.add_argument("--format", choices=["markdown", "json"], default="markdown")
    args = parser.parse_args()

    rows: list[dict[str, Any]] = []
    for entry in args.logs:
        path = Path(entry)
        if not path.exists():
            continue
        rows.append(parse_log(path))

    rows.sort(key=lambda row: row.get("dataset", ""))

    if args.format == "json":
        print(json.dumps(rows, indent=2))
    else:
        print(format_markdown(rows))


if __name__ == "__main__":
    main()
