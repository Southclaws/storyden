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

// DefaultPartSize is the default part size for multipart uploads to S3. This is
// used only when the content size is unknown and is a low value for small VMs.
const DefaultPartSize = 1024 * 1024 * 8

type s3Storer struct {
	bucket      string
	minioClient *minio.Client
}

func NewS3Storer(ctx context.Context, cfg config.Config) (Storer, error) {
	minioClient, err := minio.New(cfg.S3Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(cfg.S3AccessKey, cfg.S3SecretKey, ""),
		Region: cfg.S3Region,
		Secure: cfg.S3Secure,
	})
	if err != nil {
		log.Fatalln(err)
	}

	exists, err := minioClient.BucketExists(ctx, cfg.S3Bucket)
	if err != nil {
		return nil, err
	}
	if !exists {
		err := minioClient.MakeBucket(ctx, cfg.S3Bucket, minio.MakeBucketOptions{})
		if err != nil {
			return nil, err
		}
	}

	return &s3Storer{
		bucket:      cfg.S3Bucket,
		minioClient: minioClient,
	}, nil
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

func (s *s3Storer) Read(ctx context.Context, path string) (io.Reader, int64, error) {
	obj, err := s.minioClient.GetObject(ctx, s.bucket, path, minio.GetObjectOptions{})
	if err != nil {
		return nil, 0, fault.Wrap(err, fctx.With(ctx))
	}

	info, err := obj.Stat()
	if err != nil {
		if minio.ToErrorResponse(err).StatusCode == 404 {
			return nil, 0, fault.Wrap(err, fctx.With(ctx), ftag.With(ftag.NotFound))
		}
		return nil, 0, fault.Wrap(err, fctx.With(ctx))
	}

	return obj, info.Size, nil
}

func (s *s3Storer) Write(ctx context.Context, path string, stream io.Reader, size int64) error {
	opts := minio.PutObjectOptions{}

	if size <= 0 {
		// If size is unknown Minio needs to compute an md5 hash of the content.
		opts.SendContentMd5 = true
		opts.PartSize = DefaultPartSize
	}

	_, err := s.minioClient.PutObject(ctx, s.bucket, path, stream, size, opts)
	if err != nil {
		return fault.Wrap(err, fctx.With(ctx))
	}

	return nil
}

func (s *s3Storer) Delete(ctx context.Context, path string) error {
	err := s.minioClient.RemoveObject(ctx, s.bucket, path, minio.RemoveObjectOptions{})
	if err != nil {
		return fault.Wrap(err, fctx.With(ctx))
	}

	return nil
}

func (s *s3Storer) List(ctx context.Context, prefix string) ([]string, error) {
	doneCh := make(chan struct{})
	defer close(doneCh)

	var objects []string
	for object := range s.minioClient.ListObjects(ctx, s.bucket, minio.ListObjectsOptions{
		Prefix: prefix,
	}) {
		if object.Err != nil {
			return nil, fault.Wrap(object.Err, fctx.With(ctx))
		}
		objects = append(objects, object.Key)
	}

	return objects, nil
}
