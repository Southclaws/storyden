package oauth_writer

import (
	"context"
	"time"

	"entgo.io/ent/dialect/sql"
	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/ftag"
	"github.com/rs/xid"

	"github.com/Southclaws/storyden/app/resources/account"
	"github.com/Southclaws/storyden/app/resources/oauth"
	"github.com/Southclaws/storyden/internal/ent"
	"github.com/Southclaws/storyden/internal/ent/oauthdynamicregistrationaccesstokens"
)

var ErrDynamicRegistrationAccessTokenUnavailable = fault.New("initial Access Token is invalid, expired, revoked, or exhausted", ftag.With(ftag.PermissionDenied))

type DynamicRegistrationAccessTokenCreate struct {
	Label            string
	TokenHash        string
	CreatorAccountID account.AccountID
	ExpiresAt        time.Time
	MaxRegistrations int
}

func (w *Writer) CreateDynamicRegistrationAccessToken(ctx context.Context, input DynamicRegistrationAccessTokenCreate) (*oauth.DynamicRegistrationAccessToken, error) {
	row, err := w.db.OAuthDynamicRegistrationAccessTokens.Create().
		SetLabel(input.Label).
		SetTokenHash(input.TokenHash).
		SetCreatorAccountID(xid.ID(input.CreatorAccountID)).
		SetExpiresAt(input.ExpiresAt.UTC()).
		SetMaxRegistrations(input.MaxRegistrations).
		Save(ctx)
	if err != nil {
		return nil, wrapWriteError(ctx, err)
	}

	return oauth.MapDynamicRegistrationAccessToken(row), nil
}

func (w *Writer) RevokeDynamicRegistrationAccessToken(ctx context.Context, id oauth.DynamicRegistrationAccessTokenID, now time.Time) (bool, error) {
	updated, err := w.db.OAuthDynamicRegistrationAccessTokens.Update().
		Where(oauthdynamicregistrationaccesstokens.ID(id.XID()), oauthdynamicregistrationaccesstokens.RevokedAtIsNil()).
		SetRevokedAt(now).Save(ctx)
	if err != nil {
		return false, wrapWriteError(ctx, err)
	}

	return updated > 0, nil
}

func consumeDynamicRegistrationAccessToken(ctx context.Context, tx *ent.Tx, id oauth.DynamicRegistrationAccessTokenID) error {
	n, err := tx.OAuthDynamicRegistrationAccessTokens.Update().Where(
		oauthdynamicregistrationaccesstokens.ID(id.XID()),
		oauthdynamicregistrationaccesstokens.RevokedAtIsNil(),
		oauthdynamicregistrationaccesstokens.ExpiresAtGT(time.Now().UTC()),
		func(s *sql.Selector) {
			s.Where(sql.ColumnsLT(s.C(oauthdynamicregistrationaccesstokens.FieldRegistrationCount), s.C(oauthdynamicregistrationaccesstokens.FieldMaxRegistrations)))
		},
	).AddRegistrationCount(1).Save(ctx)
	if err != nil {
		return err
	}

	if n != 1 {
		return ErrDynamicRegistrationAccessTokenUnavailable
	}

	return nil
}
