package plugin

import (
	"context"
	"encoding/json"

	"github.com/Southclaws/storyden/internal/infrastructure/wrun"
)

const PluginDirectory = "plugins"

type Binary []byte

func (b Binary) Validate(ctx context.Context, r wrun.Runner) (*Metadata, error) {
	mb, err := r.RunOnce(ctx, b, nil)
	if err != nil {
		return nil, err
	}

	var m Metadata
	if err = json.Unmarshal(mb, &m); err != nil {
		return nil, err
	}

	return &m, nil
}

type Metadata struct {
	Name    string
	Version string
}

type Package struct {
	Metadata Metadata
	Binary   []byte
}

type Available struct {
	StoragePath string
	Loaded      *Package
	Error       error
}

type Instance struct {
	Package Package
}
