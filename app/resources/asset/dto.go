package asset

import (
	"github.com/Southclaws/storyden/internal/ent"
)

type AssetID string

type Asset struct {
	ID       AssetID
	URL      string
	MIMEType string
	Size     int
	Width    int
	Height   int
}

func FromModel(a *ent.Asset) *Asset {
	return &Asset{
		ID:       AssetID(a.ID),
		URL:      a.URL,
		MIMEType: a.Mimetype,
		Width:    a.Width,
		Height:   a.Height,
	}
}
