package asset

import (
	"github.com/Southclaws/storyden/internal/ent"
)

type Asset struct {
	URL      string
	MIMEType string
	Width    int
	Height   int
}

func FromModel(a *ent.Asset) *Asset {
	return &Asset{
		URL:      a.URL,
		MIMEType: a.Mimetype,
		Width:    a.Width,
		Height:   a.Height,
	}
}
