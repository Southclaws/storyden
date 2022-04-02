package web

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/pkg/errors"
)

func Write(w http.ResponseWriter, data interface{}) {
	bytes, err := json.Marshal(data)
	if err != nil {
		StatusInternalServerError(w, errors.Wrap(err, "failed to encode payload"))
		return
	}

	w.Header().Add("Content-Length", fmt.Sprint(len(bytes)))
	w.Write(bytes)
}
