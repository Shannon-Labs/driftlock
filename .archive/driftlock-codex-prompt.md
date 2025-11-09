
# Driftlock Demo Fix and YC Launch Prompt

**Objective:** Get the Driftlock demo operational for the YC application by implementing a basic streaming baseline and anomaly detection capability.

**Background:** The current system is partially implemented. The Go-based OpenTelemetry processor (`collector-processor/driftlockcbad`) calls the core Rust anomaly detection library (`cbad-core`), but it uses a hardcoded, single-entry baseline. This makes the anomaly detection ineffective. We need to modify the processor to dynamically build a baseline from an initial stream of data and then use that baseline to detect anomalies in subsequent data.

**Step-by-Step Instructions:**

**Part 1: Implement Streaming Baseline Logic**

1.  **Modify `collector-processor/driftlockcbad/processor.go`:**
    *   **Add new fields to the `cbadProcessor` struct:**
        *   `baselineSize int`: The number of log entries to collect for the baseline.
        *   `baselineBuffer [][]byte`: A slice to store the raw log data for the baseline.
        *   `isBaselinePeriod bool`: A flag to indicate if the processor is currently in the baseline collection phase.
    *   **Update the `newProcessor` function:**
        *   Initialize `baselineSize` to a reasonable value (e.g., 100).
        *   Initialize `isBaselinePeriod` to `true`.
        *   Initialize `baselineBuffer` as an empty slice.
    *   **Modify the `processLogs` function:**
        *   **Baseline Collection Logic:**
            *   If `isBaselinePeriod` is `true`:
                *   Append the incoming log data (`logData := p.logRecordToBytes(lr)`) to the `baselineBuffer`.
                *   Log a message indicating that the log is being added to the baseline.
                *   If `len(p.baselineBuffer) >= p.baselineSize`:
                    *   Set `isBaselinePeriod` to `false`.
                    *   Concatenate the `baselineBuffer` into the `p.baseline` field.
                    *   Log a message that the baseline collection is complete.
                *   Return immediately, do not perform anomaly detection during the baseline period.
        *   **Anomaly Detection Logic:**
            *   If `isBaselinePeriod` is `false`:
                *   Use the now-populated `p.baseline` for the `ComputeMetricsQuick` call.
                *   The rest of the function can remain as is.

**Part 2: Build the Custom Collector**

1.  **Build the Rust Core Library:**
    *   Run `cargo build --release` in the `cbad-core` directory. This is crucial for the CGO linking to work.
2.  **Build the Custom OpenTelemetry Collector:**
    *   You will need the `opentelemetry-collector-builder` tool. If you don't have it, install it:
        ```bash
        go install go.opentelemetry.io/collector/cmd/builder@latest
        ```
    *   From the root of the `driftlock` project, run the builder using the provided manifest:
        ```bash
        builder --config=./collector-processor/manifest.yaml --output-path=./otelcol-custom
        ```
    *   This will create a custom collector binary at `./otelcol-custom/otelcol-custom`.

**Part 3: Run the System**

1.  **Start Dependent Services:**
    *   The project requires PostgreSQL and Redis. Use `docker-compose` for this. There are `docker-compose.yml` files in the `.archive` directory that can be used as a reference. Create a `docker-compose.yml` in the root of the project with `redis` and `postgres` services.
    *   Run `docker-compose up -d`.
2.  **Configure the Collector:**
    *   Create a `config.yaml` for the custom collector. It should look something like this:
        ```yaml
        receivers:
          otlp:
            protocols:
              grpc:
              http:

        processors:
          driftlockcbad:
            # Add any config options from collector-processor/driftlockcbad/config.go if needed

        exporters:
          logging:
            loglevel: debug

        service:
          pipelines:
            logs:
              receivers: [otlp]
              processors: [driftlockcbad]
              exporters: [logging]
        ```
3.  **Run the Custom Collector:**
    *   `./otelcol-custom/otelcol-custom --config=./config.yaml`

**Part 4: Demonstrate Anomaly Detection**

1.  **Send Baseline Data:**
    *   Use the `generate-test-data.sh` script in the `test-data` directory, or manually send OTLP log data representing normal traffic. Send at least 100 "normal" log entries to build the baseline.
2.  **Send Anomalous Data:**
    *   After the baseline is established, send a few anomalous log entries (e.g., logs with "ERROR" severity, different structures, or unusual attributes).
3.  **Observe the Output:**
    *   The custom collector's console output (from the `logging` exporter) will show the detected anomalies with `driftlock.anomaly_detected="true"` attributes. This output is your demo.

This prompt provides a clear, actionable plan to get a working demo of the Driftlock system.
