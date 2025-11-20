# Go SDK

The official Go client for Driftlock. High-performance, type-safe anomaly detection for your Go services.

## Installation

Install the package using `go get`:

```bash
go get github.com/shannon-labs/driftlock-go
```

## Quick Start

Initialize the client and start detecting anomalies.

```go
package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/shannon-labs/driftlock-go"
)

func main() {
	// Initialize the client
	client := driftlock.NewClient(os.Getenv("DRIFTLOCK_API_KEY"))

	// Create an event
	event := driftlock.Event{
		Timestamp: time.Now(),
		Type:      "system_metric",
		Body: map[string]interface{}{
			"cpu_usage": 85.5,
			"memory":    1024,
		},
	}

	// Detect anomalies
	resp, err := client.Detect(context.Background(), &driftlock.DetectRequest{
		StreamID: "server-metrics",
		Events:   []driftlock.Event{event},
	})
	if err != nil {
		log.Fatalf("Error detecting anomalies: %v", err)
	}

	if len(resp.Anomalies) > 0 {
		fmt.Printf("Detected %d anomalies:\n", len(resp.Anomalies))
		for _, anomaly := range resp.Anomalies {
			fmt.Printf("- %s (Confidence: %.2f)\n", anomaly.Why, anomaly.Metrics.Confidence)
		}
	} else {
		fmt.Println("No anomalies detected.")
	}
}
```

## Configuration

You can configure the client using `ClientOption` functions.

```go
client := driftlock.NewClient(
	"your-api-key",
	driftlock.WithBaseURL("https://driftlock-api-o6kjgrsowq-uc.a.run.app"),
	driftlock.WithTimeout(5 * time.Second),
	driftlock.WithRetries(3),
)
```

## Struct-Based Events

You can use your own structs as event bodies. The SDK will automatically serialize them to JSON.

```go
type Transaction struct {
	Amount   float64 `json:"amount"`
	Currency string  `json:"currency"`
	UserID   string  `json:"user_id"`
}

event := driftlock.Event{
	Timestamp: time.Now(),
	Type:      "transaction",
	Body: Transaction{
		Amount:   150.00,
		Currency: "USD",
		UserID:   "user_123",
	},
}
```

## Error Handling

The SDK returns typed errors for common API issues.

```go
resp, err := client.Detect(ctx, req)
if err != nil {
	if rateLimitErr, ok := err.(*driftlock.RateLimitError); ok {
		fmt.Printf("Rate limit exceeded. Retry after %v\n", rateLimitErr.RetryAfter)
		time.Sleep(rateLimitErr.RetryAfter)
		// Retry logic...
	} else {
		log.Printf("API Error: %v", err)
	}
}
```

## Middleware Integration

The Go SDK includes middleware for common frameworks like Gin and Echo.

### Gin Middleware

```go
import "github.com/shannon-labs/driftlock-go/middleware/gin"

r := gin.Default()
r.Use(gin.DriftlockMiddleware(client))
```

### Echo Middleware

```go
import "github.com/shannon-labs/driftlock-go/middleware/echo"

e := echo.New()
e.Use(echo.DriftlockMiddleware(client))
```

## Support

For issues, please [open an issue on GitHub](https://github.com/Shannon-Labs/driftlock-go) or contact [support@driftlock.io](mailto:support@driftlock.io).
