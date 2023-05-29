package object

import (
	"context"
	"io"
	"log"

	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"
	"github.com/Southclaws/fault/ftag"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"

	"github.com/Southclaws/storyden/internal/config"
)

type s3Storer struct {
	bucket      string
	minioClient *minio.Client
}

func NewS3Storer(cfg config.Config) Storer {
	minioClient, err := minio.New(cfg.S3Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(cfg.S3AccessKey, cfg.S3SecretKey, ""),
		Region: cfg.S3Region,
		Secure: true,
	})
	if err != nil {
		log.Fatalln(err)
	}

	return &s3Storer{
		bucket:      cfg.S3Bucket,
		minioClient: minioClient,
	}
}

func (s *s3Storer) Exists(ctx context.Context, path string) (bool, error) {
	_, err := s.minioClient.StatObject(ctx, s.bucket, path, minio.GetObjectOptions{})
	if err != nil {
		// TODO: figure out if there's a way to differentiate between an error
		// and an item just not existing. Ideally we want to treat actual
		// transport errors and such as actual errors and non-existence as nil.
		return false, fault.Wrap(err, fctx.With(ctx))
	}

	return true, nil
}

func (s *s3Storer) Read(ctx context.Context, path string) (io.Reader, error) {
	obj, err := s.minioClient.GetObject(ctx, s.bucket, path, minio.GetObjectOptions{})
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	_, err = obj.Stat()
	if err != nil {
		if minio.ToErrorResponse(err).StatusCode == 404 {
			return nil, fault.Wrap(err, fctx.With(ctx), ftag.With(ftag.NotFound))
		}
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return obj, nil
}

func (s *s3Storer) Write(ctx context.Context, path string, stream io.Reader) error {
	_, err := s.minioClient.PutObject(ctx, s.bucket, path, stream, -1, minio.PutObjectOptions{
		SendContentMd5: true,
	})
	if err != nil {
		return fault.Wrap(err, fctx.With(ctx))
	}

	return nil
}
