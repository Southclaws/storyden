package mcp

import "encoding/json"

type FixedJSONRPCResponse struct {
	JSONRPC string    `json:"jsonrpc"`
	ID      RequestId `json:"id"`
	Result  any       `json:"result"`
}

func (r FixedJSONRPCResponse) MarshalJSON() ([]byte, error) {
	obj := make(map[string]interface{})
	obj["jsonrpc"] = r.JSONRPC
	obj["id"] = r.ID
	obj["result"] = r.Result
	return json.Marshal(obj)
}
