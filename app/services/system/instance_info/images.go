package instance_info

import (
	"image"
	"math"

	"github.com/Southclaws/storyden/app/resources/asset"
	"github.com/Southclaws/storyden/app/resources/settings"
	"github.com/Southclaws/storyden/internal/infrastructure/imagery"
)

type MediaServiceInfo struct {
	ImageMemoryBudget    int64
	ImageEstimatedMemory int64
}

func estimateImageMemory(limits settings.ImageLimits) int64 {
	side := max(1, int(math.Ceil(math.Sqrt(float64(limits.MaxPixels)))))
	source := &imagery.Source{
		Config: image.Config{Width: side, Height: side},
		Size:   limits.MaxSize,
	}
	preset, _ := asset.Preset(asset.Large)

	estimate, err := imagery.Estimate(source, preset.Width, preset.Height, imagery.FitMode)
	if err != nil {
		return math.MaxInt64
	}

	return estimate.Total()
}
