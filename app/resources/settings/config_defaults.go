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

	d.hydrateRateLimitDefaults(settings)

	return settings, nil
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
