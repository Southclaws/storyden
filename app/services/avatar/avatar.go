package avatar

import (
	"context"
	"image/png"
	"io"
	"path"

	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"

	"github.com/Southclaws/storyden/app/resources/account"
)

func avatarPath(aid account.AccountID) string {
	return path.Join("avatar", aid.String())
}

func (s *service) Exists(ctx context.Context, accountID account.AccountID) bool {
	exists, err := s.storage.Exists(ctx, avatarPath(accountID))
	if err != nil {
		return false // errors are ignored for now 🤠
	}

	return exists
}

func (s *service) Set(ctx context.Context, accountID account.AccountID, stream io.Reader) error {
	if err := s.storage.Write(ctx, avatarPath(accountID), stream); err != nil {
		return fault.Wrap(err, fctx.With(ctx))
	}

	return nil
}

func (s *service) Get(ctx context.Context, accountID account.AccountID) (io.Reader, error) {
	stream, err := s.storage.Read(ctx, avatarPath(accountID))
	if err != nil {
		r, w := io.Pipe()

		go func() {
			defer r.Close()

			i, err := s.avatar_gen.Generate(ctx, accountID.String())
			if err != nil {
				r.CloseWithError(err)
				return
			}

			if err := png.Encode(w, i); err != nil {
				r.CloseWithError(err)
				return
			}

			return
		}()

		return r, nil
	}

	return stream, nil
}
