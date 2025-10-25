package stream

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/your-org/driftlock/api-server/internal/models"
)

// Client represents an SSE client connection
type Client struct {
	ID       string
	Stream   chan *models.Anomaly
	Done     chan bool
	LastSeen time.Time
}

// Streamer manages SSE connections and broadcasts
type Streamer struct {
	clients    map[string]*Client
	mu         sync.RWMutex
	register   chan *Client
	unregister chan *Client
	broadcast  chan *models.Anomaly
	maxClients int
}

// NewStreamer creates a new SSE streamer
func NewStreamer(maxClients int) *Streamer {
	if maxClients == 0 {
		maxClients = 1000
	}

	s := &Streamer{
		clients:    make(map[string]*Client),
		register:   make(chan *Client, 10),
		unregister: make(chan *Client, 10),
		broadcast:  make(chan *models.Anomaly, 100),
		maxClients: maxClients,
	}

	// Start background goroutine to manage clients
	go s.run()

	return s
}

// run manages client registration, unregistration, and broadcasts
func (s *Streamer) run() {
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case client := <-s.register:
			s.mu.Lock()
			if len(s.clients) < s.maxClients {
				s.clients[client.ID] = client
				log.Printf("SSE client registered: %s (total: %d)", client.ID, len(s.clients))
			} else {
				log.Printf("SSE client limit reached, rejecting: %s", client.ID)
				close(client.Done)
			}
			s.mu.Unlock()

		case client := <-s.unregister:
			s.mu.Lock()
			if _, ok := s.clients[client.ID]; ok {
				delete(s.clients, client.ID)
				close(client.Stream)
				log.Printf("SSE client unregistered: %s (total: %d)", client.ID, len(s.clients))
			}
			s.mu.Unlock()

		case anomaly := <-s.broadcast:
			s.mu.RLock()
			for _, client := range s.clients {
				select {
				case client.Stream <- anomaly:
					client.LastSeen = time.Now()
				default:
					// Client's buffer is full, skip this message
					log.Printf("SSE client %s buffer full, skipping message", client.ID)
				}
			}
			s.mu.RUnlock()

		case <-ticker.C:
			// Clean up stale connections
			s.mu.Lock()
			now := time.Now()
			for id, client := range s.clients {
				if now.Sub(client.LastSeen) > 5*time.Minute {
					delete(s.clients, id)
					close(client.Stream)
					log.Printf("SSE client timeout: %s", id)
				}
			}
			s.mu.Unlock()
		}
	}
}

// BroadcastAnomaly sends an anomaly to all connected clients
func (s *Streamer) BroadcastAnomaly(anomaly *models.Anomaly) {
	select {
	case s.broadcast <- anomaly:
	default:
		log.Printf("Broadcast channel full, dropping anomaly: %s", anomaly.ID)
	}
}

// GetClientCount returns the number of connected clients
func (s *Streamer) GetClientCount() int {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return len(s.clients)
}

// ServeHTTP handles SSE connections at GET /v1/stream/anomalies
func (s *Streamer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Check if the client supports SSE
	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "SSE not supported", http.StatusInternalServerError)
		return
	}

	// Create new client
	clientID := fmt.Sprintf("client-%d", time.Now().UnixNano())
	client := &Client{
		ID:       clientID,
		Stream:   make(chan *models.Anomaly, 10),
		Done:     make(chan bool),
		LastSeen: time.Now(),
	}

	// Register client
	s.register <- client

	// Set SSE headers
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	// Send initial connection message
	fmt.Fprintf(w, "event: connected\ndata: {\"client_id\": \"%s\"}\n\n", clientID)
	flusher.Flush()

	// Context for handling client disconnection
	ctx := r.Context()

	// Heartbeat ticker
	ticker := time.NewTicker(15 * time.Second)
	defer ticker.Stop()

	// Stream events to client
	for {
		select {
		case <-ctx.Done():
			// Client disconnected
			s.unregister <- client
			return

		case <-client.Done:
			// Client was rejected or closed
			return

		case anomaly := <-client.Stream:
			// Send anomaly event
			data, err := json.Marshal(anomaly)
			if err != nil {
				log.Printf("Failed to marshal anomaly: %v", err)
				continue
			}

			fmt.Fprintf(w, "event: anomaly\ndata: %s\n\n", data)
			flusher.Flush()

		case <-ticker.C:
			// Send heartbeat
			fmt.Fprintf(w, "event: heartbeat\ndata: {\"timestamp\": \"%s\"}\n\n", time.Now().Format(time.RFC3339))
			flusher.Flush()
			client.LastSeen = time.Now()
		}
	}
}

// Shutdown gracefully closes all client connections
func (s *Streamer) Shutdown(ctx context.Context) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	for id, client := range s.clients {
		close(client.Stream)
		delete(s.clients, id)
	}

	log.Printf("SSE streamer shutdown complete")
	return nil
}
