package account_manage

import (
	"context"

	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"
	"github.com/Southclaws/fault/fmsg"
	"github.com/Southclaws/fault/ftag"

	"github.com/Southclaws/storyden/app/resources/account"
	"github.com/Southclaws/storyden/app/resources/account/account_querier"
	"github.com/Southclaws/storyden/app/resources/rbac"
	"github.com/Southclaws/storyden/app/services/authentication/session"
)

type Manager struct {
	accountQuery *account_querier.Querier
}

func New(accountQuery *account_querier.Querier) *Manager {
	return &Manager{
		accountQuery: accountQuery,
	}
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
