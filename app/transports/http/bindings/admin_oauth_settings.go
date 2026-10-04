package bindings

import (
	"net/url"
	"time"

	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/ftag"
	"github.com/Southclaws/opt"
	"github.com/oapi-codegen/nullable"
	"github.com/rs/xid"

	"github.com/Southclaws/storyden/app/resources/account/role"
	"github.com/Southclaws/storyden/app/resources/settings"
	"github.com/Southclaws/storyden/app/transports/http/openapi"
	"github.com/Southclaws/storyden/internal/deletable"
)

func deserialiseOAuthServiceSettings(in openapi.OAuthServiceSettings) (settings.OAuthServiceSettings, error) {
	registrationRole, err := deletable.NewMapErr(in.AutonomousRegistrationRoleId, func(value openapi.NullableIdentifier) (role.RoleID, error) {
		id, err := xid.FromString(string(value))
		return role.RoleID(id), err
	})
	if err != nil {
		return settings.OAuthServiceSettings{}, fault.Wrap(err, ftag.With(ftag.InvalidArgument))
	}

	roleID, clearRole := registrationRole.Get()
	if clearRole {
		roleID = opt.New(role.RoleID{})
	}

	mode, err := opt.MapErr(opt.NewPtr(in.AutonomousRegistrationMode), func(mode openapi.OAuthAutonomousRegistrationMode) (settings.OAuthAutonomousRegistrationMode, error) {
		return settings.NewOAuthAutonomousRegistrationMode(string(mode))
	})
	if err != nil {
		return settings.OAuthServiceSettings{}, fault.Wrap(err, ftag.With(ftag.InvalidArgument))
	}

	if in.RegistrationApprovalTtl != nil && (*in.RegistrationApprovalTtl < 1 || *in.RegistrationApprovalTtl > 9223372036) {
		return settings.OAuthServiceSettings{}, fault.New("registration approval lifetime must be between 1 and 9223372036 seconds", ftag.With(ftag.InvalidArgument))
	}

	if in.RegistrationApprovalUrl != nil && *in.RegistrationApprovalUrl != "" {
		u, err := url.Parse(*in.RegistrationApprovalUrl)
		if err != nil || (u.Scheme != "https" && u.Scheme != "http") || u.Hostname() == "" || u.User != nil || u.Fragment != "" {
			return settings.OAuthServiceSettings{}, fault.New("registration approval URL must be an absolute HTTP or HTTPS URL without credentials or a fragment", ftag.With(ftag.InvalidArgument))
		}
	}

	return settings.OAuthServiceSettings{
		AutonomousRegistrationRoleID: roleID,
		DynamicRegistrationEnabled:   opt.NewPtr(in.DynamicRegistrationEnabled),
		AutonomousRegistrationMode:   mode,
		RegistrationApprovalTTL:      opt.NewPtrMap(in.RegistrationApprovalTtl, func(seconds int) time.Duration { return time.Duration(seconds) * time.Second }),
		RegistrationApprovalURL:      opt.NewPtr(in.RegistrationApprovalUrl),
	}, nil
}

func serialiseOAuthServiceSettings(in settings.OAuthServiceSettings) openapi.OAuthServiceSettings {
	registrationRole := nullable.NewNullNullable[openapi.NullableIdentifier]()
	if id, ok := in.AutonomousRegistrationRoleID.Get(); ok && xid.ID(id) != xid.NilID() {
		registrationRole.Set(openapi.NullableIdentifier(id.String()))
	}

	return openapi.OAuthServiceSettings{
		AutonomousRegistrationRoleId: registrationRole,
		DynamicRegistrationEnabled:   in.DynamicRegistrationEnabled.Ptr(),
		AutonomousRegistrationMode: opt.Map(in.AutonomousRegistrationMode, func(mode settings.OAuthAutonomousRegistrationMode) openapi.OAuthAutonomousRegistrationMode {
			return openapi.OAuthAutonomousRegistrationMode(mode.String())
		}).Ptr(),
		RegistrationApprovalTtl: opt.Map(in.RegistrationApprovalTTL, func(d time.Duration) int { return int(d.Seconds()) }).Ptr(),
		RegistrationApprovalUrl: in.RegistrationApprovalURL.Ptr(),
	}
}
