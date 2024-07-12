package rbac

import (
	"github.com/Southclaws/fault"
	"github.com/el-mike/restrict"

	"github.com/Southclaws/storyden/app/resources/account"
	"github.com/Southclaws/storyden/app/resources/collection"
	"github.com/Southclaws/storyden/app/resources/post/reply"
	"github.com/Southclaws/storyden/app/resources/post/thread"
)

var (
	permissionCreate             = restrict.Permission{Action: ActionCreate}
	permissionRead               = restrict.Permission{Action: ActionRead}
	permissionUpdateThread       = restrict.Permission{Preset: "update_thread"}
	permissionUpdatePost         = restrict.Permission{Preset: "update_post"}
	permissionDeleteThread       = restrict.Permission{Preset: "delete_thread"}
	permissionDeletePost         = restrict.Permission{Preset: "delete_post"}
	permissionUpdateCollection   = restrict.Permission{Preset: "update_collection"}
	permissionDeleteCollection   = restrict.Permission{Preset: "delete_collection"}
	permissionSubmitToCollection = restrict.Permission{Preset: "submit_to_collection"}
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
	ResourceCollection: {
		&permissionCreate,
		&permissionRead,
		&permissionUpdateCollection,
		&permissionDeleteCollection,
		&permissionSubmitToCollection,
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

type collectionAccessCondition struct{}

func (c *collectionAccessCondition) Type() string { return "collection_access" }
func (c *collectionAccessCondition) Check(request *restrict.AccessRequest) error {
	acc := request.Subject.(*account.Account)
	col := request.Resource.(*collection.Collection)

	if col.Owner.ID == acc.ID {
		return nil
	}

	if acc.Admin {
		return nil
	}

	return restrict.NewConditionNotSatisfiedError(c, request, fault.New("Account is not the owner of the collection"))
}

var updateCollection = restrict.Permission{
	Action: ActionUpdate,
	Conditions: restrict.Conditions{
		&collectionAccessCondition{},
	},
}

type collectionSubmitCondition struct{}

func (c *collectionSubmitCondition) Type() string { return "collection_access" }
func (c *collectionSubmitCondition) Check(request *restrict.AccessRequest) error {
	acc := request.Subject.(*account.Account)
	col := request.Resource.(*collection.Collection)

	if col.Owner.ID == acc.ID {
		return nil
	}

	if acc.Admin {
		return nil
	}

	// Currently there are no rules to prevent any user from submitting anything
	// to any collection, this will change in future so for now this is a no-op.
	return nil
}

var submitToCollection = restrict.Permission{
	Action: ActionSubmit,
	Conditions: restrict.Conditions{
		&collectionSubmitCondition{},
	},
}

var deleteCollection = restrict.Permission{
	Action: ActionDelete,
	Conditions: restrict.Conditions{
		&collectionAccessCondition{},
	},
}

var defaultPolicy = &restrict.PolicyDefinition{
	Roles: restrict.Roles{
		EveryoneRole.ID: &EveryoneRole,
		OwnerRole.ID:    &OwnerRole,
	},
	PermissionPresets: restrict.PermissionPresets{
		"update_thread":        &updateThread,
		"update_post":          &updatePost,
		"delete_thread":        &deleteThread,
		"delete_post":          &deletePost,
		"update_collection":    &updateCollection,
		"submit_to_collection": &submitToCollection,
		"delete_collection":    &deleteCollection,
	},
}
