package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"time"
)

type Params struct {
	Message string `json:"message"`
	Wait    int    `json:"wait"`
}

type RPC struct {
	ID     string `json:"id"`
	Method string `json:"method"`
	Params Params `json:"params"`
}

type RPCResponse struct {
	ID     string      `json:"id"`
	Result interface{} `json:"result"`
	Error  interface{} `json:"error,omitempty"`
}

func main() {
	s := bufio.NewScanner(os.Stdin)
	s.Split(bufio.ScanLines)

	fmt.Println(`{"name":"Test Plugin","version":"2.0","id":"test","author":"local"}`)

	for s.Scan() {
		line := s.Text()

		var req RPC
		if err := json.Unmarshal([]byte(line), &req); err != nil {
			fmt.Fprintf(os.Stderr, `{"error":"invalid_json","details":%q}`+"\n", err.Error())
			continue
		}

		fmt.Println("received rpc", req.Method, req.ID, req.Params)
		spinSleep(req.ID, req.Params.Wait)

		resp := RPCResponse{
			ID: req.ID,
			Result: fmt.Sprintf("handled method %q: %s for %d",
				req.Method, req.Params.Message, req.Params.Wait),
		}

		fmt.Println("responding to rpc", req.Method, req.ID, req.Params)
		b, _ := json.Marshal(resp)
		fmt.Println(string(b))
	}

	fmt.Println("CLOSING")
}

var sink any

func spinSleep(id string, ms int) {
	start := time.Now()
	for time.Since(start) < time.Millisecond*time.Duration(ms) {
		fmt.Println("%s: waiting...", id, time.Since(start), "/", time.Millisecond*time.Duration(ms))
		sink = time.Now()
	}
}
