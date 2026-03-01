package role_writer

import "github.com/Southclaws/storyden/app/resources/account/role/role_repo"

var (
	ErrWritePermissionsNotAllowed = role_repo.ErrWritePermissionsNotAllowed
	ErrAdminPermissionsNotAllowed = role_repo.ErrAdminPermissionsNotAllowed
)

type Writer struct {
	*role_repo.Repository
}

type Mutation = role_repo.Mutation

var (
	WithName        = role_repo.WithName
	WithColour      = role_repo.WithColour
	WithPermissions = role_repo.WithPermissions
	WithMeta        = role_repo.WithMeta
)

func New(repo *role_repo.Repository) *Writer {
	return &Writer{Repository: repo}
}
