package asset

type ImageMetadata struct {
	Width  int     `json:"width,omitempty"`
	Height int     `json:"height,omitempty"`
	Preset Variant `json:"image_preset,omitempty"`
	Recipe int     `json:"image_recipe,omitempty"`
}

func (a Asset) ImageMetadata() (*ImageMetadata, error) {
	return decodeMetadata[ImageMetadata](a)
}
