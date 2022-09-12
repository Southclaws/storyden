package openapi

import "github.com/rs/xid"

// XID converts an openapi identifier to an xid. This is to work around an issue
// with the oapi-codegen generated code which generates the identifier as a new
// type instead of an alias which results in the marshal functions being hidden.
func (i Identifier) XID() xid.ID {
	v, err := xid.FromString(string(i))
	if err != nil {
		panic(err)
	}

	return v
}

// id converts any arbitrary xid.ID derivative to an *openapi.Identifier type.
func IdentifierFrom(id xid.ID) *Identifier {
	oid := Identifier(id.String())
	return &oid
}
