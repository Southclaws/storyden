package spam_checker

import (
	"compress/gzip"
	"context"
	"io"

	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"
	"golang.org/x/text/transform"
)

// Detector describes a service which can detect spam messages.
type Detector interface {
	Detect(ctx context.Context, r io.Reader) (bool, error)
}

func New() Detector {
	return &repeatedContentDetector{
		// a compression ratio of 0.0003 is very leniant and only blocks spammy
		// sequences of characters repeated hundreds or thousands of times.
		threshold: 0.0003,
	}
}

// repeatedContentDetector implements a detector with an extremely basic repeat
// tokens detector, this only works on very simple spam messages that typically
// involve repeating a word or letter hundreds of times.
type repeatedContentDetector struct {
	threshold float64
}

func (d *repeatedContentDetector) Detect(ctx context.Context, r io.Reader) (bool, error) {
	ratio, err := d.getRatio(r)
	if err != nil {
		return false, fault.Wrap(err, fctx.With(ctx))
	}

	isSpam := ratio < d.threshold

	return isSpam, nil
}

func (d *repeatedContentDetector) getRatio(r io.Reader) (float64, error) {
	treader := transform.NewReader(r, preprocessor{})
	writer := countWriter{}

	compressor := gzip.NewWriter(&writer)
	defer compressor.Close()

	originalSize, err := io.Copy(compressor, treader)
	if err != nil {
		return 0, err
	}

	compressedSize := writer.count

	return float64(compressedSize) / float64(originalSize), nil
}

// countWriter discards the bytes but counts the throughput
type countWriter struct {
	count int
}

func (b *countWriter) Write(p []byte) (int, error) {
	n := len(p)
	b.count += n
	return n, nil
}

type preprocessor struct{}

func (t preprocessor) Transform(dst, src []byte, atEOF bool) (nDst, nSrc int, err error) {
	for nSrc < len(src) {
		switch src[nSrc] {
		// strip out all whitespace and newlines, mostly because posts with
		// source code will have a lot of repeated whitespace and newlines.
		case '\n', '\r', ' ':

		default:
			if nDst < len(dst) {
				dst[nDst] = src[nSrc]
				nDst++
			}
		}
		nSrc++
	}
	return nDst, nSrc, nil
}

func (t preprocessor) Reset() {}
