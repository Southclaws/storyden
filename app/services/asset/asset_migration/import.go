package asset_migration

import (
	"context"
	"errors"
	"io"

	"github.com/Southclaws/fault/ftag"
	"github.com/Southclaws/opt"
	"github.com/rs/xid"

	"github.com/Southclaws/storyden/app/resources/account/account_repo"
	"github.com/Southclaws/storyden/app/resources/asset"
	"github.com/Southclaws/storyden/app/resources/asset/asset_image"
	migration "github.com/Southclaws/storyden/app/resources/asset/asset_migration"
	"github.com/Southclaws/storyden/app/resources/asset/asset_writer"
	jobs "github.com/Southclaws/storyden/app/resources/background_job"
	"github.com/Southclaws/storyden/app/resources/settings"
	"github.com/Southclaws/storyden/app/services/background_job"
	"github.com/Southclaws/storyden/internal/ent"
	"github.com/Southclaws/storyden/internal/mime"
)

var errReferenceChanged = errors.New("image reference changed during import")

func (m *Migrator) run(ctx context.Context, job *jobs.Job[migration.Payload]) (background_job.Result, error) {
	if job.Payload.SourceID == xid.NilID() {
		err := m.importSource(ctx, job)
		if err != nil {
			return background_job.Result{Terminal: errors.Is(err, errReferenceChanged)}, err
		}
	}

	var err error
	switch job.Payload.Target {
	case migration.Avatar:
		_, err = m.accounts.GetByID(ctx, job.Payload.AccountID)
	case migration.Icon:
		_, err = m.settings.Set(ctx, settings.Settings{IconAssetID: opt.New(job.Payload.SourceID)})
	case migration.Banner:
		_, err = m.settings.Set(ctx, settings.Settings{BannerAssetID: opt.New(job.Payload.SourceID)})
	}

	if err != nil {
		return background_job.Result{}, err
	}

	if job.Payload.Target == migration.Avatar {
		err = m.remove(ctx, job.Payload.Path)
	} else {
		for _, filename := range migration.LegacyFilenames(job.Payload.Target) {
			key := asset.BuildAssetPath(asset.NewFilepathFilename(filename))
			if err := m.remove(ctx, key); err != nil {
				return background_job.Result{}, err
			}
		}
	}

	return background_job.Result{}, err
}

func (m *Migrator) remove(ctx context.Context, key string) error {
	err := m.objects.Delete(ctx, key)
	if ftag.Get(err) == ftag.NotFound {
		return nil
	}

	return err
}

func (m *Migrator) importSource(ctx context.Context, job *jobs.Job[migration.Payload]) error {
	key := job.Payload.Path
	reader, size, err := m.objects.Read(ctx, key)
	if err != nil {
		return err
	}

	if closer, ok := reader.(io.Closer); ok {
		defer closer.Close()
	}

	mediaType, stream, err := mime.Detect(reader)
	if err != nil {
		return err
	}

	name := asset.NewFilename(string(job.Payload.Target) + "-original")
	source := &asset.Asset{
		ID:      name.GetID(),
		OwnerID: xid.ID(job.Payload.AccountID),
		Name:    name,
		Size:    int(size),
		MIME:    *mediaType,
	}

	candidate := asset.BuildAssetPath(name)
	published := false
	defer func() {
		if !published {
			if err := m.objects.Delete(context.WithoutCancel(ctx), candidate); err != nil {
				m.logger.WarnContext(ctx, "failed to remove import candidate", "path", candidate, "error", err)
			}
		}
	}()

	err = m.objects.Write(ctx, candidate, stream, size)
	if err != nil {
		return err
	}

	err = m.jobs.Checkpoint(ctx, job.State, func(tx *ent.Tx, payload *migration.Payload) error {
		return publishSource(ctx, tx, payload, source)
	})
	if err != nil {
		return err
	}

	published = true
	job.Payload.SourceID = source.ID
	if asset.ImageEligible(source.MIME.String()) {
		m.runner.Dispatch(ctx, source.ID)
	}

	return nil
}

func publishSource(
	ctx context.Context,
	tx *ent.Tx,
	payload *migration.Payload,
	source *asset.Asset,
) error {
	db := tx.Client()

	if payload.Target == migration.Avatar {
		updated, err := account_repo.SetAvatarAssetIfUnchanged(
			ctx,
			tx,
			xid.ID(payload.AccountID),
			payload.AccountUpdatedAt,
			source.ID,
		)
		if err != nil {
			return err
		}

		if !updated {
			return errReferenceChanged
		}
	}

	if err := asset_writer.AddSource(ctx, db, source); err != nil {
		return err
	}

	if asset.ImageEligible(source.MIME.String()) {
		err := asset_image.Request(ctx, db, source.ID, payload.Presets())
		if err != nil {
			return err
		}
	}

	payload.SourceID = source.ID

	return nil
}
