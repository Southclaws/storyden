package httpserver

import "net/http"

func NewRouter() *http.ServeMux {
	return http.NewServeMux()
}
