package bindings

import "net/http"

func unwrapWriter(w http.ResponseWriter) http.ResponseWriter {
	switch v := w.(type) {
	case interface{ Unwrap() http.ResponseWriter }:
		return v.Unwrap()
	default:
		return nil
	}
}
