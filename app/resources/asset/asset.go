package asset

import (
	"github.com/Southclaws/opt"
	"github.com/rs/xid"

	"github.com/Southclaws/storyden/internal/ent"
	"github.com/Southclaws/storyden/internal/mime"
)

type AssetID = xid.ID

func NewID() AssetID {
	return AssetID(xid.New())
}

type Asset struct {
	ID       AssetID
	Name     Filename
	Size     int
	MIME     mime.Type
	Metadata Metadata
	Parent   opt.Optional[Asset]
}

func Map(a *ent.Asset) *Asset {
	parent := opt.NewPtrMap(a.Edges.Parent, func(a ent.Asset) Asset { return *Map(&a) })

	return &Asset{
		ID: AssetID(a.ID),
		Name: Filename{
			id:    opt.New(a.ID),
			name:  a.Filename,
			hasID: true,
		},
		Size:     a.Size,
		MIME:     mime.New(a.MimeType),
		Metadata: a.Metadata,
		Parent:   parent,
	}
}

type Metadata map[string]any

func (m Metadata) GetWidth() float64 {
	v, ok := m["width"]
	if !ok {
		return 0.0
	}

	s, ok := v.(float64)
	if !ok {
		return 0.0
	}

	return s
}

func (m Metadata) GetHeight() float64 {
	v, ok := m["height"]
	if !ok {
		return 0.0
	}

	s, ok := v.(float64)
	if !ok {
		return 0.0
	}

	return s
}
