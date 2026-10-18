package avatar

import (
	"context"
	"image/png"
	"io"

	"github.com/Southclaws/opt"

	"github.com/Southclaws/storyden/app/resources/account"
	"github.com/Southclaws/storyden/app/resources/account/account_repo"
	"github.com/Southclaws/storyden/app/resources/asset"
	"github.com/Southclaws/storyden/app/services/asset/asset_upload"
	"github.com/Southclaws/storyden/internal/mime"
)

func (s *service) Exists(ctx context.Context, accountID account.AccountID) bool {
	owner, err := s.accounts.GetRefByID(ctx, accountID)
	return err == nil && owner.AvatarAssetID.Ok()
}

func (s *service) Set(ctx context.Context, accountID account.AccountID, stream io.Reader, size int64) error {
	source, err := s.upload.UploadImage(
		ctx,
		accountID,
		stream,
		size,
		asset.NewFilename("avatar-original"),
		asset_upload.Options{},
	)
	if err != nil {
		return err
	}

	_, err = s.accounts.Update(ctx, accountID, account_repo.SetAvatarAssetID(opt.New(source.ID)))
	return err
}

func (s *service) Delete(ctx context.Context, accountID account.AccountID) error {
	_, err := s.accounts.Update(ctx, accountID, account_repo.SetAvatarAssetID(opt.NewEmpty[asset.AssetID]()))
	return err
}

func (s *service) Get(ctx context.Context, accountID account.AccountID) (io.Reader, int64, error) {
	a, r, _, err := s.GetVariant(ctx, accountID, asset.Original)
	if err != nil {
		return nil, 0, err
	}
	return r, int64(a.Size), nil
}

func (s *service) GetVariant(
	ctx context.Context,
	accountID account.AccountID,
	variant asset.Variant,
) (*asset.Asset, io.Reader, bool, error) {
	owner, err := s.accounts.GetRefByID(ctx, accountID)
	if err != nil {
		return nil, nil, false, err
	}
	if id, ok := owner.AvatarAssetID.Get(); ok {
		return s.download.VariantByID(ctx, id, variant)
	}

	return s.generated(ctx, accountID)
}

func (s *service) generated(ctx context.Context, accountID account.AccountID) (*asset.Asset, io.Reader, bool, error) {
	r, w := io.Pipe()

	go func() {
		i, err := s.generator.Generate(ctx, accountID.String())
		if err != nil {
			w.CloseWithError(err)
			return
		}

		if err := png.Encode(w, i); err != nil {
			w.CloseWithError(err)
			return
		}

		w.Close()
	}()

	return &asset.Asset{MIME: mime.New("image/png")}, r, false, nil
}
