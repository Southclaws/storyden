package idempotency

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"sync/atomic"
	"time"

	"github.com/Southclaws/storyden/internal/ent"
	"github.com/Southclaws/storyden/internal/ent/idempotencyreceipt"
	"github.com/rs/xid"
)

const retention = 24 * time.Hour
const abandonedRetention = 30 * 24 * time.Hour

type State int

const (
	Claimed State = iota
	Completed
	InProgress
	Uncertain
	Mismatch
)

type Result struct {
	State    State
	Response string
	ID       string
	Token    string
}

type Repository struct {
	db        *ent.Client
	lastPrune atomic.Int64
}

func New(db *ent.Client) *Repository { return &Repository{db: db} }

func (r *Repository) Claim(ctx context.Context, principal, operation, key, fingerprint string) (Result, error) {
	r.pruneOccasionally(ctx)
	keyHash := digest(key)
	id := digest(principal + "\x00" + operation + "\x00" + keyHash)
	now := time.Now()

	for range 3 {
		token := xid.New().String()
		_, err := r.db.IdempotencyReceipt.Create().
			SetID(id).
			SetPrincipalID(principal).
			SetOperation(operation).
			SetKeyHash(keyHash).
			SetClaimToken(token).
			SetFingerprint(fingerprint).
			SetResponse("").
			SetExpiresAt(now.Add(retention)).
			Save(ctx)
		if err == nil {
			return Result{State: Claimed, ID: id, Token: token}, nil
		}
		if !ent.IsConstraintError(err) {
			return Result{}, err
		}

		receipt, err := r.db.IdempotencyReceipt.Get(ctx, id)
		if ent.IsNotFound(err) {
			continue
		}
		if err != nil {
			return Result{}, err
		}
		if receipt.Response != "" {
			if now.After(receipt.ExpiresAt) {
				_, err := r.db.IdempotencyReceipt.Delete().Where(
					idempotencyreceipt.IDEQ(id),
					idempotencyreceipt.ExpiresAtLT(now),
				).Exec(ctx)
				if err != nil {
					return Result{}, err
				}
				continue
			}
			if receipt.Fingerprint != fingerprint {
				return Result{State: Mismatch}, nil
			}
			return Result{State: Completed, Response: receipt.Response}, nil
		}
		if receipt.Fingerprint != fingerprint {
			return Result{State: Mismatch}, nil
		}
		if receipt.FailedAt != nil {
			return Result{State: Uncertain}, nil
		}
		if now.After(receipt.ExpiresAt) {
			return Result{State: Uncertain}, nil
		}
		return Result{State: InProgress}, nil
	}
	return Result{}, errors.New("idempotency claim changed during lookup")
}

func (r *Repository) Complete(ctx context.Context, id, token, response string) error {
	if response == "" {
		return errors.New("empty idempotency response")
	}
	count, err := r.db.IdempotencyReceipt.Update().Where(
		idempotencyreceipt.IDEQ(id),
		idempotencyreceipt.ClaimTokenEQ(token),
		idempotencyreceipt.ResponseEQ(""),
	).SetResponse(response).Save(ctx)
	if err != nil {
		return err
	}
	if count != 1 {
		return errors.New("idempotency claim was no longer active")
	}
	return nil
}

func (r *Repository) Fail(ctx context.Context, id, token string) error {
	_, err := r.db.IdempotencyReceipt.Update().Where(
		idempotencyreceipt.IDEQ(id),
		idempotencyreceipt.ClaimTokenEQ(token),
		idempotencyreceipt.ResponseEQ(""),
		idempotencyreceipt.FailedAtIsNil(),
	).SetFailedAt(time.Now()).Save(ctx)
	return err
}

func (r *Repository) pruneOccasionally(ctx context.Context) {
	now := time.Now().Unix()
	last := r.lastPrune.Load()
	if now-last < 3600 || !r.lastPrune.CompareAndSwap(last, now) {
		return
	}
	if err := r.Prune(ctx); err != nil {
		r.lastPrune.Store(last)
	}
}

func (r *Repository) Prune(ctx context.Context) error {
	now := time.Now()
	_, err := r.db.IdempotencyReceipt.Delete().Where(
		idempotencyreceipt.Or(
			idempotencyreceipt.And(idempotencyreceipt.ResponseNEQ(""), idempotencyreceipt.ExpiresAtLT(now)),
			idempotencyreceipt.And(idempotencyreceipt.ResponseEQ(""), idempotencyreceipt.CreatedAtLT(now.Add(-abandonedRetention))),
		),
	).Exec(ctx)
	return err
}

func digest(value string) string {
	sum := sha256.Sum256([]byte(value))
	return hex.EncodeToString(sum[:])
}
