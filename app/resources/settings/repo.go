package settings

import "context"

type Settings struct {
	Title        Value[string]
	Description  Value[string]
	AccentColour Value[uint32]
	Public       Value[bool]
}

type Repository interface {
	// Get returns all the current settings.
	Get(ctx context.Context) (*Settings, error)

	// SetValue and GetValue are just for internal or special use cases. They
	// work with serialised string formats rather than type safe struct values.
	SetValue(ctx context.Context, key, value string) error
	GetValue(ctx context.Context, key string) (string, error)
}
