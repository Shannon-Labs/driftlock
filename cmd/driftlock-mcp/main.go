package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/driftlock/driftlock/pkg/entropywindow"
)

// MCP Protocol Types (Simplified for MVP)
type JSONRPCRequest struct {
	JSONRPC string          `json:"jsonrpc"`
	Method  string          `json:"method"`
	Params  json.RawMessage `json:"params"`
	ID      interface{}     `json:"id"`
}

type JSONRPCResponse struct {
	JSONRPC string      `json:"jsonrpc"`
	Result  interface{} `json:"result,omitempty"`
	Error   *RPCError   `json:"error,omitempty"`
	ID      interface{} `json:"id"`
}

type RPCError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type Tool struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	InputSchema any    `json:"inputSchema"`
}

type CallToolParams struct {
	Name      string          `json:"name"`
	Arguments json.RawMessage `json:"arguments"`
}

const (
	API_URL = "https://driftlock-api-o6kjgrsowq-uc.a.run.app/v1"
)

// Tools Definition
var tools = []Tool{
	{
		Name:        "detect_anomalies",
		Description: "Analyze raw strings or JSON payloads for entropy anomalies using Driftlock.",
		InputSchema: map[string]any{
			"type": "object",
			"properties": map[string]any{
				"data": map[string]string{
					"type":        "string",
					"description": "Raw log / metric data to analyze.",
				},
				"mode": map[string]any{
					"type":        "string",
					"enum":        []string{"auto", "raw", "json"},
					"description": "auto = detect raw vs JSON, raw = entropy-only, json = send to API",
				},
				"format": map[string]any{
					"type":        "string",
					"enum":        []string{"json", "ndjson"},
					"description": "Hint for remote API requests",
				},
				"windowLines": map[string]any{
					"type":        "integer",
					"description": "Sliding baseline size for local entropy analysis",
				},
				"threshold": map[string]any{
					"type":        "number",
					"description": "Anomaly score threshold (0-1) for local analysis",
				},
			},
			"required": []string{"data"},
		},
	},
	{
		Name:        "check_system_health",
		Description: "Check the health status of the Driftlock API.",
		InputSchema: map[string]any{
			"type":       "object",
			"properties": map[string]any{},
		},
	},
}

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		line := scanner.Bytes()
		var req JSONRPCRequest
		if err := json.Unmarshal(line, &req); err != nil {
			continue // invalid json, ignore or log
		}

		var resp JSONRPCResponse
		resp.JSONRPC = "2.0"
		resp.ID = req.ID

		switch req.Method {
		case "tools/list":
			resp.Result = map[string]any{
				"tools": tools,
			}
		case "tools/call":
			var params CallToolParams
			if err := json.Unmarshal(req.Params, &params); err != nil {
				resp.Error = &RPCError{Code: -32602, Message: "Invalid params"}
				break
			}
			result, err := handleToolCall(params)
			if err != nil {
				resp.Error = &RPCError{Code: -32000, Message: err.Error()}
			} else {
				resp.Result = map[string]any{
					"content": []map[string]string{
						{"type": "text", "text": result},
					},
					"isError": false,
				}
			}
		default:
			// For notifications or unhandled methods
			continue
		}

		out, _ := json.Marshal(resp)
		fmt.Println(string(out))
	}
}

func handleToolCall(params CallToolParams) (string, error) {
	switch params.Name {
	case "check_system_health":
		resp, err := http.Get(API_URL + "/healthz")
		if err != nil {
			return "", err
		}
		defer resp.Body.Close()
		body, _ := io.ReadAll(resp.Body)
		return string(body), nil

	case "detect_anomalies":
		var args detectArgs
		if err := json.Unmarshal(params.Arguments, &args); err != nil {
			return "", fmt.Errorf("invalid arguments")
		}
		args.applyDefaults()
		if args.shouldAnalyzeLocally() {
			return runLocalAnalyzer(args)
		}
		return callRemoteAPI(args)
	}
	return "", fmt.Errorf("unknown tool")
}

type detectArgs struct {
	Data        string  `json:"data"`
	Mode        string  `json:"mode"`
	Format      string  `json:"format"`
	WindowLines int     `json:"windowLines"`
	Threshold   float64 `json:"threshold"`
}

func (d *detectArgs) applyDefaults() {
	if d.Mode == "" {
		d.Mode = "auto"
	}
	if d.WindowLines <= 0 {
		d.WindowLines = 400
	}
	if d.Threshold <= 0 {
		d.Threshold = 0.35
	}
	if d.Format == "" {
		d.Format = "json"
	}
}

func (d detectArgs) shouldAnalyzeLocally() bool {
	switch strings.ToLower(d.Mode) {
	case "raw":
		return true
	case "json":
		return false
	default:
	}
	trimmed := strings.TrimSpace(d.Data)
	if trimmed == "" {
		return true
	}
	if strings.HasPrefix(trimmed, "{") || strings.HasPrefix(trimmed, "[") {
		return false
	}
	return true
}

func runLocalAnalyzer(args detectArgs) (string, error) {
	analyzer, err := entropywindow.NewAnalyzer(entropywindow.Config{
		BaselineLines:        args.WindowLines,
		Threshold:            args.Threshold,
		CompressionAlgorithm: "zstd",
		MinLineLength:        8,
	})
	if err != nil {
		return "", err
	}
	defer analyzer.Close()

	lines := splitLines(args.Data)
	results := make([]entropywindow.Result, 0, len(lines))
	for _, line := range lines {
		res := analyzer.Process(line)
		if res.Ready && res.IsAnomaly {
			results = append(results, res)
		}
	}
	payload := map[string]any{
		"mode":          "local",
		"total_records": len(lines),
		"anomaly_count": len(results),
		"anomalies":     results,
	}
	body, _ := json.MarshalIndent(payload, "", "  ")
	return string(body), nil
}

func callRemoteAPI(args detectArgs) (string, error) {
	req, err := http.NewRequest("POST", API_URL+"/detect?algo=zstd", strings.NewReader(args.Data))
	if err != nil {
		return "", err
	}
	switch strings.ToLower(args.Format) {
	case "ndjson":
		req.Header.Set("Content-Type", "application/x-ndjson")
	default:
		req.Header.Set("Content-Type", "application/json")
	}
	if key := os.Getenv("DRIFTLOCK_API_KEY"); key != "" {
		req.Header.Set("X-Api-Key", key)
	}
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	if resp.StatusCode >= 300 {
		return "", fmt.Errorf("remote API error: %s", string(body))
	}
	return string(body), nil
}

func splitLines(payload string) []string {
	scanner := bufio.NewScanner(strings.NewReader(payload))
	scanner.Buffer(make([]byte, 0, 1024), 1024*1024)
	var out []string
	for scanner.Scan() {
		out = append(out, scanner.Text())
	}
	return out
}
