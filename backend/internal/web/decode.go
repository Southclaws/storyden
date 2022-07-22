package web

import (
	"encoding/json"
	"net/http"

	"github.com/Southclaws/qstring"
	"github.com/asaskevich/govalidator"
)

func ParseQuery(w http.ResponseWriter, r *http.Request, out interface{}) bool {
	if err := qstring.Unmarshal(r.URL.Query(), out); err != nil {
		StatusBadRequest(w, err)
		return false
	}
	if _, err := govalidator.ValidateStruct(out); err != nil {
		StatusBadRequest(w, err)
		return false
	}
	return true
}

func ParseBody(w http.ResponseWriter, r *http.Request, out interface{}) bool {
	if err := json.NewDecoder(r.Body).Decode(out); err != nil {
		StatusBadRequest(w, WithSuggestion(err,
			"Could not process request data",
			"Please try again, if the issue persists contact the support team."))
		return false
	}
	if _, err := govalidator.ValidateStruct(out); err != nil {
		StatusBadRequest(w, WithSuggestion(err,
			"Could not validate request data",
			"Please try again, if the issue persists contact the support team."))
		return false
	}
	return true
}
