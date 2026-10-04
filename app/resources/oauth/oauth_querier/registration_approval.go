package oauth_querier

import (
	"context"
	"time"

	"github.com/Southclaws/dt"
	"github.com/Southclaws/storyden/app/resources/pagination"

	"github.com/Southclaws/storyden/app/resources/oauth"
	"github.com/Southclaws/storyden/internal/ent/oauthregistrationapproval"
)

func (q *Querier) GetRegistrationApprovalByRegistrationCodeHash(ctx context.Context, hash string) (*oauth.RegistrationApproval, error) {
	row, err := q.db.OAuthRegistrationApproval.Query().Where(oauthregistrationapproval.RegistrationCodeHash(hash)).Only(ctx)
	if err != nil {
		return nil, wrapReadError(ctx, err)
	}

	return oauth.MapRegistrationApproval(row), nil
}

func (q *Querier) GetRegistrationApprovalByVerificationCodeHash(ctx context.Context, hash string) (*oauth.RegistrationApproval, error) {
	row, err := q.db.OAuthRegistrationApproval.Query().Where(oauthregistrationapproval.VerificationCodeHash(hash)).Only(ctx)
	if err != nil {
		return nil, wrapReadError(ctx, err)
	}

	return oauth.MapRegistrationApproval(row), nil
}

func (q *Querier) ListPendingRegistrationApprovals(ctx context.Context, params pagination.Parameters) (*pagination.Result[*oauth.RegistrationApproval], error) {
	query := q.db.OAuthRegistrationApproval.Query().Where(
		oauthregistrationapproval.ExpiresAtGT(time.Now()),
		oauthregistrationapproval.ApprovedAtIsNil(),
		oauthregistrationapproval.DeniedAtIsNil(),
		oauthregistrationapproval.CancelledAtIsNil(),
		oauthregistrationapproval.ConsumedAtIsNil(),
	)
	total, err := query.Clone().Count(ctx)
	if err != nil {
		return nil, wrapReadError(ctx, err)
	}

	rows, err := query.Order(oauthregistrationapproval.ByCreatedAt(), oauthregistrationapproval.ByID()).Limit(params.Limit()).Offset(params.Offset()).All(ctx)
	if err != nil {
		return nil, wrapReadError(ctx, err)
	}

	result := pagination.NewPageResult(params, total, dt.Map(rows, oauth.MapRegistrationApproval))
	return &result, nil
}
