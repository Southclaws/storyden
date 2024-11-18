package provider

import (
	"context"

	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"
	"github.com/Southclaws/fault/ftag"
	"go.uber.org/zap"

	"github.com/Southclaws/storyden/app/resources/account/authentication"
	"github.com/Southclaws/storyden/app/resources/settings"
)

var ErrIncorrectMode = fault.New("incorrect authentication mode", ftag.With(ftag.InvalidArgument))

func CheckMode(ctx context.Context, l *zap.Logger, settings *settings.SettingsRepository, requiredMode authentication.Mode) error {
	set, err := settings.Get(ctx)
	if err != nil {
		return fault.Wrap(err, fctx.With(ctx))
	}

	if set.AuthenticationMode.Or(authentication.ModeHandle) != requiredMode {
		l.Warn("authentication occurred with a different mode to the preferred mode",
			zap.String("preferred_mode", requiredMode.String()),
			zap.String("actual_mode", set.AuthenticationMode.String()),
			zap.Error(fault.Wrap(ErrIncorrectMode, fctx.With(ctx))))
	}

	return nil
}
