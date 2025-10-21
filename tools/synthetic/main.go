package main

import (
    "bytes"
    "encoding/json"
    "fmt"
    "math/rand"
    "net/http"
    "os"
    "time"
)

func main() {
    api := getenv("API_URL", "http://localhost:8080")
    n := 100
    if len(os.Args) > 1 {
        fmt.Sscanf(os.Args[1], "%d", &n)
    }
    rand.Seed(42)
    for i := 0; i < n; i++ {
        payload := map[string]any{
            "timestamp": time.Now().UTC().Format(time.RFC3339Nano),
            "value":     rand.NormFloat64(),
            "seq":       i,
        }
        b, _ := json.Marshal(payload)
        _, _ = http.Post(api+"/v1/events", "application/json", bytes.NewReader(b))
        time.Sleep(10 * time.Millisecond)
    }
    fmt.Println("sent", n, "events to", api)
}

func getenv(k, d string) string { if v := os.Getenv(k); v != "" { return v }; return d }

