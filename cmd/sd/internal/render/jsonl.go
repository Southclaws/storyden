package render

import (
	"encoding/json"
	"io"
)

// StreamJSONL writes each item as a single JSON line. The encoder writes
// directly to out, so callers can interleave many pages of results without
// buffering everything in memory.
func StreamJSONL[T any](out io.Writer, items []T) error {
	encoder := json.NewEncoder(out)
	for _, item := range items {
		if err := encoder.Encode(item); err != nil {
			return err
		}
	}
	return nil
}
