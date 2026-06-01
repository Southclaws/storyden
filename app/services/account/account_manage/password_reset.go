package account_manage

import (
	"context"
	"net/mail"

	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"
	"github.com/Southclaws/fault/fmsg"
	"github.com/Southclaws/fault/ftag"
	"github.com/Southclaws/opt"
	"github.com/rs/xid"

	"github.com/Southclaws/storyden/app/resources/account"
	"github.com/Southclaws/storyden/app/resources/audit"
	"github.com/Southclaws/storyden/app/resources/datagraph"
	"github.com/Southclaws/storyden/app/resources/rbac"
	"github.com/Southclaws/storyden/app/services/authentication/provider/password/password_reset"
	"github.com/Southclaws/storyden/app/services/authentication/session"
)

func (m *Manager) GetPasswordResetToken(ctx context.Context, targetID account.AccountID, tokenProvider *password_reset.TokenProvider) (string, error) {
	if err := session.Authorise(ctx, nil, rbac.PermissionManageAccounts); err != nil {
		return "", fault.Wrap(err, fctx.With(ctx), ftag.With(ftag.PermissionDenied))
	}

	target, err := m.accountQuery.GetByID(ctx, targetID)
	if err != nil {
		return "", fault.Wrap(err, fctx.With(ctx))
	}

	token, err := tokenProvider.GetResetToken(ctx, target.ID)
	if err != nil {
		return "", fault.Wrap(err, fctx.With(ctx))
	}

	enactedBy := session.GetOptAccountID(ctx)
	_, err = m.auditWriter.Create(
		ctx,
		audit.EventTypeAccountPasswordResetTokenIssued,
		enactedBy,
		opt.New(datagraph.Ref{
			ID:   xid.ID(target.ID),
			Kind: datagraph.KindProfile,
		}),
		map[string]any{
			"account_id": target.ID.String(),
		},
	)
	if err != nil {
		return "", fault.Wrap(err, fctx.With(ctx))
	}

	return token, nil
}

func (m *Manager) SendPasswordResetEmail(
	ctx context.Context,
	targetID account.AccountID,
	emailAddressID xid.ID,
	linkTemplate password_reset.LinkTemplate,
	emailResetter *password_reset.EmailResetter,
) error {
	if err := session.Authorise(ctx, nil, rbac.PermissionManageAccounts); err != nil {
		return fault.Wrap(err, fctx.With(ctx), ftag.With(ftag.PermissionDenied))
	}

	target, err := m.accountQuery.GetByID(ctx, targetID)
	if err != nil {
		return fault.Wrap(err, fctx.With(ctx))
	}

	var address mail.Address
	found := false
	for _, e := range target.EmailAddresses {
		if e.ID == emailAddressID {
			address = e.Email
			found = true
			break
		}
	}
	if !found {
		return fault.Wrap(
			fault.New("email address not found or does not belong to this account", ftag.With(ftag.NotFound)),
			fctx.With(ctx),
			fmsg.WithDesc("email not found", "The specified email address does not belong to this account."),
		)
	}

	if err := emailResetter.SendPasswordReset(ctx, target.ID, address, linkTemplate); err != nil {
		return fault.Wrap(err, fctx.With(ctx))
	}

	enactedBy := session.GetOptAccountID(ctx)
	_, err = m.auditWriter.Create(
		ctx,
		audit.EventTypeAccountPasswordResetEmailSent,
		enactedBy,
		opt.New(datagraph.Ref{
			ID:   xid.ID(target.ID),
			Kind: datagraph.KindProfile,
		}),
		map[string]any{
			"account_id":       target.ID.String(),
			"email_address_id": emailAddressID.String(),
		},
	)
	if err != nil {
		return fault.Wrap(err, fctx.With(ctx))
	}

	return nil
}
