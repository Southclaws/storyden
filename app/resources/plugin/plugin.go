package plugin

import (
	"context"
	"time"

	"github.com/rs/xid"

	"github.com/Southclaws/opt"
	"github.com/Southclaws/storyden/app/resources/account"
	"github.com/Southclaws/storyden/internal/ent"
	"github.com/Southclaws/storyden/lib/plugin"
)

const PluginDirectory = "plugins"

type ID xid.ID

// Record represents the database record for a plugin. It's the source of truth
// for plugins in general, it stores a copy of the manifest, who added it and
// the file path for the plugin's WASM binary application.
type Record struct {
	ID       ID
	Created  time.Time
	AddedBy  account.AccountID
	Manifest plugin.Manifest
	FilePath string

	State          ActiveState
	StateChangedAt time.Time
	StatusMessage  string
	Details        map[string]any

	// TODO: Invariant type?
	StartedAt time.Time
}

type Loaded struct {
	Record Record
	State  ActiveState
}

// Validated represents a plugin that has been read and validated from its binary.
type Validated struct {
	Metadata plugin.Manifest
	Binary   []byte
}

// Available represents a plugin that has been installed and registered in the
// database and is ready to use.
type Available struct {
	Record Record
	Loaded *Validated
}

type Binary []byte

// TODO: Restructure, move
func (b Binary) Validate(ctx context.Context, executor func(b []byte) (*plugin.Manifest, error)) (*Validated, error) {
	m, err := executor(b)
	if err != nil {
		return nil, err
	}

	return &Validated{
		Metadata: *m,
		Binary:   b,
	}, nil
}

func MapRecord(in *ent.Plugin) (*Record, error) {
	desiredActiveState, err := NewActiveState(in.ActiveState)
	if err != nil {
		return nil, err
	}

	return &Record{
		ID:       ID(in.ID),
		Created:  in.CreatedAt,
		AddedBy:  account.AccountID(in.AddedBy),
		Manifest: in.Manifest,
		FilePath: in.Path,

		State:          desiredActiveState,
		StateChangedAt: in.ActiveStateChangedAt,
		StatusMessage:  opt.NewPtr(in.StatusMessage).OrZero(),
		Details:        in.StatusDetails,
	}, nil
}
