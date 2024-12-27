package chaos

import (
	"fmt"
	"math/rand"
	"net/http"
	"time"

	"github.com/Southclaws/storyden/internal/config"
)

type Middleware struct {
	enabled  bool
	failRate float64
	slowMode time.Duration
}

func New(cfg config.Config) *Middleware {
	return &Middleware{
		enabled:  cfg.DevChaosFailRate > 0 || cfg.DevChaosSlowMode == 0,
		failRate: cfg.DevChaosFailRate,
		slowMode: cfg.DevChaosSlowMode,
	}
}

func (m *Middleware) WithChaos() func(http.Handler) http.Handler {
	if !m.enabled {
		return func(h http.Handler) http.Handler { return h }
	}

	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if m.slowMode > 0 {
				wait := time.Duration(rand.Intn(int(m.slowMode)))
				fmt.Println("[DEV_CHAOS_SLOW_MODE] waiting", wait)
				time.Sleep(wait)
			}

			if m.failRate > 0 {
				chance := rand.Float64()
				if chance < m.failRate {
					fmt.Println("[DEV_CHAOS_FAIL_RATE] crashing")
					w.WriteHeader(http.StatusInternalServerError)
					return
				}
			}

			h.ServeHTTP(w, r)
		})
	}
}
