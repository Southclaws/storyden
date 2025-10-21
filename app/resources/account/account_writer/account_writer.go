package account_writer

import "github.com/Southclaws/storyden/app/resources/account/account_repo"

type Writer struct {
	*account_repo.Repository
}

type (
	Option   = account_repo.Option
	Mutation = account_repo.Mutation
)

var (
	WithID        = account_repo.WithID
	WithAdmin     = account_repo.WithAdmin
	WithName      = account_repo.WithName
	WithKind      = account_repo.WithKind
	WithBio       = account_repo.WithBio
	WithSignature = account_repo.WithSignature
	WithInvitedBy = account_repo.WithInvitedBy

	SetHandle         = account_repo.SetHandle
	SetName           = account_repo.SetName
	SetBio            = account_repo.SetBio
	SetSignature      = account_repo.SetSignature
	SetAdmin          = account_repo.SetAdmin
	SetVerifiedStatus = account_repo.SetVerifiedStatus
	SetInterests      = account_repo.SetInterests
	SetLinks          = account_repo.SetLinks
	SetMetadata       = account_repo.SetMetadata
	SetDeleted        = account_repo.SetDeleted
)

func New(repo *account_repo.Repository) *Writer {
	return &Writer{Repository: repo}
}
