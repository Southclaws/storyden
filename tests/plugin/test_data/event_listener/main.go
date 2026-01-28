package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"sync"
	"syscall"
	"time"

	"github.com/coder/websocket"
	"github.com/rs/xid"
)

type RPCRequest struct {
	Jsonrpc string                 `json:"jsonrpc"`
	Method  string                 `json:"method"`
	Params  map[string]interface{} `json:"params,omitempty"`
	ID      int                    `json:"id"`
}

type RPCResponse struct {
	Jsonrpc string                 `json:"jsonrpc"`
	Result  map[string]interface{} `json:"result,omitempty"`
	Error   *RPCError              `json:"error,omitempty"`
	ID      int                    `json:"id"`
}

type RPCError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type Plugin struct {
	conn          *websocket.Conn
	mu            sync.Mutex
	nextID        int
	pending       map[int]chan RPCResponse
	eventHandlers map[string]func(map[string]interface{})
	outputDir     string
}

func main() {
	rpcURL := os.Getenv("STORYDEN_RPC_URL")
	if rpcURL == "" {
		log.Fatal("STORYDEN_RPC_URL not set")
	}

	outputDir := os.Getenv("OUTPUT_DIR")
	if outputDir == "" {
		cwd, err := os.Getwd()
		if err != nil {
			outputDir = "."
		} else {
			outputDir = cwd
		}
	}

	if err := os.MkdirAll(outputDir, 0o755); err != nil {
		log.Fatalf("Failed to create output directory: %v", err)
	}

	log.Printf("Connecting to %s", rpcURL)
	log.Printf("Output directory: %s", outputDir)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	conn, _, err := websocket.Dial(ctx, rpcURL, nil)
	if err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}
	defer conn.Close(websocket.StatusNormalClosure, "shutting down")

	plugin := &Plugin{
		conn:          conn,
		pending:       make(map[int]chan RPCResponse),
		eventHandlers: make(map[string]func(map[string]interface{})),
		outputDir:     outputDir,
	}

	plugin.eventHandlers["event"] = plugin.handleEvent

	go plugin.readLoop(ctx)

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	log.Println("Plugin started, waiting for events...")

	<-sigChan
	log.Println("Shutting down...")
}

func (p *Plugin) readLoop(ctx context.Context) {
	for {
		_, message, err := p.conn.Read(ctx)
		if err != nil {
			log.Printf("Read error: %v", err)
			return
		}

		var req RPCRequest
		if err := json.Unmarshal(message, &req); err == nil {
			if handler, ok := p.eventHandlers[req.Method]; ok {
				handler(req.Params)

				response := RPCResponse{
					Jsonrpc: "2.0",
					Result:  map[string]interface{}{"ok": true},
					ID:      req.ID,
				}

				respData, err := json.Marshal(response)
				if err != nil {
					log.Printf("Failed to marshal response: %v", err)
					continue
				}

				if err := p.conn.Write(ctx, websocket.MessageText, respData); err != nil {
					log.Printf("Failed to send response: %v", err)
				}
			}
			continue
		}

		var resp RPCResponse
		if err := json.Unmarshal(message, &resp); err == nil {
			p.mu.Lock()
			if ch, ok := p.pending[resp.ID]; ok {
				delete(p.pending, resp.ID)
				p.mu.Unlock()
				ch <- resp
			} else {
				p.mu.Unlock()
			}
			continue
		}

		log.Printf("Unknown message: %s", string(message))
	}
}

func (p *Plugin) handleEvent(params map[string]interface{}) {
	log.Printf("Received event: %+v", params)

	idVal, ok := params["ID"]
	if !ok {
		log.Printf("Event missing ID field")
		return
	}

	var threadID string
	switch v := idVal.(type) {
	case string:
		threadID = v
	case []interface{}:
		// ID is a byte array from xid.ID (12 bytes)
		if len(v) == 12 {
			var xidBytes [12]byte
			for i, b := range v {
				if num, ok := b.(float64); ok {
					xidBytes[i] = byte(num)
				}
			}
			// Convert raw bytes to xid string representation
			id := xid.ID(xidBytes)
			threadID = id.String()
		}
	case map[string]interface{}:
		if idStr, ok := v["id"].(string); ok {
			threadID = idStr
		}
	default:
		log.Printf("Unknown ID type: %T", idVal)
		return
	}

	if threadID == "" {
		log.Printf("Could not extract thread ID from event")
		return
	}

	// Replace the byte array ID with the string representation
	params["ID"] = threadID

	filename := filepath.Join(p.outputDir, fmt.Sprintf("%s.json", threadID))
	data, err := json.MarshalIndent(params, "", "  ")
	if err != nil {
		log.Printf("Failed to marshal event: %v", err)
		return
	}

	if err := os.WriteFile(filename, data, 0o644); err != nil {
		log.Printf("Failed to write event file: %v", err)
		return
	}

	log.Printf("Wrote event to %s", filename)
}

func (p *Plugin) send(method string, params map[string]interface{}) (map[string]interface{}, error) {
	p.mu.Lock()
	id := p.nextID
	p.nextID++
	respChan := make(chan RPCResponse, 1)
	p.pending[id] = respChan
	p.mu.Unlock()

	req := RPCRequest{
		Jsonrpc: "2.0",
		Method:  method,
		Params:  params,
		ID:      id,
	}

	data, err := json.Marshal(req)
	if err != nil {
		p.mu.Lock()
		delete(p.pending, id)
		p.mu.Unlock()
		return nil, err
	}

	if err := p.conn.Write(context.Background(), websocket.MessageText, data); err != nil {
		p.mu.Lock()
		delete(p.pending, id)
		p.mu.Unlock()
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	select {
	case resp := <-respChan:
		if resp.Error != nil {
			return nil, fmt.Errorf("RPC error: %s", resp.Error.Message)
		}
		return resp.Result, nil
	case <-ctx.Done():
		p.mu.Lock()
		delete(p.pending, id)
		p.mu.Unlock()
		return nil, fmt.Errorf("timeout waiting for response")
	}
}
