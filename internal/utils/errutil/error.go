package errutil

import (
	"context"
	"errors"

	"github.com/glebarez/go-sqlite"
)

// IsIgnored returns true for errors we just don't give a flying fuck about...
// ...in most contexts at least!
func IsIgnored(err error) bool {
	if errors.Is(err, context.Canceled) {
		return true
	}

	se := &sqlite.Error{}
	if errors.As(err, &se) {
		// "interrupted" not useful: https://www.sqlite.org/c3ref/interrupt.html
		return se.Code() == 9
	}

	return false
}
