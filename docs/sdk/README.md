# SDK Reference

Driftlock provides official SDKs for high-velocity integration. These libraries handle batching, backpressure, and authentication automatically.

## Node.js

```bash
npm install @driftlock/sdk
```

```typescript
import { Driftlock } from '@driftlock/sdk';

const radar = new Driftlock('your_api_key');

// Stream events
radar.push({
  id: 'evt_001',
  user: 'alice',
  action: 'login'
});

// Listen for verdicts
radar.on('anomaly', (verdict) => {
  console.error('Security Event:', verdict.explanation);
  // Trigger automated lockdown
});
```

## Python

```bash
pip install driftlock
```

```python
from driftlock import Radar

radar = Radar("your_api_key")

# Check a single batch synchronously
verdict = radar.scan(events)

if verdict.is_anomaly:
    print(f"Drift detected: {verdict.explanation}")
```

## Go

```bash
go get github.com/shannon-labs/driftlock-go
```

```go
import "github.com/shannon-labs/driftlock-go"

func main() {
    client := driftlock.NewClient("your_api_key")
    
    verdict, err := client.Scan(context.Background(), events)
    if err != nil {
        log.Fatal(err)
    }
    
    if verdict.IsAnomaly {
        log.Printf("Entropy Spike: %v", verdict.Variance)
    }
}
```

---

*© 2025 Shannon Labs. Built for scale.*
