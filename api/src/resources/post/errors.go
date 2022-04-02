package post

import "errors"

var (
	ErrNoTitle        = errors.New("missing title")
	ErrNoBody         = errors.New("missing body")
	ErrUnauthorised   = errors.New("unauthorised")
	ErrTagNameTooLong = errors.New("tag name too long")
)
