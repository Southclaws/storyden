package rbac

import (
	"github.com/Southclaws/fault"
	"github.com/el-mike/restrict"

	"github.com/Southclaws/storyden/app/resources/account"
	"github.com/Southclaws/storyden/app/resources/reply"
	"github.com/Southclaws/storyden/app/resources/thread"
)

var (
	permissionCreate       = restrict.Permission{Action: ActionCreate}
	permissionRead         = restrict.Permission{Action: ActionRead}
	permissionUpdateThread = restrict.Permission{Preset: "update_thread"}
	permissionUpdatePost   = restrict.Permission{Preset: "update_post"}
	permissionDeleteThread = restrict.Permission{Preset: "delete_thread"}
	permissionDeletePost   = restrict.Permission{Preset: "delete_post"}
)

var defaultGrants = restrict.GrantsMap{
	ResourceThread: {
		&permissionCreate,
		&permissionRead,
		&permissionUpdateThread,
		&permissionDeleteThread,
	},
	ResourcePost: {
		&permissionCreate,
		&permissionRead,
		&permissionUpdatePost,
		&permissionDeletePost,
	},
}

var EveryoneRole = restrict.Role{
	ID:          "everyone",
	Description: "This role applies to every person. It cannot be deleted.",
	Grants:      defaultGrants,
}

var OwnerRole = restrict.Role{
	ID:          "owner",
	Description: "The owner role controls everything. It cannot be deleted",
	Parents:     []string{EveryoneRole.ID},
	Grants:      defaultGrants,
}

type threadAccessCondition struct{}

func (c *threadAccessCondition) Type() string { return "thread_access" }
func (c *threadAccessCondition) Check(request *restrict.AccessRequest) error {
	acc := request.Subject.(*account.Account)
	thr := request.Resource.(*thread.Thread)

	if thr.Author.ID == acc.ID {
		return nil
	}

	if thr.Author.Admin {
		return nil
	}

	return restrict.NewConditionNotSatisfiedError(c, request, fault.New("Account is not the author of the thread"))
}

type postAccessCondition struct{}

func (c *postAccessCondition) Type() string { return "post_access" }
func (c *postAccessCondition) Check(request *restrict.AccessRequest) error {
	acc := request.Subject.(*account.Account)
	thr := request.Resource.(*reply.Reply)

	if thr.Author.ID == acc.ID {
		return nil
	}

	if acc.Admin {
		return nil
	}

	return restrict.NewConditionNotSatisfiedError(c, request, fault.New("Account is not the author of the post"))
}

var updateThread = restrict.Permission{
	Action: ActionUpdate,
	Conditions: restrict.Conditions{
		&threadAccessCondition{},
	},
}

var updatePost = restrict.Permission{
	Action: ActionUpdate,
	Conditions: restrict.Conditions{
		&postAccessCondition{},
	},
}

var deleteThread = restrict.Permission{
	Action: ActionDelete,
	Conditions: restrict.Conditions{
		&threadAccessCondition{},
	},
}

var deletePost = restrict.Permission{
	Action: ActionDelete,
	Conditions: restrict.Conditions{
		&postAccessCondition{},
	},
}

var defaultPolicy = &restrict.PolicyDefinition{
	Roles: restrict.Roles{
		EveryoneRole.ID: &EveryoneRole,
		OwnerRole.ID:    &OwnerRole,
	},
	PermissionPresets: restrict.PermissionPresets{
		"update_thread": &updateThread,
		"update_post":   &updatePost,
		"delete_thread": &deleteThread,
		"delete_post":   &deletePost,
	},
}
