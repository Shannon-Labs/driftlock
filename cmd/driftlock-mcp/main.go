package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
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
		Description: "Analyze a string of log/metric data for anomalies using Driftlock's CBAD engine.",
		InputSchema: map[string]any{
			"type": "object",
			"properties": map[string]any{
				"data": map[string]string{
					"type":        "string",
					"description": "The raw log or metric data (JSON/NDJSON) to analyze.",
				},
				"format": map[string]string{
					"type":        "string",
					"description": "The format of data: 'json' or 'ndjson'",
					"enum":        []string{"json", "ndjson"},
				},
			},
			"required": []string{"data"},
		},
	},
	{
		Name:        "check_system_health",
		Description: "Check the health status of the Driftlock API.",
		InputSchema: map[string]any{
			"type": "object",
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
		var args struct {
			Data   string `json:"data"`
			Format string `json:"format"`
		}
		if err := json.Unmarshal(params.Arguments, &args); err != nil {
			return "", fmt.Errorf("invalid arguments")
		}
		
		req, err := http.NewRequest("POST", API_URL+"/detect?algo=zstd", strings.NewReader(args.Data))
		if err != nil {
			return "", err
		}
		req.Header.Set("Content-Type", "application/json")
		// Note: MCP context usually implies authorized agent, but here we'd need an API key.
		// For MVP, we assume public/demo or Env var.
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
		return string(body), nil
	}
	return "", fmt.Errorf("unknown tool")
}
