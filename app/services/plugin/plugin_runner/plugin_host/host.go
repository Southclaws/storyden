package plugin_host

import (
	"context"
	"log/slog"
	"net/url"
	"strings"

	"github.com/puzpuzpuz/xsync/v4"

	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"
	"github.com/Southclaws/fault/fmsg"
	"github.com/Southclaws/fault/ftag"
	"github.com/Southclaws/storyden/app/resources/plugin"
	"github.com/Southclaws/storyden/app/resources/plugin/plugin_reader"
	"github.com/Southclaws/storyden/app/resources/robot/llm_provider"
	"github.com/Southclaws/storyden/app/resources/robot/model_ref"
	"github.com/Southclaws/storyden/app/services/plugin/plugin_runner"
	"github.com/Southclaws/storyden/app/services/plugin/plugin_runner/duplex"
	"github.com/Southclaws/storyden/app/services/plugin/plugin_runner/plugin_llmprovider"
	"github.com/Southclaws/storyden/app/services/plugin/plugin_runner/plugin_logger"
	"github.com/Southclaws/storyden/app/services/plugin/plugin_runner/plugin_session"
	"github.com/Southclaws/storyden/app/services/plugin/plugin_runner/rpc_handler"
	"github.com/Southclaws/storyden/app/services/plugin/plugin_runner/supervised_runtime"
	"github.com/Southclaws/storyden/app/services/plugin/plugin_runner/supervised_runtime/local"
	"github.com/Southclaws/storyden/app/services/plugin/plugin_runner/supervised_runtime/sprites"
	"github.com/Southclaws/storyden/internal/config"
	"github.com/Southclaws/storyden/internal/infrastructure/pubsub"
	"github.com/Southclaws/storyden/lib/plugin/rpc"
)

type Host struct {
	logger            *slog.Logger
	sessions          *xsync.Map[plugin.InstallationID, plugin_runner.Session]
	pluginReader      *plugin_reader.Reader
	modelProviders    *llm_provider.Factory
	pluginProviders   *xsync.Map[plugin.InstallationID, []model_ref.Provider]
	rpcHandlerFactory *rpc_handler.Factory
	bus               *pubsub.Bus
	serverURL         url.URL
	runtimeProvider   supervised_runtime.Provider
	runtimeCtx        context.Context
}

func New(
	ctx context.Context,
	logger *slog.Logger,
	pluginReader *plugin_reader.Reader,
	modelProviders *llm_provider.Factory,
	pluginLogger *plugin_logger.Writer,
	rpcHandlerFactory *rpc_handler.Factory,
	bus *pubsub.Bus,
	cfg config.Config,
) (plugin_runner.Host, error) {
	localDefaultURL, _ := url.Parse("http://localhost:8000")

	runtimeProviderName, err := plugin_runner.ParseRuntimeProvider(cfg.PluginRuntimeProvider)
	if err != nil {
		return nil, fault.Wrap(err, fmsg.With("failed to parse plugin runtime provider"))
	}

	var runtimeProvider supervised_runtime.Provider
	serverURL := *localDefaultURL
	switch runtimeProviderName {
	case plugin_runner.RuntimeProviderNone:
		// Plugins are disabled by configuration.
		runtimeProvider = nil

	case plugin_runner.RuntimeProviderLocal:
		runtimeProvider = local.NewProvider(
			pluginLogger,
			pluginReader,
			cfg.PluginDataPath,
			cfg.PluginMaxRestartAttempts,
			cfg.PluginMaxBackoffDuration,
			cfg.PluginRuntimeCrashThreshold,
			cfg.PluginRuntimeCrashBackoff,
		)

	case plugin_runner.RuntimeProviderSprites:
		if strings.TrimSpace(cfg.SpritesAPIKey) == "" {
			return nil, fault.New("SPRITES_API_KEY is required when PLUGIN_RUNTIME_PROVIDER=sprites")
		}
		if strings.TrimSpace(cfg.PublicAPIAddress.Host) == "" {
			return nil, fault.New("PUBLIC_API_ADDRESS is required when PLUGIN_RUNTIME_PROVIDER=sprites")
		}
		serverURL = cfg.PublicAPIAddress
		runtimeProvider = sprites.NewProvider(
			cfg.SpritesAPIKey,
			pluginLogger,
			pluginReader,
			cfg.PluginDataPath,
			cfg.PluginMaxRestartAttempts,
			cfg.PluginMaxBackoffDuration,
			cfg.PluginRuntimeCrashThreshold,
			cfg.PluginRuntimeCrashBackoff,
		)

	default:
		return nil, fault.Newf("unsupported plugin runtime provider: %s", runtimeProviderName)
	}

	r := &Host{
		logger:            logger,
		sessions:          xsync.NewMap[plugin.InstallationID, plugin_runner.Session](),
		pluginReader:      pluginReader,
		modelProviders:    modelProviders,
		pluginProviders:   xsync.NewMap[plugin.InstallationID, []model_ref.Provider](),
		rpcHandlerFactory: rpcHandlerFactory,
		bus:               bus,
		serverURL:         serverURL,
		runtimeProvider:   runtimeProvider,
		runtimeCtx:        ctx,
	}

	logger.Info("configured plugin runtime provider", slog.String("provider", runtimeProviderName.String()))
	logger.Info("configured plugin runtime server url", slog.String("server_url", serverURL.String()))

	return r, nil
}

func (h *Host) Connect(ctx context.Context, id plugin.InstallationID, conn duplex.Duplex) error {
	sess, ok := h.sessions.Load(id)
	if !ok {
		return fault.Wrap(
			fault.New("plugin session not found"),
			fctx.With(ctx),
			fmsg.With("session must be loaded before connecting"),
		)
	}

	return sess.Connect(ctx, conn)
}

func (h *Host) Load(ctx context.Context, rec plugin.Record) error {
	// Check if already loaded
	if _, exists := h.sessions.Load(rec.InstallationID); exists {
		return nil
	}

	if rec.Mode.Supervised() {
		if h.runtimeProvider == nil {
			return fault.Wrap(
				fault.New("plugins"),
				fctx.With(ctx),
				ftag.With(ftag.PermissionDenied),
				fmsg.WithDesc(
					"disabled",
					"This Storyden instance has not enabled plugins.",
				),
			)
		}

		// Load binary from database
		bin, err := h.pluginReader.LoadBinary(ctx, rec.InstallationID)
		if err != nil {
			return fault.Wrap(err, fctx.With(ctx), fmsg.With("failed to load plugin binary"))
		}

		// Validate binary
		validated, err := plugin.Binary(bin).Validate(ctx)
		if err != nil {
			return fault.Wrap(err, fctx.With(ctx), fmsg.With("failed to validate plugin"))
		}

		rpch := h.rpcHandlerFactory.New(
			h.logger.With(slog.String("plugin_id", rec.InstallationID.String())),
			rec.InstallationID,
			validated,
		)

		runtime := h.runtimeProvider.New(
			rec.InstallationID,
			bin,
			validated,
			h.logger,
			h.serverURL,
			h.runtimeCtx,
		)

		sess := plugin_session.New(
			rec.InstallationID,
			bin,
			validated,
			h.bus,
			h.logger,
			rpch,
			runtime,
		)

		h.sessions.Store(rec.InstallationID, sess)
		h.registerModelProviders(rec.InstallationID, validated, sess)
	} else {
		// External plugin - no process management, websocket-only session.
		validated := &plugin.Validated{Metadata: rec.Manifest}
		rpch := h.rpcHandlerFactory.New(
			h.logger.With(slog.String("plugin_id", rec.InstallationID.String())),
			rec.InstallationID,
			validated,
		)

		sess := plugin_session.New(
			rec.InstallationID,
			nil,
			validated,
			h.bus,
			h.logger,
			rpch,
			nil,
		)

		h.sessions.Store(rec.InstallationID, sess)
		h.registerModelProviders(rec.InstallationID, validated, sess)
	}

	return nil
}

func (h *Host) Unload(ctx context.Context, id plugin.InstallationID) error {
	sess, ok := h.sessions.Load(id)
	if !ok {
		// Idempotent unload: session may not exist for inactive plugins that
		// were never loaded into memory.
		return nil
	}

	// Deactivate before unloading.
	// - Supervised plugins: stops the managed process.
	// - External plugins: cancels active websocket connection(s).
	if err := sess.SetActiveState(ctx, plugin.ActiveStateInactive); err != nil {
		return fault.Wrap(err, fctx.With(ctx), fmsg.With("failed to deactivate plugin session"))
	}

	h.unregisterModelProviders(id)
	h.sessions.Delete(id)

	return nil
}

func (h *Host) registerModelProviders(id plugin.InstallationID, manifest *plugin.Validated, sess plugin_runner.Session) {
	providers := []model_ref.Provider{}

	for _, capability := range manifest.Metadata.Capabilities {
		declaration, ok := capability.CapabilityConfigUnion.(*rpc.RobotLLMProviderCapabilityConfig)
		if !ok {
			continue
		}

		provider := model_ref.NewProvider(declaration.ID)
		h.modelProviders.Put(plugin_llmprovider.New(provider, sess))
		providers = append(providers, provider)

		name := declaration.Name.Or(declaration.ID)
		h.logger.Info("registered plugin robot model provider",
			slog.String("plugin_id", id.String()),
			slog.String("provider", provider.String()),
			slog.String("name", name))
	}
	if len(providers) > 0 {
		h.pluginProviders.Store(id, providers)
	}
}

func (h *Host) unregisterModelProviders(id plugin.InstallationID) {
	providers, ok := h.pluginProviders.LoadAndDelete(id)
	if !ok {
		h.logger.Warn("plugin robot model providers not found during unregister",
			slog.String("plugin_id", id.String()))
		return
	}
	for _, provider := range providers {
		h.modelProviders.Delete(provider)
		h.logger.Info("unregistered plugin robot model provider",
			slog.String("plugin_id", id.String()),
			slog.String("provider", provider.String()))
	}
}

func (h *Host) GetSession(ctx context.Context, id plugin.InstallationID) (plugin_runner.Session, error) {
	sess, ok := h.sessions.Load(id)
	if !ok {
		return nil, fault.Wrap(
			fault.New("plugin session not found"),
			fctx.With(ctx),
		)
	}
	return sess, nil
}

func (h *Host) GetSessions(ctx context.Context) ([]plugin_runner.Session, error) {
	sessions := []plugin_runner.Session{}
	h.sessions.Range(func(_ plugin.InstallationID, sess plugin_runner.Session) bool {
		sessions = append(sessions, sess)
		return true
	})
	return sessions, nil
}

func (h *Host) SetServerURL(u string) {
	parsed, err := url.Parse(u)
	if err != nil {
		h.logger.Error("failed to parse server URL", slog.String("url", u), slog.Any("error", err))
		return
	}
	h.serverURL = *parsed
}
