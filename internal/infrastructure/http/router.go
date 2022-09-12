package http

import "net/http"

func NewRouter() *http.ServeMux {
	return http.NewServeMux()
}
