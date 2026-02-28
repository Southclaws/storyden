package account_repo

import (
	"context"
	"encoding/json"
	"log/slog"
	"time"

	"github.com/Southclaws/dt"
	"github.com/Southclaws/fault"
	"github.com/Southclaws/opt"
	"github.com/rs/xid"

	"github.com/Southclaws/storyden/app/resources/account"
	"github.com/Southclaws/storyden/app/resources/account/role"
	"github.com/Southclaws/storyden/app/resources/account/role/held"
	"github.com/Southclaws/storyden/app/resources/datagraph"
	"github.com/Southclaws/storyden/internal/infrastructure/cache"
)

const (
	accountCachePrefix = "account:data:"
	accountCacheTTL    = time.Hour * 24 * 365
)

type accountCache struct {
	store cache.Store
}

func newAccountCache(store cache.Store) *accountCache {
	return &accountCache{
		store: store,
	}
}

func (c *accountCache) get(ctx context.Context, id xid.ID) (*account.Account, bool) {
	val, err := c.store.Get(ctx, c.cacheKey(id))
	if err != nil {
		return nil, false
	}

	var cached cachedAccount
	if err := json.Unmarshal([]byte(val), &cached); err != nil {
		c.delete(ctx, id)
		return nil, false
	}

	parsed, err := cached.toDomain()
	if err != nil {
		c.delete(ctx, id)
		return nil, false
	}

	return parsed, true
}

func (c *accountCache) storeAccount(ctx context.Context, in *account.Account) {
	if in == nil {
		return
	}

	payload, err := json.Marshal(fromDomain(in))
	if err != nil {
		slog.Error("failed to marshal account cache payload",
			slog.String("account_id", in.ID.String()),
			slog.String("error", err.Error()),
		)
		return
	}

	if err := c.store.Set(ctx, c.cacheKey(xid.ID(in.ID)), string(payload), accountCacheTTL); err != nil {
		slog.Error("failed to store account cache payload",
			slog.String("account_id", in.ID.String()),
			slog.String("error", err.Error()),
		)
	}
}

func (c *accountCache) cacheKey(id xid.ID) string {
	return accountCachePrefix + id.String()
}

func (c *accountCache) delete(ctx context.Context, id xid.ID) {
	if err := c.store.Delete(ctx, c.cacheKey(id)); err != nil {
		slog.Error("failed to delete account cache payload",
			slog.String("account_id", id.String()),
			slog.String("error", err.Error()),
		)
	}
}

type cachedAccount struct {
	ID             string         `json:"id"`
	CreatedAt      time.Time      `json:"created_at"`
	UpdatedAt      time.Time      `json:"updated_at"`
	Handle         string         `json:"handle"`
	Name           string         `json:"name"`
	Bio            string         `json:"bio"`
	Signature      *string        `json:"signature,omitempty"`
	Kind           string         `json:"kind"`
	VerifiedStatus string         `json:"verified_status"`
	Admin          bool           `json:"admin"`
	RoleIDs        []string       `json:"role_ids,omitempty"`
	Metadata       map[string]any `json:"metadata,omitempty"`
	DeletedAt      *time.Time     `json:"deleted_at,omitempty"`
	IndexedAt      *time.Time     `json:"indexed_at,omitempty"`
}

func fromDomain(in *account.Account) cachedAccount {
	signature := opt.Map(in.Signature, datagraph.Content.HTML).Ptr()

	return cachedAccount{
		ID:             in.ID.String(),
		CreatedAt:      in.CreatedAt,
		UpdatedAt:      in.UpdatedAt,
		Handle:         in.Handle,
		Name:           in.Name,
		Bio:            in.Bio.HTML(),
		Signature:      signature,
		Kind:           in.Kind.String(),
		VerifiedStatus: in.VerifiedStatus.String(),
		Admin:          in.Admin,
		RoleIDs: dt.Map(in.Roles, func(r *held.Role) string {
			return r.ID.String()
		}),
		Metadata:  in.Metadata,
		DeletedAt: in.DeletedAt.Ptr(),
		IndexedAt: in.IndexedAt.Ptr(),
	}
}

func (in cachedAccount) toDomain() (*account.Account, error) {
	id, err := xid.FromString(in.ID)
	if err != nil {
		return nil, fault.Wrap(err)
	}

	kind, err := account.NewAccountKind(in.Kind)
	if err != nil {
		return nil, fault.Wrap(err)
	}

	verifiedStatus, err := account.NewVerifiedStatus(in.VerifiedStatus)
	if err != nil {
		// Backward compatibility for legacy cache payloads.
		if in.VerifiedStatus == "verified_email" {
			verifiedStatus = account.VerifiedStatusVerifiedEmail
		} else {
			return nil, fault.Wrap(err)
		}
	}

	bio, err := datagraph.NewRichText(in.Bio)
	if err != nil {
		return nil, fault.Wrap(err)
	}

	signature, err := opt.MapErr(opt.NewPtr(in.Signature), datagraph.NewRichText)
	if err != nil {
		return nil, fault.Wrap(err)
	}

	roles := make(held.Roles, 0, len(in.RoleIDs))
	for _, roleID := range in.RoleIDs {
		parsedID, err := xid.FromString(roleID)
		if err != nil {
			return nil, fault.Wrap(err)
		}

		roles = append(roles, &held.Role{
			Role: role.Role{
				ID: role.RoleID(parsedID),
			},
			Default: role.RoleID(parsedID) == role.DefaultRoleGuestID ||
				role.RoleID(parsedID) == role.DefaultRoleMemberID ||
				role.RoleID(parsedID) == role.DefaultRoleAdminID,
		})
	}

	return &account.Account{
		ID:             account.AccountID(id),
		CreatedAt:      in.CreatedAt,
		UpdatedAt:      in.UpdatedAt,
		Handle:         in.Handle,
		Name:           in.Name,
		Bio:            bio,
		Signature:      signature,
		Kind:           kind,
		VerifiedStatus: verifiedStatus,
		Admin:          in.Admin,
		Roles:          roles,
		Metadata:       in.Metadata,
		DeletedAt:      opt.NewPtr(in.DeletedAt),
		IndexedAt:      opt.NewPtr(in.IndexedAt),
	}, nil
}
