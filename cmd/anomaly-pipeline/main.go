package main

import (
	"context"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/gorilla/websocket"
	"github.com/segmentio/kafka-go"
)

type pipelineConfig struct {
	Brokers        []string
	EventsTopic    string
	AnomalyTopic   string
	GroupID        string
	ClientID       string
	HTTPAddr       string
	LatencyTrigger float64
}

func main() {
	cfg := loadPipelineConfig()

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	anomalies := make(chan json.RawMessage, 128)
	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()
		runDetector(ctx, cfg, anomalies)
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		runWeb(ctx, cfg, anomalies)
	}()

	<-ctx.Done()
	close(anomalies)
	wg.Wait()
}

func loadPipelineConfig() pipelineConfig {
	var (
		brokers        = flag.String("brokers", getEnv("KAFKA_BROKERS", "localhost:9092"), "Kafka broker list (comma separated)")
		eventsTopic    = flag.String("events-topic", getEnv("KAFKA_EVENTS_TOPIC", "otlp-events"), "Kafka topic to consume OTLP events from")
		anomalyTopic   = flag.String("anomaly-topic", getEnv("KAFKA_ANOMALY_TOPIC", "anomaly-events"), "Kafka topic to publish anomalies to")
		groupID        = flag.String("group-id", getEnv("KAFKA_GROUP_ID", "anomaly-detector"), "Kafka consumer group ID")
		clientID       = flag.String("client-id", getEnv("KAFKA_CLIENT_ID", "anomaly-pipeline"), "Kafka client ID")
		httpAddr       = flag.String("http", getEnv("ANOMALY_HTTP_ADDR", ":8090"), "HTTP address for anomaly dashboard")
		latencyTrigger = flag.Float64("latency-threshold", getEnvFloat("ANOMALY_LATENCY_THRESHOLD", 1000.0), "Latency threshold (ms) for metrics anomalies")
	)
	flag.Parse()

	return pipelineConfig{
		Brokers:        parseBrokers(*brokers),
		EventsTopic:    *eventsTopic,
		AnomalyTopic:   *anomalyTopic,
		GroupID:        *groupID,
		ClientID:       *clientID,
		HTTPAddr:       *httpAddr,
		LatencyTrigger: *latencyTrigger,
	}
}

func parseBrokers(raw string) []string {
	fields := strings.Split(raw, ",")
	out := make([]string, 0, len(fields))
	for _, f := range fields {
		if trimmed := strings.TrimSpace(f); trimmed != "" {
			out = append(out, trimmed)
		}
	}
	return out
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func getEnvFloat(key string, fallback float64) float64 {
	if v := os.Getenv(key); v != "" {
		if parsed, err := strconv.ParseFloat(v, 64); err == nil {
			return parsed
		}
	}
	return fallback
}

func runDetector(ctx context.Context, cfg pipelineConfig, anomalies chan<- json.RawMessage) {
	if len(cfg.Brokers) == 0 {
		log.Println("detector: no brokers configured; skipping")
		return
	}

	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers: cfg.Brokers,
		GroupID: cfg.GroupID,
		Topic:   cfg.EventsTopic,
		Logger:  kafka.LoggerFunc(func(string, ...interface{}) {}),
	})
	defer reader.Close()

	writer := kafka.NewWriter(kafka.WriterConfig{
		Brokers: cfg.Brokers,
		Topic:   cfg.AnomalyTopic,
		Async:   false,
	})
	defer writer.Close()

	log.Printf("detector: consuming %s, publishing anomalies to %s", cfg.EventsTopic, cfg.AnomalyTopic)

	for {
		msg, err := reader.ReadMessage(ctx)
		if err != nil {
			if errors.Is(err, context.Canceled) {
				return
			}
			log.Printf("detector: read error: %v", err)
			continue
		}

		payload, detected, err := detectAnomaly(msg.Value, cfg.LatencyTrigger)
		if err != nil {
			log.Printf("detector: decode error: %v", err)
			continue
		}
		if !detected {
			continue
		}

		if err := writer.WriteMessages(ctx, kafka.Message{Value: payload}); err != nil {
			log.Printf("detector: failed to publish anomaly: %v", err)
		}

		select {
		case anomalies <- payload:
		default:
			log.Println("detector: anomaly channel full, dropping in-memory broadcast")
		}
	}
}

func detectAnomaly(data []byte, latencyThreshold float64) (json.RawMessage, bool, error) {
	var event map[string]interface{}
	if err := json.Unmarshal(data, &event); err != nil {
		return nil, false, err
	}

	rawType, _ := event["type"].(string)
	dataMap, _ := event["data"].(map[string]interface{})
	if dataMap == nil {
		return nil, false, nil
	}

	var reasons []string

	switch strings.ToLower(rawType) {
	case "log":
		severity := strings.ToUpper(asString(dataMap["severity"]))
		body := strings.ToLower(asString(dataMap["body"]))

		if severity == "ERROR" || severity == "FATAL" || severity == "PANIC" {
			reasons = append(reasons, fmt.Sprintf("log severity %s", severity))
		}
		if strings.Contains(body, "panic") || strings.Contains(body, "out of memory") {
			reasons = append(reasons, "log contains panic signature")
		}
	case "metric":
		if hist, ok := dataMap["histogram"].([]interface{}); ok {
			for _, item := range hist {
				bucket, ok := item.(map[string]interface{})
				if !ok {
					continue
				}
				if sum, ok := asFloat(bucket["sum"]); ok && sum > latencyThreshold {
					reasons = append(reasons, fmt.Sprintf("histogram sum %.2fms", sum))
					break
				}
				if max, ok := asFloat(bucket["max"]); ok && max > latencyThreshold {
					reasons = append(reasons, fmt.Sprintf("histogram max %.2fms", max))
					break
				}
			}
		}
	}

	if len(reasons) == 0 {
		return nil, false, nil
	}

	anomaly := map[string]interface{}{
		"detected_at": time.Now().UTC().Format(time.RFC3339Nano),
		"reasons":     reasons,
		"event":       event,
	}

	payload, err := json.Marshal(anomaly)
	if err != nil {
		return nil, false, err
	}

	return json.RawMessage(payload), true, nil
}

func asString(v interface{}) string {
	if s, ok := v.(string); ok {
		return s
	}
	return ""
}

func asFloat(v interface{}) (float64, bool) {
	switch value := v.(type) {
	case float64:
		return value, true
	case float32:
		return float64(value), true
	case json.Number:
		f, err := value.Float64()
		return f, err == nil
	default:
		return 0, false
	}
}

var wsUpgrader = websocket.Upgrader{CheckOrigin: func(*http.Request) bool { return true }}

func runWeb(ctx context.Context, cfg pipelineConfig, anomalies <-chan json.RawMessage) {
	hub := newWebHub()
	go hub.process(ctx, anomalies)

	mux := http.NewServeMux()
	mux.HandleFunc("/", serveDashboard)
	mux.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		hub.handleWebsocket(ctx, w, r)
	})

	server := &http.Server{Addr: cfg.HTTPAddr, Handler: mux}

	go func() {
		<-ctx.Done()
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		_ = server.Shutdown(shutdownCtx)
	}()

	log.Printf("web: anomaly dashboard listening on %s", cfg.HTTPAddr)
	if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Printf("web: server error: %v", err)
	}
}

type webHub struct {
	mu      sync.Mutex
	clients map[*websocket.Conn]struct{}
	history []json.RawMessage
}

func newWebHub() *webHub {
	return &webHub{
		clients: make(map[*websocket.Conn]struct{}),
		history: make([]json.RawMessage, 0, 64),
	}
}

func (h *webHub) process(ctx context.Context, anomalies <-chan json.RawMessage) {
	for {
		select {
		case msg, ok := <-anomalies:
			if !ok {
				h.closeAll()
				return
			}
			h.appendHistory(msg)
			h.broadcast(msg)
		case <-ctx.Done():
			h.closeAll()
			return
		}
	}
}

func (h *webHub) appendHistory(msg json.RawMessage) {
	h.mu.Lock()
	defer h.mu.Unlock()

	h.history = append(h.history, msg)
	if len(h.history) > 100 {
		h.history = append([]json.RawMessage(nil), h.history[len(h.history)-100:]...)
	}
}

func (h *webHub) broadcast(msg json.RawMessage) {
	h.mu.Lock()
	defer h.mu.Unlock()

	for conn := range h.clients {
		if err := conn.WriteMessage(websocket.TextMessage, msg); err != nil {
			conn.Close()
			delete(h.clients, conn)
		}
	}
}

func (h *webHub) closeAll() {
	h.mu.Lock()
	defer h.mu.Unlock()

	for conn := range h.clients {
		_ = conn.WriteControl(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, "server shutdown"), time.Now().Add(2*time.Second))
		conn.Close()
	}
	h.clients = make(map[*websocket.Conn]struct{})
}

func (h *webHub) handleWebsocket(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	conn, err := wsUpgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("websocket: upgrade failed: %v", err)
		return
	}

	conn.SetReadLimit(1024)
	conn.SetReadDeadline(time.Now().Add(60 * time.Second))
	conn.SetPongHandler(func(string) error {
		conn.SetReadDeadline(time.Now().Add(60 * time.Second))
		return nil
	})

	h.mu.Lock()
	for _, msg := range h.history {
		if err := conn.WriteMessage(websocket.TextMessage, msg); err != nil {
			h.mu.Unlock()
			conn.Close()
			return
		}
	}
	h.clients[conn] = struct{}{}
	h.mu.Unlock()

	go func() {
		pingTicker := time.NewTicker(30 * time.Second)
		defer pingTicker.Stop()
		for {
			select {
			case <-ctx.Done():
				return
			case <-pingTicker.C:
				conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
				if err := conn.WriteMessage(websocket.PingMessage, nil); err != nil {
					return
				}
			}
		}
	}()

	go func() {
		defer func() {
			h.mu.Lock()
			delete(h.clients, conn)
			h.mu.Unlock()
			conn.Close()
		}()
		for {
			if _, _, err := conn.ReadMessage(); err != nil {
				return
			}
		}
	}()
}

const dashboardHTML = `<!doctype html>
<html lang="en">
  <head>
    <meta charset="utf-8" />
    <meta name="viewport" content="width=device-width, initial-scale=1" />
    <title>DriftLock Anomalies</title>
    <style>
      body { margin: 0; font-family: system-ui, -apple-system, BlinkMacSystemFont, "Segoe UI", sans-serif; background: #0b0f17; color: #e6eefb; }
      header { padding: 16px 24px; border-bottom: 1px solid #243046; }
      main { padding: 24px; display: grid; gap: 16px; }
      .card { background: #101726; border: 1px solid #1e2a3d; border-radius: 12px; padding: 16px 20px; box-shadow: 0 4px 18px rgba(0,0,0,0.45); }
      .meta { color: #8ca0c8; font-size: 12px; margin-bottom: 8px; }
      pre { background: #121b2d; padding: 12px; border-radius: 8px; overflow-x: auto; font-size: 12px; }
      .reason { display: inline-block; margin-right: 12px; padding: 4px 8px; border-radius: 6px; background: rgba(255,99,132,0.2); color: #ff9eb5; font-size: 12px; }
      #status { font-size: 12px; color: #8ca0c8; }
    </style>
  </head>
  <body>
    <header>
      <h1 style="margin:0;font-size:20px;">DriftLock Anomaly Stream</h1>
      <p id="status">Connecting…</p>
    </header>
    <main id="anomalies"></main>
    <script>
      const target = document.getElementById('anomalies');
      const status = document.getElementById('status');
      function render(item) {
        const card = document.createElement('div');
        card.className = 'card';
        const when = new Date(item.detected_at || Date.now()).toLocaleString();
        const reasons = (item.reasons || []).map(r => '<span class="reason">' + r + '</span>').join('');
        card.innerHTML = '<div class="meta">' + when + '</div>' + reasons + '<pre>' + JSON.stringify(item.event, null, 2) + '</pre>';
        target.prepend(card);
        while (target.children.length > 50) target.removeChild(target.lastChild);
      }
      function connect() {
        const protocol = location.protocol === 'https:' ? 'wss' : 'ws';
        const ws = new WebSocket(protocol + '://' + location.host + '/ws');
        ws.onopen = () => { status.textContent = 'Live'; };
        ws.onclose = () => { status.textContent = 'Disconnected, retrying…'; setTimeout(connect, 2000); };
        ws.onerror = () => { status.textContent = 'Error'; ws.close(); };
        ws.onmessage = (evt) => {
          try { const data = JSON.parse(evt.data); render(data); }
          catch (_) { console.error('invalid anomaly payload', evt.data); }
        };
      }
      connect();
    </script>
  </body>
</html>`

func serveDashboard(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	_, _ = w.Write([]byte(dashboardHTML))
}
