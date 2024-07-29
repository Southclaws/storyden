package asset

import (
	"context"

	"github.com/Southclaws/fault"
	"github.com/Southclaws/opt"
	"github.com/rs/xid"

	"github.com/Southclaws/storyden/app/resources/account"
	"github.com/Southclaws/storyden/internal/ent"
)

var errInvalidFormat = fault.New("invalid format")

type Repository interface {
	Add(ctx context.Context,
		owner account.AccountID,
		filename Filename,
		url string,
	) (*Asset, error)

	Get(ctx context.Context, id Filename) (*Asset, error)
	GetByID(ctx context.Context, id AssetID) (*Asset, error)

	Remove(ctx context.Context, owner account.AccountID, id Filename) error
}

type AssetID = xid.ID

func NewID() AssetID {
	return AssetID(xid.New())
}

type Asset struct {
	ID       AssetID
	Name     Filename
	URL      string
	Size     int
	Metadata Metadata
}

func FromModel(a *ent.Asset) *Asset {
	return &Asset{
		ID: AssetID(a.ID),
		Name: Filename{
			id:    opt.New(a.ID),
			name:  a.Filename,
			hasID: true,
		},
		URL:      a.URL,
		Metadata: a.Metadata,
	}
}

type Metadata map[string]any

func (m Metadata) GetMIMEType() string {
	v, ok := m["mime_type"]
	if !ok {
		return ""
	}

	s, ok := v.(string)
	if !ok {
		return ""
	}

	return s
}

func (m Metadata) GetWidth() float64 {
	v, ok := m["mime_type"]
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
