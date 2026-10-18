package account_repo

import (
	"context"
	"time"

	"entgo.io/ent/dialect/sql"
	"github.com/rs/xid"

	"github.com/Southclaws/storyden/app/resources/asset"
	"github.com/Southclaws/storyden/internal/ent"
	ent_account "github.com/Southclaws/storyden/internal/ent/account"
)

func SetAvatarAssetIfUnchanged(
	ctx context.Context,
	tx *ent.Tx,
	id xid.ID,
	updatedAt time.Time,
	assetID asset.AssetID,
) (bool, error) {
	// Lock without rewriting the legacy timestamp before comparing its decoded value.
	updated, err := tx.Account.Update().
		Where(
			ent_account.ID(id),
			ent_account.AvatarAssetIDIsNil(),
			ent_account.DeletedAtIsNil(),
		).
		Modify(func(update *sql.UpdateBuilder) {
			update.Set(ent_account.FieldUpdatedAt, sql.Expr(ent_account.FieldUpdatedAt))
		}).
		Save(ctx)
	if err != nil || updated == 0 {
		return false, err
	}

	current, err := tx.Account.Query().
		Where(ent_account.ID(id)).
		Select(ent_account.FieldUpdatedAt).
		Only(ctx)
	if err != nil {
		return false, err
	}
	if !current.UpdatedAt.Equal(updatedAt) {
		return false, nil
	}

	err = tx.Account.UpdateOneID(id).SetAvatarAssetID(assetID).Exec(ctx)

	return err == nil, err
}
