package asset

import "fmt"

type Variant string

const (
	Original        Variant = "original"
	Small           Variant = "s"
	Medium          Variant = "m"
	Large           Variant = "l"
	ThumbnailSmall  Variant = "thumb_s"
	ThumbnailMedium Variant = "thumb_m"
	ThumbnailLarge  Variant = "thumb_l"
	ImageRecipe             = 1
)

type ImagePreset struct {
	Name          Variant
	Width, Height int
	Square, PNG   bool
	Bit           int64
}

var ImagePresets = []ImagePreset{
	{Name: Small, Width: 320, Height: 320, Bit: 1 << 0},
	{Name: Medium, Width: 960, Height: 960, Bit: 1 << 1},
	{Name: Large, Width: 1920, Height: 1920, Bit: 1 << 2},

	{Name: ThumbnailSmall, Width: 64, Height: 64, Square: true, Bit: 1 << 3},
	{Name: ThumbnailMedium, Width: 128, Height: 128, Square: true, Bit: 1 << 4},
	{Name: ThumbnailLarge, Width: 256, Height: 256, Square: true, Bit: 1 << 5},

	{Name: "icon_32", Width: 32, Height: 32, Square: true, PNG: true, Bit: 1 << 6},
	{Name: "icon_120", Width: 120, Height: 120, Square: true, PNG: true, Bit: 1 << 7},
	{Name: "icon_152", Width: 152, Height: 152, Square: true, PNG: true, Bit: 1 << 8},
	{Name: "icon_167", Width: 167, Height: 167, Square: true, PNG: true, Bit: 1 << 9},
	{Name: "icon_180", Width: 180, Height: 180, Square: true, PNG: true, Bit: 1 << 10},
	{Name: "icon_192", Width: 192, Height: 192, Square: true, PNG: true, Bit: 1 << 13},
	{Name: "icon_512", Width: 512, Height: 512, Square: true, PNG: true, Bit: 1 << 11},
	{Name: "banner", Width: 1200, Height: 630, PNG: true, Bit: 1 << 12},
}

const (
	EagerImagePresets  int64 = 1 | 8 | 16 | 32
	IconImagePresets   int64 = 64 | 128 | 256 | 512 | 1024 | 2048 | 8192
	BannerImagePresets int64 = 4096
)

func Preset(v Variant) (ImagePreset, bool) {
	for _, p := range ImagePresets {
		if p.Name == v {
			return p, true
		}
	}
	return ImagePreset{}, false
}

func ParseVariant(s string) (Variant, error) {
	if s == "" || s == string(Original) {
		return Original, nil
	}

	v := Variant(s)
	p, ok := Preset(v)
	if !ok || p.PNG {
		return "", fmt.Errorf("unknown image variant %q", s)
	}

	return v, nil
}

func ImageEligible(mimeType string) bool {
	switch mimeType {
	case "image/jpeg", "image/png", "image/webp":
		return true
	default:
		return false
	}
}
