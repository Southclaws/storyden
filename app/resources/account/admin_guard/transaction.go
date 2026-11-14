package admin_guard

import (
	"context"

	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqljson"
	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"
	"github.com/Southclaws/fault/fmsg"
	"github.com/Southclaws/fault/ftag"
	"github.com/rs/xid"

	"github.com/Southclaws/storyden/app/resources/account/role"
	"github.com/Southclaws/storyden/app/resources/rbac"
	"github.com/Southclaws/storyden/internal/ent"
	ent_account "github.com/Southclaws/storyden/internal/ent/account"
	"github.com/Southclaws/storyden/internal/ent/predicate"
	ent_role "github.com/Southclaws/storyden/internal/ent/role"
	ent_setting "github.com/Southclaws/storyden/internal/ent/setting"
)

var ErrLastAdministrator = fault.New("cannot remove the last administrator",
	ftag.With(ftag.InvalidArgument),
	fmsg.With("At least one active account must retain administrator access."))

type Transaction struct {
	*ent.Tx
	hadAdministrator bool
}

func Begin(ctx context.Context, db *ent.Client) (*Transaction, error) {
	tx, err := db.Tx(ctx)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	err = tx.Setting.Create().
		SetID("administrator_guard").
		SetValue("").
		OnConflictColumns(ent_setting.FieldID).
		UpdateNewValues().
		Exec(ctx)
	if err != nil {
		_ = tx.Rollback()
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	hadAdministrator, err := hasAdministrator(ctx, tx.Client())
	if err != nil {
		_ = tx.Rollback()
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return &Transaction{Tx: tx, hadAdministrator: hadAdministrator}, nil
}

func (t *Transaction) Commit(ctx context.Context, beforeCommit ...func() error) error {
	if t.hadAdministrator {
		hasAdministrator, err := hasAdministrator(ctx, t.Client())
		if err != nil {
			return fault.Wrap(err, fctx.With(ctx))
		}

		if !hasAdministrator {
			return ErrLastAdministrator
		}
	}

	for _, prepare := range beforeCommit {
		if err := prepare(); err != nil {
			return err
		}
	}

	return t.Tx.Commit()
}

func hasAdministrator(ctx context.Context, db *ent.Client) (bool, error) {
	permission := predicate.Role(func(s *sql.Selector) {
		s.Where(sqljson.ValueContains(ent_role.FieldPermissions, rbac.PermissionAdministrator.String()))
	})

	memberAdministrator, err := db.Role.Query().
		Where(ent_role.ID(xid.ID(role.DefaultRoleMemberID)), permission).
		Exist(ctx)
	if err != nil {
		return false, err
	}

	accounts := db.Account.Query().Where(ent_account.DeletedAtIsNil())
	if !memberAdministrator {
		accounts.Where(ent_account.Or(
			ent_account.Admin(true),
			ent_account.HasRolesWith(permission),
		))
	}

	return accounts.Exist(ctx)
}
