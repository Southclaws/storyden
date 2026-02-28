package role_assign

import (
	"context"
	"encoding/json"
	"log/slog"
	"time"

	"github.com/Southclaws/dt"
	"github.com/Southclaws/fault"
	"github.com/rs/xid"

	"github.com/Southclaws/storyden/internal/ent"
)

const (
	cachePrefix = "account:role_ids:"
	cacheTTL    = time.Hour * 24 * 365
)

func (w *Assignment) getRoleAssignmentsCache(ctx context.Context, accountID xid.ID) ([]*ent.AccountRoles, bool) {
	raw, err := w.store.Get(ctx, w.key(accountID))
	if err != nil {
		return nil, false
	}

	var payload []cachedAssignment
	if err := json.Unmarshal([]byte(raw), &payload); err == nil {
		assignments := make([]*ent.AccountRoles, 0, len(payload))
		for _, in := range payload {
			parsedID, err := xid.FromString(in.RoleID)
			if err != nil {
				_ = w.deleteRoleIDsCache(ctx, accountID)
				return nil, false
			}

			assignments = append(assignments, &ent.AccountRoles{
				ID:        xid.New(),
				CreatedAt: in.CreatedAt,
				AccountID: accountID,
				RoleID:    parsedID,
				Badge:     &in.Badge,
			})
		}

		return assignments, true
	}

	// Backward compatibility with legacy payloads that only stored role IDs.
	var legacy []string
	if err := json.Unmarshal([]byte(raw), &legacy); err != nil {
		_ = w.deleteRoleIDsCache(ctx, accountID)
		return nil, false
	}

	assignments := make([]*ent.AccountRoles, 0, len(legacy))
	for _, in := range legacy {
		parsedID, err := xid.FromString(in)
		if err != nil {
			_ = w.deleteRoleIDsCache(ctx, accountID)
			return nil, false
		}

		assignments = append(assignments, &ent.AccountRoles{
			ID:        xid.New(),
			AccountID: accountID,
			RoleID:    parsedID,
		})
	}

	return assignments, true
}

func (w *Assignment) storeRoleAssignmentsCache(ctx context.Context, accountID xid.ID, assignments []*ent.AccountRoles) error {
	payload, err := json.Marshal(dt.Map(assignments, func(in *ent.AccountRoles) cachedAssignment {
		return cachedAssignment{
			RoleID:    in.RoleID.String(),
			CreatedAt: in.CreatedAt,
			Badge:     in.Badge != nil && *in.Badge,
		}
	}))
	if err != nil {
		slog.Error("failed to marshal assignment cache payload",
			slog.String("account_id", accountID.String()),
			slog.String("error", err.Error()),
		)
		return fault.Wrap(err)
	}

	if err := w.store.Set(ctx, w.key(accountID), string(payload), cacheTTL); err != nil {
		slog.Error("failed to store assignment cache payload",
			slog.String("account_id", accountID.String()),
			slog.String("error", err.Error()),
		)
		return fault.Wrap(err)
	}

	return nil
}

func (w *Assignment) invalidateRoleIDsCache(ctx context.Context, accountID xid.ID) error {
	return w.deleteRoleIDsCache(ctx, accountID)
}

func (w *Assignment) key(accountID xid.ID) string {
	return cachePrefix + accountID.String()
}

func (w *Assignment) deleteRoleIDsCache(ctx context.Context, accountID xid.ID) error {
	if err := w.store.Delete(ctx, w.key(accountID)); err != nil {
		slog.Error("failed to delete assignment cache payload",
			slog.String("account_id", accountID.String()),
			slog.String("error", err.Error()),
		)
		return fault.Wrap(err)
	}

	return nil
}

type cachedAssignment struct {
	RoleID    string    `json:"role_id"`
	CreatedAt time.Time `json:"created_at"`
	Badge     bool      `json:"badge"`
}
