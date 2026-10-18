package asset_upload

import (
	"context"
	"io"
	"strings"

	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"
	"github.com/Southclaws/fault/ftag"
	"github.com/Southclaws/opt"
	"github.com/rs/xid"

	"github.com/Southclaws/storyden/app/resources/account"
	"github.com/Southclaws/storyden/app/resources/asset"
	"github.com/Southclaws/storyden/app/resources/asset/asset_writer"
	"github.com/Southclaws/storyden/app/services/background_job/job_dispatcher"
	"github.com/Southclaws/storyden/internal/infrastructure/object"
	"github.com/Southclaws/storyden/internal/mime"
)

type Uploader struct {
	assets     *asset_writer.Writer
	objects    object.Storer
	dispatcher job_dispatcher.Dispatcher
}

func New(
	assets *asset_writer.Writer,
	objects object.Storer,
	dispatcher job_dispatcher.Dispatcher,
) *Uploader {
	return &Uploader{
		assets:     assets,
		objects:    objects,
		dispatcher: dispatcher,
	}
}

type Options struct {
	ParentID     opt.Optional[asset.AssetID]
	ImagePresets int64
}

func (s *Uploader) UploadImage(
	ctx context.Context,
	accountID account.AccountID,
	stream io.Reader,
	size int64,
	name asset.Filename,
	opts Options,
) (*asset.Asset, error) {
	mediaType, reader, err := mime.Detect(stream)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}
	if !strings.HasPrefix(mediaType.String(), "image/") {
		return nil, fault.Wrap(fault.New("image required"), ftag.With(ftag.InvalidArgument))
	}

	return s.Upload(ctx, accountID, reader, size, name, opts)
}

func (s *Uploader) Upload(
	ctx context.Context,
	accountID account.AccountID,
	stream io.Reader,
	size int64,
	name asset.Filename,
	opts Options,
) (*asset.Asset, error) {
	mediaType, reader, err := mime.Detect(stream)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	counter := &countingReader{reader: reader}

	if err := s.objects.Write(ctx, asset.BuildAssetPath(name), counter, size); err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	stored := int(counter.bytesRead)
	var requested int64
	if asset.ImageEligible(mediaType.String()) {
		requested = asset.EagerImagePresets | opts.ImagePresets
	}
	processing := asset_writer.WithImagePresets(requested)

	a, err := func() (asset *asset.Asset, err error) {
		if parentID, ok := opts.ParentID.Get(); ok {
			return s.assets.AddVersion(ctx, xid.ID(accountID), name, stored, *mediaType, parentID, processing)
		}

		return s.assets.Add(ctx, xid.ID(accountID), name, stored, *mediaType, processing)
	}()
	if err != nil {
		s.objects.Delete(context.WithoutCancel(ctx), asset.BuildAssetPath(name))
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	if requested != 0 {
		s.dispatcher.Dispatch(ctx, a.ID)
	}

	return a, nil
}

type countingReader struct {
	reader    io.Reader
	bytesRead int64
}

func (c *countingReader) Read(p []byte) (int, error) {
	n, err := c.reader.Read(p)
	c.bytesRead += int64(n)

	return n, err
}
