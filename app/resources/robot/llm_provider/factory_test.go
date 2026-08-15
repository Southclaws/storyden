package llm_provider

import (
	"context"
	"database/sql"
	"errors"
	"io"
	"log/slog"
	"sync"
	"testing"
	"time"

	"entgo.io/ent/dialect"
	entsql "entgo.io/ent/dialect/sql"
	"github.com/Southclaws/opt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/fx/fxtest"
	"google.golang.org/adk/v2/model"

	"github.com/Southclaws/storyden/app/resources/robot/model_ref"
	"github.com/Southclaws/storyden/app/resources/robot/robot_model_cache"
	"github.com/Southclaws/storyden/app/resources/settings"
	"github.com/Southclaws/storyden/internal/config"
	"github.com/Southclaws/storyden/internal/ent"
	"github.com/Southclaws/storyden/internal/ent/enttest"
	"github.com/Southclaws/storyden/internal/infrastructure/cache/cachetest"

	_ "github.com/glebarez/go-sqlite"
)

var (
	testProvider = model_ref.NewProvider("test")
	modelAlpha   = model_ref.NewModel("alpha")
	modelBeta    = model_ref.NewModel("beta")
)

func TestFactoryRefreshProviderModelsWithKeyReplacesCache(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	factory, cacheRepo, provider := newTestFactory(t)

	provider.setModels(modelInfo(modelAlpha), modelInfo(modelBeta))
	models, err := factory.RefreshProviderModelsWithKey(ctx, testProvider, "key-1")
	require.NoError(t, err)
	require.Len(t, models, 2)
	assert.Equal(t, "key-1", provider.lastAPIKey())

	_, err = cacheRepo.GetModel(ctx, ref(modelAlpha))
	require.NoError(t, err)
	_, err = cacheRepo.GetModel(ctx, ref(modelBeta))
	require.NoError(t, err)

	provider.setModels(modelInfo(modelBeta))
	models, err = factory.RefreshProviderModelsWithKey(ctx, testProvider, "key-2")
	require.NoError(t, err)
	require.Len(t, models, 1)
	assert.Equal(t, "key-2", provider.lastAPIKey())

	_, err = cacheRepo.GetModel(ctx, ref(modelAlpha))
	assert.Error(t, err)
	_, err = cacheRepo.GetModel(ctx, ref(modelBeta))
	assert.NoError(t, err)
	assert.Equal(t, 2, provider.listCalls())

	provider.setModels()
	models, err = factory.RefreshProviderModelsWithKey(ctx, testProvider, "key-3")
	require.NoError(t, err)
	assert.Empty(t, models)

	_, err = cacheRepo.GetModel(ctx, ref(modelBeta))
	assert.Error(t, err)
}

func TestFactoryEnsureModelAvailableUsesFreshCacheWithoutProviderValidation(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	factory, cacheRepo, provider := newTestFactory(t)

	require.NoError(t, cacheRepo.UpsertProviderModels(ctx, testProvider, []model_ref.Info{modelInfo(modelAlpha)}))
	require.NoError(t, cacheRepo.SetRefreshed(ctx, testProvider, time.Now()))
	provider.validateErr = errors.New("provider should not be asked")

	err := factory.EnsureModelAvailable(ctx, ref(modelAlpha))
	require.NoError(t, err)
	assert.Equal(t, 0, provider.validateCalls())
	assert.Equal(t, 0, provider.listCalls())
}

func TestFactoryEnsureModelAvailableRejectsMissingModelWhenCacheIsFresh(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	factory, cacheRepo, provider := newTestFactory(t)

	require.NoError(t, cacheRepo.UpsertProviderModels(ctx, testProvider, []model_ref.Info{modelInfo(modelAlpha)}))
	require.NoError(t, cacheRepo.SetRefreshed(ctx, testProvider, time.Now()))
	provider.validateErr = errors.New("not found upstream")

	err := factory.EnsureModelAvailable(ctx, ref(modelBeta))
	require.Error(t, err)
	assert.Contains(t, err.Error(), "not available in the test model cache")
	assert.Equal(t, 1, provider.validateCalls())
	assert.Equal(t, 0, provider.listCalls())
}

func TestFactoryEnsureModelAvailableRefreshesStaleCacheForMissingModel(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	factory, cacheRepo, provider := newTestFactory(t)

	require.NoError(t, cacheRepo.UpsertProviderModels(ctx, testProvider, []model_ref.Info{modelInfo(modelAlpha)}))
	require.NoError(t, cacheRepo.SetRefreshed(ctx, testProvider, time.Now().Add(-ModelCacheTTL-time.Minute)))
	provider.validateErr = errors.New("not found before refresh")
	provider.setModels(modelInfo(modelAlpha), modelInfo(modelBeta))

	err := factory.EnsureModelAvailable(ctx, ref(modelBeta))
	require.NoError(t, err)

	_, err = cacheRepo.GetModel(ctx, ref(modelBeta))
	assert.NoError(t, err)
	assert.Equal(t, 1, provider.validateCalls())
	assert.Equal(t, 1, provider.listCalls())
}

func TestFactoryRefreshProviderModelsWithKeyRecordsErrorsWithoutReplacingCache(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	factory, cacheRepo, provider := newTestFactory(t)

	require.NoError(t, cacheRepo.UpsertProviderModels(ctx, testProvider, []model_ref.Info{modelInfo(modelAlpha)}))
	require.NoError(t, cacheRepo.SetRefreshed(ctx, testProvider, time.Now()))
	provider.listErr = errors.New("provider unavailable")

	models, err := factory.RefreshProviderModelsWithKey(ctx, testProvider, "key")
	require.Error(t, err)
	assert.Nil(t, models)

	_, err = cacheRepo.GetModel(ctx, ref(modelAlpha))
	assert.NoError(t, err)

	status, err := cacheRepo.GetStatus(ctx, testProvider)
	require.NoError(t, err)
	lastErr, ok := status.LastError.Get()
	require.True(t, ok)
	assert.Equal(t, "provider unavailable", lastErr)
}

func TestFactoryEnsureModelAvailableRefreshesStaleCachedModelBestEffort(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	factory, cacheRepo, provider := newTestFactory(t)

	require.NoError(t, cacheRepo.UpsertProviderModels(ctx, testProvider, []model_ref.Info{modelInfo(modelAlpha)}))
	require.NoError(t, cacheRepo.SetRefreshed(ctx, testProvider, time.Now().Add(-ModelCacheTTL-time.Minute)))
	provider.setModels(modelInfo(modelBeta))

	err := factory.EnsureModelAvailable(ctx, ref(modelAlpha))
	require.NoError(t, err)

	require.Eventually(t, func() bool {
		_, alphaErr := cacheRepo.GetModel(ctx, ref(modelAlpha))
		_, betaErr := cacheRepo.GetModel(ctx, ref(modelBeta))
		return alphaErr != nil && betaErr == nil
	}, time.Second, 10*time.Millisecond)
	assert.Equal(t, 1, provider.listCalls())
}

func newTestFactory(t *testing.T) (*Factory, *robot_model_cache.Repository, *fakeProvider) {
	t.Helper()

	ctx := context.Background()
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))
	sqlDB, err := sql.Open("sqlite", "file:"+t.Name()+"?mode=memory&cache=shared&_pragma=foreign_keys(1)")
	require.NoError(t, err)
	t.Cleanup(func() { _ = sqlDB.Close() })

	db := enttest.NewClient(t, enttest.WithOptions(ent.Driver(entsql.OpenDB(dialect.SQLite, sqlDB))))
	t.Cleanup(func() { _ = db.Close() })

	store := cachetest.New()
	cacheRepo := robot_model_cache.New(db, store, logger)

	settingsRepo, err := settings.New(ctx, fxtest.NewLifecycle(t), logger, db, config.Config{})
	require.NoError(t, err)
	_, err = settingsRepo.Set(ctx, settings.Settings{
		Services: opt.New(settings.ServiceSettings{
			Robots: opt.New(settings.RobotServiceSettings{
				Enabled: opt.New(true),
				Providers: opt.New(map[string]settings.RobotProviderSettings{
					testProvider.String(): {
						Enabled: opt.New(true),
						APIKey:  opt.New("test-key"),
					},
				}),
			}),
		}),
	})
	require.NoError(t, err)

	provider := &fakeProvider{requiresAPIKey: true}
	factory := New(logger, settingsRepo, cacheRepo)
	factory.Put(provider)

	return factory, cacheRepo, provider
}

func ref(model model_ref.Model) model_ref.ModelRef {
	return model_ref.ModelRef{
		Provider: testProvider,
		Model:    model,
	}
}

func modelInfo(model model_ref.Model) model_ref.Info {
	return model_ref.Info{Ref: ref(model)}
}

type fakeProvider struct {
	mu             sync.Mutex
	requiresAPIKey bool
	apiKey         string
	models         []model_ref.Info
	listErr        error
	validateErr    error
	listCount      int
	validateCount  int
}

func (p *fakeProvider) Provider() model_ref.Provider { return testProvider }

func (p *fakeProvider) RequiresAPIKey() bool { return p.requiresAPIKey }

func (p *fakeProvider) Configure(config Config) {
	p.mu.Lock()
	defer p.mu.Unlock()

	p.apiKey = config.APIKey
}

func (p *fakeProvider) ListModels(ctx context.Context) ([]model_ref.Info, error) {
	p.mu.Lock()
	defer p.mu.Unlock()

	p.listCount++
	if p.listErr != nil {
		return nil, p.listErr
	}

	out := make([]model_ref.Info, len(p.models))
	copy(out, p.models)
	return out, nil
}

func (p *fakeProvider) GetADKModelLLM(ctx context.Context, ref model_ref.ModelRef) (model.LLM, error) {
	return nil, nil
}

func (p *fakeProvider) ValidateModel(ctx context.Context, ref model_ref.ModelRef) error {
	p.mu.Lock()
	defer p.mu.Unlock()

	p.validateCount++
	return p.validateErr
}

func (p *fakeProvider) setModels(models ...model_ref.Info) {
	p.mu.Lock()
	defer p.mu.Unlock()

	p.models = models
}

func (p *fakeProvider) lastAPIKey() string {
	p.mu.Lock()
	defer p.mu.Unlock()

	return p.apiKey
}

func (p *fakeProvider) listCalls() int {
	p.mu.Lock()
	defer p.mu.Unlock()

	return p.listCount
}

func (p *fakeProvider) validateCalls() int {
	p.mu.Lock()
	defer p.mu.Unlock()

	return p.validateCount
}
