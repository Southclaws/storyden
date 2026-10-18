package asset_image

import (
	"context"
	"errors"

	"github.com/rs/xid"

	"github.com/Southclaws/storyden/app/resources/asset"
	"github.com/Southclaws/storyden/app/resources/background_job"
	"github.com/Southclaws/storyden/internal/ent"
	ent_asset "github.com/Southclaws/storyden/internal/ent/asset"
)

const Kind = "asset_image"

type Repository struct {
	db   *ent.Client
	jobs *background_job.Repository
}

func New(db *ent.Client) *Repository {
	return &Repository{
		db:   db,
		jobs: background_job.New(db),
	}
}

func Request(ctx context.Context, db *ent.Client, id xid.ID, mask int64) error {
	return New(db).Request(ctx, id, mask)
}

func (r *Repository) Request(ctx context.Context, id xid.ID, mask int64) error {
	payload := Payload{Requested: mask}
	if err := r.jobs.Ensure(ctx, id, Kind, id.String(), payload); err != nil {
		return err
	}

	return r.jobs.Request(ctx, id, func(payload *Payload) error {
		payload.Requested |= mask
		return nil
	}, Payload.IsComplete)
}

func (r *Repository) Get(ctx context.Context, id xid.ID) (*background_job.Job[Payload], error) {
	return r.jobs.Get[Payload](ctx, id)
}

func (r *Repository) Variants(ctx context.Context, id xid.ID) ([]*asset.Asset, error) {
	rows, err := r.db.Asset.Query().Where(ent_asset.ParentAssetIDEQ(id)).All(ctx)
	if err != nil {
		return nil, err
	}
	out := make([]*asset.Asset, 0, len(rows))
	for _, row := range rows {
		if !asset.ImageEligible(row.MimeType) {
			continue
		}
		variant := asset.Map(row)
		metadata, err := variant.ImageMetadata()
		if err != nil {
			return nil, err
		}
		if metadata.Recipe != asset.ImageRecipe {
			continue
		}
		if _, ok := asset.Preset(metadata.Preset); ok {
			out = append(out, variant)
		}
	}
	return out, nil
}

func (r *Repository) Checkpoint(
	ctx context.Context,
	job *background_job.Job[Payload],
	preset asset.ImagePreset,
	write func(*ent.Tx) error,
) error {
	return r.jobs.Checkpoint(ctx, job.State, func(tx *ent.Tx, payload *Payload) error {
		if payload.Completed&preset.Bit != 0 {
			return errors.New("image preset already completed")
		}

		if err := write(tx); err != nil {
			return err
		}
		payload.Completed |= preset.Bit
		return nil
	})
}

func (r *Repository) Dimensions(ctx context.Context, id xid.ID, width, height int) error {
	row, err := r.db.Asset.Get(ctx, id)
	if err != nil {
		return err
	}
	metadata := asset.Metadata{RawMessage: row.Metadata}
	updated, err := metadata.Merge(asset.ImageMetadata{Width: width, Height: height})
	if err != nil {
		return err
	}

	return r.db.Asset.UpdateOneID(id).SetMetadata(updated.RawMessage).Exec(ctx)
}

func (r *Repository) SkipPreset(
	ctx context.Context,
	job *background_job.Job[Payload],
	preset asset.ImagePreset,
	reason string,
) error {
	return r.jobs.UpdatePayload(ctx, job.ID, job.LeaseToken, func(payload *Payload) error {
		if payload.Skipped == nil {
			payload.Skipped = make(map[asset.Variant]string)
		}
		payload.Skipped[preset.Name] = reason
		payload.Completed |= preset.Bit
		return nil
	})
}
