package remote

import (
	"context"
	"time"

	"entgo.io/ent/dialect/sql"
	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"
	"github.com/Southclaws/fault/ftag"
	"github.com/rs/xid"

	"github.com/Southclaws/storyden/app/resources/account"
	"github.com/Southclaws/storyden/internal/ent"
	ent_flow "github.com/Southclaws/storyden/internal/ent/oauthremoteauthorisationflow"
	ent_connection "github.com/Southclaws/storyden/internal/ent/oauthremoteconnection"
)

type Repository struct {
	db *ent.Client
}

func New(db *ent.Client) *Repository {
	return &Repository{db: db}
}

func (r *Repository) ListConnections(ctx context.Context) ([]Connection, error) {
	rows, err := r.db.OAuthRemoteConnection.Query().
		Order(ent_connection.ByCreatedAt(sql.OrderDesc())).
		All(ctx)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}
	out := make([]Connection, 0, len(rows))
	for _, row := range rows {
		out = append(out, MapConnection(row))
	}
	return out, nil
}

func (r *Repository) GetConnection(ctx context.Context, id ConnectionID) (Connection, error) {
	row, err := r.db.OAuthRemoteConnection.Get(ctx, xid.ID(id))
	if err != nil {
		if ent.IsNotFound(err) {
			err = fault.Wrap(err, ftag.With(ftag.NotFound))
		}
		return Connection{}, fault.Wrap(err, fctx.With(ctx))
	}
	return MapConnection(row), nil
}

func (r *Repository) CreateConnection(ctx context.Context, in ConnectionCreate) (Connection, error) {
	status := in.Status
	if status == "" {
		status = StatusPending
	}
	row, err := r.db.OAuthRemoteConnection.Create().
		SetResourceURL(in.ResourceURL).
		SetResource(in.Resource).
		SetResourceName(in.ResourceName).
		SetProtectedResourceMetadata(in.ProtectedResourceMetadata).
		SetAuthorizationServer(in.AuthorizationServer).
		SetAuthorizationServerMetadata(in.AuthorizationServerMetadata).
		SetMode(toEntMode(in.Mode)).
		SetStatus(toEntStatus(status)).
		SetClientID(in.ClientID).
		SetClientSecret(in.ClientSecret).
		SetAuthorizationEndpoint(in.AuthorizationEndpoint).
		SetTokenEndpoint(in.TokenEndpoint).
		SetRegistrationEndpoint(in.RegistrationEndpoint).
		SetTokenEndpointAuthMethod(in.TokenEndpointAuthMethod).
		SetRedirectUris(in.RedirectURIs).
		SetRedirectURI(in.RedirectURI).
		SetScope(in.Scope).
		SetAddedBy(xid.ID(in.AddedBy)).
		Save(ctx)
	if err != nil {
		if ent.IsConstraintError(err) {
			return Connection{}, fault.Wrap(err, fctx.With(ctx), ftag.With(ftag.AlreadyExists))
		}
		return Connection{}, fault.Wrap(err, fctx.With(ctx))
	}
	return MapConnection(row), nil
}

func (r *Repository) DeleteUnconnectedConnectionByIdentity(ctx context.Context, resourceURL, authorizationServer string, addedBy account.AccountID) (int, error) {
	affected, err := r.db.OAuthRemoteConnection.Delete().
		Where(
			ent_connection.ResourceURLEQ(resourceURL),
			ent_connection.AuthorizationServerEQ(authorizationServer),
			ent_connection.AddedByEQ(xid.ID(addedBy)),
			ent_connection.Or(
				ent_connection.AccessTokenIsNil(),
				ent_connection.AccessTokenEQ(""),
			),
		).
		Exec(ctx)
	if err != nil {
		return 0, fault.Wrap(err, fctx.With(ctx))
	}
	return affected, nil
}

func (r *Repository) HasConnectedConnectionByIdentity(ctx context.Context, resourceURL, authorizationServer string, addedBy account.AccountID) (bool, error) {
	exists, err := r.db.OAuthRemoteConnection.Query().
		Where(
			ent_connection.ResourceURLEQ(resourceURL),
			ent_connection.AuthorizationServerEQ(authorizationServer),
			ent_connection.AddedByEQ(xid.ID(addedBy)),
			ent_connection.AccessTokenNotNil(),
			ent_connection.AccessTokenNEQ(""),
		).
		Exist(ctx)
	if err != nil {
		return false, fault.Wrap(err, fctx.With(ctx))
	}
	return exists, nil
}

func (r *Repository) CreateFlow(ctx context.Context, connectionID ConnectionID, stateHash string, verifier string, redirectURI string, expiresAt time.Time) (Flow, error) {
	row, err := r.db.OAuthRemoteAuthorisationFlow.Create().
		SetConnectionID(xid.ID(connectionID)).
		SetStateHash(stateHash).
		SetPkceVerifier(verifier).
		SetRedirectURI(redirectURI).
		SetExpiresAt(expiresAt).
		Save(ctx)
	if err != nil {
		return Flow{}, fault.Wrap(err, fctx.With(ctx))
	}
	return MapFlow(row), nil
}

func (r *Repository) GetFlowByStateHash(ctx context.Context, stateHash string) (Flow, error) {
	row, err := r.db.OAuthRemoteAuthorisationFlow.Query().
		Where(ent_flow.StateHashEQ(stateHash)).
		WithConnection().
		Only(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			err = fault.Wrap(err, ftag.With(ftag.NotFound))
		}
		return Flow{}, fault.Wrap(err, fctx.With(ctx))
	}
	return MapFlow(row), nil
}

func (r *Repository) ConsumeFlow(ctx context.Context, id FlowID, now time.Time) error {
	err := r.db.OAuthRemoteAuthorisationFlow.UpdateOneID(xid.ID(id)).
		SetConsumedAt(now).
		Exec(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			err = fault.Wrap(err, ftag.With(ftag.NotFound))
		}
		return fault.Wrap(err, fctx.With(ctx))
	}
	return nil
}

func (r *Repository) ClaimFlow(ctx context.Context, id FlowID, now time.Time) (bool, error) {
	affected, err := r.db.OAuthRemoteAuthorisationFlow.Update().
		Where(
			ent_flow.ID(xid.ID(id)),
			ent_flow.ConsumedAtIsNil(),
			ent_flow.ExpiresAtGT(now),
		).
		SetConsumedAt(now).
		Save(ctx)
	if err != nil {
		return false, fault.Wrap(err, fctx.With(ctx))
	}
	return affected == 1, nil
}

func (r *Repository) MarkError(ctx context.Context, id ConnectionID, message string) error {
	return fault.Wrap(
		r.db.OAuthRemoteConnection.UpdateOneID(xid.ID(id)).
			SetStatus(ent_connection.StatusError).
			SetLastError(message).
			ClearTokenRefreshStartedAt().
			Exec(ctx),
		fctx.With(ctx),
	)
}

func (r *Repository) ClaimTokenRefresh(ctx context.Context, id ConnectionID, now time.Time, staleBefore time.Time) (bool, error) {
	affected, err := r.db.OAuthRemoteConnection.Update().
		Where(
			ent_connection.IDEQ(xid.ID(id)),
			ent_connection.Or(
				ent_connection.TokenRefreshStartedAtIsNil(),
				ent_connection.TokenRefreshStartedAtLT(staleBefore),
			),
		).
		SetTokenRefreshStartedAt(now).
		Save(ctx)
	if err != nil {
		return false, fault.Wrap(err, fctx.With(ctx))
	}
	return affected == 1, nil
}

func (r *Repository) ReleaseTokenRefresh(ctx context.Context, id ConnectionID) error {
	err := r.db.OAuthRemoteConnection.UpdateOneID(xid.ID(id)).
		ClearTokenRefreshStartedAt().
		Exec(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			err = fault.Wrap(err, ftag.With(ftag.NotFound))
		}
		return fault.Wrap(err, fctx.With(ctx))
	}
	return nil
}

func (r *Repository) StoreTokens(ctx context.Context, id ConnectionID, in TokenUpdate) (Connection, error) {
	update := r.db.OAuthRemoteConnection.UpdateOneID(xid.ID(id)).
		SetStatus(ent_connection.StatusConnected).
		ClearLastError().
		ClearTokenRefreshStartedAt().
		SetAccessToken(in.AccessToken).
		SetTokenType(in.TokenType)
	if in.RefreshToken != "" {
		update.SetRefreshToken(in.RefreshToken)
	}
	if in.TokenExpiry != nil {
		update.SetTokenExpiry(*in.TokenExpiry)
	} else {
		update.ClearTokenExpiry()
	}
	if in.Scope != "" {
		update.SetScope(in.Scope)
	}
	row, err := update.Save(ctx)
	if err != nil {
		return Connection{}, fault.Wrap(err, fctx.With(ctx))
	}
	return MapConnection(row), nil
}
