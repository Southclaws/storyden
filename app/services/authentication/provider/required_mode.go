package provider

import (
	"context"
	"log/slog"

	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"
	"github.com/Southclaws/fault/ftag"

	"github.com/Southclaws/storyden/app/resources/account/authentication"
	"github.com/Southclaws/storyden/app/resources/settings"
)

var ErrIncorrectMode = fault.New("incorrect authentication mode", ftag.With(ftag.InvalidArgument))

func CheckMode(ctx context.Context, logger *slog.Logger, settings *settings.SettingsRepository, requiredMode authentication.Mode) error {
	set, err := settings.Get(ctx)
	if err != nil {
		return fault.Wrap(err, fctx.With(ctx))
	}

	if set.AuthenticationMode.Or(authentication.ModeHandle) != requiredMode {
		logger.Warn("authentication occurred with a different mode to the preferred mode",
			slog.String("preferred_mode", requiredMode.String()),
			slog.String("actual_mode", set.AuthenticationMode.String()),
			slog.String("error", fault.Wrap(ErrIncorrectMode, fctx.With(ctx)).Error()))
	}

	return nil
}
