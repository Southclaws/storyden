package web

import (
	"context"
	"errors"
	"net/http"

	"go.uber.org/zap"
)

// Error represents the object returned from the API for any form of problem
// encountered during a request.
//
// Errors contain technical information as well as optional human-readable
// error messages and suggestions for solutions.
type Error struct {
	Message    string `json:"message,omitempty"`   // human readable message
	Suggestion string `json:"suggested,omitempty"` // suggestion for fix
	Error      string `json:"error,omitempty"`     // internal error string
}

// StatusNotFound writes a pretty error
func StatusNotFound(w http.ResponseWriter, err error) {
	w.WriteHeader(http.StatusNotFound)
	if err == nil {
		errToWriter(w, errors.New("not found"))
	} else {
		errToWriter(w, err)
	}
}

// StatusInternalServerError writes a pretty error and logs
func StatusInternalServerError(w http.ResponseWriter, err error) {
	if errors.Is(err, context.Canceled) {
		return
	}
	zap.L().Error("internal error", zap.Error(err))
	w.WriteHeader(http.StatusInternalServerError)
	errToWriter(w, errors.New("something went wrong but the details have been omitted from this error for security reasons"))
}

// StatusUnauthorized writes a pretty error
func StatusUnauthorized(w http.ResponseWriter, err error) {
	w.WriteHeader(http.StatusUnauthorized)
	errToWriter(w, err)
}

// StatusNotAcceptable writes a pretty error
func StatusNotAcceptable(w http.ResponseWriter, err error) {
	w.WriteHeader(http.StatusNotAcceptable)
	errToWriter(w, err)
}

// StatusBadRequest writes a pretty error
func StatusBadRequest(w http.ResponseWriter, err error) {
	w.WriteHeader(http.StatusBadRequest)
	errToWriter(w, err)
}

// HumanReadable is an error with a human readable description for displaying
// on UI to end-users. It also contains an optional suggested solution field.
type HumanReadable struct {
	err     error
	desc    string
	suggest string
}

// Error implements error
func (h *HumanReadable) Error() string {
	return h.err.Error()
}

// WithDescription wraps an error with a human readable description.
func WithDescription(err error, desc string) error {
	return &HumanReadable{err, desc, ""}
}

// WithSuggestion wraps an error with a human readable description and solution.
func WithSuggestion(err error, desc, suggest string) error {
	return &HumanReadable{err, desc, suggest}
}

func errToWriter(w http.ResponseWriter, err error) {
	if err == nil {
		err = errors.New("Unknown or unspecified error")
	}

	var message string
	var suggest string
	herr, ok := err.(*HumanReadable)
	if ok {
		message = herr.desc
		suggest = herr.suggest
	}

	Write(w, Error{
		Message:    message,
		Suggestion: suggest,
		Error:      err.Error(),
	})
}
