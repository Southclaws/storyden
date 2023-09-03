package settings

import (
	"context"

	"github.com/Southclaws/opt"
)

type Settings struct {
	Title        Value[string]
	Description  Value[string]
	AccentColour Value[string]
	Public       Value[bool]
}

type Partial struct {
	Title        opt.Optional[string]
	Description  opt.Optional[string]
	AccentColour opt.Optional[string]
	Public       opt.Optional[bool]
}

type Repository interface {
	// Init initialises with defaults if there are no settings.
	Init(ctx context.Context) error

	// Get returns all the current settings.
	Get(ctx context.Context) (*Settings, error)
	Set(ctx context.Context, s Partial) (*Settings, error)

	// SetValue and GetValue are just for internal or special use cases. They
	// work with serialised string formats rather than type safe struct values.
	SetValue(ctx context.Context, key, value string) error
	GetValue(ctx context.Context, key string) (string, error)
}
