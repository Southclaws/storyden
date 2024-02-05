package datagraph

import "github.com/rs/xid"

type Indexable interface {
	GetID() xid.ID
	GetKind() Kind
	GetName() string
	GetSlug() string
	GetDesc() string
	GetText() string
	GetProps() any
}
