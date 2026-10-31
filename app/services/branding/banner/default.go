package banner

import _ "embed"

// A string keeps the immutable PNG in read-only storage; strings.NewReader
// reads its binary bytes without copying the image.
//
//go:embed storyden-mountains-1200x630.png
var defaultBanner string
