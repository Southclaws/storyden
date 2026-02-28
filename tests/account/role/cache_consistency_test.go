package role_test

import (
	"context"
	"errors"
	"strings"
	"testing"
	"time"

	"github.com/rs/xid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/fx"

	"github.com/Southclaws/storyden/app/resources/account/account_querier"
	"github.com/Southclaws/storyden/app/resources/account/account_ref"
	"github.com/Southclaws/storyden/app/resources/account/account_writer"
	"github.com/Southclaws/storyden/app/resources/account/role"
	"github.com/Southclaws/storyden/app/resources/account/role/held"
	"github.com/Southclaws/storyden/app/resources/account/role/role_assign"
	"github.com/Southclaws/storyden/app/resources/account/role/role_hydrate"
	"github.com/Southclaws/storyden/app/resources/account/role/role_repo"
	"github.com/Southclaws/storyden/app/resources/rbac"
	"github.com/Southclaws/storyden/app/resources/seed"
	"github.com/Southclaws/storyden/internal/ent"
	ent_account "github.com/Southclaws/storyden/internal/ent/account"
	ent_accountroles "github.com/Southclaws/storyden/internal/ent/accountroles"
	ent_role "github.com/Southclaws/storyden/internal/ent/role"
	"github.com/Southclaws/storyden/internal/infrastructure/cache"
	"github.com/Southclaws/storyden/internal/integration"
	"github.com/Southclaws/storyden/internal/integration/e2e"
)

func TestAccountRefCacheRefreshesAfterRoleAssignment(t *testing.T) {
	t.Parallel()

	integration.Test(t, nil, e2e.Setup(), fx.Invoke(func(
		lc fx.Lifecycle,
		root context.Context,
		aw *account_writer.Writer,
		accountQuery *account_querier.Querier,
		roleRepo *role_repo.Repository,
		assign *role_assign.Assignment,
	) {
		lc.Append(fx.StartHook(func() {
			r := require.New(t)
			a := assert.New(t)

			accountCtx, acc := e2e.WithAccount(root, aw, seed.Account_004_Loki)

			_, err := accountQuery.GetRefByID(accountCtx, acc.ID)
			r.NoError(err)

			customRole, err := roleRepo.Create(accountCtx, "cache-refresh-"+xid.New().String(), "red", rbac.PermissionList{})
			r.NoError(err)

			err = assign.UpdateRoles(accountCtx, account_ref.ID(acc.ID), role_assign.Add(customRole.ID))
			r.NoError(err)

			updated, err := accountQuery.GetRefByID(accountCtx, acc.ID)
			r.NoError(err)
			a.True(hasRoleID(updated.Roles, customRole.ID), "account cache should refresh role membership on cache hit")
		}))
	}))
}

func TestHydratorRefreshesDeletedRoleAssignments(t *testing.T) {
	t.Parallel()

	integration.Test(t, nil, e2e.Setup(), fx.Invoke(func(
		lc fx.Lifecycle,
		root context.Context,
		aw *account_writer.Writer,
		db *ent.Client,
		roleRepo *role_repo.Repository,
		assign *role_assign.Assignment,
		hydrator *role_hydrate.Hydrator,
	) {
		lc.Append(fx.StartHook(func() {
			r := require.New(t)
			a := assert.New(t)

			accountCtx, acc := e2e.WithAccount(root, aw, seed.Account_004_Loki)
			accountID := xid.ID(acc.ID)

			customRole, err := roleRepo.Create(accountCtx, "deleted-role-"+xid.New().String(), "red", rbac.PermissionList{})
			r.NoError(err)

			err = assign.UpdateRoles(accountCtx, account_ref.ID(acc.ID), role_assign.Add(customRole.ID))
			r.NoError(err)

			beforeDelete, err := db.Account.Query().Where(ent_account.ID(accountID)).Only(accountCtx)
			r.NoError(err)
			r.NoError(hydrator.HydrateRoleEdges(accountCtx, beforeDelete))

			err = roleRepo.Delete(accountCtx, customRole.ID)
			r.NoError(err)

			afterDelete, err := db.Account.Query().Where(ent_account.ID(accountID)).Only(accountCtx)
			r.NoError(err)
			r.NoError(hydrator.HydrateRoleEdges(accountCtx, afterDelete))

			idsByAccount, _, err := assign.ResolveRoleIDs(accountCtx, []xid.ID{accountID})
			r.NoError(err)

			a.NotContains(idsByAccount[accountID], xid.ID(customRole.ID), "deleted role should be removed from assignment cache")
		}))
	}))
}

func TestRoleAssignmentRollsBackOnCacheWriteFailure(t *testing.T) {
	t.Parallel()

	integration.Test(
		t,
		nil,
		e2e.Setup(),
		fx.Decorate(func(store cache.Store) cache.Store {
			return &failingCacheStore{
				base:            store,
				failSetPrefixes: []string{"account:role_ids:"},
			}
		}),
		fx.Invoke(func(
			lc fx.Lifecycle,
			root context.Context,
			db *ent.Client,
			roleRepo *role_repo.Repository,
			assign *role_assign.Assignment,
		) {
			lc.Append(fx.StartHook(func() {
				r := require.New(t)
				a := assert.New(t)

				accountID := xid.New()
				_, err := db.Account.Create().
					SetID(accountID).
					SetHandle("rollback-assign-" + accountID.String()).
					SetName("Rollback Assign").
					Save(root)
				r.NoError(err)

				customRole, err := roleRepo.Create(root, "rollback-assign-"+xid.New().String(), "red", rbac.PermissionList{})
				r.NoError(err)

				err = assign.UpdateRoles(root, account_ref.ID(accountID), role_assign.Add(customRole.ID))
				r.Error(err)

				exists, err := db.AccountRoles.Query().
					Where(
						ent_accountroles.AccountIDEQ(accountID),
						ent_accountroles.RoleIDEQ(xid.ID(customRole.ID)),
					).
					Exist(root)
				r.NoError(err)
				a.False(exists, "assignment mutation should rollback when cache write fails")
			}))
		}),
	)
}

func TestRoleDeleteRollsBackOnCacheInvalidationFailure(t *testing.T) {
	t.Parallel()

	integration.Test(
		t,
		nil,
		e2e.Setup(),
		fx.Decorate(func(store cache.Store) cache.Store {
			return &failingCacheStore{
				base:            store,
				failSetPrefixes: []string{"role:list"},
			}
		}),
		fx.Invoke(func(
			lc fx.Lifecycle,
			root context.Context,
			roleRepo *role_repo.Repository,
			db *ent.Client,
		) {
			lc.Append(fx.StartHook(func() {
				r := require.New(t)
				a := assert.New(t)

				customRole, err := roleRepo.Create(root, "rollback-delete-"+xid.New().String(), "red", rbac.PermissionList{})
				r.NoError(err)

				err = roleRepo.Delete(root, customRole.ID)
				r.Error(err)

				exists, err := db.Role.Query().Where(ent_role.ID(xid.ID(customRole.ID))).Exist(root)
				r.NoError(err)
				a.True(exists, "role delete should rollback when cache invalidation fails")
			}))
		}),
	)
}

func hasRoleID(roles held.Roles, id role.RoleID) bool {
	for _, heldRole := range roles {
		if heldRole.ID == id {
			return true
		}
	}

	return false
}

type failingCacheStore struct {
	base cache.Store

	failSetPrefixes    []string
	failDeletePrefixes []string
}

func (s *failingCacheStore) Get(ctx context.Context, key string) (string, error) {
	return s.base.Get(ctx, key)
}

func (s *failingCacheStore) Set(ctx context.Context, key string, object string, ttl time.Duration) error {
	if hasPrefix(key, s.failSetPrefixes) {
		return errors.New("forced cache set failure")
	}

	return s.base.Set(ctx, key, object, ttl)
}

func (s *failingCacheStore) Delete(ctx context.Context, key string) error {
	if hasPrefix(key, s.failDeletePrefixes) {
		return errors.New("forced cache delete failure")
	}

	return s.base.Delete(ctx, key)
}

func (s *failingCacheStore) HIncrBy(ctx context.Context, key string, field string, incr int64) (int, error) {
	return s.base.HIncrBy(ctx, key, field, incr)
}

func (s *failingCacheStore) HGetAll(ctx context.Context, key string) (map[string]string, error) {
	return s.base.HGetAll(ctx, key)
}

func (s *failingCacheStore) HDel(ctx context.Context, key string, field string) error {
	return s.base.HDel(ctx, key, field)
}

func (s *failingCacheStore) Expire(ctx context.Context, key string, expiration time.Duration) error {
	return s.base.Expire(ctx, key, expiration)
}

func hasPrefix(key string, prefixes []string) bool {
	for _, prefix := range prefixes {
		if strings.HasPrefix(key, prefix) {
			return true
		}
	}

	return false
}
