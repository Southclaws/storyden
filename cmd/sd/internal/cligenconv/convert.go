// Package cligenconv converts values between the HTTP API client's openapi
// types and cligen's OpenCLI-schema-generated output types. The two are
// independently generated from the same wire contract, so a JSON round-trip
// is a safe, generic way to move between them without hand-written field
// mapping for every command.
package cligenconv

import "encoding/json"

func Convert[T any](v any) (T, error) {
	var out T

	data, err := json.Marshal(v)
	if err != nil {
		return out, err
	}

	if err := json.Unmarshal(data, &out); err != nil {
		return out, err
	}

	return out, nil
}
