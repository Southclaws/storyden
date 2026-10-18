package asset_image

import (
	"context"
	"errors"
	"image"
	"io"
	"log/slog"
	"os"
	"sync/atomic"

	"github.com/Southclaws/storyden/app/resources/asset"
	jobs "github.com/Southclaws/storyden/app/resources/asset/asset_image"
	"github.com/Southclaws/storyden/app/resources/asset/asset_querier"
	"github.com/Southclaws/storyden/app/resources/asset/asset_writer"
	job_resource "github.com/Southclaws/storyden/app/resources/background_job"
	"github.com/Southclaws/storyden/app/services/background_job"
	"github.com/Southclaws/storyden/internal/config"
	"github.com/Southclaws/storyden/internal/ent"
	"github.com/Southclaws/storyden/internal/infrastructure/imagery"
	"github.com/Southclaws/storyden/internal/infrastructure/object"
	"github.com/Southclaws/storyden/internal/mime"
)

type Processor struct {
	jobs    *jobs.Repository
	assets  *asset_querier.Querier
	objects object.Storer
	logger  *slog.Logger
	cfg     config.Config

	memory atomic.Int64
}

func New(
	imageJobs *jobs.Repository,
	assets *asset_querier.Querier,
	objects object.Storer,
	runner *background_job.Runner,
	logger *slog.Logger,
	cfg config.Config,
) (*Processor, error) {
	if cfg.ImageMemoryBudget <= 0 {
		cfg.ImageMemoryBudget = 256 << 20
	}
	if cfg.ImageMaxSourceBytes <= 0 {
		cfg.ImageMaxSourceBytes = 32 << 20
	}
	if cfg.ImageMaxSourcePixels <= 0 {
		cfg.ImageMaxSourcePixels = 12000000
	}

	p := &Processor{
		jobs:    imageJobs,
		assets:  assets,
		objects: objects,
		logger:  logger,
		cfg:     cfg,
	}

	err := runner.Register(jobs.Kind, max(1, cfg.ImageWorkers), background_job.Handler[jobs.Payload]{
		Run:      p.run,
		Complete: jobs.Payload.IsComplete,
	})
	if err != nil {
		return nil, err
	}

	return p, nil
}

func (p *Processor) run(ctx context.Context, job *job_resource.Job[jobs.Payload]) (background_job.Result, error) {
	err := p.process(ctx, job)
	if errors.Is(err, errMemoryBusy) {
		return background_job.Result{Deferred: true}, err
	}

	terminal := errors.Is(err, imagery.ErrUnsupported) || errors.Is(err, imagery.ErrTooLarge)
	return background_job.Result{Terminal: terminal}, err
}

func (p *Processor) reserve(bytes int64) bool {
	for {
		reserved := p.memory.Load()
		if bytes > p.cfg.ImageMemoryBudget-reserved {
			return false
		}
		if p.memory.CompareAndSwap(reserved, reserved+bytes) {
			return true
		}
	}
}

func closeReader(reader io.Reader) {
	if closer, ok := reader.(io.Closer); ok {
		closer.Close()
	}
}

var errMemoryBusy = errors.New("image processing memory is currently reserved")

func (p *Processor) process(ctx context.Context, job *job_resource.Job[jobs.Payload]) error {
	source, err := p.assets.GetByID(ctx, job.ID)
	if err != nil {
		return err
	}

	img, reserved, admitted, err := p.decode(ctx, source, job.Payload.Requested & ^job.Payload.Completed)
	if err != nil {
		return err
	}
	defer p.memory.Add(-reserved)

	if err = p.jobs.Dimensions(ctx, source.ID, img.Bounds().Dx(), img.Bounds().Dy()); err != nil {
		return err
	}

	for _, preset := range asset.ImagePresets {
		if job.Payload.Requested & ^job.Payload.Completed & preset.Bit == 0 {
			continue
		}
		if admitted&preset.Bit == 0 {
			if err = p.jobs.SkipPreset(ctx, job, preset, "memory budget exceeded"); err != nil {
				return err
			}
			continue
		}
		if err = ctx.Err(); err != nil {
			return err
		}
		if err = p.generate(ctx, job, source, img, preset); err != nil {
			return err
		}
	}
	return nil
}

func (p *Processor) decode(
	ctx context.Context,
	source *asset.Asset,
	requested int64,
) (image.Image, int64, int64, error) {
	r, _, err := p.objects.Read(ctx, asset.BuildAssetPath(source.Name))
	if err != nil {
		return nil, 0, 0, err
	}
	defer closeReader(r)

	input, err := imagery.Open(r, p.cfg.ImageMaxSourceBytes, p.cfg.ImageMaxSourcePixels)
	if err != nil {
		return nil, 0, 0, err
	}
	defer input.Close()

	width, height := 0, 0
	var admitted int64
	for _, preset := range asset.ImagePresets {
		if requested&preset.Bit == 0 {
			continue
		}

		estimate := imagery.Estimate(input, preset.Width, preset.Height)
		if estimate > p.cfg.ImageMemoryBudget {
			continue
		}

		admitted |= preset.Bit
		width = max(width, preset.Width)
		height = max(height, preset.Height)
	}

	if admitted == 0 {
		return nil, 0, 0, imagery.ErrTooLarge
	}
	estimate := imagery.Estimate(input, width, height)
	if estimate > p.cfg.ImageMemoryBudget {
		return nil, 0, 0, imagery.ErrTooLarge
	}
	if !p.reserve(estimate) {
		return nil, 0, 0, errMemoryBusy
	}

	img, err := input.Decode()
	if err != nil {
		p.memory.Add(-estimate)
		return nil, 0, 0, errors.Join(imagery.ErrUnsupported, err)
	}
	return img, estimate, admitted, nil
}

func (p *Processor) generate(
	ctx context.Context,
	job *job_resource.Job[jobs.Payload],
	source *asset.Asset,
	img image.Image,
	preset asset.ImagePreset,
) error {
	var resized *image.NRGBA
	if preset.PNG {
		resized = imagery.Fill(img, preset.Width, preset.Height)
	} else {
		resized = imagery.Resize(img, preset.Width, preset.Height, preset.Square)
	}
	out, err := os.CreateTemp("", "storyden-variant-*")
	if err != nil {
		return err
	}
	defer func() {
		out.Close()
		os.Remove(out.Name())
	}()

	if err = imagery.Encode(out, resized, preset.PNG); err != nil {
		return err
	}
	size, err := out.Seek(0, io.SeekCurrent)
	if err != nil {
		return err
	}
	if size > p.cfg.ImageMaxSourceBytes {
		return imagery.ErrTooLarge
	}
	if _, err = out.Seek(0, io.SeekStart); err != nil {
		return err
	}
	ext := "webp"
	if preset.PNG {
		ext = "png"
	}
	name := asset.NewFilename(string(preset.Name) + "." + ext)
	bounds := resized.Bounds()
	mediaType := mime.New("image/" + ext)
	metadata, err := asset.NewMetadata(mediaType, asset.ImageMetadata{
		Width:  bounds.Dx(),
		Height: bounds.Dy(),
		Preset: preset.Name,
		Recipe: asset.ImageRecipe,
	})
	if err != nil {
		return err
	}
	variant := &asset.Asset{
		Name:     name,
		Size:     int(size),
		MIME:     mediaType,
		Metadata: metadata,
	}

	path := asset.BuildAssetPath(name)
	if err = p.objects.Write(ctx, path, out, size); err != nil {
		p.objects.Delete(context.WithoutCancel(ctx), path)
		return err
	}

	err = p.jobs.Checkpoint(ctx, job, preset, func(tx *ent.Tx) error {
		return asset_writer.AddVariant(ctx, tx.Client(), source, variant)
	})
	if err != nil {
		p.objects.Delete(context.WithoutCancel(ctx), path)
		return err
	}
	p.logger.Debug("image variant generated", "asset_id", source.ID, "preset", preset.Name, "bytes", size)
	return nil
}
