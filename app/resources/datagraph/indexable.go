package datagraph

import (
	"github.com/rs/xid"

	"github.com/Southclaws/storyden/app/resources/asset"
)

type Indexable interface {
	GetID() xid.ID
	GetKind() Kind
	GetName() string
	GetSlug() string
	GetDesc() string
	GetText() string
	GetProps() any
}

// NodeReference is a general structure that refers to a datagraph node of Kind.
// TODO: distinguish between a reference and a hydrated instance.
type NodeReference struct {
	ID          xid.ID
	Score       float64
	Kind        Kind
	Name        string
	Slug        string
	Description string
	Asset       *asset.Asset
}

type NodeReferenceList []*NodeReference
