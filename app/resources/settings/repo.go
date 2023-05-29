package settings

import "context"

type Settings struct {
	Title        Value[string]
	Description  Value[string]
	AccentColour Value[string]
	Public       Value[bool]
}

type Repository interface {
	// Init initialises with defaults if there are no settings.
	Init(ctx context.Context) error

	// Get returns all the current settings.
	Get(ctx context.Context) (*Settings, error)

	// SetValue and GetValue are just for internal or special use cases. They
	// work with serialised string formats rather than type safe struct values.
	SetValue(ctx context.Context, key, value string) error
	GetValue(ctx context.Context, key string) (string, error)
}
