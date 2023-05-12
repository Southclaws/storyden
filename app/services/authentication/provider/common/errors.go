package common

import "github.com/Southclaws/fault"

var (
	ErrAccountAlreadyExists = fault.New("account already exists")
	ErrNotFound             = fault.New("account not found")
)
