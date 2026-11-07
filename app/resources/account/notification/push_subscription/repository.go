package push_subscription

import (
	"context"

	"entgo.io/ent/dialect/sql"
	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"
	"github.com/Southclaws/fault/ftag"
	"github.com/rs/xid"

	"github.com/Southclaws/storyden/app/resources/account"
	"github.com/Southclaws/storyden/internal/config"
	"github.com/Southclaws/storyden/internal/ent"
	entaccount "github.com/Southclaws/storyden/internal/ent/account"
	"github.com/Southclaws/storyden/internal/ent/pushsubscription"
	"github.com/Southclaws/storyden/internal/infrastructure/webpush"
)

type Subscription struct {
	ID xid.ID
	webpush.Subscription
}

type Repository struct {
	db               *ent.Client
	maxSubscriptions int
}

func New(db *ent.Client, cfg config.Config) (*Repository, error) {
	limit := cfg.WebPushMaxSubscriptions
	if !cfg.WebPushEnabled {
		limit = 5
	}

	if limit < 1 {
		return nil, fault.New("WEB_PUSH_MAX_SUBSCRIPTIONS must be positive when Web Push is enabled")
	}

	return &Repository{
		db:               db,
		maxSubscriptions: limit,
	}, nil
}

func (r *Repository) Upsert(
	ctx context.Context,
	owner account.AccountID,
	input webpush.SubscriptionInput,
) (*Subscription, error) {
	subscription, err := webpush.ParseSubscription(input)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx), ftag.With(ftag.InvalidArgument))
	}

	var row *ent.PushSubscription
	err = ent.WithTx(ctx, r.db, func(tx *ent.Tx) error {
		// Serialize registration and eviction for this account across replicas.
		updated, err := tx.Account.Update().
			Where(entaccount.IDEQ(xid.ID(owner))).
			Modify(func(update *sql.UpdateBuilder) {
				update.Set(entaccount.FieldUpdatedAt, sql.Expr(entaccount.FieldUpdatedAt))
			}).
			Save(ctx)
		if err != nil {
			return err
		}

		if updated == 0 {
			return fault.New("account not found", ftag.With(ftag.NotFound))
		}

		if err := tx.PushSubscription.Create().
			SetAccountID(xid.ID(owner)).
			SetEndpoint(input.Endpoint).
			SetP256dh(input.P256DH).
			SetAuth(input.Auth).
			OnConflictColumns(pushsubscription.FieldEndpoint).
			Ignore().
			Exec(ctx); err != nil {
			return err
		}

		existing, err := tx.PushSubscription.Query().
			Where(pushsubscription.EndpointEQ(input.Endpoint)).
			Only(ctx)
		if err != nil {
			return err
		}

		if existing.AccountID != xid.ID(owner) {
			return fault.New("push subscription belongs to another account", ftag.With(ftag.AlreadyExists))
		}

		row, err = tx.PushSubscription.UpdateOneID(existing.ID).
			SetP256dh(input.P256DH).
			SetAuth(input.Auth).
			Save(ctx)
		if err != nil {
			return err
		}

		return r.evictOldest(ctx, tx, owner, row.ID)
	})
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return &Subscription{
		ID:           row.ID,
		Subscription: *subscription,
	}, nil
}

func (r *Repository) evictOldest(
	ctx context.Context,
	tx *ent.Tx,
	owner account.AccountID,
	keep xid.ID,
) error {
	ids, err := tx.PushSubscription.Query().
		Where(
			pushsubscription.AccountIDEQ(xid.ID(owner)),
			pushsubscription.IDNEQ(keep),
		).
		Order(
			ent.Desc(pushsubscription.FieldCreatedAt),
			ent.Desc(pushsubscription.FieldID),
		).
		Offset(r.maxSubscriptions - 1).
		IDs(ctx)
	if err != nil {
		return err
	}

	if len(ids) == 0 {
		return nil
	}

	_, err = tx.PushSubscription.Delete().
		Where(
			pushsubscription.AccountIDEQ(xid.ID(owner)),
			pushsubscription.IDIn(ids...),
		).
		Exec(ctx)

	return err
}

func (r *Repository) List(ctx context.Context, owner account.AccountID) ([]Subscription, error) {
	rows, err := r.db.PushSubscription.Query().
		Where(pushsubscription.AccountIDEQ(xid.ID(owner))).
		All(ctx)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	result := make([]Subscription, len(rows))
	for i, row := range rows {
		subscription, err := webpush.ParseSubscription(webpush.SubscriptionInput{
			Endpoint: row.Endpoint,
			P256DH:   row.P256dh,
			Auth:     row.Auth,
		})
		if err != nil {
			return nil, fault.Wrap(err, fctx.With(ctx))
		}

		result[i] = Subscription{
			ID:           row.ID,
			Subscription: *subscription,
		}
	}
	return result, nil
}

func (r *Repository) Delete(ctx context.Context, owner account.AccountID, id xid.ID) error {
	count, err := r.db.PushSubscription.Delete().
		Where(
			pushsubscription.IDEQ(id),
			pushsubscription.AccountIDEQ(xid.ID(owner)),
		).
		Exec(ctx)
	if err != nil {
		return fault.Wrap(err, fctx.With(ctx))
	}

	if count == 0 {
		return fault.New("push subscription not found", fctx.With(ctx), ftag.With(ftag.NotFound))
	}
	return nil
}
