package storyden

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"sync/atomic"
	"testing"
	"time"

	"github.com/coder/websocket"
	"github.com/rs/xid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/Southclaws/opt"
	"github.com/Southclaws/storyden/app/resources/account"
	"github.com/Southclaws/storyden/app/resources/post"
	"github.com/Southclaws/storyden/lib/plugin/rpc"
)

func TestRPCRequestToPlugin_UnmarshalsPing(t *testing.T) {
	messageID := xid.New()
	message := fmt.Sprintf(`{
		"jsonrpc": "2.0",
		"id": "%s",
		"method": "ping"
	}`, messageID.String())

	var req rpc.HostToPluginRequest
	err := json.Unmarshal([]byte(message), &req)

	require.NoError(t, err)
	require.NotNil(t, req.HostToPluginRequestUnion)

	ping, ok := req.HostToPluginRequestUnion.(*rpc.RPCRequestPing)
	require.True(t, ok, "expected *rpc.RPCRequestPing, got %T", req.HostToPluginRequestUnion)
	assert.Equal(t, messageID, ping.ID)
	assert.Equal(t, "ping", ping.Method)
}

func TestRPCRequestToPlugin_UnmarshalsEvent(t *testing.T) {
	// Create an actual event using the codegen types
	messageID := xid.New()
	threadID := post.ID(xid.New())
	event := rpc.RPCRequestEvent{
		ID:      messageID,
		Jsonrpc: "2.0",
		Method:  "event",
		Params: rpc.EventPayload{
			EventPayloadUnion: &rpc.EventThreadPublished{
				Event: "EventThreadPublished",
				ID:    threadID,
			},
		},
	}

	// Marshal it to JSON
	data, err := json.Marshal(event)
	require.NoError(t, err)

	assert.Contains(t, string(data), messageID.String())
	assert.Contains(t, string(data), threadID.String())

	// Now unmarshal it back through the union type
	var req rpc.HostToPluginRequest
	err = json.Unmarshal(data, &req)
	require.NoError(t, err)
	require.NotNil(t, req.HostToPluginRequestUnion)

	// Type switch to verify it's the right type
	eventReq, ok := req.HostToPluginRequestUnion.(*rpc.RPCRequestEvent)
	require.True(t, ok, "expected *rpc.RPCRequestEvent, got %T", req.HostToPluginRequestUnion)
	assert.Equal(t, messageID, eventReq.ID)
	assert.Equal(t, "event", eventReq.Method)

	// Verify the event payload
	threadEvent, ok := eventReq.Params.EventPayloadUnion.(*rpc.EventThreadPublished)
	require.True(t, ok, "expected *rpc.EventThreadPublished, got %T", eventReq.Params.EventPayloadUnion)
	assert.Equal(t, threadID, threadEvent.ID)
}

func TestRPCRequestToPlugin_FailsOnUnknownMethod(t *testing.T) {
	message := `{
		"jsonrpc": "2.0",
		"id": "unknown-123",
		"method": "unknown_method"
	}`

	var req rpc.HostToPluginRequest
	err := json.Unmarshal([]byte(message), &req)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "unknown type")
}

func TestRPCRequestToPlugin_FailsOnMissingMethod(t *testing.T) {
	message := `{
		"jsonrpc": "2.0",
		"id": "no-method-123"
	}`

	var req rpc.HostToPluginRequest
	err := json.Unmarshal([]byte(message), &req)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "missing discriminator field")
}

func TestRPCRequestToPlugin_TypeSwitchOnPing(t *testing.T) {
	messageID := xid.New()
	message := fmt.Sprintf(`{
		"jsonrpc": "2.0",
		"id": "%s",
		"method": "ping"
	}`, messageID)

	var req rpc.HostToPluginRequest
	err := json.Unmarshal([]byte(message), &req)
	require.NoError(t, err)
	require.NotNil(t, req.HostToPluginRequestUnion)

	switch r := req.HostToPluginRequestUnion.(type) {
	case *rpc.RPCRequestPing:
		assert.Equal(t, messageID, r.ID)
		assert.Equal(t, "ping", r.Method)
	default:
		t.Fatalf("expected *rpc.RPCRequestPing, got %T", r)
	}
}

func TestRPCRequestToPlugin_TypeSwitchOnEvent(t *testing.T) {
	// Create and marshal an event using codegen types
	messageID := xid.New()
	accountID := account.AccountID(xid.New())

	event := rpc.RPCRequestEvent{
		ID:      messageID,
		Jsonrpc: "2.0",
		Method:  "event",
		Params: rpc.EventPayload{
			EventPayloadUnion: &rpc.EventAccountCreated{
				Event: "EventAccountCreated",
				ID:    accountID,
			},
		},
	}

	data, err := json.Marshal(event)
	require.NoError(t, err)

	var req rpc.HostToPluginRequest
	err = json.Unmarshal(data, &req)
	require.NoError(t, err)
	require.NotNil(t, req.HostToPluginRequestUnion)

	// Use type switch like handleMessage does
	switch r := req.HostToPluginRequestUnion.(type) {
	case *rpc.RPCRequestEvent:
		assert.Equal(t, messageID, r.ID)
		assert.Equal(t, "event", r.Method)

		// Verify we can type switch on the event payload
		switch evt := r.Params.EventPayloadUnion.(type) {
		case *rpc.EventAccountCreated:
			assert.Equal(t, accountID, evt.ID)
		default:
			t.Fatalf("expected *rpc.EventAccountCreated, got %T", evt)
		}
	default:
		t.Fatalf("expected *rpc.RPCRequestEvent, got %T", r)
	}
}

func TestConfigureDispatchAndAck(t *testing.T) {
	configureReceived := make(chan map[string]any, 1)
	ackReceived := make(chan struct{}, 1)

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		conn, err := websocket.Accept(w, r, nil)
		if err != nil {
			return
		}
		defer conn.CloseNow()

		reqID := xid.New()
		req := rpc.RPCRequestConfigure{
			ID:      reqID,
			Jsonrpc: "2.0",
			Method:  "configure",
			Params: map[string]any{
				"name":    "configured",
				"enabled": true,
			},
		}

		reqBytes, err := json.Marshal(req)
		require.NoError(t, err)
		require.NoError(t, conn.Write(r.Context(), websocket.MessageText, reqBytes))

		_, responseBytes, err := conn.Read(r.Context())
		require.NoError(t, err)

		var response rpc.HostToPluginResponse
		require.NoError(t, json.Unmarshal(responseBytes, &response))
		require.Equal(t, reqID, response.ID)

		configureResponse, ok := response.Result.HostToPluginResponseUnionUnion.(*rpc.RPCResponseConfigure)
		require.True(t, ok)
		require.True(t, configureResponse.Ok)
		ackReceived <- struct{}{}

		_ = conn.Close(websocket.StatusNormalClosure, "done")
	}))
	defer srv.Close()

	t.Setenv("STORYDEN_RPC_URL", wsURLWithToken(srv.URL, "sdprt_test_configure_ack_123456789012"))

	pl, err := New(context.Background())
	require.NoError(t, err)

	pl.OnConfigure(func(_ context.Context, config map[string]any) error {
		configureReceived <- config
		return nil
	})

	runDone := make(chan error, 1)
	go func() {
		runDone <- pl.Run(context.Background())
	}()

	select {
	case cfg := <-configureReceived:
		assert.Equal(t, "configured", cfg["name"])
		assert.Equal(t, true, cfg["enabled"])
	case <-time.After(2 * time.Second):
		t.Fatal("configure callback was not called")
	}

	select {
	case <-ackReceived:
	case <-time.After(2 * time.Second):
		t.Fatal("configure response ack was not received by test server")
	}

	select {
	case err := <-runDone:
		require.NoError(t, err)
	case <-time.After(2 * time.Second):
		t.Fatal("run did not return")
	}
}

func TestConfigureHandlerErrorRespondsNotOK(t *testing.T) {
	ackReceived := make(chan bool, 1)

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		conn, err := websocket.Accept(w, r, nil)
		if err != nil {
			return
		}
		defer conn.CloseNow()

		reqID := xid.New()
		req := rpc.RPCRequestConfigure{
			ID:      reqID,
			Jsonrpc: "2.0",
			Method:  "configure",
			Params: map[string]any{
				"reject": true,
			},
		}

		reqBytes, err := json.Marshal(req)
		require.NoError(t, err)
		require.NoError(t, conn.Write(r.Context(), websocket.MessageText, reqBytes))

		_, responseBytes, err := conn.Read(r.Context())
		require.NoError(t, err)

		var response rpc.HostToPluginResponse
		require.NoError(t, json.Unmarshal(responseBytes, &response))
		configureResponse, ok := response.Result.HostToPluginResponseUnionUnion.(*rpc.RPCResponseConfigure)
		require.True(t, ok)
		ackReceived <- configureResponse.Ok

		_ = conn.Close(websocket.StatusNormalClosure, "done")
	}))
	defer srv.Close()

	t.Setenv("STORYDEN_RPC_URL", wsURLWithToken(srv.URL, "sdprt_test_configure_error_12345678901"))

	pl, err := New(context.Background())
	require.NoError(t, err)

	pl.OnConfigure(func(_ context.Context, _ map[string]any) error {
		return errors.New("reject config")
	})

	runDone := make(chan error, 1)
	go func() {
		runDone <- pl.Run(context.Background())
	}()

	select {
	case ok := <-ackReceived:
		assert.False(t, ok)
	case <-time.After(2 * time.Second):
		t.Fatal("configure response ack was not received by test server")
	}

	select {
	case err := <-runDone:
		require.NoError(t, err)
	case <-time.After(2 * time.Second):
		t.Fatal("run did not return")
	}
}

func TestRPCResponseBase_DistinguishesFromRequest(t *testing.T) {
	messageID := xid.New()
	responseMessage := fmt.Sprintf(`{
		"jsonrpc": "2.0",
		"id": "%s",
		"result": {
			"method": "get_config",
			"config": {
				"key1": "value1"
			}
		}
	}`, messageID)

	// Should fail to unmarshal as a request (no "method" field at top level)
	var req rpc.HostToPluginRequest
	err := json.Unmarshal([]byte(responseMessage), &req)
	assert.Error(t, err, "response should not unmarshal as request")

	// Should succeed to unmarshal as a response
	var base rpc.JsonRpcResponse
	err = json.Unmarshal([]byte(responseMessage), &base)
	require.NoError(t, err)
	assert.Equal(t, messageID, base.ID)
	assert.Equal(t, "2.0", base.Jsonrpc)
}

func TestShutdownBeforeRunReturnsImmediately(t *testing.T) {
	t.Setenv("STORYDEN_RPC_URL", "ws://localhost:12345/rpc")

	pl, err := New(context.Background())
	require.NoError(t, err)

	start := time.Now()
	err = pl.Shutdown()
	require.NoError(t, err)
	assert.Less(t, time.Since(start), 100*time.Millisecond)
}

func TestShutdownClosesConnectionAndStopsRun(t *testing.T) {
	connected := make(chan struct{})

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		conn, err := websocket.Accept(w, r, nil)
		if err != nil {
			return
		}
		defer conn.Close(websocket.StatusNormalClosure, "server closing")

		select {
		case <-connected:
		default:
			close(connected)
		}

		for {
			_, _, err := conn.Read(r.Context())
			if err != nil {
				return
			}
		}
	}))
	defer srv.Close()

	t.Setenv("STORYDEN_RPC_URL", toWSURL(srv.URL))

	pl, err := New(context.Background())
	require.NoError(t, err)

	done := make(chan error, 1)
	go func() {
		done <- pl.Run(context.Background())
	}()

	select {
	case <-connected:
	case <-time.After(2 * time.Second):
		t.Fatal("plugin did not connect")
	}
	require.Eventually(t, func() bool {
		return pl.getConn() != nil
	}, time.Second, 10*time.Millisecond)

	start := time.Now()
	err = pl.Shutdown()
	require.NoError(t, err)
	assert.Less(t, time.Since(start), 1*time.Second)

	select {
	case <-done:
	case <-time.After(2 * time.Second):
		t.Fatal("run did not return after shutdown")
	}
}

func TestExternalReconnectsOnAbnormalDisconnect(t *testing.T) {
	var connCount atomic.Int32
	reconnected := make(chan struct{}, 1)

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		conn, err := websocket.Accept(w, r, nil)
		if err != nil {
			return
		}
		defer conn.CloseNow()

		n := connCount.Add(1)
		if n == 1 {
			// Abrupt close without close frame (exceptional).
			conn.CloseNow()
			return
		}

		select {
		case reconnected <- struct{}{}:
		default:
		}

		for {
			_, _, err := conn.Read(r.Context())
			if err != nil {
				return
			}
		}
	}))
	defer srv.Close()

	t.Setenv("STORYDEN_RPC_URL", wsURLWithToken(srv.URL, "sdprt_test_external_token_1234567890123"))

	pl, err := New(context.Background())
	require.NoError(t, err)

	runCtx, cancel := context.WithCancel(context.Background())
	defer cancel()

	done := make(chan error, 1)
	go func() {
		done <- pl.Run(runCtx)
	}()

	select {
	case <-reconnected:
	case <-time.After(3 * time.Second):
		t.Fatal("plugin did not reconnect after abnormal disconnect")
	}

	require.NoError(t, pl.Shutdown())

	select {
	case err := <-done:
		require.NoError(t, err)
	case <-time.After(2 * time.Second):
		t.Fatal("run did not stop after shutdown")
	}

	assert.GreaterOrEqual(t, connCount.Load(), int32(2))
}

func TestExternalDoesNotReconnectOnNormalServerClose(t *testing.T) {
	var connCount atomic.Int32

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		conn, err := websocket.Accept(w, r, nil)
		if err != nil {
			return
		}
		defer conn.CloseNow()

		connCount.Add(1)
		_ = conn.Close(websocket.StatusNormalClosure, "inactive")
	}))
	defer srv.Close()

	t.Setenv("STORYDEN_RPC_URL", wsURLWithToken(srv.URL, "sdprt_test_external_token_abcdefabcdef"))

	pl, err := New(context.Background())
	require.NoError(t, err)

	err = pl.Run(context.Background())
	require.NoError(t, err)

	// Give some time in case an unwanted reconnect attempt is made.
	time.Sleep(250 * time.Millisecond)
	assert.Equal(t, int32(1), connCount.Load())
}

func TestExternalReconnectsOnServiceRestartClose(t *testing.T) {
	var connCount atomic.Int32
	reconnected := make(chan struct{}, 1)

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		conn, err := websocket.Accept(w, r, nil)
		if err != nil {
			return
		}
		defer conn.CloseNow()

		n := connCount.Add(1)
		if n == 1 {
			_ = conn.Close(websocket.StatusServiceRestart, "restarting")
			return
		}

		select {
		case reconnected <- struct{}{}:
		default:
		}

		for {
			_, _, err := conn.Read(r.Context())
			if err != nil {
				return
			}
		}
	}))
	defer srv.Close()

	t.Setenv("STORYDEN_RPC_URL", wsURLWithToken(srv.URL, "sdprt_test_external_token_service_restart_123"))

	pl, err := New(context.Background())
	require.NoError(t, err)

	runCtx, cancel := context.WithCancel(context.Background())
	defer cancel()

	done := make(chan error, 1)
	go func() {
		done <- pl.Run(runCtx)
	}()

	select {
	case <-reconnected:
	case <-time.After(3 * time.Second):
		t.Fatal("plugin did not reconnect after service restart close")
	}

	require.NoError(t, pl.Shutdown())
	select {
	case err := <-done:
		require.NoError(t, err)
	case <-time.After(2 * time.Second):
		t.Fatal("run did not stop after shutdown")
	}
}

func TestExternalDoesNotReconnectOnPolicyViolationClose(t *testing.T) {
	var connCount atomic.Int32

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		conn, err := websocket.Accept(w, r, nil)
		if err != nil {
			return
		}
		defer conn.CloseNow()

		connCount.Add(1)
		_ = conn.Close(websocket.StatusPolicyViolation, "token invalid")
	}))
	defer srv.Close()

	t.Setenv("STORYDEN_RPC_URL", wsURLWithToken(srv.URL, "sdprt_test_external_token_policy_123456789"))

	pl, err := New(context.Background())
	require.NoError(t, err)

	err = pl.Run(context.Background())
	require.NoError(t, err)

	time.Sleep(250 * time.Millisecond)
	assert.Equal(t, int32(1), connCount.Load())
}

func TestExternalDoesNotRetryOnAuthFailure(t *testing.T) {
	var reqCount atomic.Int32

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		reqCount.Add(1)
		http.Error(w, "unauthorized", http.StatusUnauthorized)
	}))
	defer srv.Close()

	t.Setenv("STORYDEN_RPC_URL", wsURLWithToken(srv.URL, "sdprt_test_external_token_deadbeefdeadbeef"))

	pl, err := New(context.Background())
	require.NoError(t, err)

	err = pl.Run(context.Background())
	require.Error(t, err)

	time.Sleep(250 * time.Millisecond)
	assert.Equal(t, int32(1), reqCount.Load())
}

func TestSendMatchesOutOfOrderResponses(t *testing.T) {
	type rpcRequest struct {
		id     xid.ID
		method string
	}

	received := make(chan []rpcRequest, 1)

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		conn, err := websocket.Accept(w, r, nil)
		if err != nil {
			return
		}
		defer conn.CloseNow()

		requests := make([]rpcRequest, 0, 2)
		for range 2 {
			_, message, err := conn.Read(r.Context())
			require.NoError(t, err)

			var req rpc.PluginToHostRequest
			require.NoError(t, json.Unmarshal(message, &req))
			require.NotNil(t, req.PluginToHostRequestUnion)

			switch v := req.PluginToHostRequestUnion.(type) {
			case *rpc.RPCRequestAccessGet:
				requests = append(requests, rpcRequest{id: v.ID, method: v.Method})
			case *rpc.RPCRequestGetConfig:
				requests = append(requests, rpcRequest{id: v.ID, method: v.Method})
			default:
				t.Fatalf("unexpected request type: %T", v)
			}
		}

		require.Len(t, requests, 2)

		sendResponse := func(req rpcRequest) {
			var result rpc.PluginToHostResponseUnionUnion
			switch req.method {
			case "access_get":
				result = &rpc.RPCResponseAccessGet{
					ID:      req.id,
					Jsonrpc: "2.0",
					Method:  opt.New("access_get"),
					Result: rpc.RPCResponseAccessGetResult{
						AccessKey: "sdbak_out_of_order",
						APIBaseURL: url.URL{
							Scheme: "http",
							Host:   "localhost:8000",
						},
					},
				}
			case "get_config":
				result = &rpc.RPCResponseGetConfig{
					Method: "get_config",
					Config: map[string]any{
						"name": "cfg",
					},
				}
			default:
				t.Fatalf("unexpected request method: %s", req.method)
			}

			resp := rpc.PluginToHostResponse{
				ID:      req.id,
				Jsonrpc: "2.0",
				Result: rpc.PluginToHostResponseUnion{
					PluginToHostResponseUnionUnion: result,
				},
			}

			b, err := json.Marshal(resp)
			require.NoError(t, err)
			require.NoError(t, conn.Write(r.Context(), websocket.MessageText, b))
		}

		// Reverse send order from receive order.
		sendResponse(requests[1])
		sendResponse(requests[0])
		received <- requests

		_ = conn.Close(websocket.StatusNormalClosure, "done")
	}))
	defer srv.Close()

	t.Setenv("STORYDEN_RPC_URL", wsURLWithToken(srv.URL, "sdprt_test_send_order_1234567890123"))

	pl, err := New(context.Background())
	require.NoError(t, err)

	runDone := make(chan error, 1)
	go func() {
		runDone <- pl.Run(context.Background())
	}()

	require.Eventually(t, func() bool {
		return pl.getConn() != nil
	}, 2*time.Second, 10*time.Millisecond)

	type sendResult struct {
		resp rpc.PluginToHostResponseUnionUnion
		err  error
	}
	accessCh := make(chan sendResult, 1)
	configCh := make(chan sendResult, 1)

	go func() {
		resp, err := pl.Send(context.Background(), rpc.RPCRequestAccessGet{
			Jsonrpc: "2.0",
			Method:  "access_get",
		})
		accessCh <- sendResult{resp: resp, err: err}
	}()

	go func() {
		resp, err := pl.Send(context.Background(), rpc.RPCRequestGetConfig{
			Jsonrpc: "2.0",
			Method:  "get_config",
		})
		configCh <- sendResult{resp: resp, err: err}
	}()

	select {
	case <-received:
	case <-time.After(2 * time.Second):
		t.Fatal("server did not receive both requests")
	}

	accessResult := <-accessCh
	require.NoError(t, accessResult.err)
	accessResp, ok := accessResult.resp.(*rpc.RPCResponseAccessGet)
	require.True(t, ok, "expected access_get response, got %T", accessResult.resp)
	assert.Equal(t, "sdbak_out_of_order", accessResp.Result.AccessKey)

	configResult := <-configCh
	require.NoError(t, configResult.err)
	configResp, ok := configResult.resp.(*rpc.RPCResponseGetConfig)
	require.True(t, ok, "expected get_config response, got %T", configResult.resp)
	assert.Equal(t, "cfg", configResp.Config["name"])

	select {
	case err := <-runDone:
		require.NoError(t, err)
	case <-time.After(2 * time.Second):
		t.Fatal("run did not return")
	}
}

func TestConfigureHandlersRunConcurrently(t *testing.T) {
	slowID := xid.New()
	fastID := xid.New()
	responseOrder := make(chan []xid.ID, 1)

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		conn, err := websocket.Accept(w, r, nil)
		if err != nil {
			return
		}
		defer conn.CloseNow()

		slowReq := rpc.RPCRequestConfigure{
			ID:      slowID,
			Jsonrpc: "2.0",
			Method:  "configure",
			Params: map[string]any{
				"name": "slow",
			},
		}
		fastReq := rpc.RPCRequestConfigure{
			ID:      fastID,
			Jsonrpc: "2.0",
			Method:  "configure",
			Params: map[string]any{
				"name": "fast",
			},
		}

		slowBytes, err := json.Marshal(slowReq)
		require.NoError(t, err)
		fastBytes, err := json.Marshal(fastReq)
		require.NoError(t, err)

		require.NoError(t, conn.Write(r.Context(), websocket.MessageText, slowBytes))
		require.NoError(t, conn.Write(r.Context(), websocket.MessageText, fastBytes))

		order := make([]xid.ID, 0, 2)
		for range 2 {
			_, responseBytes, err := conn.Read(r.Context())
			require.NoError(t, err)

			var response rpc.HostToPluginResponse
			require.NoError(t, json.Unmarshal(responseBytes, &response))
			configureResponse, ok := response.Result.HostToPluginResponseUnionUnion.(*rpc.RPCResponseConfigure)
			require.True(t, ok, "expected configure response, got %T", response.Result.HostToPluginResponseUnionUnion)
			require.True(t, configureResponse.Ok)
			order = append(order, response.ID)
		}

		responseOrder <- order
		_ = conn.Close(websocket.StatusNormalClosure, "done")
	}))
	defer srv.Close()

	t.Setenv("STORYDEN_RPC_URL", wsURLWithToken(srv.URL, "sdprt_test_configure_concurrency_123456"))

	pl, err := New(context.Background())
	require.NoError(t, err)

	pl.OnConfigure(func(_ context.Context, config map[string]any) error {
		name, _ := config["name"].(string)
		if name == "slow" {
			time.Sleep(300 * time.Millisecond)
		}
		return nil
	})

	runDone := make(chan error, 1)
	go func() {
		runDone <- pl.Run(context.Background())
	}()

	var order []xid.ID
	select {
	case order = <-responseOrder:
	case <-time.After(3 * time.Second):
		t.Fatal("did not receive configure responses")
	}

	require.Len(t, order, 2)
	assert.Equal(t, fastID, order[0], "fast configure should respond before slow configure")
	assert.Equal(t, slowID, order[1], "slow configure should respond second")

	select {
	case err := <-runDone:
		require.NoError(t, err)
	case <-time.After(2 * time.Second):
		t.Fatal("run did not return")
	}
}

func TestMarshalRequestWithGeneratedIDAlwaysInjectsID(t *testing.T) {
	id, data, err := marshalRequestWithGeneratedID(rpc.RPCRequestAccessGet{
		Jsonrpc: "2.0",
		Method:  "access_get",
	})
	require.NoError(t, err)
	require.False(t, id.IsNil())

	var body map[string]any
	require.NoError(t, json.Unmarshal(data, &body))
	require.Equal(t, id.String(), body["id"])

	id2, data2, err := marshalRequestWithGeneratedID(rpc.RPCRequestAccessGet{
		ID:      xid.New(),
		Jsonrpc: "2.0",
		Method:  "access_get",
	})
	require.NoError(t, err)
	require.False(t, id2.IsNil())

	var body2 map[string]any
	require.NoError(t, json.Unmarshal(data2, &body2))
	require.Equal(t, id2.String(), body2["id"])
}

func TestWithDefaultTimeoutSetsTimeoutWhenMissing(t *testing.T) {
	ctx, cancel := withDefaultTimeout(context.Background())
	defer cancel()

	deadline, ok := ctx.Deadline()
	require.True(t, ok)

	remaining := time.Until(deadline)
	assert.Greater(t, remaining, 28*time.Second)
	assert.LessOrEqual(t, remaining, 30*time.Second)
}

func TestWithDefaultTimeoutPreservesExistingDeadline(t *testing.T) {
	originalCtx, originalCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer originalCancel()

	ctx, cancel := withDefaultTimeout(originalCtx)
	defer cancel()

	originalDeadline, ok := originalCtx.Deadline()
	require.True(t, ok)

	deadline, ok := ctx.Deadline()
	require.True(t, ok)
	assert.Equal(t, originalDeadline, deadline)
}

func TestSanitizesRPCURLInErrors(t *testing.T) {
	u, err := url.Parse("ws://localhost:8000/rpc?token=sdprt_secret_token_123456")
	require.NoError(t, err)

	pl := &Plugin{
		rpcURL: u,
	}

	endpointURL := pl.rpcEndpointURL()
	assert.Equal(t, "ws://localhost:8000/rpc", endpointURL)
	assert.NotContains(t, endpointURL, "?token=")

	sanitizedErr := pl.sanitizeError(errors.New(`failed handshake to http://localhost:8000/rpc?token=sdprt_secret_token_123456`))
	assert.NotContains(t, sanitizedErr, "sdprt_secret_token_123456")
	assert.NotContains(t, sanitizedErr, "?token=")
	assert.Contains(t, sanitizedErr, "http://localhost:8000/rpc")
}

func toWSURL(httpURL string) string {
	if strings.HasPrefix(httpURL, "https://") {
		return "wss://" + strings.TrimPrefix(httpURL, "https://")
	}
	return "ws://" + strings.TrimPrefix(httpURL, "http://")
}

func wsURLWithToken(httpURL string, token string) string {
	return toWSURL(httpURL) + "/rpc?token=" + token
}
