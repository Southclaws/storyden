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
	OwnerID  xid.ID
	Name     Filename
	Size     int
	MIME     mime.Type
	Metadata Metadata
	Parent   opt.Optional[Asset]
}

func Map(a *ent.Asset) *Asset {
	parent := opt.NewPtrMap(a.Edges.Parent, func(a ent.Asset) Asset { return *Map(&a) })

	return &Asset{
		ID:      AssetID(a.ID),
		OwnerID: a.AccountID,
		Name: Filename{
			id:    opt.New(a.ID),
			name:  a.Filename,
			hasID: true,
		},
		Size:     a.Size,
		MIME:     mime.New(a.MimeType),
		Metadata: Metadata{RawMessage: a.Metadata},
		Parent:   parent,
	}
}
