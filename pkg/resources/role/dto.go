package role

import "github.com/rs/xid"

type RoleID xid.ID

func (u RoleID) String() string { return xid.ID(u).String() }

type Role struct {
	ID          RoleID
	Name        string
	Colour      string
	Permissions []string
}
