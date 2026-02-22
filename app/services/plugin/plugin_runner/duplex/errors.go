package duplex

import (
	"errors"
	"fmt"
)

var (
	ErrClosed      = errors.New("duplex closed")
	ErrRejected    = errors.New("duplex rejected")
	ErrUnavailable = errors.New("duplex unavailable")
	ErrFailed      = errors.New("duplex failed")
)

type Error struct {
	Kind   error
	Reason string
}

func NewError(kind error, reason string) error {
	return Error{
		Kind:   kind,
		Reason: reason,
	}
}

func (e Error) Error() string {
	if e.Reason == "" {
		return e.Kind.Error()
	}
	return fmt.Sprintf("%s: %s", e.Kind, e.Reason)
}

func (e Error) Is(target error) bool {
	return target == e.Kind
}

func IsExpectedDisconnect(err error) bool {
	return errors.Is(err, ErrClosed)
}
