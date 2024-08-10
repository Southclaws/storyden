package datagraph

import (
	"github.com/rs/xid"

	"github.com/Southclaws/storyden/app/resources/asset"
	"github.com/Southclaws/storyden/app/resources/content"
)

type (
	Identifiable interface{ GetID() xid.ID }             // Has a unique ID
	Slugged      interface{ GetSlug() string }           // Has a URL slug for web browser access
	Named        interface{ GetName() string }           // Has a renderable display name
	Described    interface{ GetDesc() string }           // Has a short description of some sort
	WithContent  interface{ GetContent() content.Rich }  // Has long-form rich-text content
	WithProps    interface{ GetProps() map[string]any }  // Has arbitrary metadata
	WithAssets   interface{ GetAssets() []*asset.Asset } // Has media assets
)

// Addressable describes a type that can be uniquely identified via either an ID
// or a slug and also posses a human-readable display name.
type Addressable interface {
	Identifiable
	Slugged
	Named
}

// Publishable describes a type which can be published to platform audiences via
// APIs and it is uniquely addressable as well as contains rich text content.
type Publishable interface {
	Addressable
	Described
	WithContent
	WithAssets
}

// Item describes an object which exists in the "datagraph", an abstract concept
// which is formed of a graph of content which may reference each other such as
// discussion posts, blog posts, newsletters, library nodes, profiles, and more.
// It is a generic type which can be narrowed via `GetKind` or a type switch.
type Item interface {
	GetKind() Kind
	Publishable
	WithProps
}

type ItemList []Item

// Ref is a non-hydrated type to express a result type from semdex operations
// such as searching or recommendations. It can be hydrated into an Item using
// the Kind field to select a relevant resource querier to find the full object.
type Ref struct {
	ID        xid.ID
	Kind      Kind
	Relevance float64
}

type RefList []*Ref
