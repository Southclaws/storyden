package datagraph

import "github.com/rs/xid"

type Indexable interface {
	GetID() xid.ID
	GetKind() Kind
	GetName() string
	GetText() string
	GetProps() any
}
