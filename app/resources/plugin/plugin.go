package plugin

import (
	"context"
	"encoding/json"
	"time"

	"github.com/rs/xid"

	"github.com/Southclaws/storyden/app/resources/account"
	"github.com/Southclaws/storyden/internal/ent"
	"github.com/Southclaws/storyden/internal/infrastructure/wrun"
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
	Status   Status
	AddedBy  account.AccountID
	Manifest plugin.Manifest
	FilePath string
}

// Loaded represents a plugin that has been read and validated from its binary.
type Loaded struct {
	Metadata plugin.Manifest
	Binary   []byte
}

// Available represents a plugin that has been installed and registered in the
// database and is ready to use.
type Available struct {
	Record Record
	Loaded *Loaded
}

type Binary []byte

func (b Binary) Validate(ctx context.Context, r wrun.Runner) (*Loaded, error) {
	mb, err := r.RunOnce(ctx, b, nil)
	if err != nil {
		return nil, err
	}

	var m plugin.Manifest
	if err = json.Unmarshal(mb, &m); err != nil {
		return nil, err
	}

	return &Loaded{
		Metadata: m,
		Binary:   b,
	}, nil
}

func MapRecord(in *ent.Plugin) (*Record, error) {
	status, err := MapStatus(in)
	if err != nil {
		return nil, err
	}

	return &Record{
		ID:       ID(in.ID),
		Created:  in.CreatedAt,
		Status:   *status,
		AddedBy:  account.AccountID(in.AddedBy),
		Manifest: in.Manifest,
		FilePath: in.Path,
	}, nil
}
