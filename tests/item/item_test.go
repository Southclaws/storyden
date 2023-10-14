package item_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestMain(m *testing.M) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Hello, client")
	}))
	defer ts.Close()

	fmt.Println("DONE", ts.URL)
}
