package imagery

import (
	"image"
	"log/slog"
	"math"
)

const (
	sourcePixelBytes        = 16
	outputPixelBytes        = 24
	codecOverhead           = 24 << 20
	resizePixelBytes        = 4 * 8
	contributionBytes       = 16
	contributionSourceBytes = 4 * contributionBytes
	contributionOutputBytes = 8 * contributionBytes
)

type MemoryEstimate struct {
	Source        int64 `json:"source"`
	Output        int64 `json:"output"`
	Encoded       int64 `json:"encoded"`
	Resize        int64 `json:"resize"`
	Contributions int64 `json:"contributions"`
	Codec         int64 `json:"codec"`
	Safety        int64 `json:"safety"`
}

func (e MemoryEstimate) Total() int64 {
	return e.Source + e.Output + e.Encoded + e.Resize + e.Contributions + e.Codec + e.Safety
}

func (e MemoryEstimate) LogValue() slog.Value {
	return slog.GroupValue(
		slog.Int64("source", e.Source),
		slog.Int64("output", e.Output),
		slog.Int64("encoded", e.Encoded),
		slog.Int64("resize", e.Resize),
		slog.Int64("contributions", e.Contributions),
		slog.Int64("codec", e.Codec),
		slog.Int64("safety", e.Safety),
		slog.Int64("total", e.Total()),
	)
}

func Estimate(s *Source, width, height int, mode ResizeMode) (MemoryEstimate, error) {
	if s == nil || s.Config.Width <= 0 || s.Config.Height <= 0 || s.Size < 0 {
		return MemoryEstimate{}, ErrTooLarge
	}
	if width <= 0 || height <= 0 {
		return MemoryEstimate{}, ErrTooLarge
	}
	if mode != FitMode && mode != ThumbnailMode && mode != FillMode {
		return MemoryEstimate{}, ErrUnsupported
	}

	sourceBytes, err := memoryProduct(int64(s.Config.Width), int64(s.Config.Height), sourcePixelBytes)
	if err != nil {
		return MemoryEstimate{}, err
	}
	encodedBytes, err := memoryProduct(s.Size, 2)
	if err != nil {
		return MemoryEstimate{}, err
	}

	base := MemoryEstimate{
		Source:  sourceBytes,
		Encoded: encodedBytes,
		Codec:   codecOverhead,
	}

	// Orientation is applied after admission. Reserve for the more expensive
	// orientation without combining buffers from two different transforms.
	var largest MemoryEstimate
	for _, bounds := range []image.Rectangle{
		image.Rect(0, 0, s.Config.Width, s.Config.Height),
		image.Rect(0, 0, s.Config.Height, s.Config.Width),
	} {
		estimate, err := estimateResize(base, bounds, width, height, mode)
		if err != nil {
			return MemoryEstimate{}, err
		}
		if estimate.Total() > largest.Total() {
			largest = estimate
		}
	}

	return largest, nil
}

func estimateResize(
	base MemoryEstimate,
	bounds image.Rectangle,
	width, height int,
	mode ResizeMode,
) (MemoryEstimate, error) {
	crop, destination := resizeGeometry(bounds, width, height, mode)
	if crop.Empty() || destination.Empty() {
		return MemoryEstimate{}, ErrTooLarge
	}

	if max(crop.Dx(), crop.Dy(), destination.Dx(), destination.Dy()) > math.MaxInt32 {
		return MemoryEstimate{}, ErrTooLarge
	}

	output, err := memoryProduct(int64(destination.Dx()), int64(destination.Dy()), outputPixelBytes)
	if err != nil {
		return MemoryEstimate{}, err
	}

	// x/image/draw v0.47.0 scales horizontally first and allocates one
	// [4]float64 per destination column and cropped source row.
	workspacePixels, err := memoryProduct(int64(destination.Dx()), int64(crop.Dy()))
	if err != nil || workspacePixels > math.MaxInt32 {
		return MemoryEstimate{}, ErrTooLarge
	}
	resize, err := memoryProduct(workspacePixels, resizePixelBytes)
	if err != nil {
		return MemoryEstimate{}, err
	}

	sourceAxes, err := memorySum(int64(crop.Dx()), int64(crop.Dy()))
	if err != nil {
		return MemoryEstimate{}, err
	}
	sourceContributions, err := memoryProduct(sourceAxes, contributionSourceBytes)
	if err != nil {
		return MemoryEstimate{}, err
	}
	outputAxes, err := memorySum(int64(destination.Dx()), int64(destination.Dy()))
	if err != nil {
		return MemoryEstimate{}, err
	}
	outputContributions, err := memoryProduct(outputAxes, contributionOutputBytes)
	if err != nil {
		return MemoryEstimate{}, err
	}
	contributions, err := memorySum(sourceContributions, outputContributions)
	if err != nil || contributions/contributionBytes > math.MaxInt32 {
		return MemoryEstimate{}, ErrTooLarge
	}

	base.Output = output
	base.Resize = resize
	base.Contributions = contributions
	working, err := memorySum(base.Source, base.Output, base.Encoded, base.Resize, base.Contributions, base.Codec)
	if err != nil {
		return MemoryEstimate{}, err
	}

	// Keep the measured allowance for codec buffer growth and intermediates
	// retained between sequential variants until the corpus supports reducing it.
	base.Safety = working / 2
	if _, err := memorySum(working, base.Safety); err != nil {
		return MemoryEstimate{}, err
	}

	return base, nil
}

func memoryProduct(values ...int64) (int64, error) {
	product := int64(1)
	for _, value := range values {
		if value < 0 || value > 0 && product > math.MaxInt64/value {
			return 0, ErrTooLarge
		}
		product *= value
	}

	return product, nil
}

func memorySum(values ...int64) (int64, error) {
	var sum int64
	for _, value := range values {
		if value < 0 || sum > math.MaxInt64-value {
			return 0, ErrTooLarge
		}
		sum += value
	}

	return sum, nil
}
