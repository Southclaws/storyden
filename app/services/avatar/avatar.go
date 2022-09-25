package avatar

import (
	"context"
	"io"
	"path"

	"github.com/Southclaws/fault/errctx"

	"github.com/Southclaws/storyden/app/resources/account"
)

func (s *service) Set(ctx context.Context, accountID account.AccountID, stream io.Reader) error {
	if err := s.storage.Write(ctx, path.Join("avatar", accountID.String()), stream); err != nil {
		return errctx.Wrap(err, ctx)
	}

	return nil
}

func (s *service) Get(ctx context.Context, accountID account.AccountID) (io.Reader, error) {
	stream, err := s.storage.Read(ctx, path.Join("avatar", accountID.String()))
	if err != nil {
		return nil, errctx.Wrap(err, ctx)
	}

	return stream, nil
}
