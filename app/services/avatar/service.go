package avatar

import (
	"context"
	"io"

	"go.uber.org/fx"

	"github.com/Southclaws/storyden/app/resources/account"
	"github.com/Southclaws/storyden/app/resources/account/account_repo"
	"github.com/Southclaws/storyden/app/resources/asset"
	"github.com/Southclaws/storyden/app/services/asset/asset_download"
	"github.com/Southclaws/storyden/app/services/asset/asset_upload"
	"github.com/Southclaws/storyden/app/services/avatar_gen"
)

type Service interface {
	Exists(ctx context.Context, accountID account.AccountID) bool
	Set(ctx context.Context, accountID account.AccountID, stream io.Reader, size int64) error
	Delete(ctx context.Context, accountID account.AccountID) error
	Get(ctx context.Context, accountID account.AccountID) (io.Reader, int64, error)
	GetVariant(context.Context, account.AccountID, asset.Variant) (*asset.Asset, io.Reader, bool, error)
}

func Build() fx.Option {
	return fx.Provide(New)
}

type service struct {
	accounts  *account_repo.Repository
	upload    *asset_upload.Uploader
	download  *asset_download.Downloader
	generator avatar_gen.AvatarGenerator
}

func New(
	accounts *account_repo.Repository,
	upload *asset_upload.Uploader,
	download *asset_download.Downloader,
	generator avatar_gen.AvatarGenerator,
) Service {
	return &service{
		accounts:  accounts,
		upload:    upload,
		download:  download,
		generator: generator,
	}
}
