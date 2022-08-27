package role

import "github.com/google/uuid"

type RoleID uuid.UUID

func (u RoleID) String() string { return uuid.UUID(u).String() }

type Role struct {
	ID          RoleID
	Name        string
	Colour      string
	Permissions []string
}
