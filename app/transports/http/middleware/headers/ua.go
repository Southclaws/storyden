package headers

import (
	"context"
	"log/slog"
	"net/http"
	"sync/atomic"

	"go.uber.org/fx"

	"github.com/getkin/kin-openapi/routers/gorillamux"

	"github.com/Southclaws/storyden/app/resources/settings"
	"github.com/Southclaws/storyden/app/services/reqinfo"
	"github.com/Southclaws/storyden/app/transports/http/openapi"
	"github.com/Southclaws/storyden/internal/infrastructure/pubsub"
	"github.com/Southclaws/storyden/lib/plugin/rpc"
)

type Middleware struct {
	clientIPConfig atomic.Value
	settingsRepo   *settings.SettingsRepository
	logger         *slog.Logger
}

func New(
	ctx context.Context,
	lc fx.Lifecycle,
	settingsRepo *settings.SettingsRepository,
	bus *pubsub.Bus,
	logger *slog.Logger,
) *Middleware {
	m := &Middleware{
		settingsRepo: settingsRepo,
		logger:       logger,
	}
	m.clientIPConfig.Store(defaultClientIPConfiguration())

	lc.Append(fx.StartHook(func(hctx context.Context) error {
		m.reloadClientIPConfiguration(hctx)

		_, err := pubsub.Subscribe(ctx, bus, "headers.client_ip_settings_updated", func(ctx context.Context, evt *rpc.EventSettingsUpdated) error {
			m.reloadClientIPConfiguration(ctx)
			return nil
		})
		return err
	}))

	return m
}

// WithHeaderContext stores in the request context header info.
func (m *Middleware) WithHeaderContext() func(next http.Handler) http.Handler {
	spec, err := openapi.GetSwagger()
	if err != nil {
		panic(err)
	}

	gr, err := gorillamux.NewRouter(spec)
	if err != nil {
		panic(err)
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()

			opid := "unknown"
			rt, _, err := gr.FindRoute(r)
			if err == nil && rt.Operation != nil {
				opid = rt.Operation.OperationID
			}

			cfg := m.currentClientIPConfiguration()
			clientAddress := m.clientAddressWithConfig(r, cfg)
			ssrClientAddress := m.ssrClientAddress(r, clientAddress)

			newctx := reqinfo.WithRequestInfo(ctx, r, opid, clientAddress, ssrClientAddress)

			r = r.WithContext(newctx)

			next.ServeHTTP(w, r)
		})
	}
}
