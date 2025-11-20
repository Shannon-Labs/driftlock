# IoT Sensor Monitoring

Detect equipment failure and predictive maintenance needs using Driftlock.

## Scenario

You manage a fleet of industrial machines (e.g., wind turbines, manufacturing robots). Each machine sends telemetry data (temperature, vibration, RPM) every second. You want to detect when a machine starts behaving abnormally *before* it fails.

## Implementation

### 1. Data Structure

We will track sensor readings.

```json
{
  "machine_id": "turbine-001",
  "timestamp": "2025-01-01T12:00:00Z",
  "metrics": {
    "temperature_c": 45.2,
    "vibration_hz": 12.5,
    "rpm": 1200,
    "power_output_kw": 350
  },
  "status": "running"
}
```

### 2. Python Client (Edge Device)

Run this on your edge gateway or cloud ingestion service.

```python
import time
import random
from driftlock import DriftlockClient

client = DriftlockClient(api_key="...")

def read_sensors():
    # Simulate reading from hardware
    return {
        "temperature_c": 45 + random.uniform(-1, 1),
        "vibration_hz": 12 + random.uniform(-0.5, 0.5),
        "rpm": 1200 + random.randint(-10, 10)
    }

async def monitor_machine(machine_id):
    while True:
        data = read_sensors()
        
        # Detect anomalies
        result = await client.detect(
            stream_id=f"machine-{machine_id}",
            events=[{
                "type": "telemetry",
                "body": data
            }]
        )
        
        if result.anomalies:
            anomaly = result.anomalies[0]
            print(f"[WARNING] Machine {machine_id} anomaly: {anomaly.why}")
            
            # Trigger emergency shutdown if confidence is very high
            if anomaly.metrics.confidence > 0.99:
                shutdown_machine(machine_id)
                
        time.sleep(1)
```

## Use Cases

### Predictive Maintenance

A bearing failure often starts with a subtle increase in vibration and temperature.
- **Traditional Threshold**: "Alert if Temp > 80°C". By the time it hits 80°C, the machine might already be damaged.
- **Driftlock**: Detects that Temp is 48°C when it's normally 45°C, combined with a slight shift in vibration frequency. This "multivariate" anomaly is detected days before catastrophic failure.

### Mode Detection

Machines have different operating modes (Idle, Startup, Full Load).
- Driftlock learns these modes automatically.
- It won't flag "Low RPM" as an anomaly during startup if that's the normal pattern for startup.
- It *will* flag "High RPM" during Idle.

## Bandwidth Optimization

IoT devices often have limited bandwidth.

1. **Batching**: Send data every minute instead of every second.
2. **Compression**: Driftlock supports GZIP.
3. **Edge Processing**: If you have a powerful edge gateway, you can run a local anomaly detection model (contact Enterprise sales for Edge/On-Prem deployment).
