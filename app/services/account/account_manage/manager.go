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
	"github.com/Southclaws/storyden/app/resources/account/account_querier"
	"github.com/Southclaws/storyden/app/resources/account/account_writer"
	"github.com/Southclaws/storyden/app/resources/audit/audit_writer"
	"github.com/Southclaws/storyden/app/resources/rbac"
	"github.com/Southclaws/storyden/app/services/account/account_email"
	"github.com/Southclaws/storyden/app/services/account/account_update"
	"github.com/Southclaws/storyden/app/services/authentication/session"
)

type Manager struct {
	accountQuery *account_querier.Querier
	accountWrite *account_writer.Writer
	accountEmail *account_email.Manager
	updater      *account_update.Updater
	auditWriter  *audit_writer.Writer
}

func New(
	accountQuery *account_querier.Querier,
	accountWrite *account_writer.Writer,
	accountEmail *account_email.Manager,
	updater *account_update.Updater,
	auditWriter *audit_writer.Writer,
) *Manager {
	return &Manager{
		accountQuery: accountQuery,
		accountWrite: accountWrite,
		accountEmail: accountEmail,
		updater:      updater,
		auditWriter:  auditWriter,
	}
}

type InitialProps struct {
	Handle         string
	Name           opt.Optional[string]
	Bio            opt.Optional[string]
	Signature      opt.Optional[string]
	Interests      opt.Optional[[]xid.ID]
	Links          opt.Optional[[]account.ExternalLink]
	Admin          opt.Optional[bool]
	EmailAddress   opt.Optional[mail.Address]
	VerifiedStatus opt.Optional[account.VerifiedStatus]
	Meta           opt.Optional[map[string]any]
}

func (m *Manager) Create(ctx context.Context, props InitialProps) (*account.AccountWithEdges, error) {
	if err := session.Authorise(ctx, nil, rbac.PermissionManageAccounts); err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx), ftag.With(ftag.PermissionDenied))
	}

	if err := account.ValidateHandle(ctx, props.Handle); err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	if admin, ok := props.Admin.Get(); ok && admin {
		if err := session.Authorise(ctx, nil, rbac.PermissionAdministrator); err != nil {
			return nil, fault.Wrap(err, fctx.With(ctx), ftag.With(ftag.PermissionDenied))
		}
	}

	options := []account_writer.Option{}
	if name, ok := props.Name.Get(); ok {
		options = append(options, account_writer.WithName(name))
	}
	if bio, ok := props.Bio.Get(); ok {
		options = append(options, account_writer.WithBioString(bio))
	}
	if signature, ok := props.Signature.Get(); ok {
		options = append(options, account_writer.WithSignatureString(signature))
	}
	if interests, ok := props.Interests.Get(); ok {
		options = append(options, account_writer.WithInterests(interests))
	}
	if links, ok := props.Links.Get(); ok {
		options = append(options, account_writer.WithLinks(links))
	}
	if admin, ok := props.Admin.Get(); ok {
		options = append(options, account_writer.WithAdmin(admin))
	}
	if status, ok := props.VerifiedStatus.Get(); ok {
		options = append(options, account_writer.WithVerifiedStatus(status))
	}
	if meta, ok := props.Meta.Get(); ok {
		options = append(options, account_writer.WithMetadata(meta))
	}

	acc, err := m.accountWrite.Create(ctx, props.Handle, options...)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	if address, ok := props.EmailAddress.Get(); ok {
		if _, err := m.accountEmail.AddUnverified(ctx, acc.ID, address); err != nil {
			return nil, fault.Wrap(err, fctx.With(ctx))
		}

		acc, err = m.accountQuery.GetByID(ctx, acc.ID)
		if err != nil {
			return nil, fault.Wrap(err, fctx.With(ctx))
		}
	}

	return acc, nil
}

func (m *Manager) Update(ctx context.Context, targetID account.AccountID, params account_update.Partial) (*account.AccountWithEdges, error) {
	if err := session.Authorise(ctx, nil, rbac.PermissionManageAccounts); err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx), ftag.With(ftag.PermissionDenied))
	}

	if params.Admin.Ok() {
		if err := session.Authorise(ctx, nil, rbac.PermissionAdministrator); err != nil {
			return nil, fault.Wrap(err, fctx.With(ctx), ftag.With(ftag.PermissionDenied))
		}
	}

	return m.updater.Update(ctx, targetID, params)
}

// GetByID retrieves an account by ID with proper authorization checks.
// Supports hierarchical permissions:
// - Owner can always view their own account
// - ADMINISTRATOR can view any account (including other administrators)
// - VIEW_ACCOUNTS can view non-administrator accounts only
func (m *Manager) GetByID(ctx context.Context, targetID account.AccountID) (*account.AccountWithEdges, error) {
	// Get caller's account from session (already loaded in middleware)
	callerID, err := session.GetAccountID(ctx)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	caller, err := m.accountQuery.GetByID(ctx, callerID)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	// Get target account
	target, err := m.accountQuery.GetByID(ctx, targetID)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	// Check authorization
	if err := m.canViewAccount(ctx, caller, target); err != nil {
		return nil, err
	}

	return target, nil
}

// canViewAccount checks if the caller is authorized to view the target account.
// Returns an error if the caller does not have sufficient permissions.
func (m *Manager) canViewAccount(ctx context.Context, caller, target *account.AccountWithEdges) error {
	// Owner can always view their own account
	if caller.ID == target.ID {
		return nil
	}

	callerPerms := caller.Roles.Permissions()
	targetPerms := target.Roles.Permissions()

	// Administrator can view any account (including other administrators)
	if callerPerms.HasAny(rbac.PermissionAdministrator) {
		return nil
	}

	// VIEW_ACCOUNTS can view non-administrator accounts only
	if callerPerms.HasAny(rbac.PermissionViewAccounts) {
		if targetPerms.HasAny(rbac.PermissionAdministrator) {
			return fault.Wrap(
				fault.New("cannot view administrator account", ftag.With(ftag.PermissionDenied)),
				fctx.With(ctx),
				fmsg.WithDesc("insufficient permissions", "You do not have permission to view detailed information about this account."),
			)
		}
		return nil
	}

	// No valid permission to view this account
	return fault.Wrap(
		fault.New("insufficient permissions", ftag.With(ftag.PermissionDenied)),
		fctx.With(ctx),
		fmsg.WithDesc("insufficient permissions", "You do not have permission to view detailed information about this account."),
	)
}
