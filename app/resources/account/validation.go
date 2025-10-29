package account

import (
	"context"

	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"
	"github.com/Southclaws/fault/fmsg"
	"github.com/Southclaws/fault/ftag"

	"github.com/Southclaws/storyden/app/resources/mark"
)

var errInvalidHandle = fault.New("invalid handle")

// ValidateHandle checks if a handle meets the requirements:
// - Must be a valid slug (lowercase, letters, numbers, hyphens, underscores)
// - Must be 30 characters or less
func ValidateHandle(ctx context.Context, handle string) error {
	if !mark.IsSlug(handle) {
		return fault.Wrap(
			errInvalidHandle,
			fctx.With(ctx),
			ftag.With(ftag.InvalidArgument),
			fmsg.WithDesc("invalid handle format", "Handle must be lowercase letters, numbers, hyphens, and underscores only"),
		)
	}

	if len(handle) > 30 {
		return fault.Wrap(
			errInvalidHandle,
			fctx.With(ctx),
			ftag.With(ftag.InvalidArgument),
			fmsg.WithDesc("handle too long", "Handle must be 30 characters or less"),
		)
	}

	return nil
}
