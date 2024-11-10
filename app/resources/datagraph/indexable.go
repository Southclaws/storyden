package datagraph

import (
	"time"

	"github.com/rs/xid"

	"github.com/Southclaws/fault"
	"github.com/Southclaws/opt"
	"github.com/Southclaws/storyden/app/resources/asset"
)

var ErrInvalidReferenceScheme = fault.New("invalid reference scheme")

type OptAsset = opt.Optional[asset.Asset]

type (
	Identifiable interface{ GetID() xid.ID }             // Has a unique ID
	Slugged      interface{ GetSlug() string }           // Has a URL slug for web browser access
	Named        interface{ GetName() string }           // Has a renderable display name
	Described    interface{ GetDesc() string }           // Has a short description of some sort
	WithContent  interface{ GetContent() Content }       // Has long-form rich-text content
	WithProps    interface{ GetProps() map[string]any }  // Has arbitrary metadata
	WithAssets   interface{ GetAssets() []*asset.Asset } // Has media assets
	WithCreated  interface{ GetCreated() time.Time }     // Has a creation timestamp
	WithUpdated  interface{ GetUpdated() time.Time }     // Has an update timestamp
	WithCover    interface{ GetCover() OptAsset }        // Has a cover image
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
	WithCreated
	WithUpdated
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

// ItemRef describes a type which knows its ID and kind, but nothing else.
type ItemRef interface {
	Identifiable
	GetKind() Kind
}

type ItemList []Item
