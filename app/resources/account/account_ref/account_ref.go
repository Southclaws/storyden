package account_ref

import "github.com/rs/xid"

type ID xid.ID

func (u ID) String() string { return xid.ID(u).String() }
