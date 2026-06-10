package oauth_writer

import (
	"context"
	"time"

	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"
	"github.com/Southclaws/fault/ftag"
	"github.com/Southclaws/opt"
	"github.com/rs/xid"

	"github.com/Southclaws/storyden/app/resources/account"
	"github.com/Southclaws/storyden/app/resources/oauth"
	"github.com/Southclaws/storyden/internal/ent"
	"github.com/Southclaws/storyden/internal/ent/oauthauthorisationcode"
	"github.com/Southclaws/storyden/internal/ent/oauthauthorisationrequest"
	"github.com/Southclaws/storyden/internal/ent/oauthclient"
	"github.com/Southclaws/storyden/internal/ent/oauthdeviceauthorisation"
	"github.com/Southclaws/storyden/internal/ent/oauthrefreshtoken"
)

type Writer struct {
	db *ent.Client
}

func New(db *ent.Client) *Writer {
	return &Writer{db: db}
}

type ClientCreate struct {
	AccountID               opt.Optional[account.AccountID]
	ClientID                string
	ClientSecretHash        opt.Optional[string]
	Name                    string
	Type                    oauth.ClientType
	ScopePolicy             opt.Optional[oauth.ScopePolicy]
	TokenEndpointAuthMethod opt.Optional[string]
	PKCERequired            opt.Optional[bool]
	RedirectURIs            []string
	AllowedScopes           []string
	AllowedGrants           []string
}

type ClientUpdate struct {
	Name             opt.Optional[string]
	ClientSecretHash opt.Optional[string]
	ScopePolicy      opt.Optional[oauth.ScopePolicy]
	RedirectURIs     opt.Optional[[]string]
	AllowedScopes    opt.Optional[[]string]
	AllowedGrants    opt.Optional[[]string]
}

type AuthorisationCodeCreate struct {
	ClientID            oauth.ClientID
	AccountID           account.AccountID
	CodeHash            string
	RedirectURI         string
	Scope               string
	CodeChallenge       string
	CodeChallengeMethod string
	ExpiresAt           time.Time
}

type AuthorisationRequestCreate struct {
	ClientID            oauth.ClientID
	AccountID           account.AccountID
	RequestIDHash       string
	RedirectURI         string
	Scope               string
	State               opt.Optional[string]
	CodeChallenge       string
	CodeChallengeMethod string
	ExpiresAt           time.Time
}

type DeviceAuthorisationCreate struct {
	ClientID            oauth.ClientID
	DeviceCodeHash      string
	UserCodeHash        string
	UserCodeDisplay     string
	Scope               string
	ExpiresAt           time.Time
	PollIntervalSeconds int
}

type RefreshTokenCreate struct {
	ClientID  oauth.ClientID
	AccountID account.AccountID
	TokenHash string
	Scope     string
	ExpiresAt time.Time
}

func (w *Writer) CreateClient(ctx context.Context, input ClientCreate) (*oauth.Client, error) {
	create := w.db.OAuthClient.Create().
		SetClientID(input.ClientID).
		SetName(input.Name).
		SetType(oauthclient.Type(input.Type.String())).
		SetRedirectUris(input.RedirectURIs).
		SetAllowedScopes(input.AllowedScopes).
		SetAllowedGrants(input.AllowedGrants)

	input.AccountID.Call(func(accountID account.AccountID) {
		create.SetAccountID(xid.ID(accountID))
	})
	input.ClientSecretHash.Call(func(hash string) {
		create.SetClientSecretHash(hash)
	})
	input.ScopePolicy.Call(func(policy oauth.ScopePolicy) {
		create.SetScopePolicy(oauthclient.ScopePolicy(policy.String()))
	})
	input.TokenEndpointAuthMethod.Call(func(method string) {
		create.SetTokenEndpointAuthMethod(method)
	})
	input.PKCERequired.Call(func(required bool) {
		create.SetPkceRequired(required)
	})

	row, err := create.Save(ctx)
	if err != nil {
		return nil, wrapWriteError(ctx, err)
	}

	return oauth.MapClient(row), nil
}

func (w *Writer) UpdateClient(ctx context.Context, id oauth.ClientID, input ClientUpdate) (*oauth.Client, error) {
	update := w.db.OAuthClient.UpdateOneID(id.XID())

	input.Name.Call(func(name string) {
		update.SetName(name)
	})
	input.ClientSecretHash.Call(func(hash string) {
		update.SetClientSecretHash(hash)
	})
	input.ScopePolicy.Call(func(policy oauth.ScopePolicy) {
		update.SetScopePolicy(oauthclient.ScopePolicy(policy.String()))
	})
	input.RedirectURIs.Call(func(uris []string) {
		update.SetRedirectUris(uris)
	})
	input.AllowedScopes.Call(func(scopes []string) {
		update.SetAllowedScopes(scopes)
	})
	input.AllowedGrants.Call(func(grants []string) {
		update.SetAllowedGrants(grants)
	})

	row, err := update.Save(ctx)
	if err != nil {
		return nil, wrapWriteError(ctx, err)
	}

	return oauth.MapClient(row), nil
}

func (w *Writer) DeleteClient(ctx context.Context, id oauth.ClientID) error {
	err := ent.WithTx(ctx, w.db, func(tx *ent.Tx) error {
		clientID := id.XID()

		if _, err := tx.OAuthAuthorisationCode.Delete().
			Where(oauthauthorisationcode.ClientID(clientID)).
			Exec(ctx); err != nil {
			return err
		}

		if _, err := tx.OAuthAuthorisationRequest.Delete().
			Where(oauthauthorisationrequest.ClientID(clientID)).
			Exec(ctx); err != nil {
			return err
		}

		if _, err := tx.OAuthDeviceAuthorisation.Delete().
			Where(oauthdeviceauthorisation.ClientID(clientID)).
			Exec(ctx); err != nil {
			return err
		}

		if _, err := tx.OAuthRefreshToken.Delete().
			Where(oauthrefreshtoken.ClientID(clientID)).
			Exec(ctx); err != nil {
			return err
		}

		return tx.OAuthClient.DeleteOneID(clientID).Exec(ctx)
	})
	if err != nil {
		return wrapWriteError(ctx, err)
	}

	return nil
}

func (w *Writer) CreateAuthorisationCode(ctx context.Context, input AuthorisationCodeCreate) (*oauth.AuthorisationCode, error) {
	row, err := w.db.OAuthAuthorisationCode.Create().
		SetClientID(input.ClientID.XID()).
		SetAccountID(xid.ID(input.AccountID)).
		SetCodeHash(input.CodeHash).
		SetRedirectURI(input.RedirectURI).
		SetScope(input.Scope).
		SetCodeChallenge(input.CodeChallenge).
		SetCodeChallengeMethod(oauthauthorisationcode.CodeChallengeMethod(input.CodeChallengeMethod)).
		SetExpiresAt(input.ExpiresAt).
		Save(ctx)
	if err != nil {
		return nil, wrapWriteError(ctx, err)
	}

	return oauth.MapAuthorisationCode(row), nil
}

func (w *Writer) CreateAuthorisationRequest(ctx context.Context, input AuthorisationRequestCreate) (*oauth.AuthorisationRequest, error) {
	create := w.db.OAuthAuthorisationRequest.Create().
		SetClientID(input.ClientID.XID()).
		SetAccountID(xid.ID(input.AccountID)).
		SetRequestIDHash(input.RequestIDHash).
		SetRedirectURI(input.RedirectURI).
		SetScope(input.Scope).
		SetCodeChallenge(input.CodeChallenge).
		SetCodeChallengeMethod(oauthauthorisationrequest.CodeChallengeMethod(input.CodeChallengeMethod)).
		SetExpiresAt(input.ExpiresAt)

	input.State.Call(func(state string) {
		create.SetState(state)
	})

	row, err := create.Save(ctx)
	if err != nil {
		return nil, wrapWriteError(ctx, err)
	}

	return oauth.MapAuthorisationRequest(row), nil
}

func (w *Writer) ApproveAuthorisationRequestAndCreateCode(ctx context.Context, id oauth.AuthorisationRequestID, code AuthorisationCodeCreate, approvedAt time.Time) (bool, error) {
	var approved bool

	err := ent.WithTx(ctx, w.db, func(tx *ent.Tx) error {
		updated, err := tx.OAuthAuthorisationRequest.Update().
			Where(
				oauthauthorisationrequest.IDEQ(id.XID()),
				oauthauthorisationrequest.ApprovedAtIsNil(),
				oauthauthorisationrequest.DeniedAtIsNil(),
			).
			SetApprovedAt(approvedAt).
			Save(ctx)
		if err != nil {
			return err
		}
		if updated == 0 {
			approved = false
			return nil
		}

		_, err = tx.OAuthAuthorisationCode.Create().
			SetClientID(code.ClientID.XID()).
			SetAccountID(xid.ID(code.AccountID)).
			SetCodeHash(code.CodeHash).
			SetRedirectURI(code.RedirectURI).
			SetScope(code.Scope).
			SetCodeChallenge(code.CodeChallenge).
			SetCodeChallengeMethod(oauthauthorisationcode.CodeChallengeMethod(code.CodeChallengeMethod)).
			SetExpiresAt(code.ExpiresAt).
			Save(ctx)
		if err != nil {
			return err
		}

		approved = true
		return nil
	})
	if err != nil {
		return false, wrapWriteError(ctx, err)
	}

	return approved, nil
}

func (w *Writer) DenyAuthorisationRequest(ctx context.Context, id oauth.AuthorisationRequestID, deniedAt time.Time) (bool, error) {
	updated, err := w.db.OAuthAuthorisationRequest.Update().
		Where(
			oauthauthorisationrequest.IDEQ(id.XID()),
			oauthauthorisationrequest.ApprovedAtIsNil(),
			oauthauthorisationrequest.DeniedAtIsNil(),
		).
		SetDeniedAt(deniedAt).
		Save(ctx)
	if err != nil {
		return false, wrapWriteError(ctx, err)
	}

	return updated > 0, nil
}

func (w *Writer) ConsumeAuthorisationCode(ctx context.Context, id oauth.AuthorisationCodeID, consumedAt time.Time) (bool, error) {
	updated, err := w.db.OAuthAuthorisationCode.Update().
		Where(
			oauthauthorisationcode.IDEQ(id.XID()),
			oauthauthorisationcode.ConsumedAtIsNil(),
		).
		SetConsumedAt(consumedAt).
		Save(ctx)
	if err != nil {
		return false, wrapWriteError(ctx, err)
	}

	return updated > 0, nil
}

func (w *Writer) CreateDeviceAuthorisation(ctx context.Context, input DeviceAuthorisationCreate) (*oauth.DeviceAuthorisation, error) {
	row, err := w.db.OAuthDeviceAuthorisation.Create().
		SetClientID(input.ClientID.XID()).
		SetDeviceCodeHash(input.DeviceCodeHash).
		SetUserCodeHash(input.UserCodeHash).
		SetUserCodeDisplay(input.UserCodeDisplay).
		SetScope(input.Scope).
		SetExpiresAt(input.ExpiresAt).
		SetPollIntervalSeconds(input.PollIntervalSeconds).
		Save(ctx)
	if err != nil {
		return nil, wrapWriteError(ctx, err)
	}

	return oauth.MapDeviceAuthorisation(row), nil
}

func (w *Writer) RecordDeviceAuthorisationPoll(ctx context.Context, id oauth.DeviceAuthorisationID, polledAt time.Time, pollIntervalSeconds int) error {
	_, err := w.db.OAuthDeviceAuthorisation.UpdateOneID(id.XID()).
		SetLastPolledAt(polledAt).
		SetPollIntervalSeconds(pollIntervalSeconds).
		Save(ctx)
	if err != nil {
		return wrapWriteError(ctx, err)
	}

	return nil
}

func (w *Writer) ClaimDeviceAuthorisation(ctx context.Context, id oauth.DeviceAuthorisationID, accountID account.AccountID) (bool, error) {
	updated, err := w.db.OAuthDeviceAuthorisation.Update().
		Where(
			oauthdeviceauthorisation.IDEQ(id.XID()),
			oauthdeviceauthorisation.ApprovedAtIsNil(),
			oauthdeviceauthorisation.DeniedAtIsNil(),
			oauthdeviceauthorisation.ConsumedAtIsNil(),
		).
		Where(
			oauthdeviceauthorisation.Or(
				oauthdeviceauthorisation.ClaimedByAccountIDIsNil(),
				oauthdeviceauthorisation.ClaimedByAccountID(xid.ID(accountID)),
			),
		).
		SetClaimedByAccountID(xid.ID(accountID)).
		Save(ctx)
	if err != nil {
		return false, wrapWriteError(ctx, err)
	}

	return updated > 0, nil
}

func (w *Writer) ApproveDeviceAuthorisation(ctx context.Context, id oauth.DeviceAuthorisationID, accountID account.AccountID, scope string, approvedAt time.Time) (bool, error) {
	updated, err := w.db.OAuthDeviceAuthorisation.Update().
		Where(
			oauthdeviceauthorisation.IDEQ(id.XID()),
			oauthdeviceauthorisation.ApprovedAtIsNil(),
			oauthdeviceauthorisation.DeniedAtIsNil(),
			oauthdeviceauthorisation.ConsumedAtIsNil(),
		).
		SetApprovedByAccountID(xid.ID(accountID)).
		SetScope(scope).
		SetApprovedAt(approvedAt).
		Save(ctx)
	if err != nil {
		return false, wrapWriteError(ctx, err)
	}

	return updated > 0, nil
}

func (w *Writer) DenyDeviceAuthorisation(ctx context.Context, id oauth.DeviceAuthorisationID, deniedAt time.Time) (bool, error) {
	updated, err := w.db.OAuthDeviceAuthorisation.Update().
		Where(
			oauthdeviceauthorisation.IDEQ(id.XID()),
			oauthdeviceauthorisation.ApprovedAtIsNil(),
			oauthdeviceauthorisation.DeniedAtIsNil(),
			oauthdeviceauthorisation.ConsumedAtIsNil(),
		).
		SetDeniedAt(deniedAt).
		Save(ctx)
	if err != nil {
		return false, wrapWriteError(ctx, err)
	}

	return updated > 0, nil
}

func (w *Writer) ConsumeDeviceAuthorisation(ctx context.Context, id oauth.DeviceAuthorisationID, consumedAt time.Time) (bool, error) {
	updated, err := w.db.OAuthDeviceAuthorisation.Update().
		Where(
			oauthdeviceauthorisation.IDEQ(id.XID()),
			oauthdeviceauthorisation.ConsumedAtIsNil(),
		).
		SetConsumedAt(consumedAt).
		Save(ctx)
	if err != nil {
		return false, wrapWriteError(ctx, err)
	}

	return updated > 0, nil
}

func (w *Writer) DeleteExpiredDeviceAuthorisations(ctx context.Context, now time.Time) (int, error) {
	deleted, err := w.db.OAuthDeviceAuthorisation.Delete().
		Where(oauthdeviceauthorisation.ExpiresAtLT(now)).
		Exec(ctx)
	if err != nil {
		return 0, wrapWriteError(ctx, err)
	}

	return deleted, nil
}

func (w *Writer) DeleteExpiredAuthorisationRequests(ctx context.Context, now time.Time) (int, error) {
	deleted, err := w.db.OAuthAuthorisationRequest.Delete().
		Where(oauthauthorisationrequest.ExpiresAtLT(now)).
		Exec(ctx)
	if err != nil {
		return 0, wrapWriteError(ctx, err)
	}

	return deleted, nil
}

func (w *Writer) CreateRefreshToken(ctx context.Context, input RefreshTokenCreate) (*oauth.RefreshToken, error) {
	row, err := w.db.OAuthRefreshToken.Create().
		SetClientID(input.ClientID.XID()).
		SetAccountID(xid.ID(input.AccountID)).
		SetTokenHash(input.TokenHash).
		SetScope(input.Scope).
		SetExpiresAt(input.ExpiresAt).
		Save(ctx)
	if err != nil {
		return nil, wrapWriteError(ctx, err)
	}

	return oauth.MapRefreshToken(row), nil
}

func (w *Writer) RevokeRefreshToken(ctx context.Context, id oauth.RefreshTokenID, revokedAt time.Time, replacedBy opt.Optional[oauth.RefreshTokenID]) (bool, error) {
	update := w.db.OAuthRefreshToken.Update().
		Where(
			oauthrefreshtoken.IDEQ(id.XID()),
			oauthrefreshtoken.RevokedAtIsNil(),
		).
		SetRevokedAt(revokedAt).
		SetLastUsedAt(revokedAt)

	replacedBy.Call(func(replacement oauth.RefreshTokenID) {
		update.SetReplacedByTokenID(replacement.XID())
	})

	updated, err := update.Save(ctx)
	if err != nil {
		return false, wrapWriteError(ctx, err)
	}

	return updated > 0, nil
}

func (w *Writer) ConsumeRefreshToken(ctx context.Context, id oauth.RefreshTokenID, revokedAt time.Time) (bool, error) {
	updated, err := w.db.OAuthRefreshToken.Update().
		Where(
			oauthrefreshtoken.IDEQ(id.XID()),
			oauthrefreshtoken.RevokedAtIsNil(),
		).
		SetRevokedAt(revokedAt).
		SetLastUsedAt(revokedAt).
		Save(ctx)
	if err != nil {
		return false, wrapWriteError(ctx, err)
	}

	return updated > 0, nil
}

func (w *Writer) SetRefreshTokenReplacement(ctx context.Context, id oauth.RefreshTokenID, replacement oauth.RefreshTokenID) error {
	_, err := w.db.OAuthRefreshToken.UpdateOneID(id.XID()).
		SetReplacedByTokenID(replacement.XID()).
		Save(ctx)
	if err != nil {
		return wrapWriteError(ctx, err)
	}

	return nil
}

func wrapWriteError(ctx context.Context, err error) error {
	if ent.IsNotFound(err) {
		return fault.Wrap(err, fctx.With(ctx), ftag.With(ftag.NotFound))
	}

	return fault.Wrap(err, fctx.With(ctx), ftag.With(ftag.Internal))
}
