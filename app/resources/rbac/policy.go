package rbac

import "github.com/el-mike/restrict"

type Action = string

const (
	ActionCreate Action = "create"
	ActionRead   Action = "read"
	ActionUpdate Action = "update"
	ActionDelete Action = "delete"
)

type Resource = string

const (
	ResourceAccount Resource = "account"
	ResourceThread  Resource = "thread"
	ResourcePost    Resource = "post"
)

var (
	permissionCreate      = restrict.Permission{Action: ActionCreate}
	permissionRead        = restrict.Permission{Action: ActionRead}
	permissionUpdateOwned = restrict.Permission{Preset: "updateOwned"}
	permissionDelete      = restrict.Permission{
		Action: "delete",
		Conditions: restrict.Conditions{
			&restrict.EmptyCondition{
				ID: "deleteActive",
				Value: &restrict.ValueDescriptor{
					Source: restrict.ResourceField,
					Field:  "Active",
				},
			},
		},
	}
)

var EveryoneRole = restrict.Role{
	ID:          "everyone",
	Description: "This role applies to every person. It cannot be deleted.",
	Grants: restrict.GrantsMap{
		ResourceThread: {
			&permissionCreate,
			&permissionRead,
			&permissionUpdateOwned,
			&permissionDelete,
		},
	},
}

var OwnerRole = restrict.Role{
	ID:          "owner",
	Description: "The owner role controls everything. It cannot be deleted",
	Parents:     []string{EveryoneRole.ID},
	Grants: restrict.GrantsMap{
		EveryoneRole.ID: {
			&restrict.Permission{Action: ActionUpdate},
		},
	},
}

var updateOwned = restrict.Permission{
	Action: ActionUpdate,
	Conditions: restrict.Conditions{
		&restrict.EqualCondition{
			ID: "isOwner",
			Left: &restrict.ValueDescriptor{
				Source: restrict.ResourceField,
				Field:  "CreatedBy",
			},
			Right: &restrict.ValueDescriptor{
				Source: restrict.SubjectField,
				Field:  "ID",
			},
		},
	},
}

var defaultPolicy = &restrict.PolicyDefinition{
	Roles: restrict.Roles{
		EveryoneRole.ID: &EveryoneRole,
		OwnerRole.ID:    &OwnerRole,
	},
	PermissionPresets: restrict.PermissionPresets{
		"updateOwned": &updateOwned,
	},
}
