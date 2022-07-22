package web

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/schema"
	"github.com/pkg/errors"
)

var decoder = schema.NewDecoder()

func DecodeBody(r *http.Request, v interface{}) error {
	switch r.Header.Get("Content-Type") {
	case "application/json":
		return json.NewDecoder(r.Body).Decode(v)

	case "application/x-www-form-urlencoded":
		err := r.ParseForm()
		if err != nil {
			return errors.Wrap(err, "failed to parse form")
		}
		return decoder.Decode(v, r.PostForm)
	}
	return errors.Errorf("cannot decode content type %s", r.Header.Get("Content-Type"))
}
