package bindings

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"time"

	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"
	"github.com/Southclaws/fault/fmsg"
	"github.com/Southclaws/fault/ftag"

	"github.com/Southclaws/storyden/app/resources/idempotency"
	"github.com/Southclaws/storyden/app/transports/http/openapi"
)

func idempotentCreate[T any](
	ctx context.Context,
	receipts *idempotency.Repository,
	principal, operation string,
	key *openapi.IdempotencyKey,
	request any,
	create func() (T, error),
	recheck func(T) error,
) (T, error) {
	var zero T
	if key == nil {
		return create()
	}
	if len(*key) == 0 || len(*key) > 128 {
		return zero, fault.New("invalid idempotency key", ftag.With(ftag.InvalidArgument))
	}
	for _, c := range *key {
		if c < '!' || c > '~' {
			return zero, fault.New("invalid idempotency key", ftag.With(ftag.InvalidArgument))
		}
	}

	data, err := json.Marshal(request)
	if err != nil {
		return zero, err
	}
	if len(data) > 1<<20 {
		return zero, fault.New("idempotent request exceeds one megabyte", ftag.With(ftag.InvalidArgument))
	}
	sum := sha256.Sum256(data)
	claim, err := receipts.Claim(ctx, principal, operation, string(*key), hex.EncodeToString(sum[:]))
	if err != nil {
		return zero, err
	}
	switch claim.State {
	case idempotency.Completed:
		var response T
		if err := json.Unmarshal([]byte(claim.Response), &response); err != nil {
			return zero, err
		}
		if err := recheck(response); err != nil {
			return zero, err
		}
		return response, nil
	case idempotency.Mismatch:
		return zero, fault.New("idempotency key was used with a different request", ftag.With(ftag.AlreadyExists), fctx.With(ctx, "idempotency_reason", "key_mismatch"))
	case idempotency.InProgress:
		return zero, fault.New("idempotent request is still in progress", ftag.With(ftag.AlreadyExists), fctx.With(ctx, "idempotency_reason", "in_progress"))
	case idempotency.Uncertain:
		return zero, fault.New("idempotent request outcome is uncertain; reconcile before using a new key", ftag.With(ftag.AlreadyExists), fctx.With(ctx, "idempotency_reason", "uncertain"))
	case idempotency.Claimed:
	default:
		return zero, errors.New("invalid idempotency claim state")
	}

	markFailed := func() {
		failureCtx, cancel := context.WithTimeout(context.WithoutCancel(ctx), 5*time.Second)
		defer cancel()
		_ = receipts.Fail(failureCtx, claim.ID, claim.Token)
	}
	response, err := create()
	if err != nil {
		markFailed()
		return zero, err
	}
	encoded, err := json.Marshal(response)
	if err != nil {
		markFailed()
		return zero, err
	}
	if len(encoded) > 2<<20 {
		markFailed()
		return zero, errors.New("idempotent response exceeds two megabytes")
	}
	completionCtx, cancel := context.WithTimeout(context.WithoutCancel(ctx), 5*time.Second)
	defer cancel()
	if err := receipts.Complete(completionCtx, claim.ID, claim.Token, string(encoded)); err != nil {
		markFailed()
		return zero, fault.Wrap(err, fmsg.With("create succeeded but idempotency response could not be recorded"))
	}
	return response, nil
}
