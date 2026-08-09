package oauth_querier

import (
	"context"

	"github.com/Southclaws/dt"
	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"
	"github.com/Southclaws/fault/ftag"
	"github.com/rs/xid"

	"github.com/Southclaws/storyden/app/resources/account"
	"github.com/Southclaws/storyden/app/resources/oauth"
	"github.com/Southclaws/storyden/app/resources/pagination"
	"github.com/Southclaws/storyden/internal/ent"
	"github.com/Southclaws/storyden/internal/ent/oauthauthorisationcode"
	"github.com/Southclaws/storyden/internal/ent/oauthauthorisationrequest"
	"github.com/Southclaws/storyden/internal/ent/oauthclient"
	"github.com/Southclaws/storyden/internal/ent/oauthdeviceauthorisation"
	"github.com/Southclaws/storyden/internal/ent/oauthrefreshtoken"
)

type Querier struct {
	db *ent.Client
}

func New(db *ent.Client) *Querier {
	return &Querier{db: db}
}

func (q *Querier) GetClient(ctx context.Context, id oauth.ClientID) (*oauth.Client, error) {
	row, err := q.db.OAuthClient.Get(ctx, id.XID())
	if err != nil {
		return nil, wrapReadError(ctx, err)
	}

	return oauth.MapClient(row), nil
}

func (q *Querier) GetClientByClientID(ctx context.Context, clientID string) (*oauth.Client, error) {
	row, err := q.db.OAuthClient.Query().
		Where(oauthclient.ClientID(clientID)).
		Only(ctx)
	if err != nil {
		return nil, wrapReadError(ctx, err)
	}

	return oauth.MapClient(row), nil
}

func (q *Querier) ListClients(ctx context.Context) ([]*oauth.Client, error) {
	rows, err := q.db.OAuthClient.Query().All(ctx)
	if err != nil {
		return nil, wrapReadError(ctx, err)
	}

	return dt.Map(rows, oauth.MapClient), nil
}

func (q *Querier) ListClientsByAccount(ctx context.Context, accountID account.AccountID) ([]*oauth.Client, error) {
	rows, err := q.db.OAuthClient.Query().
		Where(oauthclient.AccountID(xid.ID(accountID))).
		All(ctx)
	if err != nil {
		return nil, wrapReadError(ctx, err)
	}

	return dt.Map(rows, oauth.MapClient), nil
}

func (q *Querier) ListDeviceAuthorisations(ctx context.Context) ([]*oauth.DeviceAuthorisation, error) {
	rows, err := q.db.OAuthDeviceAuthorisation.Query().All(ctx)
	if err != nil {
		return nil, wrapReadError(ctx, err)
	}

	return dt.Map(rows, oauth.MapDeviceAuthorisation), nil
}

func (q *Querier) GetAuthorisationCodeByCodeHash(ctx context.Context, codeHash string) (*oauth.AuthorisationCode, error) {
	row, err := q.db.OAuthAuthorisationCode.Query().
		Where(oauthauthorisationcode.CodeHash(codeHash)).
		Only(ctx)
	if err != nil {
		return nil, wrapReadError(ctx, err)
	}

	return oauth.MapAuthorisationCode(row), nil
}

func (q *Querier) GetAuthorisationRequestByRequestIDHash(ctx context.Context, requestIDHash string) (*oauth.AuthorisationRequest, error) {
	row, err := q.db.OAuthAuthorisationRequest.Query().
		Where(oauthauthorisationrequest.RequestIDHash(requestIDHash)).
		Only(ctx)
	if err != nil {
		return nil, wrapReadError(ctx, err)
	}

	return oauth.MapAuthorisationRequest(row), nil
}

func (q *Querier) GetDeviceAuthorisationByDeviceCodeHash(ctx context.Context, deviceCodeHash string) (*oauth.DeviceAuthorisation, error) {
	row, err := q.db.OAuthDeviceAuthorisation.Query().
		Where(oauthdeviceauthorisation.DeviceCodeHash(deviceCodeHash)).
		Only(ctx)
	if err != nil {
		return nil, wrapReadError(ctx, err)
	}

	return oauth.MapDeviceAuthorisation(row), nil
}

func (q *Querier) GetDeviceAuthorisationByUserCodeHash(ctx context.Context, userCodeHash string) (*oauth.DeviceAuthorisation, error) {
	row, err := q.db.OAuthDeviceAuthorisation.Query().
		Where(oauthdeviceauthorisation.UserCodeHash(userCodeHash)).
		Only(ctx)
	if err != nil {
		return nil, wrapReadError(ctx, err)
	}

	return oauth.MapDeviceAuthorisation(row), nil
}

func (q *Querier) GetRefreshToken(ctx context.Context, id oauth.RefreshTokenID) (*oauth.RefreshToken, error) {
	row, err := q.db.OAuthRefreshToken.Get(ctx, id.XID())
	if err != nil {
		return nil, wrapReadError(ctx, err)
	}

	return oauth.MapRefreshToken(row), nil
}

func (q *Querier) GetRefreshTokenByTokenHash(ctx context.Context, tokenHash string) (*oauth.RefreshToken, error) {
	row, err := q.db.OAuthRefreshToken.Query().
		Where(oauthrefreshtoken.TokenHash(tokenHash)).
		Only(ctx)
	if err != nil {
		return nil, wrapReadError(ctx, err)
	}

	return oauth.MapRefreshToken(row), nil
}

func (q *Querier) ListRefreshTokens(ctx context.Context, params pagination.Parameters) (*pagination.Result[*oauth.RefreshToken], error) {
	return q.pageRefreshTokens(ctx, params, q.db.OAuthRefreshToken.Query())
}

func (q *Querier) ListRefreshTokensByAccount(ctx context.Context, accountID account.AccountID, params pagination.Parameters) (*pagination.Result[*oauth.RefreshToken], error) {
	return q.pageRefreshTokens(ctx, params, q.db.OAuthRefreshToken.Query().
		Where(oauthrefreshtoken.AccountID(xid.ID(accountID))))
}

func (q *Querier) pageRefreshTokens(ctx context.Context, params pagination.Parameters, query *ent.OAuthRefreshTokenQuery) (*pagination.Result[*oauth.RefreshToken], error) {
	total, err := query.Clone().Count(ctx)
	if err != nil {
		return nil, wrapReadError(ctx, err)
	}

	rows, err := query.
		WithClient().
		// id breaks ties, created_at alone leaves rows sharing a timestamp free
		// to move between adjacent offset pages
		Order(
			ent.Desc(oauthrefreshtoken.FieldCreatedAt),
			ent.Desc(oauthrefreshtoken.FieldID),
		).
		Limit(params.Limit()).
		Offset(params.Offset()).
		All(ctx)
	if err != nil {
		return nil, wrapReadError(ctx, err)
	}

	result := pagination.NewPageResult(params, total, dt.Map(rows, oauth.MapRefreshToken))

	return &result, nil
}

func wrapReadError(ctx context.Context, err error) error {
	if ent.IsNotFound(err) {
		return fault.Wrap(err, fctx.With(ctx), ftag.With(ftag.NotFound))
	}

	return fault.Wrap(err, fctx.With(ctx), ftag.With(ftag.Internal))
}
