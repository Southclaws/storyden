package provider

import (
	"context"
	"fmt"

	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"
	"github.com/Southclaws/fault/fmsg"
	"github.com/Southclaws/fault/ftag"

	"github.com/Southclaws/storyden/app/resources/account/authentication"
	"github.com/Southclaws/storyden/app/resources/settings"
)

var ErrIncorrectMode = fault.New("incorrect authentication mode", ftag.With(ftag.InvalidArgument))

func CheckMode(ctx context.Context, settings *settings.SettingsRepository, requiredMode authentication.Mode) error {
	set, err := settings.Get(ctx)
	if err != nil {
		return fault.Wrap(err, fctx.With(ctx))
	}

	if set.AuthenticationMode.Or(authentication.ModeHandle) != requiredMode {
		message := fmt.Sprintf("This site is configured for %s authentication only", set.AuthenticationMode.Or(authentication.ModeHandle))
		return fault.Wrap(ErrIncorrectMode, fctx.With(ctx), fmsg.WithDesc("required mode", message))
	}

	return nil
}
