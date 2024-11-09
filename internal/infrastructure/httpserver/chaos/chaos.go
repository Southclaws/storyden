package chaos

import (
	"fmt"
	"math/rand"
	"net/http"
	"time"

	"github.com/Southclaws/storyden/internal/config"
)

func WithChaos(cfg config.Config) func(http.Handler) http.Handler {
	failRate := cfg.DevChaosFailRate
	slowMode := cfg.DevChaosSlowMode

	disabled := failRate == 0 && slowMode == 0

	if disabled {
		return func(h http.Handler) http.Handler { return h }
	}

	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if slowMode > 0 {
				wait := time.Duration(rand.Intn(int(slowMode)))
				fmt.Println("[DEV_CHAOS_SLOW_MODE] waiting", wait)
				time.Sleep(wait)
			}

			if failRate > 0 {
				chance := rand.Float64()
				if chance < failRate {
					fmt.Println("[DEV_CHAOS_FAIL_RATE] crashing")
					w.WriteHeader(http.StatusInternalServerError)
					return
				}
			}

			h.ServeHTTP(w, r)
		})
	}
}
