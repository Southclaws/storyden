package storyden

import (
	"context"
	"fmt"
	"net/http"

	"github.com/Southclaws/storyden/app/transports/http/openapi"
	"github.com/Southclaws/storyden/lib/plugin/rpc"
)

func (p *Plugin) GetAccess(ctx context.Context) (rpc.RPCResponseAccessGetResult, error) {
	req := rpc.RPCRequestAccessGet{
		Jsonrpc: "2.0",
		Method:  "access_get",
	}

	resp, err := p.Send(ctx, req)
	if err != nil {
		return rpc.RPCResponseAccessGetResult{}, err
	}

	typed, ok := resp.(*rpc.RPCResponseAccessGet)
	if !ok {
		return rpc.RPCResponseAccessGetResult{}, fmt.Errorf("unexpected RPC response type: %T", resp)
	}

	if methodErr, ok := typed.Error.Get(); ok {
		if msg, ok := methodErr.Message.Get(); ok {
			return rpc.RPCResponseAccessGetResult{}, fmt.Errorf("access_get error: %s", msg)
		}
		return rpc.RPCResponseAccessGetResult{}, fmt.Errorf("access_get error")
	}

	return typed.Result, nil
}

func (p *Plugin) BuildAPIClient(ctx context.Context) (*openapi.ClientWithResponses, error) {
	access, err := p.GetAccess(ctx)
	if err != nil {
		return nil, err
	}

	url := fmt.Sprintf("%s/api", &access.APIBaseURL)
	authHeader := "Bearer " + access.AccessKey

	return openapi.NewClientWithResponses(
		url,
		openapi.WithRequestEditorFn(func(_ context.Context, req *http.Request) error {
			req.Header.Set("Authorization", authHeader)
			return nil
		}),
	)
}
