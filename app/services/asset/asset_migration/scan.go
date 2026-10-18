package asset_migration

import (
	"context"
	"path"
	"strings"

	"github.com/rs/xid"

	"github.com/Southclaws/storyden/app/resources/account"
	"github.com/Southclaws/storyden/app/resources/asset"
	migration "github.com/Southclaws/storyden/app/resources/asset/asset_migration"
	"github.com/Southclaws/storyden/app/resources/rbac"
)

func (m *Migrator) scan(ctx context.Context) error {
	if err := m.scanAvatars(ctx); err != nil {
		return err
	}

	for _, target := range []migration.Target{migration.Icon, migration.Banner} {
		if err := m.scanBranding(ctx, target); err != nil {
			return err
		}
	}

	return nil
}

func (m *Migrator) scanAvatars(ctx context.Context) error {
	files, err := m.objects.List(ctx, "avatar/")
	if err != nil {
		return err
	}

	for _, file := range files {
		name := strings.TrimPrefix(file, "avatar/")
		id, err := xid.FromString(name)
		if err != nil {
			continue
		}

		owner, err := m.accounts.GetRefByID(ctx, account.AccountID(id))
		if err != nil {
			m.logger.WarnContext(ctx, "legacy avatar owner lookup failed", "account_id", id, "error", err)
			continue
		}

		if owner.AvatarAssetID.Ok() {
			continue
		}

		payload := migration.Payload{
			Target:           migration.Avatar,
			Path:             path.Join("avatar", owner.ID.String()),
			AccountID:        owner.ID,
			AccountUpdatedAt: owner.UpdatedAt,
		}

		err = m.jobs.Ensure(ctx, xid.New(), migration.Kind, payload.Key(), payload)
		if err != nil {
			return err
		}
	}

	return nil
}

func (m *Migrator) scanBranding(ctx context.Context, target migration.Target) error {
	for _, filename := range migration.LegacyFilenames(target) {
		key := asset.BuildAssetPath(asset.NewFilepathFilename(filename))
		exists, err := m.objects.Exists(ctx, key)
		if err != nil {
			return err
		}

		if !exists {
			continue
		}

		// assets need an owner, so set the first settings/admin manager.
		owners, err := m.accounts.ListByHeldPermission(ctx, rbac.PermissionManageSettings, rbac.PermissionAdministrator)
		if err != nil {
			return err
		}

		if len(owners) == 0 {
			m.logger.WarnContext(ctx, "legacy branding import has no owner", "path", key)

			return nil
		}

		owner := owners[0]

		payload := migration.Payload{
			Target:    target,
			Path:      key,
			AccountID: owner.ID,
		}

		return m.jobs.Ensure(ctx, xid.New(), migration.Kind, payload.Key(), payload)
	}

	return nil
}
