package openapi

import "encoding/json"

// Marshal PluginGet responses through Plugin so union fields are preserved.
func (response PluginGet200JSONResponse) MarshalJSON() ([]byte, error) {
	return json.Marshal(Plugin(response.PluginGetOKJSONResponse))
}
