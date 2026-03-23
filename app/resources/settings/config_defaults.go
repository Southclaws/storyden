package settings

import (
	"github.com/Southclaws/fault"
	"github.com/Southclaws/opt"

	"github.com/Southclaws/storyden/internal/ent"
)

// hydrateConfigDefaults takes an *ent.Setting, maps it to *Settings, and
// injects config.Config defaults for any values not set in the database.
func (d *SettingsRepository) hydrateConfigDefaults(in *ent.Setting) (*Settings, error) {
	settings, err := mapSettings(in)
	if err != nil {
		return nil, fault.Wrap(err)
	}

	d.hydrateClientIPDefaults(settings)
	d.hydrateRateLimitDefaults(settings)

	return settings, nil
}

func (d *SettingsRepository) hydrateClientIPDefaults(settings *Settings) {
	services, ok := settings.Services.Get()
	if !ok {
		services = ServiceSettings{}
	}

	clientIP, ok := services.ClientIP.Get()
	if !ok {
		services.ClientIP = opt.New(ClientIPServiceSettings{
			ClientIPMode:   opt.New(ClientIPModeRemoteAddr),
			ClientIPHeader: opt.New("X-Real-IP"),
		})
		settings.Services = opt.New(services)
		return
	}

	if !clientIP.ClientIPMode.Ok() {
		clientIP.ClientIPMode = opt.New(ClientIPModeRemoteAddr)
	}
	if !clientIP.ClientIPHeader.Ok() {
		clientIP.ClientIPHeader = opt.New("X-Real-IP")
	}

	services.ClientIP = opt.New(clientIP)
	settings.Services = opt.New(services)
}

func (d *SettingsRepository) hydrateRateLimitDefaults(settings *Settings) {
	services, ok := settings.Services.Get()
	if !ok {
		settings.Services = opt.New(ServiceSettings{
			RateLimit: opt.New(RateLimitServiceSettings{
				RateLimit:          opt.New(d.config.RateLimit),
				RateLimitPeriod:    opt.New(d.config.RateLimitPeriod),
				RateLimitBucket:    opt.New(d.config.RateLimitBucket),
				RateLimitGuestCost: opt.New(d.config.RateLimitGuestCost),
			}),
		})
		return
	}

	rateLimit, ok := services.RateLimit.Get()
	if !ok {
		services.RateLimit = opt.New(RateLimitServiceSettings{
			RateLimit:          opt.New(d.config.RateLimit),
			RateLimitPeriod:    opt.New(d.config.RateLimitPeriod),
			RateLimitBucket:    opt.New(d.config.RateLimitBucket),
			RateLimitGuestCost: opt.New(d.config.RateLimitGuestCost),
		})
		settings.Services = opt.New(services)
		return
	}

	if !rateLimit.RateLimit.Ok() {
		rateLimit.RateLimit = opt.New(d.config.RateLimit)
	}

	if !rateLimit.RateLimitPeriod.Ok() {
		rateLimit.RateLimitPeriod = opt.New(d.config.RateLimitPeriod)
	}

	if !rateLimit.RateLimitBucket.Ok() {
		rateLimit.RateLimitBucket = opt.New(d.config.RateLimitBucket)
	}

	if !rateLimit.RateLimitGuestCost.Ok() {
		rateLimit.RateLimitGuestCost = opt.New(d.config.RateLimitGuestCost)
	}

	services.RateLimit = opt.New(rateLimit)
	settings.Services = opt.New(services)
}
