package oauth_querier

import (
	"context"

	"github.com/Southclaws/dt"
	"github.com/Southclaws/storyden/app/resources/oauth"
	"github.com/Southclaws/storyden/app/resources/pagination"
	"github.com/Southclaws/storyden/internal/ent/oauthdynamicregistrationaccesstokens"
)

func (q *Querier) GetDynamicRegistrationAccessToken(ctx context.Context, id oauth.DynamicRegistrationAccessTokenID) (*oauth.DynamicRegistrationAccessToken, error) {
	row, err := q.db.OAuthDynamicRegistrationAccessTokens.Get(ctx, id.XID())
	if err != nil {
		return nil, wrapReadError(ctx, err)
	}

	return oauth.MapDynamicRegistrationAccessToken(row), nil
}

func (q *Querier) GetDynamicRegistrationAccessTokenByHash(ctx context.Context, hash string) (*oauth.DynamicRegistrationAccessToken, error) {
	row, err := q.db.OAuthDynamicRegistrationAccessTokens.Query().Where(oauthdynamicregistrationaccesstokens.TokenHash(hash)).Only(ctx)
	if err != nil {
		return nil, wrapReadError(ctx, err)
	}

	return oauth.MapDynamicRegistrationAccessToken(row), nil
}

func (q *Querier) ListDynamicRegistrationAccessTokens(ctx context.Context, params pagination.Parameters) (*pagination.Result[*oauth.DynamicRegistrationAccessToken], error) {
	query := q.db.OAuthDynamicRegistrationAccessTokens.Query()
	total, err := query.Clone().Count(ctx)
	if err != nil {
		return nil, wrapReadError(ctx, err)
	}

	rows, err := query.Order(oauthdynamicregistrationaccesstokens.ByCreatedAt(), oauthdynamicregistrationaccesstokens.ByID()).Limit(params.Limit()).Offset(params.Offset()).All(ctx)
	if err != nil {
		return nil, wrapReadError(ctx, err)
	}

	result := pagination.NewPageResult(params, total, dt.Map(rows, oauth.MapDynamicRegistrationAccessToken))
	return &result, nil
}
