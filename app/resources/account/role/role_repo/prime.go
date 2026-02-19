package role_repo

import (
	"context"
	"log/slog"

	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"
	"github.com/rs/xid"
)

func (h *Repository) PrimeAssignmentsForAccount(ctx context.Context, accountID xid.ID) error {
	if _, ok := h.cachedAssignments(ctx, accountID); ok {
		return nil
	}

	assignments, err := h.assignmentsFromDB(ctx, accountID)
	if err != nil {
		return fault.Wrap(err, fctx.With(ctx))
	}

	if err := h.storeAssignments(ctx, accountID, assignments); err != nil {
		if recoveryErr := h.recoverFromCacheWriteFailure(ctx, "prime.prime_assignments_for_account.store_assignments account_id="+accountID.String(), err); recoveryErr != nil {
			slog.Error("role repository cache write recovery failed after priming account assignments",
				slog.String("account_id", accountID.String()),
				slog.String("cache_write_error", err.Error()),
				slog.String("recovery_error", recoveryErr.Error()))
		}
	}

	return nil
}
