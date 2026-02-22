package main

//go:generate ./package.nu

import (
	"context"
	"encoding/json"
	"log"
	"os"
	"os/signal"
	"time"

	"github.com/coder/websocket"
)

type RPCRequest struct {
	Jsonrpc string                 `json:"jsonrpc"`
	Method  string                 `json:"method"`
	Params  map[string]interface{} `json:"params"`
	ID      int                    `json:"id"`
}

type RPCResponse struct {
	Jsonrpc string                 `json:"jsonrpc"`
	ID      int                    `json:"id"`
	Result  map[string]interface{} `json:"result,omitempty"`
	Error   *RPCError              `json:"error,omitempty"`
}

type RPCError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	url := os.Getenv("STORYDEN_RPC_URL")
	if url == "" {
		url = "ws://localhost:8000/rpc/"
	}

	log.Printf("connecting to %s", url)

	conn, _, err := websocket.Dial(ctx, url, nil)
	if err != nil {
		log.Fatalf("failed to connect: %v", err)
	}
	defer conn.CloseNow()

	log.Println("connected successfully")

	registerReq := RPCRequest{
		Jsonrpc: "2.0",
		Method:  "plugin.register",
		Params: map[string]interface{}{
			"plugin_id": "example-plugin",
		},
		ID: 1,
	}

	registerData, err := json.Marshal(registerReq)
	if err != nil {
		log.Fatalf("failed to marshal register request: %v", err)
	}

	if err := conn.Write(ctx, websocket.MessageText, registerData); err != nil {
		log.Fatalf("failed to send register request: %v", err)
	}

	msgType, data, err := conn.Read(ctx)
	if err != nil {
		log.Fatalf("failed to read register response: %v", err)
	}

	if msgType == websocket.MessageText {
		var resp RPCResponse
		if err := json.Unmarshal(data, &resp); err != nil {
			log.Fatalf("failed to unmarshal response: %v", err)
		}
		log.Printf("registered: %+v", resp.Result)
	}

	go func() {
		<-ctx.Done()
		log.Println("shutting down")
		conn.Close(websocket.StatusNormalClosure, "shutdown")
	}()

	for {
		select {
		case <-ctx.Done():
			return
		case <-time.After(5 * time.Second):
			log.Println("heartbeat")
		}
	}
}
