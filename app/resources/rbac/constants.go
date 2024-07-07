package rbac

type Action = string

const (
	ActionCreate Action = "create"
	ActionRead   Action = "read"
	ActionUpdate Action = "update"
	ActionSubmit Action = "submit"
	ActionDelete Action = "delete"
)

type Resource = string

const (
	ResourceAccount    Resource = "account"
	ResourceThread     Resource = "thread"
	ResourcePost       Resource = "post"
	ResourceCollection Resource = "collection"
)
