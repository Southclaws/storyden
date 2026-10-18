package asset

type ThemeAssetKind string

const (
	ThemeAssetKindStylesheet ThemeAssetKind = "stylesheet"
	ThemeAssetKindScript     ThemeAssetKind = "script"
)

type ThemeMetadata struct {
	Purpose   string         `json:"purpose"`
	Kind      ThemeAssetKind `json:"theme_kind"`
	Integrity string         `json:"theme_integrity"`
}

func (k ThemeAssetKind) MIME() string {
	if k == ThemeAssetKindStylesheet {
		return "text/css"
	}
	return "application/javascript"
}

func (a Asset) ThemeMetadata() (*ThemeMetadata, error) {
	return decodeMetadata[ThemeMetadata](a)
}
