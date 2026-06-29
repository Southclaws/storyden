package workspaceprovider

import (
	"context"
	"sort"

	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"
	"github.com/Southclaws/fault/ftag"
	"github.com/puzpuzpuz/xsync/v4"
	"go.uber.org/fx"
	adksession "google.golang.org/adk/session"

	robotresource "github.com/Southclaws/storyden/app/resources/robot"
	"github.com/Southclaws/storyden/app/services/semdex/robot/workspaceprovider/local"
	spriteprovider "github.com/Southclaws/storyden/app/services/semdex/robot/workspaceprovider/sprites"
	workspacecap "github.com/Southclaws/storyden/app/services/semdex/robot/workspaceprovider/workspace"
	"github.com/Southclaws/storyden/app/services/semdex/robot/workspacestate"
)

var ErrProviderNotFound = fault.New("workspace provider not found")

type Provider interface {
	Provider() robotresource.WorkspaceProvider
	Mount(ctx context.Context, instance *robotresource.WorkspaceInstance) (map[string]any, error)
	Open(ctx context.Context, mount robotresource.WorkspaceMount) (Workspace, error)
	Cleanup(ctx context.Context, instance *robotresource.WorkspaceInstance) error
}

type Workspace = workspacecap.Workspace
type ListOptions = workspacecap.ListOptions
type FileInfo = workspacecap.FileInfo
type ReadFileResult = workspacecap.ReadFileResult
type WriteFileResult = workspacecap.WriteFileResult
type SearchMatch = workspacecap.SearchMatch
type SearchResult = workspacecap.SearchResult
type CommandSpec = workspacecap.CommandSpec
type CommandResult = workspacecap.CommandResult

type ProviderInfo struct {
	Provider robotresource.WorkspaceProvider
	Name     string
}

type Registry struct {
	providers *xsync.Map[robotresource.WorkspaceProvider, Provider]
}

func Build() fx.Option {
	return fx.Options(
		fx.Provide(NewRegistry, local.New, spriteprovider.New),
		fx.Invoke(registerLocal, registerSprites),
	)
}

func NewRegistry() *Registry {
	return &Registry{
		providers: xsync.NewMap[robotresource.WorkspaceProvider, Provider](),
	}
}

func (r *Registry) Register(provider Provider) {
	r.providers.Store(provider.Provider(), provider)
}

func (r *Registry) List() []ProviderInfo {
	providers := []ProviderInfo{}

	r.providers.Range(func(key robotresource.WorkspaceProvider, provider Provider) bool {
		providers = append(providers, ProviderInfo{
			Provider: key,
			Name:     string(key),
		})

		return true
	})

	sort.Slice(providers, func(i, j int) bool {
		return providers[i].Provider < providers[j].Provider
	})

	return providers
}

func (r *Registry) Get(ctx context.Context, provider robotresource.WorkspaceProvider) (Provider, error) {
	p, ok := r.providers.Load(provider)
	if !ok {
		return nil, fault.Wrap(
			fault.Newf("workspace provider not found: %s", provider),
			fctx.With(ctx),
			ftag.With(ftag.InvalidArgument),
		)
	}

	return p, nil
}

type stateContext interface {
	State() adksession.State
}

func WorkspaceFromState(ctx context.Context, registry *Registry) (Workspace, error) {
	if registry == nil {
		return nil, fault.Wrap(fault.New("workspace provider registry is not configured"), fctx.With(ctx), ftag.With(ftag.InvalidArgument))
	}

	stateProvider, ok := ctx.(stateContext)
	if !ok {
		return nil, fault.Wrap(fault.New("tool context does not expose session state"), fctx.With(ctx), ftag.With(ftag.InvalidArgument))
	}

	state := map[string]any{}
	for key, value := range stateProvider.State().All() {
		state[key] = value
	}

	mount, ok := workspacestate.MountFromState(state).Get()
	if !ok {
		return nil, fault.Wrap(fault.New("tool requires an active robot workspace"), fctx.With(ctx), ftag.With(ftag.InvalidArgument))
	}

	provider, err := registry.Get(ctx, mount.Provider)
	if err != nil {
		return nil, err
	}

	return provider.Open(ctx, mount)
}

func registerLocal(registry *Registry, provider *local.Provider) {
	registry.Register(provider)
}

func registerSprites(registry *Registry, provider *spriteprovider.Provider) {
	if provider.Enabled() {
		registry.Register(provider)
	}
}
