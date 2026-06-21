package robot_model_cache

import (
	"context"
	"encoding/json"
	"log/slog"
	"strings"
	"time"

	"entgo.io/ent/dialect/sql"
	"github.com/Southclaws/dt"
	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"
	"github.com/Southclaws/opt"

	"github.com/Southclaws/storyden/app/resources/robot/model_ref"
	"github.com/Southclaws/storyden/internal/ent"
	ent_robot_provider_model "github.com/Southclaws/storyden/internal/ent/robotprovidermodel"
	"github.com/Southclaws/storyden/internal/infrastructure/cache"
)

const (
	statusCachePrefix = "robots:provider-model-cache-status:"
	statusCacheTTL    = 7 * 24 * time.Hour
)

type Repository struct {
	db     *ent.Client
	store  cache.Store
	logger *slog.Logger
}

type providerStatusCache struct {
	LastRefreshedAt *time.Time `json:"last_refreshed_at,omitempty"`
	LastError       *string    `json:"last_error,omitempty"`
}

func New(db *ent.Client, store cache.Store, logger *slog.Logger) *Repository {
	return &Repository{db: db, store: store, logger: logger}
}

func (r *Repository) ListProviderModels(ctx context.Context, provider model_ref.Provider) ([]model_ref.Info, error) {
	rows, err := r.db.RobotProviderModel.Query().
		Where(ent_robot_provider_model.ProviderEQ(provider.String())).
		Order(ent_robot_provider_model.ByName(sql.OrderAsc())).
		All(ctx)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return dt.Map(rows, mapModel), nil
}

func (r *Repository) ListActiveModels(ctx context.Context, providers []model_ref.Provider) ([]model_ref.Info, error) {
	query := r.db.RobotProviderModel.Query().
		Order(ent_robot_provider_model.ByProvider(sql.OrderAsc()), ent_robot_provider_model.ByName(sql.OrderAsc()))

	if len(providers) > 0 {
		query.Where(ent_robot_provider_model.ProviderIn(dt.Map(providers, func(provider model_ref.Provider) string {
			return provider.String()
		})...))
	}

	rows, err := query.All(ctx)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return dt.Map(rows, mapModel), nil
}

func (r *Repository) GetModel(ctx context.Context, ref model_ref.ModelRef) (model_ref.Info, error) {
	row, err := r.db.RobotProviderModel.Query().
		Where(
			ent_robot_provider_model.ProviderEQ(ref.Provider.String()),
			ent_robot_provider_model.NameEQ(ref.Model.String()),
		).
		Only(ctx)
	if err != nil {
		return model_ref.Info{}, fault.Wrap(err, fctx.With(ctx))
	}

	return mapModel(row), nil
}

func (r *Repository) UpsertProviderModels(ctx context.Context, provider model_ref.Provider, models []model_ref.Info) error {
	now := time.Now()
	seen := make([]string, 0, len(models))

	for _, model := range models {
		name := model.Model().String()
		seen = append(seen, name)

		existing, err := r.db.RobotProviderModel.Query().
			Where(
				ent_robot_provider_model.ProviderEQ(provider.String()),
				ent_robot_provider_model.NameEQ(name),
			).
			Only(ctx)
		if ent.IsNotFound(err) {
			_, err = r.db.RobotProviderModel.Create().
				SetProvider(provider.String()).
				SetName(name).
				SetRaw(model.Raw).
				SetLastSeenAt(now).
				Save(ctx)
			if err != nil {
				return fault.Wrap(err, fctx.With(ctx))
			}
			continue
		}
		if err != nil {
			return fault.Wrap(err, fctx.With(ctx))
		}

		_, err = r.db.RobotProviderModel.UpdateOne(existing).
			SetRaw(model.Raw).
			SetLastSeenAt(now).
			Save(ctx)
		if err != nil {
			return fault.Wrap(err, fctx.With(ctx))
		}
	}

	deleteQuery := r.db.RobotProviderModel.Delete().
		Where(ent_robot_provider_model.ProviderEQ(provider.String()))
	if len(seen) > 0 {
		deleteQuery.Where(ent_robot_provider_model.NameNotIn(seen...))
	}
	if _, err := deleteQuery.Exec(ctx); err != nil {
		return fault.Wrap(err, fctx.With(ctx))
	}

	return nil
}

func (r *Repository) GetStatus(ctx context.Context, provider model_ref.Provider) (model_ref.CacheStatus, error) {
	key := r.statusCacheKey(provider)
	val, err := r.store.Get(ctx, key)
	if err == nil {
		var status providerStatusCache
		decodeErr := json.Unmarshal([]byte(val), &status)
		if decodeErr == nil {
			return mapStatus(provider, status), nil
		}

		r.logger.Warn("failed to decode robot model provider status cache", slog.String("provider", provider.String()), slog.String("key", key), slog.String("error", decodeErr.Error()))
		_ = r.store.Delete(ctx, key)
	} else if !isCacheMiss(err) {
		r.logger.Warn("failed to read robot model provider status cache", slog.String("provider", provider.String()), slog.String("key", key), slog.String("error", err.Error()))
	}

	row, err := r.db.RobotProviderModel.Query().
		Where(
			ent_robot_provider_model.ProviderEQ(provider.String()),
		).
		Order(ent_robot_provider_model.ByLastSeenAt(sql.OrderDesc())).
		First(ctx)
	if ent.IsNotFound(err) {
		return model_ref.CacheStatus{Provider: provider}, nil
	}
	if err != nil {
		return model_ref.CacheStatus{}, fault.Wrap(err, fctx.With(ctx))
	}

	return model_ref.CacheStatus{
		Provider:        provider,
		LastRefreshedAt: opt.New(row.LastSeenAt),
	}, nil
}

func (r *Repository) SetRefreshed(ctx context.Context, provider model_ref.Provider, refreshedAt time.Time) error {
	return r.setStatus(ctx, provider, providerStatusCache{
		LastRefreshedAt: &refreshedAt,
	})
}

func (r *Repository) SetError(ctx context.Context, provider model_ref.Provider, message string) error {
	current, _ := r.GetStatus(ctx, provider)
	status := providerStatusCache{
		LastRefreshedAt: current.LastRefreshedAt.Ptr(),
		LastError:       &message,
	}

	return r.setStatus(ctx, provider, status)
}

func (r *Repository) setStatus(ctx context.Context, provider model_ref.Provider, status providerStatusCache) error {
	body, err := json.Marshal(status)
	if err != nil {
		return fault.Wrap(err, fctx.With(ctx))
	}

	return r.store.Set(ctx, r.statusCacheKey(provider), string(body), statusCacheTTL)
}

func (r *Repository) statusCacheKey(provider model_ref.Provider) string {
	return statusCachePrefix + provider.String()
}

func isCacheMiss(err error) bool {
	return err != nil && strings.EqualFold(err.Error(), "not found")
}

func mapModel(in *ent.RobotProviderModel) model_ref.Info {
	return model_ref.Info{
		Ref: model_ref.ModelRef{
			Provider: model_ref.NewProvider(in.Provider),
			Model:    model_ref.NewModel(in.Name),
		},
		Raw:        in.Raw,
		LastSeenAt: in.LastSeenAt,
	}
}

func mapStatus(provider model_ref.Provider, in providerStatusCache) model_ref.CacheStatus {
	return model_ref.CacheStatus{
		Provider:        provider,
		LastRefreshedAt: opt.NewPtr(in.LastRefreshedAt),
		LastError:       opt.NewPtr(in.LastError),
	}
}
