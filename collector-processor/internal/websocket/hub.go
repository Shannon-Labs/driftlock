package websocket

import (
	"log"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

// MessageType represents different types of WebSocket messages
type MessageType string

const (
	MessageTypeProgress MessageType = "progress"
	MessageTypeComplete MessageType = "complete"
	MessageTypeError    MessageType = "error"
	MessageTypeStatus   MessageType = "status"
	MessageTypePing     MessageType = "ping"
)

// WebSocketMessage represents a message sent over WebSocket
type WebSocketMessage struct {
	Type      MessageType `json:"type"`
	Timestamp time.Time   `json:"timestamp"`
	Data      interface{} `json:"data"`
	RequestID string      `json:"request_id,omitempty"`
}

// ProgressData represents progress information
type ProgressData struct {
	ProcessedLines int     `json:"processed_lines"`
	TotalLines     int     `json:"total_lines"`
	Percentage     float64 `json:"percentage"`
	CurrentLine    string  `json:"current_line,omitempty"`
	Speed          int     `json:"speed,omitempty"` // lines per second
	ETA            string  `json:"eta,omitempty"`
}

// CompleteData represents completion information
type CompleteData struct {
	TotalEvents    int                    `json:"total_events"`
	AnomalyCount   int                    `json:"anomaly_count"`
	ProcessingTime string                 `json:"processing_time"`
	FileInfo       map[string]interface{} `json:"file_info"`
	Anomalies      []interface{}          `json:"anomalies,omitempty"`
}

// ErrorData represents error information
type ErrorData struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Details string `json:"details,omitempty"`
}

// StatusData represents status information
type StatusData struct {
	ActiveConnections int               `json:"active_connections"`
	ProcessingJobs    map[string]string `json:"processing_jobs"`
	ServerUptime      string            `json:"server_uptime"`
}

// Client represents a WebSocket client connection
type Client struct {
	ID        string
	Conn      *websocket.Conn
	Send      chan WebSocketMessage
	RequestID string
	Hub       *Hub
}

// Hub maintains the set of active clients and broadcasts messages to them
type Hub struct {
	// Registered clients
	clients map[*Client]bool

	// Inbound messages from the clients
	broadcast chan WebSocketMessage

	// Register requests from the clients
	register chan *Client

	// Unregister requests from clients
	unregister chan *Client

	// Mutex for thread-safe operations
	mutex sync.RWMutex

	// Start time for uptime calculation
	startTime time.Time

	// Active processing jobs by request ID
	processingJobs map[string]string
}

// NewHub creates a new WebSocket hub
func NewHub() *Hub {
	return &Hub{
		broadcast:      make(chan WebSocketMessage, 256),
		register:       make(chan *Client),
		unregister:     make(chan *Client),
		clients:        make(map[*Client]bool),
		startTime:      time.Now(),
		processingJobs: make(map[string]string),
	}
}

// Run starts the hub
func (h *Hub) Run() {
	log.Printf("WebSocket hub started at %s", h.startTime)

	for {
		select {
		case client := <-h.register:
			h.mutex.Lock()
			h.clients[client] = true
			h.mutex.Unlock()
			log.Printf("WebSocket client connected: %s (total: %d)", client.ID, len(h.clients))

			// Send welcome message
			h.sendToClient(client, WebSocketMessage{
				Type:      MessageTypeStatus,
				Timestamp: time.Now(),
				Data: StatusData{
					ActiveConnections: len(h.clients),
					ProcessingJobs:    h.getProcessingJobs(),
					ServerUptime:      h.getUptime(),
				},
			})

		case client := <-h.unregister:
			h.mutex.Lock()
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.Send)
				log.Printf("WebSocket client disconnected: %s (total: %d)", client.ID, len(h.clients))
			}
			h.mutex.Unlock()

		case message := <-h.broadcast:
			h.mutex.RLock()
			for client := range h.clients {
				// Send to all clients or specific request ID
				if message.RequestID == "" || message.RequestID == client.RequestID {
					select {
					case client.Send <- message:
					default:
						// Client channel is full, remove client
						close(client.Send)
						delete(h.clients, client)
					}
				}
			}
			h.mutex.RUnlock()
		}
	}
}

// SendProgress sends a progress update to clients
func (h *Hub) SendProgress(requestID string, data ProgressData) {
	message := WebSocketMessage{
		Type:      MessageTypeProgress,
		Timestamp: time.Now(),
		RequestID: requestID,
		Data:      data,
	}

	select {
	case h.broadcast <- message:
	default:
		log.Printf("WebSocket broadcast channel is full, dropping progress message")
	}
}

// SendComplete sends a completion message to clients
func (h *Hub) SendComplete(requestID string, data CompleteData) {
	message := WebSocketMessage{
		Type:      MessageTypeComplete,
		Timestamp: time.Now(),
		RequestID: requestID,
		Data:      data,
	}

	select {
	case h.broadcast <- message:
	default:
		log.Printf("WebSocket broadcast channel is full, dropping complete message")
	}

	// Remove from processing jobs
	h.mutex.Lock()
	delete(h.processingJobs, requestID)
	h.mutex.Unlock()
}

// SendError sends an error message to clients
func (h *Hub) SendError(requestID string, data ErrorData) {
	message := WebSocketMessage{
		Type:      MessageTypeError,
		Timestamp: time.Now(),
		RequestID: requestID,
		Data:      data,
	}

	select {
	case h.broadcast <- message:
	default:
		log.Printf("WebSocket broadcast channel is full, dropping error message")
	}

	// Remove from processing jobs
	h.mutex.Lock()
	delete(h.processingJobs, requestID)
	h.mutex.Unlock()
}

// RegisterProcessingJob registers a new processing job
func (h *Hub) RegisterProcessingJob(requestID, jobType string) {
	h.mutex.Lock()
	h.processingJobs[requestID] = jobType
	h.mutex.Unlock()
}

// GetProcessingJobs returns current processing jobs
func (h *Hub) GetProcessingJobs() map[string]string {
	h.mutex.RLock()
	defer h.mutex.RUnlock()

	jobs := make(map[string]string)
	for k, v := range h.processingJobs {
		jobs[k] = v
	}
	return jobs
}

// GetStats returns hub statistics
func (h *Hub) GetStats() map[string]interface{} {
	h.mutex.RLock()
	defer h.mutex.RUnlock()

	return map[string]interface{}{
		"active_connections": len(h.clients),
		"processing_jobs":    len(h.processingJobs),
		"uptime":             h.getUptime(),
		"started_at":         h.startTime,
	}
}

// sendToClient sends a message to a specific client
func (h *Hub) sendToClient(client *Client, message WebSocketMessage) {
	select {
	case client.Send <- message:
	default:
		log.Printf("Client channel is full, removing client %s", client.ID)
		close(client.Send)
		delete(h.clients, client)
	}
}

// getProcessingJobs returns a copy of processing jobs
func (h *Hub) getProcessingJobs() map[string]string {
	jobs := make(map[string]string)
	for k, v := range h.processingJobs {
		jobs[k] = v
	}
	return jobs
}

// getUptime returns the server uptime as a formatted string
func (h *Hub) getUptime() string {
	uptime := time.Since(h.startTime)
	return uptime.String()
}

// BroadcastStatus sends status updates to all connected clients
func (h *Hub) BroadcastStatus() {
	h.mutex.RLock()
	status := StatusData{
		ActiveConnections: len(h.clients),
		ProcessingJobs:    h.getProcessingJobs(),
		ServerUptime:      h.getUptime(),
	}
	h.mutex.RUnlock()

	message := WebSocketMessage{
		Type:      MessageTypeStatus,
		Timestamp: time.Now(),
		Data:      status,
	}

	select {
	case h.broadcast <- message:
	default:
		log.Printf("WebSocket broadcast channel is full, dropping status message")
	}
}

// StartPingLoop starts a ping loop to keep connections alive
func (h *Hub) StartPingLoop(interval time.Duration) {
	ticker := time.NewTicker(interval)
	go func() {
		for range ticker.C {
			h.mutex.RLock()
			for client := range h.clients {
				select {
				case client.Send <- WebSocketMessage{
					Type:      MessageTypePing,
					Timestamp: time.Now(),
					Data:      map[string]interface{}{"ping": "pong"},
				}:
				default:
					// Client channel is full, remove client
					close(client.Send)
					delete(h.clients, client)
					h.mutex.RUnlock()
					h.mutex.Lock()
					break
				}
			}
			if len(h.clients) > 0 {
				h.mutex.RUnlock()
			} else {
				h.mutex.RUnlock()
			}
		}
	}()
}
