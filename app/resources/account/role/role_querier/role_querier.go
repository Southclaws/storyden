package role_querier

import (
	"context"

	"github.com/Southclaws/storyden/app/resources/account/role"
	"github.com/Southclaws/storyden/app/resources/account/role/role_repo"
)

type Querier struct {
	repo *role_repo.Repository
}

func New(repo *role_repo.Repository) *Querier {
	return &Querier{repo: repo}
}

func (q *Querier) Get(ctx context.Context, id role.RoleID) (*role.Role, error) {
	return q.repo.Get(ctx, id)
}

func (q *Querier) GetMany(ctx context.Context, ids ...role.RoleID) (map[role.RoleID]*role.Role, error) {
	return q.repo.GetMany(ctx, ids...)
}

func (q *Querier) List(ctx context.Context) (role.Roles, error) {
	return q.repo.List(ctx)
}

func (q *Querier) GetMemberRole(ctx context.Context) (*role.Role, error) {
	return q.repo.GetMemberRole(ctx)
}

func (q *Querier) GetGuestRole(ctx context.Context) (*role.Role, error) {
	return q.repo.GetGuestRole(ctx)
}

func (q *Querier) GetAdminRole(ctx context.Context) (*role.Role, error) {
	return q.repo.GetAdminRole(ctx)
}
